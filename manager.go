package sago

import (
	"context"
	"fmt"
	"strings"

	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagolog"
	"git.coryptex.com/lib/sago/sagomsg"
	"git.coryptex.com/lib/sago/zeebe"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/pkg/errors"
)

type SagaManager interface {
	create(uniqueId string, data SagaData, extVars map[string]interface{}) error
	subscribeToReplyChannel()
	deployProcess(path string) error
	registerJobWorkers() error
}

func NewSagaManager(
	saga Saga,
	zb zbc.Client,
	sagaInstanceRepository SagaInstanceRepository,
	commandProducer sagocmd.CommandProducer,
	messageConsumer sagomsg.MessageConsumer,
	sagaCommandProducer *SagaCommandProducer,
) SagaManager {
	return &sagaManager{
		saga:                   saga,
		zb:                     zb,
		sagaInstanceRepository: sagaInstanceRepository,
		commandProducer:        commandProducer,
		messageConsumer:        messageConsumer,
		sagaCommandProducer:    sagaCommandProducer,
	}
}

type sagaManager struct {
	// TODO: sagaLockManager -> lack of isolation
	saga                   Saga
	zb                     zbc.Client
	sagaInstanceRepository SagaInstanceRepository
	commandProducer        sagocmd.CommandProducer
	messageConsumer        sagomsg.MessageConsumer
	sagaCommandProducer    *SagaCommandProducer
}

func (sm *sagaManager) create(uniqueId string, data SagaData, vars map[string]interface{}) error {
	const op string = "sago.manager.create"

	dataSerd, err := data.Marshal()
	if err != nil {
		return errors.Wrapf(err, "%s: %s:%s\n%+v", op, sm.getSagaType(), uniqueId, data)
	}

	sagolog.Log(sagolog.DEBUG, fmt.Sprintf(
		"%s: creating %s:%s saga with data: %s",
		op, sm.getSagaType(), uniqueId, string(dataSerd),
	))

	sagaInstance := NewSagaInstance(
		uniqueId, sm.getSagaType(), "started", "",
		dataSerd, map[string]string{},
	)

	sagaID, err := sm.sagaInstanceRepository.Save(*sagaInstance)
	if err != nil {
		return errors.Wrapf(err, "failed to store sagaInstance of %s saga\n", sm.getSagaType())
	}
	// do not need to set saga id because it's same as uniqueID
	// sagaInstance.SetID(sagaID)

	// go sm.saga.OnStarting(sagaID, dataSerd)

	// prepare zeebe flow initial variables
	if vars == nil {
		vars = make(map[string]interface{})
	}
	removeReservedVariableKeys(vars)
	vars[ZB_SAGO_KEY] = buildSagoKey(sm.getSagaType(), sagaID)
	vars[ZB_SAGA_ID] = sagaID
	vars[ZB_SAGA_TYPE] = sm.getSagaType()

	req, err := sm.zb.NewPublishMessageCommand().
		MessageName(sm.getSagaType()).
		CorrelationKey(buildSagoKey(sm.getSagaType(), sagaID)).
		VariablesFromMap(vars)
	if err != nil {
		return errors.Wrapf(err,
			"failed to create start saga message %s:%s\n", sm.getSagaType(), sagaID)
	}

	_, err = req.Send(context.Background())
	if err != nil {
		return errors.Wrapf(err,
			"failed to send start saga message %s:%s\n", sm.getSagaType(), sagaID)
	}

	sagolog.Log(sagolog.INFO, fmt.Sprintf(
		"saga instance %s:%s started",
		sm.getSagaType(), sagaID,
	))
	return nil
}

func (sm *sagaManager) deployProcess(path string) error {
	err := zeebe.DeployProcess(sm.zb, path)
	if err != nil {
		return errors.Wrapf(err,
			"failed to deploy process in %s for %s\n",
			path, sm.getSagaType(),
		)
	}
	sagolog.Log(sagolog.DEBUG, fmt.Sprintf("%s process deployed", path))
	return nil
}

func (sm *sagaManager) registerJobWorkers() error {
	def, err := sm.getStateDefinition()
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to register job workers for %s saga\n",
			sm.getSagaType(),
		)
	}
	sagolog.Log(sagolog.DEBUG, "start registering jobs")
	for step := range def.stepsName() {
		go sm.zb.NewJobWorker().JobType(step).Handler(sm.handleJob).Open()
		sagolog.Log(sagolog.DEBUG, fmt.Sprint("job registered for", step))
	}
	return nil
}

func (sm *sagaManager) handleJob(client worker.JobClient, job entities.Job) {
	const op string = "sago.manager.handleJob"

	jobKey, jobType := job.GetKey(), job.GetType()

	sagolog.Log(sagolog.DEBUG,
		fmt.Sprintf("%s: handling job %s:%d", op, jobType, jobKey),
	)

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		msg := fmt.Sprintf(
			"%s: failing job %s:%d, reason: failed to get variables of job\n%v",
			op, jobType, jobKey, err,
		)

		sagolog.Log(sagolog.ERROR, msg)
		zeebe.FailJob(client, job, msg)
		return
	}

	// fetch sagaInstance from DB
	sagaID := variables[ZB_SAGA_ID].(string)
	instance, err := sm.sagaInstanceRepository.Find(sm.getSagaType(), sagaID)
	if err != nil {
		msg := fmt.Sprintf(
			"%s: failing job %s:%d, reason: failed to get %s:%s saga instance\n%v",
			op, jobType, jobKey, sm.getSagaType(), sagaID, err,
		)

		sagolog.Log(sagolog.ERROR, msg)
		zeebe.FailJob(client, job, msg)
		return
	}
	// err has been checked in registerJobWorkers step
	def, _ := sm.getStateDefinition()

	// get saga step related to this job
	step := def.step(jobType)
	if step == nil {
		msg := fmt.Sprintf(
			"%s: failing job %s:%d, reason: saga step is nil",
			op, jobType, jobKey,
		)

		sagolog.Log(sagolog.ERROR, msg)
		zeebe.FailJob(client, job, msg)
		return
	}

	// copy variables and delete reserved keys
	vars := make(map[string]interface{})
	for k, v := range variables {
		vars[k] = v
	}
	removeReservedVariableKeys(vars)

	// get command of this step and send it to related participant
	cmd, err := step.Command(instance.SerializedSagaData(), vars)
	if err != nil {
		msg := fmt.Sprintf("%s: failed to get step's command\n%v", op, err)

		sagolog.Log(sagolog.ERROR, msg)
		zeebe.FailJob(client, job, msg)
		return
	}
	lastReqID, err := sm.sagaCommandProducer.sendCommands(
		sm.getSagaType(), sagaID,
		sm.makeSagaReplyChannel(),
		[]sagocmd.Command{cmd},
	)
	if err != nil {
		msg := fmt.Sprintf(
			"%s: failing job %s:%d, reason: failed to send commands\n%v",
			op, jobType, jobKey, err,
		)

		sagolog.Log(sagolog.ERROR, msg)
		zeebe.FailJob(client, job, msg)
		return
	}

	// set lastRequestID and state then update saga instance in DB
	instance.SetLastRequestID(lastReqID)
	instance.SetStateName(jobType)
	err = sm.sagaInstanceRepository.Update(*instance)
	if err != nil {
		msg := fmt.Sprintf(
			"%s: failing job %s:%d, reason: failed to update %s:%s saga instance\n%v",
			op, jobType, jobKey, sm.getSagaType(), sagaID, err,
		)

		sagolog.Log(sagolog.ERROR, msg)
		zeebe.FailJob(client, job, msg)
		return
	}

	// send complete command to zb
	_, err = sm.zb.NewCompleteJobCommand().
		JobKey(jobKey).Send(context.Background())

	if err != nil {
		sagolog.Log(sagolog.ERROR,
			fmt.Sprintf(
				"%s: failing job %s:%d, reason: failed to send complete job command\n%v",
				op, jobType, jobKey, err,
			),
		)
		return
	}

	sagolog.Log(sagolog.DEBUG,
		fmt.Sprintf("%s: job %s:%d completed", op, jobType, jobKey),
	)
}

func (sm *sagaManager) subscribeToReplyChannel() {
	sm.messageConsumer.Subscribe(
		fmt.Sprintf("%s-consumer", sm.saga.SagaType()),
		[]string{sm.makeSagaReplyChannel()},
		sm.handleMessage,
	)
}

func (sm *sagaManager) handleMessage(msg sagomsg.Message) error {
	const op string = "sago.manager.handleMessage"

	sagolog.Log(sagolog.DEBUG, fmt.Sprintf("%s: handling message:\n%+v", op, msg))

	if msg.HasHeader(REPLY_SAGA_ID) {
		err := sm.handleReply(msg)
		if err != nil {
			return errors.Wrap(err, op)
		}
		return nil
	}
	return errors.Errorf("%s: doesn't know what to do with:\n%+v", op, msg)
}

func (sm *sagaManager) handleReply(msg sagomsg.Message) error {
	const op string = "sago.manager.handleReply"

	if !sm.isReplyForThisSagaType(msg) {
		return errors.Errorf("%s: reply %v is not for %s saga type.", op, msg, sm.getSagaType())
	}

	sagolog.Log(sagolog.DEBUG, fmt.Sprintf("%s: handling reply\n%+v", op, msg))
	// header existence checked before
	sagaID, _ := msg.RequiredHeader(REPLY_SAGA_ID)
	sagaType, _ := msg.RequiredHeader(REPLY_SAGA_TYPE)
	replyCmdName, err := msg.RequiredHeader(sagocmd.REPLY_TYPE)
	if err != nil {
		return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
	}

	sagaInstance, err := sm.sagaInstanceRepository.Find(sagaType, sagaID)
	if err != nil {
		return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
	}

	sagolog.Log(sagolog.DEBUG, fmt.Sprintf(
		"current state of %s:%s saga is %s",
		sagaType, sagaID, sagaInstance.StateName(),
	))

	sagaDefinition, err := sm.getStateDefinition()
	if err != nil {
		return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
	}

	stepName := strings.TrimSuffix(replyCmdName, "Reply")
	step := sagaDefinition.step(stepName)
	if step == nil {
		return errors.Errorf("%s.%s:%s: there is no %s step defined",
			op, sagaType, sagaID, stepName,
		)
	}

	result := "failed"
	isSuccessful := step.IsSuccessfulReply(msg)
	if isSuccessful {
		result = "success"
	}

	sagolog.Log(sagolog.DEBUG,
		fmt.Sprintf("%s %s.%s: %s result was: %s",
			op, sagaType, sagaID, replyCmdName, result,
		),
	)

	// call business logic callback to handle reply
	handler := step.GetReplyHandler(msg)
	if handler != nil {
		sagaData, err := handler(sagaInstance.SerializedSagaData(), msg.Payload(), isSuccessful)
		if err != nil {
			return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
		}
		data, err := sagaData.Marshal()
		if err != nil {
			return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
		}
		sagaInstance.SetSerializedSagaData(data)
		err = sm.sagaInstanceRepository.Update(*sagaInstance)
		if err != nil {
			return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
		}
	}

	err = zeebe.PublishMessage(
		context.Background(),
		sm.zb, replyCmdName, buildSagoKey(sagaType, sagaID),
		map[string]interface{}{
			stepName + "Result": result,
		},
	)

	if err != nil {
		return errors.Wrapf(err, "%s %s:%s", op, sagaType, sagaID)
	}

	sagolog.Log(sagolog.DEBUG,
		fmt.Sprintf("%s %s.%s: %s handled",
			op, sagaType, sagaID, replyCmdName,
		),
	)
	return nil
}

func (sm *sagaManager) getSagaType() string {
	return sm.saga.SagaType()
}

func (sm *sagaManager) makeSagaReplyChannel() string {
	return sm.getSagaType() + "-reply"
}

func (sm *sagaManager) getStateDefinition() (SagaDefinition, error) {
	def := sm.saga.SagaDefinition()
	if def == nil {
		return nil, errors.New("saga definition should not be nil")
	}
	return def, nil
}

func (sm *sagaManager) isReplyForThisSagaType(msg sagomsg.Message) bool {
	return strings.EqualFold(msg.Header(REPLY_SAGA_TYPE), sm.getSagaType())
}

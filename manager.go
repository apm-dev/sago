package sago

import (
	"context"
	"fmt"
	"log"
	"strings"

	"git.coryptex.com/lib/sago/sagocmd"
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
	dataSerd := data.Marshal()

	sagaInstance := NewSagaInstance(
		uniqueId, sm.getSagaType(), "started", "",
		dataSerd, map[string]string{},
	)

	sagaID, err := sm.sagaInstanceRepository.Save(*sagaInstance)
	if err != nil {
		return errors.Wrapf(err, "failed to store sagaInstance of %s saga\n", sm.getSagaType())
	}
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
	log.Println("start registering jobs")
	for step := range def.stepsName() {
		go sm.zb.NewJobWorker().JobType(step).Handler(sm.handleJob).Open()
		log.Println("job registered for", step)
	}
	return nil
}

func (sm *sagaManager) handleJob(client worker.JobClient, job entities.Job) {
	jobKey, jobType := job.GetKey(), job.GetType()
	log.Printf("handle job %s:%d\n", jobType, jobKey)

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job,
			fmt.Sprintf("failed to get variables of job %s:%d\nerr: %v\n", jobType, jobKey, err))
		return
	}

	// fetch sagaInstance from DB
	sagaID := variables[ZB_SAGA_ID].(string)
	instance, err := sm.sagaInstanceRepository.Find(sm.getSagaType(), sagaID)
	if err != nil {
		zeebe.FailJob(client, job,
			fmt.Sprintf("failed to fetch instance from db of saga %s:%s to handle job %s:%d\nerr: %v",
				sm.getSagaType(), sagaID, jobType, jobKey, err))
		return
	}
	// err has been checked in registerJobWorkers step
	def, _ := sm.getStateDefinition()

	// get saga step related to this job
	step := def.step(jobType)
	if step == nil {
		zeebe.FailJob(client, job,
			fmt.Sprintf("there is no step for %s:%d job\n", jobType, jobKey))
		return
	}

	// copy variables and delete reserved keys
	vars := make(map[string]interface{})
	for k, v := range variables {
		vars[k] = v
	}
	removeReservedVariableKeys(vars)

	// get command of this step and send it to related participant
	cmd := step.Command(instance.SerializedSagaData(), vars)
	lastReqID, err := sm.sagaCommandProducer.sendCommands(
		sm.getSagaType(), sagaID,
		sm.makeSagaReplyChannel(),
		[]sagocmd.Command{cmd},
	)
	if err != nil {
		zeebe.FailJob(client, job,
			fmt.Sprintf("failed to handle job %s:%d\nerr: %v\n", jobType, jobKey, err))
		return
	}

	// set lastRequestID and state then update saga instance in DB
	instance.SetLastRequestID(lastReqID)
	instance.SetStateName(jobType)
	err = sm.sagaInstanceRepository.Update(*instance)
	if err != nil {
		zeebe.FailJob(client, job,
			fmt.Sprintf("failed to update sagaInstance of %s:%s saga for %s:%d job\nerr: %v\n",
				sm.getSagaType(), sagaID, jobType, jobKey, err))
		return
	}

	// send complete command to zb
	_, err = sm.zb.NewCompleteJobCommand().
		JobKey(jobKey).Send(context.Background())

	if err != nil {
		log.Printf("failed to send %s:%d job complete request\nerr: %v\n", jobType, jobKey, err)
		return
	}
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
	// TODO log
	log.Printf("handle message invoked %+v", msg)
	if msg.HasHeader(REPLY_SAGA_ID) {
		err := sm.handleReply(msg)
		if err != nil {
			log.Println(errors.Wrap(err, op))
			return err
		}
		return nil
	}
	// TODO log
	err := errors.Errorf("%s: handleMessage doesn't know what to do with: %+v", op, msg)
	log.Println(err)
	return err
}

func (sm *sagaManager) handleReply(msg sagomsg.Message) error {
	const op string = "sago.manager.handleReply"
	// TODO implement
	if !sm.isReplyForThisSagaType(msg) {
		return errors.Errorf("%s: reply %v is not for %s saga type.", op, msg, sm.getSagaType())
	}
	// TODO log
	log.Printf("%s: Handle reply %+v", op, msg)
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

	log.Printf("Current state of %s:%s saga is %s", sagaType, sagaID, sagaInstance.StateName())

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

	// call business logic callback to handle reply
	handler := step.GetReplyHandler(msg)
	if handler != nil {
		sagaData, err := handler(sagaInstance.SerializedSagaData(), msg.Payload(), isSuccessful)
		if err != nil {
			return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
		}
		sagaInstance.SetSerializedSagaData(sagaData.Marshal())
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
		return errors.Wrapf(err, "%s.%s:%s", op, sagaType, sagaID)
	}
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

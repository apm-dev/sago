package sago

import (
	"github.com/apm-dev/sago/sagocmd"
	"github.com/apm-dev/sago/sagomsg"
	"github.com/apm-dev/sago/zeebe"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/pkg/errors"
)

type SagaManager interface {
	create(data SagaData) error
	subscribeToReplyChannel()
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

func (sm *sagaManager) create(data SagaData) error {
	dataSerd := data.Marshal()

	sagaInstance := NewSagaInstance(
		"", sm.getSagaType(), "started", "",
		dataSerd, map[string]string{},
	)

	sagaID, err := sm.sagaInstanceRepository.Save(*sagaInstance)
	if err != nil {
		return errors.Wrapf(err, "failed to store sagaInstance of %s saga\n", sm.getSagaType())
	}
	sagaInstance.SetID(sagaID)

	go sm.saga.OnStarting(sagaID, dataSerd)

	req, err := sm.zb.NewPublishMessageCommand().
		MessageName(sm.getSagaType()).
		CorrelationKey(sagaID).
		VariablesFromMap(map[string]interface{}{
			ZB_SAGA_ID:   sagaID,
			ZB_SAGA_TYPE: sm.getSagaType(),
		})
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

	// get command of this step and send it to related participant
	cmd := step.Command(instance.SerializedSagaData())
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
	_, err = sm.zb.NewCompleteJobCommand().JobKey(jobKey).Send(context.Background())
	/* if err != nil {
		log.Printf("failed to create %s:%s job complete command\nerr: %v\n", jobType, jobKey, err)
		return
	}
	_, err = req.Send(context.Background()) */
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

func (sm *sagaManager) handleMessage(msg sagomsg.Message) {
	// TODO log
	log.Printf("handle message invoked %+v", msg)
	if msg.HasHeader(REPLY_SAGA_ID) {
		sm.handleReply(msg)
	} else {
		// TODO log
		log.Printf("handleMessage doesn't know what to do with: %+v", msg)
	}
}

func (sm *sagaManager) handleReply(msg sagomsg.Message) {
	// TODO implement
	if !sm.isReplyForThisSagaType(msg) {
		return
	}
	// TODO log
	log.Printf("Handle reply %+v", msg)
	// header existence checked before
	sagaID, _ := msg.RequiredHeader(REPLY_SAGA_ID)
	sagaType, _ := msg.RequiredHeader(REPLY_SAGA_TYPE)
	replyCmdName, err := msg.RequiredHeader(sagocmd.REPLY_TYPE)
	if err != nil {
		log.Printf("handleReply doesn't know what to do with %+v msg without %s header\n",
			msg, sagocmd.REPLY_TYPE)
		return
	}

	sagaInstance, err := sm.sagaInstanceRepository.Find(sagaType, sagaID)
	if err != nil {
		log.Printf("failed to get sagaInstance of %s:%s saga\nerr: %v\n",
			sagaID, sagaType, err)
		return
	}

	log.Printf("Current state of %s:%s saga is %s", sagaType, sagaID, sagaInstance.StateName())

	sagaDefinition, err := sm.getStateDefinition()
	if err != nil {
		log.Printf(
			"failed to get definition of %s:%s saga\nerr: %v\n",
			sagaID, sagaType, err,
		)
		return
	}

	stepName := strings.TrimSuffix(replyCmdName, "Reply")
	step := sagaDefinition.step(stepName)
	if step == nil {
		log.Printf("there is no %s step defined for %s:%s saga\n", stepName, sagaType, sagaID)
		return
	}

	result := "failed"
	isSuccessful := step.IsSuccessfulReply(msg)
	if isSuccessful {
		result = "success"
	}

	// call business logic callback to handle reply
	handler := step.GetReplyHandler(msg)
	if handler != nil {
		sagaData := handler(sagaInstance.SerializedSagaData(), msg.Payload(), isSuccessful)
		sagaInstance.SetSerializedSagaData(sagaData.Marshal())
		err = sm.sagaInstanceRepository.Update(*sagaInstance)
		if err != nil {
			log.Printf("failed to update sagaInstance of %s:%s saga\nerr: %v\n",
				sagaID, sagaType, err)
			return
		}
	}

	err = zeebe.PublishMessage(
		context.Background(),
		sm.zb, replyCmdName, sagaID,
		map[string]interface{}{
			stepName + "Result": result,
		},
	)

	if err != nil {
		log.Printf(
			"failed to publish zb message %s of saga %s:%s\nerr: %v\n",
			replyCmdName, sagaType, sagaID, err,
		)
		return
	}
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
		return nil, errors.New("state machine should not be nil")
	}
	return def, nil
}

func (sm *sagaManager) isReplyForThisSagaType(msg sagomsg.Message) bool {
	return strings.EqualFold(msg.Header(REPLY_SAGA_TYPE), sm.getSagaType())
}

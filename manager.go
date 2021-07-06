package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type SagaManager interface {
	Create(data SagaData) (*SagaInstance, error)
	SubscribeToReplyChannel()
}

func NewSagaManager(
	saga Saga,
	sagaInstanceRepository SagaInstanceRepository,
	commandProducer commands.CommandProducer,
	messageConsumer messaging.MessageConsumer,
	sagaCommandProducer SagaCommandProducer,
) SagaManager {
	return &sagaManager{
		saga:                   saga,
		sagaInstanceRepository: sagaInstanceRepository,
		commandProducer:        commandProducer,
		messageConsumer:        messageConsumer,
		sagaCommandProducer:    sagaCommandProducer,
	}
}

type sagaManager struct {
	// TODO: sagaLockManager -> lack of isolation
	saga                   Saga
	sagaInstanceRepository SagaInstanceRepository
	commandProducer        commands.CommandProducer
	messageConsumer        messaging.MessageConsumer
	sagaCommandProducer    SagaCommandProducer
}

func (sm *sagaManager) Create(data SagaData) (*SagaInstance, error) {
	// dataserd, err := serializeSagaData(data)
	// if err != nil {
	// 	return nil, err
	// }

	dataSerd := data.Marshal()

	sagaInstance := NewSagaInstance(
		"", sm.getSagaType(), "????", "",
		dataSerd, map[string]string{},
	)

	sagaID, err := sm.sagaInstanceRepository.Save(*sagaInstance)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't store sagaInstance")
	}
	sagaInstance.SetID(sagaID)

	def, err := sm.getStateDefinition()
	if err != nil {
		return nil, err
	}

	// serData, err := proto.Marshal(data)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Couldn't marshal sagaData")
	// }

	actions := def.Start(dataSerd)

	sm.processActions(sagaID, sagaInstance, dataSerd, actions)

	return sagaInstance, nil
}

func (sm *sagaManager) SubscribeToReplyChannel() {
	sm.messageConsumer.Subscribe(
		fmt.Sprintf("%s-consumer", sm.saga.SagaType()),
		[]string{sm.makeSagaReplyChannel()},
		sm.handleMessage,
	)
}

func (sm *sagaManager) handleMessage(msg messaging.Message) {
	// TODO log
	log.Printf("handle message invoked %+v", msg)
	if msg.HasHeader(REPLY_SAGA_ID) {
		sm.handleReply(msg)
	} else {
		// TODO log
		log.Printf("Handle message doesn't know what to do with: %+v", msg)
	}
}

func (sm *sagaManager) handleReply(msg messaging.Message) {
	// TODO implement
	if !sm.isReplyForThisSagaType(msg) {
		return
	}
	// TODO log
	log.Printf("Handle reply %+v", msg)

	// header existence checked before
	sagaID, _ := msg.RequiredHeader(REPLY_SAGA_ID)
	sagaType, _ := msg.RequiredHeader(REPLY_SAGA_TYPE)

	sagaInstance, err := sm.sagaInstanceRepository.Find(sagaType, sagaID)
	if err != nil {
		log.Printf("There is no sagaInstance for id: %s, type: %s", sagaID, sagaType)
		return
	}

	currentState := sagaInstance.StateName()

	log.Printf("Current state %s", currentState)

	sagaDefinition, err := sm.getStateDefinition()
	if err != nil {
		log.Printf(
			"Error while getting definition for saga id:%s, type:%s \n Error: %v",
			sagaID, sagaType, err,
		)
		return
	}

	actions, err := sagaDefinition.HandleReply(
		currentState,
		sagaInstance.SerializedSagaData(),
		msg,
	)
	if err != nil {
		log.Printf("Couldn't handle reply, err: %v", err)
		return
	}

	err = sm.processActions(
		sagaID, sagaInstance,
		sagaInstance.SerializedSagaData(),
		actions,
	)

	if err != nil {
		log.Printf("Couldn't process actions, err: %v", err)
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
		return nil, errors.New("state machine can not be nil")
	}
	return def, nil
}

func (sm *sagaManager) processActions(sagaID string, sagaInstance *SagaInstance, sagaData []byte, actions *SagaActions) error {

	lastRequestID := sm.sagaCommandProducer.sendCommands(
		sm.getSagaType(),
		sagaID,
		sm.makeSagaReplyChannel(),
		actions.Commands(),
	)

	sagaInstance.SetLastRequestID(lastRequestID)

	sm.updateState(sagaInstance, actions)

	if updatedSagaData := actions.UpdatedSagaData(); updatedSagaData != nil {
		/* serd, err := serializeSagaData(updatedSagaData)
		if err != nil {
			return err
		} */
		sagaInstance.SetSerializedSagaData(updatedSagaData)
	}

	if actions.IsEndState() {
		sm.performEndStateActions(sagaID, sagaInstance, actions.IsCompensating(), sagaData)
	}

	sm.sagaInstanceRepository.Update(*sagaInstance)
	return nil
}

func (sm *sagaManager) updateState(si *SagaInstance, actions *SagaActions) {
	updatedState := actions.UpdatedState()
	if updatedState != "" {
		si.SetStateName(updatedState)
		si.SetEndState(actions.IsEndState())
		si.SetCompensating(actions.IsCompensating())
	}
}

func (sm *sagaManager) performEndStateActions(sagaID string, si *SagaInstance, compensating bool, sagaData []byte) {
	// TODO implement me, this is for releasing(unlock) resources
	if compensating {
		sm.saga.OnSagaRolledBack(sagaID, sagaData)
	} else {
		sm.saga.OnSagaCompletedSuccessfully(sagaID, sagaData)
	}
}

func (sm *sagaManager) isReplyForThisSagaType(msg messaging.Message) bool {
	return strings.EqualFold(msg.Header(REPLY_SAGA_TYPE), sm.getSagaType())
}

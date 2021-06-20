package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type SagaManager interface {
	Create(data interface{}) *SagaInstance
	SubscribeToReplyChannel()
}

type sagaManager struct {
	// TODO: sagaLockManager -> lack of isolation
	saga                   Saga
	sagaInstanceRepository SagaInstanceRepository
	commandProducer        commands.CommandProducer
	messageConsumer        messaging.MessageConsumer
	sagaCommandProducer    SagaCommandProducer
}

func (sm *sagaManager) Create(data proto.Message) (*SagaInstance, error) {
	dataserd, err := serializeSagaData(data)
	if err != nil {
		return nil, err
	}

	sagaInstance := NewSagaInstance(
		"", sm.getSagaType(), "????", "",
		dataserd, map[string]string{},
	)

	sagaID := sm.sagaInstanceRepository.Save(*sagaInstance)
	sagaInstance.SetID(sagaID)

	def, err := sm.getStateDefinition()
	if err != nil {
		return nil, err
	}

	actions := def.Start(data)
	sm.processActions(sagaID, sagaInstance, data, actions)

	return sagaInstance, nil
}

func (sm *sagaManager) SubscribeToReplyChannel() {
	panic("implement me")
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

func (sm *sagaManager) processActions(sagaID string, sagaInstance *SagaInstance, sagaData proto.Message, actions *SagaActions) error {

	lastRequestID := sm.sagaCommandProducer.sendCommands(
		sm.getSagaType(),
		sagaID,
		sm.makeSagaReplyChannel(),
		actions.Commands(),
	)

	sagaInstance.SetLastRequestID(lastRequestID)

	sm.updateState(sagaInstance, actions)

	if updatedSagaData := actions.UpdatedSagaData(); updatedSagaData != nil {
		serd, err := serializeSagaData(updatedSagaData)
		if err != nil {
			// TODO log
			return err
		}
		sagaInstance.SetSerializedSagaData(serd)
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

func (sm *sagaManager) performEndStateActions(sagaID string, si *SagaInstance, compensating bool, sagaData proto.Message) {
	// TODO implement me, this is for releasing(unlock) resources
	if compensating {
		sm.saga.OnSagaRolledBack(sagaID, sagaData)
	} else {
		sm.saga.OnSagaCompletedSuccessfully(sagaID, sagaData)
	}
}

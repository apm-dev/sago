package sago

import (
	"apm-dev/sago/messaging"

	"google.golang.org/protobuf/proto"
)

type SagaDefinition interface {
	Start(sagaData proto.Message) *SagaActions
	HandleReply(currentState string, sagaData []byte, message messaging.Message) *SagaActions
}

type sagaDefinition struct {
	sagaSteps []SagaStep
}

func NewSagaDefinition(stps []SagaStep) SagaDefinition {
	return &sagaDefinition{stps}
}

func (sd *sagaDefinition) Start(sagaData proto.Message) *SagaActions {
	// TODO implement
	panic("implement me")
}

func (sd *sagaDefinition) HandleReply(currentState string, sagaData []byte, message messaging.Message) *SagaActions {
	// TODO implement
	panic("implement me")
}

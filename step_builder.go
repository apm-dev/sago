package sago

import "google.golang.org/protobuf/proto"

type StepBuilder struct {
	parent *SagaDefinitionBuilder
}

func NewStepBuilder(sdb *SagaDefinitionBuilder) *StepBuilder {
	return &StepBuilder{sdb}
}

func (b *StepBuilder) InvokeParticipant(ce CommandEndpoint, cmdProvider func(sagaData []byte, vars map[string]interface{}) (proto.Message, error)) *ParticipantInvocationStepBuilder {
	return NewParticipantInvocationStepBuilder(b.parent).withAction(ce, cmdProvider)
}

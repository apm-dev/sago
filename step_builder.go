package sago

type StepBuilder struct {
	parent *SagaDefinitionBuilder
}

func NewStepBuilder(sdb *SagaDefinitionBuilder) *StepBuilder {
	return &StepBuilder{sdb}
}

func (b *StepBuilder) InvokeParticipant(ce CommandEndpoint, cmdProvider func(sagaData []byte) []byte) *ParticipantInvocationStepBuilder {
	return NewParticipantInvocationStepBuilder(b.parent).WithAction(ce, cmdProvider)
}

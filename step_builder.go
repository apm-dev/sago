package sago

import "apm-dev/sago/commands"

type StepBuilder struct {
	parent *SagaDefinitionBuilder
}

func NewStepBuilder(sdb *SagaDefinitionBuilder) *StepBuilder {
	return &StepBuilder{sdb}
}

func (b *StepBuilder) InvokeParticipant(cmdProvider func() commands.Command) *ParticipantInvocationStepBuilder {
	return NewParticipantInvocationStepBuilder(b.parent).WithAction(cmdProvider)
}

func (b *StepBuilder) WithCompensation(cmdProvider func() commands.Command) *ParticipantInvocationStepBuilder {
	return NewParticipantInvocationStepBuilder(b.parent).WithCompensation(cmdProvider)
}

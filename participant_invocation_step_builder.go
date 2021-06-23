package sago

import (
	"apm-dev/sago/commands"
)

type ParticipantInvocationStepBuilder struct {
	parent                    *SagaDefinitionBuilder
	action                    *ParticipantInvocation
	compensation              *ParticipantInvocation
	actionReplyHandlers       map[string]func([]byte)
	compensationReplyHandlers map[string]func([]byte)
}

func NewParticipantInvocationStepBuilder(parent *SagaDefinitionBuilder) *ParticipantInvocationStepBuilder {
	return &ParticipantInvocationStepBuilder{
		parent:                    parent,
		actionReplyHandlers:       make(map[string]func([]byte)),
		compensationReplyHandlers: make(map[string]func([]byte)),
	}
}

func (b *ParticipantInvocationStepBuilder) WithAction(cmdProvider func() commands.Command) *ParticipantInvocationStepBuilder {
	b.action = NewParticipantInvocation(cmdProvider)
	return b
}

func (b *ParticipantInvocationStepBuilder) WithCompensation(cmdProvider func() commands.Command) *ParticipantInvocationStepBuilder {
	b.compensation = NewParticipantInvocation(cmdProvider)
	return b
}

func (b *ParticipantInvocationStepBuilder) OnReply(name string, handler func([]byte)) *ParticipantInvocationStepBuilder {
	if b.compensation != nil {
		b.compensationReplyHandlers[name] = handler
	} else {
		b.actionReplyHandlers[name] = handler
	}
	return b
}

func (b *ParticipantInvocationStepBuilder) Step() *StepBuilder {
	b.addStep()
	return NewStepBuilder(b.parent)
}

func (b *ParticipantInvocationStepBuilder) Build() SagaDefinition {
	b.addStep()
	return b.parent.Build()
}

func (b *ParticipantInvocationStepBuilder) addStep() {
	b.parent.AddStep(NewParticipantInvocationStep(
		b.action, b.compensation, b.actionReplyHandlers, b.compensationReplyHandlers,
	))
}

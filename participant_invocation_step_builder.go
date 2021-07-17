package sago

type ParticipantInvocationStepBuilder struct {
	parent                    *SagaDefinitionBuilder
	action                    *ParticipantInvocation
	compensation              *ParticipantInvocation
	actionReplyHandlers       map[string]func(data []byte, msg []byte)
	compensationReplyHandlers map[string]func(data []byte, msg []byte)
}

func NewParticipantInvocationStepBuilder(parent *SagaDefinitionBuilder) *ParticipantInvocationStepBuilder {
	return &ParticipantInvocationStepBuilder{
		parent:                    parent,
		actionReplyHandlers:       make(map[string]func(data []byte, msg []byte)),
		compensationReplyHandlers: make(map[string]func(data []byte, msg []byte)),
	}
}

func (b *ParticipantInvocationStepBuilder) WithAction(cmdEndpoint CommandEndpoint, cmdProvider func() []byte) *ParticipantInvocationStepBuilder {
	b.action = NewParticipantInvocation(cmdEndpoint, cmdProvider)
	return b
}

func (b *ParticipantInvocationStepBuilder) WithCompensation(cmdEndpoint CommandEndpoint, cmdProvider func() []byte) *ParticipantInvocationStepBuilder {
	b.compensation = NewParticipantInvocation(cmdEndpoint, cmdProvider)
	return b
}

func (b *ParticipantInvocationStepBuilder) OnReply(name string, handler func(data []byte, msg []byte)) *ParticipantInvocationStepBuilder {
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

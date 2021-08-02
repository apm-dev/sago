package sago

type ParticipantInvocationStepBuilder struct {
	parent              *SagaDefinitionBuilder
	action              *ParticipantInvocation
	actionReplyHandlers map[string]func(data, msg []byte, successful bool) SagaData
}

func NewParticipantInvocationStepBuilder(parent *SagaDefinitionBuilder) *ParticipantInvocationStepBuilder {
	return &ParticipantInvocationStepBuilder{
		parent:              parent,
		actionReplyHandlers: make(map[string]func(data, msg []byte, successful bool) SagaData),
	}
}

func (b *ParticipantInvocationStepBuilder) withAction(cmdEndpoint CommandEndpoint, cmdProvider func([]byte) []byte) *ParticipantInvocationStepBuilder {
	b.action = NewParticipantInvocation(cmdEndpoint, cmdProvider)
	return b
}

func (b *ParticipantInvocationStepBuilder) OnReply(handler func(data, msg []byte, successful bool) SagaData) *ParticipantInvocationStepBuilder {
	// b.actionReplyHandlers[common.StructName(reply)] = handler
	b.actionReplyHandlers[b.action.cmdEndpoint.CommandName()+"Reply"] = handler
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
	b.parent.AddStep(
		b.action.cmdEndpoint.CommandName(),
		NewParticipantInvocationStep(b.action, b.actionReplyHandlers),
	)
}

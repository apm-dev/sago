package sago

import "google.golang.org/protobuf/proto"

type ParticipantInvocationStepBuilder struct {
	parent              *SagaDefinitionBuilder
	action              *ParticipantInvocation
	actionReplyHandlers map[string]func(data, msg []byte, successful bool) (SagaData, error)
}

func NewParticipantInvocationStepBuilder(parent *SagaDefinitionBuilder) *ParticipantInvocationStepBuilder {
	return &ParticipantInvocationStepBuilder{
		parent:              parent,
		actionReplyHandlers: make(map[string]func(data, msg []byte, successful bool) (SagaData, error)),
	}
}

func (b *ParticipantInvocationStepBuilder) withAction(cmdEndpoint CommandEndpoint, cmdProvider func([]byte, map[string]interface{}) (proto.Message, error)) *ParticipantInvocationStepBuilder {
	b.action = NewParticipantInvocation(cmdEndpoint, cmdProvider)
	return b
}

func (b *ParticipantInvocationStepBuilder) OnReply(handler func(data, msg []byte, successful bool) (SagaData, error)) *ParticipantInvocationStepBuilder {
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

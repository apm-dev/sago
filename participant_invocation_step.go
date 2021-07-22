package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
	"log"
)

type ParticipantInvocationStep struct {
	participantInvocation *ParticipantInvocation
	actionReplyHandlers   map[string]func(data, msg []byte) SagaData
}

func NewParticipantInvocationStep(
	participantInvocation *ParticipantInvocation,
	actionReplyHandlers map[string]func(data, msg []byte) SagaData,
) *ParticipantInvocationStep {
	return &ParticipantInvocationStep{
		participantInvocation: participantInvocation,
		actionReplyHandlers:   actionReplyHandlers,
	}
}

func (stp *ParticipantInvocationStep) getParticipantInvocation() *ParticipantInvocation {
	return stp.participantInvocation
}

func (stp *ParticipantInvocationStep) IsSuccessfulReply(msg messaging.Message) bool {
	return stp.getParticipantInvocation().isSuccessfulReply(msg)
}

func (stp *ParticipantInvocationStep) GetReplyHandler(msg messaging.Message) func(data, msg []byte) SagaData {
	replyType, err := msg.RequiredHeader(commands.REPLY_TYPE)
	if err != nil {
		log.Printf("failed to get reply handler ->\nmsg: %v\nerr: %v\n", msg, err)
		return nil
	}
	return stp.actionReplyHandlers[replyType]
}

/* func (stp *ParticipantInvocationStep) MakeStepOutcome(data []byte, compensating bool) StepOutcome {
	cmd := stp.getParticipantInvocation(compensating).makeCommandToSend()
	return NewRemoteStepOutcome([]commands.Command{cmd})
} */

func (stp *ParticipantInvocationStep) Command(sagaData []byte) commands.Command {
	return stp.getParticipantInvocation().makeCommandToSend(sagaData)
}

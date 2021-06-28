package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
	"log"
)

type ParticipantInvocationStep struct {
	participantInvocation     *ParticipantInvocation
	compensation              *ParticipantInvocation
	actionReplyHandlers       map[string]func(data []byte, msg []byte)
	compensationReplyHandlers map[string]func(data []byte, msg []byte)
}

func NewParticipantInvocationStep(
	participantInvocation *ParticipantInvocation,
	compensation *ParticipantInvocation,
	actionReplyHandlers map[string]func(data []byte, msg []byte),
	compensationReplyHandlers map[string]func(data []byte, msg []byte),
) *ParticipantInvocationStep {
	return &ParticipantInvocationStep{
		participantInvocation:     participantInvocation,
		compensation:              compensation,
		actionReplyHandlers:       actionReplyHandlers,
		compensationReplyHandlers: compensationReplyHandlers,
	}
}

func (stp *ParticipantInvocationStep) getParticipantInvocation(compensating bool) *ParticipantInvocation {
	if compensating {
		return stp.compensation
	}
	return stp.participantInvocation
}

func (stp *ParticipantInvocationStep) IsSuccessfulReply(compensating bool, msg messaging.Message) bool {
	return stp.getParticipantInvocation(compensating).isSuccessfulReply(msg)
}

func (stp *ParticipantInvocationStep) GetReplyHandler(msg messaging.Message, compensating bool) func(data []byte, msg []byte) {
	replyType, err := msg.RequiredHeader(commands.REPLY_TYPE)
	if err != nil {
		// TODO log
		log.Print(err)
		return nil
	}
	if compensating {
		return stp.compensationReplyHandlers[replyType]
	}
	return stp.actionReplyHandlers[replyType]
}

/* func (stp *ParticipantInvocationStep) MakeStepOutcome(data []byte, compensating bool) StepOutcome {
	cmd := stp.getParticipantInvocation(compensating).makeCommandToSend()
	return NewRemoteStepOutcome([]commands.Command{cmd})
} */

func (stp *ParticipantInvocationStep) Command(compensating bool) *Command {
	return stp.getParticipantInvocation(compensating).makeCommandToSend()
}

func (stp *ParticipantInvocationStep) HasAction() bool {
	return stp.participantInvocation != nil
}

func (stp *ParticipantInvocationStep) HasCompensation() bool {
	return stp.compensation != nil
}

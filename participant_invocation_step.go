package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
	"log"
)

type ParticipantInvocationStep struct {
	participantInvocation     *ParticipantInvocation
	compensation              *ParticipantInvocation
	actionReplyHandlers       map[string]func([]byte)
	compensationReplyHandlers map[string]func([]byte)
}

func NewParticipantInvocationStep(
	participantInvocation *ParticipantInvocation,
	compensation *ParticipantInvocation,
	actionReplyHandlers map[string]func([]byte),
	compensationReplyHandlers map[string]func([]byte),
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

func (stp *ParticipantInvocationStep) GetReplyHandler(msg messaging.Message, compensating bool) func(msg []byte) {
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

// func (stp *ParticipantInvocationStep) MakeStepOutcome(data proto.Message, compensating bool) {}

func (stp *ParticipantInvocationStep) HasAction() bool {
	return stp.participantInvocation != nil
}

func (stp *ParticipantInvocationStep) HasCompensation() bool {
	return stp.compensation != nil
}

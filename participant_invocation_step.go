package sago

import (
	"log"

	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagomsg"
)

type ParticipantInvocationStep struct {
	participantInvocation *ParticipantInvocation
	actionReplyHandlers   map[string]func(data, msg []byte, successful bool) (SagaData, error)
}

func NewParticipantInvocationStep(
	participantInvocation *ParticipantInvocation,
	actionReplyHandlers map[string]func(data, msg []byte, successful bool) (SagaData, error),
) *ParticipantInvocationStep {
	return &ParticipantInvocationStep{
		participantInvocation: participantInvocation,
		actionReplyHandlers:   actionReplyHandlers,
	}
}

func (stp *ParticipantInvocationStep) getParticipantInvocation() *ParticipantInvocation {
	return stp.participantInvocation
}

func (stp *ParticipantInvocationStep) IsSuccessfulReply(msg sagomsg.Message) bool {
	return stp.getParticipantInvocation().isSuccessfulReply(msg)
}

func (stp *ParticipantInvocationStep) GetReplyHandler(msg sagomsg.Message) func(data, msg []byte, successful bool) (SagaData, error) {
	replyType, err := msg.RequiredHeader(sagocmd.REPLY_TYPE)
	if err != nil {
		log.Printf("failed to get reply handler ->\nmsg: %v\nerr: %v\n", msg, err)
		return nil
	}
	return stp.actionReplyHandlers[replyType]
}

/* func (stp *ParticipantInvocationStep) MakeStepOutcome(data []byte, compensating bool) StepOutcome {
	cmd := stp.getParticipantInvocation(compensating).makeCommandToSend()
	return NewRemoteStepOutcome([]sagocmd.Command{cmd})
} */

func (stp *ParticipantInvocationStep) Command(sagaData []byte, vars map[string]interface{}) sagocmd.Command {
	return stp.getParticipantInvocation().makeCommandToSend(sagaData, vars)
}

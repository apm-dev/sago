package sago

import (
	"log"
	"strings"

	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagomsg"
	"github.com/pkg/errors"
)

type ParticipantInvocation struct {
	cmdEndpoint CommandEndpoint
	cmdProvider func(data []byte, vars map[string]interface{}) ([]byte, error)
}

func NewParticipantInvocation(cmdEndpoint CommandEndpoint, cmdProvider func(data []byte, vars map[string]interface{}) ([]byte, error)) *ParticipantInvocation {
	return &ParticipantInvocation{cmdEndpoint, cmdProvider}
}

func (pi *ParticipantInvocation) isSuccessfulReply(msg sagomsg.Message) bool {
	val, err := msg.RequiredHeader(sagocmd.REPLY_OUTCOME)
	if err != nil {
		log.Printf("failed to check message successfulness\nmsg: %v\n", msg)
		return false
	}
	return strings.EqualFold(val, string(sagocmd.SUCCESS))
}

func (pi *ParticipantInvocation) makeCommandToSend(sagaData []byte, vars map[string]interface{}) (sagocmd.Command, error) {
	const op string = "sago.participant_invocation.makeCommandToSend"

	payload, err := pi.cmdProvider(sagaData, vars)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return NewCommand(
		pi.cmdEndpoint.CommandName(),
		pi.cmdEndpoint.Channel(),
		payload,
		map[string]string{},
	) , nil
}

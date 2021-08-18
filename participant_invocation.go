package sago

import (
	"log"
	"strings"

	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagomsg"
)

type ParticipantInvocation struct {
	cmdEndpoint CommandEndpoint
	cmdProvider func(data []byte, vars map[string]interface{}) []byte
}

func NewParticipantInvocation(cmdEndpoint CommandEndpoint, cmdProvider func(data []byte, vars map[string]interface{}) []byte) *ParticipantInvocation {
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

func (pi *ParticipantInvocation) makeCommandToSend(sagaData []byte, vars map[string]interface{}) sagocmd.Command {
	return NewCommand(
		pi.cmdEndpoint.CommandName(),
		pi.cmdEndpoint.Channel(),
		pi.cmdProvider(sagaData, vars),
		map[string]string{},
	)
}

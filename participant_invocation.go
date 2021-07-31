package sago

import (
	"github.com/apm-dev/sago/sagocmd"
	"github.com/apm-dev/sago/sagomsg"
	"log"
	"strings"
)

type ParticipantInvocation struct {
	cmdEndpoint CommandEndpoint
	cmdProvider func(data []byte) []byte
}

func NewParticipantInvocation(cmdEndpoint CommandEndpoint, cmdProvider func(data []byte) []byte) *ParticipantInvocation {
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

func (pi *ParticipantInvocation) makeCommandToSend(sagaData []byte) sagocmd.Command {
	return NewCommand(
		pi.cmdEndpoint.CommandName(),
		pi.cmdEndpoint.Channel(),
		pi.cmdProvider(sagaData),
		map[string]string{},
	)
}

package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
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

func (pi *ParticipantInvocation) isSuccessfulReply(msg messaging.Message) bool {
	val, err := msg.RequiredHeader(commands.REPLY_OUTCOME)
	if err != nil {
		log.Printf("failed to check message successfulness\nmsg: %v\n", msg)
		return false
	}
	return strings.EqualFold(val, string(commands.SUCCESS))
}

func (pi *ParticipantInvocation) makeCommandToSend(sagaData []byte) commands.Command {
	return NewCommand(
		pi.cmdEndpoint.CommandName(),
		pi.cmdEndpoint.Channel(),
		pi.cmdProvider(sagaData),
		map[string]string{},
	)
}

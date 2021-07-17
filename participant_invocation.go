package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
	"log"
	"strings"
)

type ParticipantInvocation struct {
	cmdEndpoint CommandEndpoint
	cmdProvider func() []byte
}

func NewParticipantInvocation(cmdEndpoint CommandEndpoint, cmdProvider func() []byte) *ParticipantInvocation {
	return &ParticipantInvocation{cmdEndpoint, cmdProvider}
}

func (pi *ParticipantInvocation) isSuccessfulReply(msg messaging.Message) bool {
	val, err := msg.RequiredHeader(commands.REPLY_OUTCOME)
	if err != nil {
		// TODO: log
		log.Println(err)
		return false
	}
	return strings.EqualFold(val, string(commands.SUCCESS))
}

func (pi *ParticipantInvocation) makeCommandToSend() commands.Command {
	return NewCommand(
		pi.cmdEndpoint.CommandName(),
		pi.cmdEndpoint.Channel(),
		pi.cmdProvider(),
		map[string]string{},
	)
}

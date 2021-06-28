package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
	"log"
	"strings"
)

type ParticipantInvocation struct {
	cmdProvider func() *Command
}

func NewParticipantInvocation(cmdProvider func() *Command) *ParticipantInvocation {
	return &ParticipantInvocation{cmdProvider}
}

func (pi *ParticipantInvocation) isSuccessfulReply(msg messaging.Message) bool {
	val, err := msg.RequiredHeader(commands.REPLY_OUTCOME)
	if err != nil {
		log.Print(err)
		return false
	}
	return strings.EqualFold(val, string(commands.SUCCESS))
}

func (pi *ParticipantInvocation) makeCommandToSend() *Command {
	return pi.cmdProvider()
}

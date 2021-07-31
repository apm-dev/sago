package sagocmd

import (
	"github.com/apm-dev/sago/sagomsg"
)

type CommandReplyOutcome string

const (
	FAILURE  CommandReplyOutcome = "failed"
	SUCCESS CommandReplyOutcome = "success"
)

type Command interface {
	GetName() string
	GetPayload() []byte
	GetChannel() string
	GetExtraHeaders() map[string]string
}

func MakeMessage(channel, replyTo string, cmd Command, headers map[string]string) sagomsg.Message {
	b := sagomsg.WithPayload(cmd.GetPayload())
	b.WithExtraHeaders("", headers)
	b.WithHeader(DESTINATION, channel)
	b.WithHeader(COMMAND_TYPE, cmd.GetName())
	b.WithHeader(REPLY_TO, replyTo)
	return b.Build()
}

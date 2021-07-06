package commands

import "apm-dev/sago/messaging"

type CommandReplyOutcome string

const (
	FAILED  CommandReplyOutcome = "failed"
	SUCCESS CommandReplyOutcome = "success"
)

type Command interface {
	GetName() string
	GetPayload() []byte
}

func MakeMessage(channel, replyTo string, cmd Command, headers map[string]string) messaging.Message {
	b := messaging.WithPayload(cmd.GetPayload())
	b.WithExtraHeaders("", headers)
	b.WithHeader(DESTINATION, channel)
	b.WithHeader(COMMAND_TYPE, cmd.GetName())
	b.WithHeader(REPLY_TO, replyTo)
	return b.Build()
}

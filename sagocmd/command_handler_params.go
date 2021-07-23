package sagocmd

import (
	"apm-dev/sago/sagomsg"
	"strings"
)

type CommandHandlerParams struct {
	command             []byte
	correlationHeaders  map[string]string
	defaultReplyChannel string
}

func NewCommandReplyHandlerParams(msg sagomsg.Message) *CommandHandlerParams {
	return &CommandHandlerParams{
		command:             msg.Payload(),
		correlationHeaders:  getCorrelationHeaders(msg.Headers()),
		defaultReplyChannel: msg.Header(REPLY_TO),
	}
}

func getCorrelationHeaders(headers map[string]string) map[string]string {
	cheaders := make(map[string]string)
	for key, value := range headers {
		if strings.HasPrefix(key, COMMAND_HEADER_PREFIX) {
			cheaders[InReply(key)] = value
		}
	}
	cheaders[IN_REPLY_TO] = headers[sagomsg.ID]
	return cheaders
}

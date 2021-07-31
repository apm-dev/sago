package sagocmd

import (
	"log"
	"strings"

	"github.com/apm-dev/sago/sagomsg"
)

type CommandHandler struct {
	channel string
	cmdType string
	handler func([]byte) sagomsg.Message
}

func NewCommandHandler(channel, cmdType string, handler func([]byte) sagomsg.Message) *CommandHandler {
	return &CommandHandler{
		channel: channel,
		cmdType: cmdType,
		handler: handler,
	}
}

func (h *CommandHandler) Channel() string {
	return h.channel
}

func (h *CommandHandler) CommandType() string {
	return h.cmdType
}

func (h *CommandHandler) Handler() func([]byte) sagomsg.Message {
	return h.handler
}

func (h *CommandHandler) InvokeHandler(cm *CommandMessage) sagomsg.Message {
	return h.handler(cm.Command())
}

func (h *CommandHandler) Handles(msg sagomsg.Message) bool {
	return h.commandTypeMatches(msg)
}

func (h *CommandHandler) commandTypeMatches(msg sagomsg.Message) bool {
	cmdType, err := msg.RequiredHeader(COMMAND_TYPE)
	if err != nil {
		log.Println("there is no", cmdType, "header in message to handle command")
		return false
	}
	return strings.EqualFold(h.cmdType, cmdType)
}

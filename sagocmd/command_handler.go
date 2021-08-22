package sagocmd

import (
	"fmt"
	"strings"

	"git.coryptex.com/lib/sago/sagolog"
	"git.coryptex.com/lib/sago/sagomsg"
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
	const op string = "sagocmd.command_handler.commandTypeMatches"
	cmdType, err := msg.RequiredHeader(COMMAND_TYPE)
	if err != nil {
		sagolog.Log(sagolog.WARN, fmt.Sprintf("%s: %v", op, err))
		return false
	}
	return strings.EqualFold(h.cmdType, cmdType)
}

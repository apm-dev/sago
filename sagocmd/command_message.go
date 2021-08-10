package sagocmd

import "git.coryptex.com/lib/sago/sagomsg"

type CommandMessage struct {
	messageID          string
	command            []byte
	correlationHeaders map[string]string
	message            sagomsg.Message
}

func NewCommandMessage(
	messageID string, command []byte,
	correlationHeaders map[string]string, message sagomsg.Message,
) *CommandMessage {
	return &CommandMessage{
		messageID:          messageID,
		command:            command,
		correlationHeaders: correlationHeaders,
		message:            message,
	}
}

func (cm *CommandMessage) MessageID() string {
	return cm.messageID
}

func (cm *CommandMessage) Command() []byte {
	return cm.command
}

func (cm *CommandMessage) CorrelationHeaders() map[string]string {
	return cm.correlationHeaders
}

func (cm *CommandMessage) Message() sagomsg.Message {
	return cm.message
}

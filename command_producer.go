package sago

import "apm-dev/sago/commands"

type SagaCommandProducer struct {
	cmdProducer commands.CommandProducer
}

func NewSagaCommandProducer(cp commands.CommandProducer) *SagaCommandProducer {
	if cp == nil {
		panic("command producer should not be nil")
	}
	return &SagaCommandProducer{cp}
}

func (cp *SagaCommandProducer) sendCommands(sagaType, sagaID, sagaReplyChannel string, commands []Command) string {
	var msgID string
	for _, cmd := range commands {
		headers := make(map[string]string)
		for k, v := range cmd.ExtraHeaders {
			headers[k] = v
		}
		headers[SAGA_TYPE] = sagaType
		headers[SAGA_ID] = sagaID
		msgID = cp.cmdProducer.Send(cmd.Channel, sagaReplyChannel, cmd, headers)
	}
	return msgID
}

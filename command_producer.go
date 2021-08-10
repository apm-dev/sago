package sago

import (
	"git.coryptex.com/lib/sago/sagocmd"

	"github.com/pkg/errors"
)

type SagaCommandProducer struct {
	cmdProducer sagocmd.CommandProducer
}

func NewSagaCommandProducer(cp sagocmd.CommandProducer) *SagaCommandProducer {
	if cp == nil {
		panic("command producer should not be nil")
	}
	return &SagaCommandProducer{cp}
}

func (cp *SagaCommandProducer) sendCommands(sagaType, sagaID, sagaReplyChannel string, commands []sagocmd.Command) (string, error) {
	var msgID string
	for _, cmd := range commands {
		headers := make(map[string]string)
		for k, v := range cmd.GetExtraHeaders() {
			headers[k] = v
		}
		headers[SAGA_TYPE] = sagaType
		headers[SAGA_ID] = sagaID
		var err error
		msgID, err = cp.cmdProducer.Send(cmd.GetChannel(), sagaReplyChannel, cmd, headers)
		if err != nil {
			return "", errors.Wrapf(err,
				"failed to send command %s of saga %s:%s on %s channel\n",
				cmd.GetName(), sagaType, sagaID, cmd.GetChannel(),
			)
		}
	}
	return msgID, nil
}

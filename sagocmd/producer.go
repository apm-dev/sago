package sagocmd

import (
	"github.com/apm-dev/sago/sagomsg"

	"github.com/pkg/errors"
)

// TODO: implement
type CommandProducer interface {
	Send(channel, replyTo string, cmd Command, headers map[string]string) (string, error)
}

type commandProducerImpl struct {
	msgProducer sagomsg.MessageProducer
}

func NewCommandProducerImpl(mp sagomsg.MessageProducer) CommandProducer {
	return &commandProducerImpl{
		msgProducer: mp,
	}
}

func (p *commandProducerImpl) Send(channel, replyTo string, cmd Command, headers map[string]string) (string, error) {
	msg := MakeMessage(channel, replyTo, cmd, headers)
	err := p.msgProducer.Send(channel, msg)

	if err != nil {
		return "", errors.Wrapf(err,
			"Couldn't send the command: %s to %s",
			cmd.GetName(), channel,
		)
	}
	id, err := msg.ID()
	if err != nil {
		return "", errors.Wrapf(err,
			"sent message doesn't have ID, cmd: %s, channel: %s",
			cmd.GetName(), channel,
		)
	}
	return id, nil
}

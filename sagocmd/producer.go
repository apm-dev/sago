package sagocmd

import (
	"git.coryptex.com/lib/sago/sagomsg"

	"github.com/pkg/errors"
)

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
	const op string = "sagocmd.producer.Send"

	msg := MakeMessage(channel, replyTo, cmd, headers)
	err := p.msgProducer.Send(channel, msg)

	if err != nil {
		return "", errors.Wrapf(err,
			"%s: failed to send %s command to %s",
			op, cmd.GetName(), channel,
		)
	}
	id, err := msg.ID()
	if err != nil {
		return "", errors.Wrapf(err,
			"%s: sent %s command message on %s channel doesn't have ID",
			op, cmd.GetName(), channel,
		)
	}
	return id, nil
}

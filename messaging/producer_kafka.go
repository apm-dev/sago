package messaging

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	wmsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type MessageProducerKafkaImpl struct {
	pub wmsg.Publisher
}

func NewMessageProducerKafkaImpl(brokers []string) (*MessageProducerKafkaImpl, error) {
	p := MessageProducerKafkaImpl{}
	var err error
	p.pub, err = kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kafka message producer")
	}
	return &p, nil
}

func (p *MessageProducerKafkaImpl) Send(destination string, msg Message) error {
	prepareMessageHeaders(msg, destination)

	kmsg := wmsg.NewMessage(
		watermill.NewUUID(),
		msg.Payload(),
	)
	kmsg.Metadata = msg.Headers()

	err := p.pub.Publish(destination, kmsg)
	if err != nil {
		return errors.Wrapf(err,
			"failed to send %v message on %s channel using kafka.\nerr: %v",
			msg, destination, err,
		)
	}

	msg.SetHeader(ID, kmsg.UUID)
	return nil
}

package messaging

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/pkg/errors"
)

var (
	// TODO: log
	logger = watermill.NewStdLogger(false, false)
)

type MessageConsumerKafkaImpl struct {
	id       string
	consumer *kafka.Subscriber
}

func NewMessageConsumerKafkaImpl(brokers []string) MessageConsumer {
	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:     brokers,
			Unmarshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	
	return &MessageConsumerKafkaImpl{
		id:       watermill.NewUUID(),
		consumer: sub,
	}
}

func (c *MessageConsumerKafkaImpl) Subscribe(subscriberID string, channels []string, handler func(m Message)) MessageSubscription {
	for _, ch := range channels {
		msgs, err := c.consumer.Subscribe(context.Background(), ch)
		if err != nil {
			panic(errors.Wrapf(err,
				"failed to subscribe on kafka topic %s, subscriberID: %s \nerr: %s",
				ch, subscriberID, err,
			))
		}
		go func() {
			for msg := range msgs {
				log.Printf("message received, id: %s", msg.UUID)
				// TODO: store msg in db and detect message duplication
				handler(NewMessage(msg.Payload, msg.Metadata))
				msg.Ack()
			}
		}()
	}
	return &KafkaSubscription{
		close: c.Close,
	}
}

func (c *MessageConsumerKafkaImpl) ID() string {
	return c.id
}

func (c *MessageConsumerKafkaImpl) Close() error {
	err := c.consumer.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close kafka subscription")
	}
	return nil
}

type KafkaSubscription struct {
	close func() error
}

func (s *KafkaSubscription) Unsubscribe() error {
	return s.close()
}

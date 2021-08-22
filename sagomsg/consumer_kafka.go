package sagomsg

import (
	"context"
	"fmt"

	"git.coryptex.com/lib/sago/sagolog"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/pkg/errors"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

type MessageConsumerKafkaImpl struct {
	id  string
	sub *kafka.Subscriber
}

func NewMessageConsumerKafkaImpl(brokers []string) MessageConsumer {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	// equivalent of auto.offset.reset: earliest
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "sago-consumer-group",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return &MessageConsumerKafkaImpl{
		id:  watermill.NewUUID(),
		sub: sub,
	}
}

func (c *MessageConsumerKafkaImpl) Subscribe(subscriberID string, channels []string, handler func(m Message) error) MessageSubscription {
	const op string = "sagomsg.consumer_kafka.Subscribe"

	for _, ch := range channels {
		if ch == "" {
			sagolog.Log(sagolog.WARN,
				fmt.Sprintf("%s: could not subscribe on empty channel name", op),
			)
			continue
		}

		sagolog.Log(sagolog.DEBUG,
			fmt.Sprintf("%s: subscribing to %s channel", op, ch),
		)

		msgs, err := c.sub.Subscribe(context.Background(), ch)
		if err != nil {
			panic(errors.Wrapf(err,
				"%s: failed to subscribe on topic %s, subscriberID: %s\nerr: %v",
				op, ch, subscriberID, err,
			))
		}
		go func(ch string) {
			for msg := range msgs {
				sagolog.Log(sagolog.DEBUG,
					fmt.Sprintf("%s: message %s:%s received", op, ch, msg.UUID),
				)
				// TODO: store msg in db and detect message duplication
				err := handler(NewMessage(msg.Payload, msg.Metadata))
				if err != nil {
					sagolog.Log(sagolog.ERROR,
						fmt.Sprintf("%s: failed to handle message %s:%s\n%v",
							op, ch, msg.UUID, err),
					)
					msg.Nack()
					continue
				}
				sagolog.Log(sagolog.DEBUG,
					fmt.Sprintf("%s: sending ack for %s:%s message", op, ch, msg.UUID),
				)
				msg.Ack()
			}
			sagolog.Log(sagolog.DEBUG,
				fmt.Sprintf("%s: stop subscribing on %s", op, ch),
			)
		}(ch)
	}
	return &KafkaSubscription{
		close: c.Close,
	}
}

func (c *MessageConsumerKafkaImpl) ID() string {
	return c.id
}

func (c *MessageConsumerKafkaImpl) Close() error {
	err := c.sub.Close()
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

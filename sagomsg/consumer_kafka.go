package sagomsg

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/pkg/errors"
)

var (
	// TODO: log
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
	for _, ch := range channels {
		if ch == "" {
			log.Println("could not subscribe on empty channel")
			continue
		}
		log.Println("registering kafka subscriber for", ch)

		msgs, err := c.sub.Subscribe(context.Background(), ch)
		if err != nil {
			panic(errors.Wrapf(err,
				"failed to subscribe on kafka topic %s, subscriberID: %s\nerr: %v\n",
				ch, subscriberID, err,
			))
		}
		go func(ch string) {
			log.Println(ch, "kafka subscribing started")
			for msg := range msgs {
				log.Printf("message received, id: %s", msg.UUID)
				// TODO: store msg in db and detect message duplication
				err := handler(NewMessage(msg.Payload, msg.Metadata))
				if err != nil {
					msg.Nack()
					continue
				}
				msg.Ack()
			}
			log.Println(ch, "kafka subscribing closed")
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

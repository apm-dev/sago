package messaging

type MessageConsumer interface {
	Subscribe(subscriberID string, channels []string, handler func(m Message)) MessageSubscription
	ID() string
	Close()
}

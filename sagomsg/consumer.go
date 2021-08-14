package sagomsg

type MessageConsumer interface {
	Subscribe(subscriberID string, channels []string, handler func(m Message) error) MessageSubscription
	ID() string
	Close() error
}

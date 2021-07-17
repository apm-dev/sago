package messaging

type MessageSubscription interface {
	Unsubscribe() error
}

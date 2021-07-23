package sagomsg

type MessageSubscription interface {
	Unsubscribe() error
}

package messaging

type MessageProducer interface {
	Send(destination string, msg Message)
}
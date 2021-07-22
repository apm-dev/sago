package messaging

import "time"

type MessageProducer interface {
	Send(destination string, msg Message) error
}

func prepareMessageHeaders(msg Message, destination string) {
	msg.SetHeader(DESTINATION, destination)
	msg.SetHeader(DATE, time.Now().UTC().Format(time.RFC1123))
}

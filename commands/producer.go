package commands

import "apm-dev/sago/messaging"

// TODO: implement
type CommandProducer interface {
	Send(channel, replyTo string, cmd Command, headers map[string]string) string
}

type CommandProducerImpl struct {
	msgProducer messaging.MessageProducer
}


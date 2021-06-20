package commands

// TODO: implement
type CommandProducer interface {
	Send(channel, replyTo string, cmd Command, headers map[string]string) string
}

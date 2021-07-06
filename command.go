package sago

import (
	"apm-dev/sago/commands"
)

type Command struct {
	commands.Command
	Name         string
	Channel      string
	Payload      []byte
	ReplyStructs []interface{}
	ExtraHeaders map[string]string
}

func (c *Command) GetPayload() []byte {
	return c.Payload
}

func (c *Command) GetName() string {
	return c.Name
}

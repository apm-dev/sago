package sago

import (
	"reflect"
	"strings"
)

type CommandEndpoint struct {
	command interface{}
	channel string
}

func NewCommandEndpoint(cmd interface{}, channel string) CommandEndpoint {
	return CommandEndpoint{cmd, channel}
}

func (c *CommandEndpoint) Command() interface{} {
	return c.command
}

func (c *CommandEndpoint) CommandName() string {
	n := strings.Split(reflect.TypeOf(c.command).String(), ".")
	return n[len(n)-1]
}

func (c *CommandEndpoint) Channel() string {
	return c.channel
}

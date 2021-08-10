package sago

import "git.coryptex.com/lib/sago/common"

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
	return common.StructName(c.command)
}

func (c *CommandEndpoint) Channel() string {
	return c.channel
}

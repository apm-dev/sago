package sago

import (
	"git.coryptex.com/lib/sago/common"
	"google.golang.org/protobuf/proto"
)

type CommandEndpoint struct {
	command proto.Message
	channel string
}

func NewCommandEndpoint(cmd proto.Message, channel string) CommandEndpoint {
	return CommandEndpoint{cmd, channel}
}

func (c *CommandEndpoint) Command() proto.Message {
	return c.command
}

func (c *CommandEndpoint) CommandName() string {
	return common.StructName(c.command)
}

func (c *CommandEndpoint) Channel() string {
	return c.channel
}

package sago

import (
	"apm-dev/sago/commands"

	"google.golang.org/protobuf/proto"
)

type Command struct {
	commands.Command
	Name         string
	Channel      string
	Payload      []proto.Message
	ReplyStructs []interface{}
	ExtraHeaders map[string]string
}

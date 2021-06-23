package sago

import (
	"apm-dev/sago/commands"

	"google.golang.org/protobuf/proto"
)

const (
	SAGA_TYPE = commands.COMMAND_HEADER_PREFIX + "saga_type"
	SAGA_ID   = commands.COMMAND_HEADER_PREFIX + "saga_id"
)

type Command struct {
	commands.Command
	Name         string
	Channel      string
	Payload      []proto.Message
	ReplyStructs []interface{}
	ExtraHeaders map[string]string
}

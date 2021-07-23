package sago

import (
	"apm-dev/sago/commands"
)

const (
	SAGA_TYPE = commands.COMMAND_HEADER_PREFIX + "saga_type"
	SAGA_ID   = commands.COMMAND_HEADER_PREFIX + "saga_id"

	REPLY_SAGA_TYPE    = commands.COMMAND_REPLY_PREFIX + "saga_type"
	REPLY_SAGA_ID      = commands.COMMAND_REPLY_PREFIX + "saga_id"
	REPLY_LOCKED       = "saga-locked-target"
)

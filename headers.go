package sago

import (
	"apm-dev/sago/sagocmd"
)

const (
	SAGA_TYPE = sagocmd.COMMAND_HEADER_PREFIX + "saga_type"
	SAGA_ID   = sagocmd.COMMAND_HEADER_PREFIX + "saga_id"

	REPLY_SAGA_TYPE = sagocmd.COMMAND_REPLY_PREFIX + "saga_type"
	REPLY_SAGA_ID   = sagocmd.COMMAND_REPLY_PREFIX + "saga_id"
	REPLY_LOCKED    = "saga-locked-target"
)

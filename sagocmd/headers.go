package sagocmd

import "strings"

const (
	COMMAND_HEADER_PREFIX = "command_"
	COMMAND_TYPE          = COMMAND_HEADER_PREFIX + "type"
	RESOURCE              = COMMAND_HEADER_PREFIX + "resource"
	DESTINATION           = COMMAND_HEADER_PREFIX + "_destination"
	COMMAND_REPLY_PREFIX  = "commandreply_"
	REPLY_TO              = COMMAND_HEADER_PREFIX + "reply_to"

	REPLY_TYPE    = "reply_type"
	IN_REPLY_TO   = "reply_to_message_id"
	REPLY_OUTCOME = "reply_outcome-type"
)

func InReply(header string) string {
	if strings.HasPrefix(header, COMMAND_HEADER_PREFIX) {
		return strings.ReplaceAll(header, COMMAND_HEADER_PREFIX, COMMAND_REPLY_PREFIX)
	}
	return ""
}

package sagocmd

import (
	"git.coryptex.com/lib/sago/common"
	"git.coryptex.com/lib/sago/sagomsg"
)

var ReplyBuilder = &CommandHandlerReplyBuilder{}

type CommandHandlerReplyBuilder struct {
}

func (b *CommandHandlerReplyBuilder) with(cmd interface{}, reply []byte, outcome CommandReplyOutcome) sagomsg.Message {
	return sagomsg.WithPayload(reply).
		WithHeader(REPLY_OUTCOME, string(outcome)).
		WithHeader(REPLY_TYPE, common.StructName(cmd)+"Reply").
		Build()
}

func (b *CommandHandlerReplyBuilder) WithSuccess(cmd interface{}, reply []byte) sagomsg.Message {
	return b.with(cmd, reply, SUCCESS)
}

func (b *CommandHandlerReplyBuilder) WithEmptySuccess(cmd interface{}) sagomsg.Message {
	return b.WithSuccess(cmd, nil)
}

func (b *CommandHandlerReplyBuilder) WithFailure(cmd interface{}, reply []byte) sagomsg.Message {
	return b.with(cmd, reply, FAILURE)
}

func (b *CommandHandlerReplyBuilder) WithEmptyFailure(cmd interface{}) sagomsg.Message {
	return b.WithFailure(cmd, nil)
}

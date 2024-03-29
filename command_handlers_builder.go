package sago

import (
	"git.coryptex.com/lib/sago/common"
	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagomsg"
)

type SagaCommandHandlersBuilder struct {
	channel  string
	handlers []sagocmd.CommandHandler
}

func NewSagaCommandHandlersBuilder(channel string) *SagaCommandHandlersBuilder {
	return &SagaCommandHandlersBuilder{
		channel:  channel,
		handlers: make([]sagocmd.CommandHandler, 0),
	}
}

func (b *SagaCommandHandlersBuilder) OnMessage(cmd interface{}, handler func([]byte) sagomsg.Message) *SagaCommandHandlersBuilder {
	b.handlers = append(b.handlers,
		*sagocmd.NewCommandHandler(b.channel, common.StructName(cmd), handler))
	return b
}

func (b *SagaCommandHandlersBuilder) Build() *sagocmd.CommandHandlers {
	return sagocmd.NewCommandHandlers(b.handlers)
}

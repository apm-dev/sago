package sagocmd

import (
	"github.com/apm-dev/sago/sagomsg"
)

type CommandDispatcherFactory struct {
	mc sagomsg.MessageConsumer
	mp sagomsg.MessageProducer
}

func NewCommandDispatcherFactory(
	mc sagomsg.MessageConsumer, mp sagomsg.MessageProducer,
) *CommandDispatcherFactory {
	return &CommandDispatcherFactory{
		mc: mc,
		mp: mp,
	}
}

func (f *CommandDispatcherFactory) Make(
	cmdDispatcherID string, target *CommandHandlers,
) *CommandDispatcher {
	return NewCommandDispatcher(cmdDispatcherID, target, f.mc, f.mp)
}

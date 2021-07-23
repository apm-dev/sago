package sagocmd

import (
	"apm-dev/sago/sagomsg"
	"log"
)

type CommandDispatcher struct {
	cmdDispatcherID string
	cmdHandlers     *CommandHandlers
	mc              sagomsg.MessageConsumer
	mp              sagomsg.MessageProducer
}

func NewCommandDispatcher(
	cmdDispatcherID string, cmdHandlers *CommandHandlers,
	mc sagomsg.MessageConsumer, mp sagomsg.MessageProducer,
) *CommandDispatcher {
	return &CommandDispatcher{
		cmdDispatcherID: cmdDispatcherID,
		cmdHandlers:     cmdHandlers,
		mc:              mc,
		mp:              mp,
	}
}

func (d *CommandDispatcher) Initialize() {
	d.mc.Subscribe(d.cmdDispatcherID,
		d.cmdHandlers.Channels(),
		d.MessageHandler,
	)
}

func (d *CommandDispatcher) MessageHandler(msg sagomsg.Message) {
	log.Printf("received message %s %v\n", d.cmdDispatcherID, msg)

	cmdHandler := d.cmdHandlers.FindTargetMethod(msg)
	if cmdHandler == nil {
		log.Printf("no method for %v\n", msg)
		return
	}

	// TODO: implement
}

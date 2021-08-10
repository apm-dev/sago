package sagocmd

import (
	"git.coryptex.com/lib/sago/sagomsg"
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
		d.handleMessage,
	)
}

func (d *CommandDispatcher) handleMessage(msg sagomsg.Message) {
	log.Printf("received message %s %v\n", d.cmdDispatcherID, msg)

	cmdHandler := d.cmdHandlers.FindTargetMethod(msg)
	if cmdHandler == nil {
		log.Printf("no method for %v\n", msg)
		return
	}

	cmdHandlerParams := NewCommandHandlerParams(msg)

	msgid, err := msg.ID()
	if err != nil {
		log.Printf("message doesn't have ID, msg: %v\b", msg)
		return
	}

	cmdmsg := NewCommandMessage(
		msgid, cmdHandlerParams.Command(),
		cmdHandlerParams.CorrelationHeaders(), msg,
	)

	replyMsg := cmdHandler.InvokeHandler(cmdmsg)

	if cmdHandlerParams.DefaultReplyChannel() == "" {
		log.Printf(
			"no %s header in command message so there is no destination to send reply to",
			REPLY_TO)
		return
	}

	if replyMsg == nil {
		log.Printf("nil reply - not publishing - command message: %v\n", cmdmsg)
		return
	}
	err = d.sendReply(
		cmdHandlerParams.CorrelationHeaders(),
		replyMsg,
		cmdHandlerParams.DefaultReplyChannel(),
	)
	if err != nil {
		log.Printf("failed to send reply ->\nreply:%v\nerr: %v\n", replyMsg, err)
		return
	}
}

func (d *CommandDispatcher) sendReply(headers map[string]string, reply sagomsg.Message, channel string) error {
	return d.mp.Send(
		channel,
		sagomsg.WithMessage(reply).
			WithExtraHeaders("", headers).
			Build(),
	)
}

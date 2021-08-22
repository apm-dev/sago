package sagocmd

import (
	"fmt"

	"git.coryptex.com/lib/sago/sagolog"
	"git.coryptex.com/lib/sago/sagomsg"
	"github.com/pkg/errors"
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

func (d *CommandDispatcher) handleMessage(msg sagomsg.Message) error {
	const op string = "sagocmd.command_dispatcher.handleMessage"

	sagolog.Log(sagolog.DEBUG,
		fmt.Sprintf("%s: message received %s %v", op, d.cmdDispatcherID, msg),
	)

	cmdHandler := d.cmdHandlers.FindTargetMethod(msg)
	if cmdHandler == nil {
		return errors.Errorf("%s: no handler for %+v", op, msg)
	}

	cmdHandlerParams := NewCommandHandlerParams(msg)

	msgid, err := msg.ID()
	if err != nil {
		return errors.Wrapf(err, "%s: message doesn't have ID, msg: %+v", op, msg)
	}

	cmdmsg := NewCommandMessage(
		msgid, cmdHandlerParams.Command(),
		cmdHandlerParams.CorrelationHeaders(), msg,
	)

	replyMsg := cmdHandler.InvokeHandler(cmdmsg)

	if cmdHandlerParams.DefaultReplyChannel() == "" {
		return errors.Errorf(
			"%s: no %s header in command message so there is no destination to send reply to",
			op, REPLY_TO,
		)
	}

	if replyMsg == nil {
		return errors.Errorf("%s: nil reply - not publishing - command message: %v", op, cmdmsg)
	}
	err = d.sendReply(
		cmdHandlerParams.CorrelationHeaders(),
		replyMsg,
		cmdHandlerParams.DefaultReplyChannel(),
	)
	if err != nil {
		return errors.Wrapf(err, op)
	}
	return nil
}

func (d *CommandDispatcher) sendReply(headers map[string]string, reply sagomsg.Message, channel string) error {
	const op string = "sagocmd.command_dispatcher.sendReply"

	err := d.mp.Send(
		channel,
		sagomsg.WithMessage(reply).
			WithExtraHeaders("", headers).
			Build(),
	)
	if err != nil {
		return errors.Wrapf(err, op)
	}
	return nil
}

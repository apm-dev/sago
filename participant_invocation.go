package sago

import (
	"fmt"
	"strings"

	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagolog"
	"git.coryptex.com/lib/sago/sagomsg"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type ParticipantInvocation struct {
	cmdEndpoint CommandEndpoint
	cmdProvider func(data []byte, vars map[string]interface{}) (proto.Message, error)
}

func NewParticipantInvocation(cmdEndpoint CommandEndpoint, cmdProvider func(data []byte, vars map[string]interface{}) (proto.Message, error)) *ParticipantInvocation {
	return &ParticipantInvocation{cmdEndpoint, cmdProvider}
}

func (pi *ParticipantInvocation) isSuccessfulReply(msg sagomsg.Message) bool {
	const op string = "sago.participant_invocation.isSuccessfulReply"
	val, err := msg.RequiredHeader(sagocmd.REPLY_OUTCOME)
	if err != nil {
		sagolog.Log(sagolog.WARN,
			fmt.Sprintf("%s: %v", op, err),
		)
		return false
	}
	return strings.EqualFold(val, string(sagocmd.SUCCESS))
}

func (pi *ParticipantInvocation) makeCommandToSend(sagaData []byte, vars map[string]interface{}) (sagocmd.Command, error) {
	const op string = "sago.participant_invocation.makeCommandToSend"

	cmd, err := pi.cmdProvider(sagaData, vars)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	payload, err := proto.Marshal(cmd)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return NewCommand(
		pi.cmdEndpoint.CommandName(),
		pi.cmdEndpoint.Channel(),
		payload,
		map[string]string{},
	), nil
}

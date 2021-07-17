package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
)

type SagaStep interface {
	//   StepOutcome makeStepOutcome(Data data, boolean compensating);
	IsSuccessfulReply(compensating bool, message messaging.Message) bool
	GetReplyHandler(msg messaging.Message, compensating bool) func(data, msg []byte)
	// MakeStepOutcome(data []byte, compensating bool) StepOutcome
	Command(compensating bool) commands.Command
	HasAction() bool
	HasCompensation() bool
}

package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
)

type SagaStep interface {
	//   StepOutcome makeStepOutcome(Data data, boolean compensating);
	IsSuccessfulReply(message messaging.Message) bool
	GetReplyHandler(msg messaging.Message) func(data, msg []byte) SagaData
	// MakeStepOutcome(data []byte, compensating bool) StepOutcome
	// Command(compensating bool) commands.Command
	Command(sagaData []byte) commands.Command
	// HasAction() bool
	// HasCompensation() bool
}

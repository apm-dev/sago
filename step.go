package sago

import (
	"apm-dev/sago/sagocmd"
	"apm-dev/sago/sagomsg"
)

type SagaStep interface {
	//   StepOutcome makeStepOutcome(Data data, boolean compensating);
	IsSuccessfulReply(message sagomsg.Message) bool
	GetReplyHandler(msg sagomsg.Message) func(data, msg []byte) SagaData
	// MakeStepOutcome(data []byte, compensating bool) StepOutcome
	// Command(compensating bool) sagocmd.Command
	Command(sagaData []byte) sagocmd.Command
	// HasAction() bool
	// HasCompensation() bool
}

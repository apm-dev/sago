package sago

import (
	"apm-dev/sago/messaging"
)

type SagaStep interface {
	//   StepOutcome makeStepOutcome(Data data, boolean compensating);
	IsSuccessfulReply(compensating bool, message messaging.Message) bool
	GetReplyHandler(msg messaging.Message, compensating bool) func(msg []byte)
	// MakeStepOutcome(data proto.Message, compensating bool)
	HasAction() bool
	HasCompensation() bool
}

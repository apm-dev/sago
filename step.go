package sago

import (
	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagomsg"
)

type SagaStep interface {
	//   StepOutcome makeStepOutcome(Data data, boolean compensating);
	IsSuccessfulReply(message sagomsg.Message) bool
	GetReplyHandler(msg sagomsg.Message) func(data, msg []byte, successful bool) (SagaData, error)
	// MakeStepOutcome(data []byte, compensating bool) StepOutcome
	// Command(compensating bool) sagocmd.Command
	Command(sagaData []byte, vars map[string]interface{}) sagocmd.Command
	// HasAction() bool
	// HasCompensation() bool
}

package sago

import "apm-dev/sago/sagocmd"

type StepOutcome interface {
	Visit(cmdConsumer func([]sagocmd.Command))
}

type RemoteStepOutcome struct {
	commandsToSend []sagocmd.Command
}

func NewRemoteStepOutcome(cmds []sagocmd.Command) *RemoteStepOutcome {
	return &RemoteStepOutcome{
		commandsToSend: cmds,
	}
}

func (r *RemoteStepOutcome) Visit(cmdConsumer func([]sagocmd.Command)) {
	cmdConsumer(r.commandsToSend)
}

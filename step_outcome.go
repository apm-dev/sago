package sago

import "apm-dev/sago/commands"

type StepOutcome interface {
	Visit(cmdConsumer func([]commands.Command))
}

type RemoteStepOutcome struct {
	commandsToSend []commands.Command
}

func NewRemoteStepOutcome(cmds []commands.Command) *RemoteStepOutcome {
	return &RemoteStepOutcome{
		commandsToSend: cmds,
	}
}

func (r *RemoteStepOutcome) Visit(cmdConsumer func([]commands.Command)) {
	cmdConsumer(r.commandsToSend)
}

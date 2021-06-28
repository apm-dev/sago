package sago

type StepOutcome interface {
	Visit(cmdConsumer func([]Command))
}

type RemoteStepOutcome struct {
	commandsToSend []Command
}

func NewRemoteStepOutcome(cmds []Command) *RemoteStepOutcome {
	return &RemoteStepOutcome{
		commandsToSend: cmds,
	}
}

func (r *RemoteStepOutcome) Visit(cmdConsumer func([]Command)) {
	cmdConsumer(r.commandsToSend)
}

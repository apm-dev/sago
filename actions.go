package sago

type SagaActions struct {
	commands        []Command
	updatedSagaData []byte
	updatedState    string
	endState        bool
	compensating    bool
}

func NewSagaActions(
	commands []Command,
	updatedSagaData []byte,
	updatedState string,
	endState, compensating bool) *SagaActions {
	return &SagaActions{
		commands, updatedSagaData, updatedState,
		endState, compensating,
	}
}

func (sa *SagaActions) Commands() []Command { return sa.commands }

func (sa *SagaActions) UpdatedSagaData() []byte { return sa.updatedSagaData }

func (sa *SagaActions) UpdatedState() string { return sa.updatedState }

func (sa *SagaActions) IsEndState() bool { return sa.endState }

func (sa *SagaActions) IsCompensating() bool { return sa.compensating }

type SagaActionsBuilder struct {
	commands        []Command
	updatedSagaData []byte
	updatedState    string
	endState        bool
	compensating    bool
}

func NewSagaActionsBuilder() *SagaActionsBuilder {
	return &SagaActionsBuilder{
		commands: []Command{},
	}
}

func (b *SagaActionsBuilder) Build() *SagaActions {
	return NewSagaActions(
		b.commands,
		b.updatedSagaData, b.updatedState,
		b.endState, b.compensating,
	)
}

func (b *SagaActionsBuilder) WithCommand(cmd Command) *SagaActionsBuilder {
	b.commands = append(b.commands, cmd)
	return b
}

func (b *SagaActionsBuilder) WithCommands(cmds []Command) *SagaActionsBuilder {
	for _, cmd := range cmds {
		b.WithCommand(cmd)
	}
	return b
}

func (b *SagaActionsBuilder) WithUpdatedSagaData(data []byte) *SagaActionsBuilder {
	b.updatedSagaData = data
	return b
}

func (b *SagaActionsBuilder) WithUpdatedState(state string) *SagaActionsBuilder {
	b.updatedState = state
	return b
}

func (b *SagaActionsBuilder) WithIsEndState(endState bool) *SagaActionsBuilder {
	b.endState = endState
	return b
}

func (b *SagaActionsBuilder) WithIsCompensating(compensating bool) *SagaActionsBuilder {
	b.compensating = compensating
	return b
}


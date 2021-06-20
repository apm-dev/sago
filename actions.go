package sago

import "google.golang.org/protobuf/proto"

type SagaActions struct {
	commands        []Command
	updatedSagaData proto.Message
	updatedState    string
	endState        bool
	compensating    bool
}

func NewSagaActions(
	commands []Command,
	updatedSagaData proto.Message,
	updatedState string,
	endState, compensating bool) *SagaActions {
	return &SagaActions{
		commands, updatedSagaData, updatedState,
		endState, compensating,
	}
}

func (sa *SagaActions) Commands() []Command { return sa.commands }

func (sa *SagaActions) UpdatedSagaData() proto.Message { return sa.updatedSagaData }

func (sa *SagaActions) UpdatedState() string { return sa.updatedState }

func (sa *SagaActions) IsEndState() bool { return sa.endState }

func (sa *SagaActions) IsCompensating() bool { return sa.compensating }

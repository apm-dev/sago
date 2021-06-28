package sago

import "encoding/json"

type sagaExecutionState struct {
	CurrentlyExecuting int
	Compensating       bool
	EndState           bool
}

func NewSagaExecutionState(currentlyExecuting int, compensating bool) *sagaExecutionState {
	return &sagaExecutionState{
		CurrentlyExecuting: currentlyExecuting,
		Compensating:       compensating,
	}
}

func (ses *sagaExecutionState) GetCurrentlyExecuting() int {
	return ses.CurrentlyExecuting
}

func (ses *sagaExecutionState) SetCurrentlyExecuting(currentlyExecuting int) {
	ses.CurrentlyExecuting = currentlyExecuting
}

func (ses *sagaExecutionState) IsCompensating() bool {
	return ses.Compensating
}

func (ses *sagaExecutionState) SetCompensating(compensating bool) {
	ses.Compensating = compensating
}

func (ses *sagaExecutionState) IsEndState() bool {
	return ses.EndState
}

func (ses *sagaExecutionState) SetEndState(endState bool) {
	ses.EndState = endState
}

func (ses *sagaExecutionState) NextState(size int) *sagaExecutionState {
	step := ses.CurrentlyExecuting + size
	if ses.Compensating {
		step = ses.CurrentlyExecuting - size
	}
	return &sagaExecutionState{CurrentlyExecuting: step, Compensating: ses.Compensating}
}

func (ses *sagaExecutionState) StartCompensating() *sagaExecutionState {
	return NewSagaExecutionState(ses.CurrentlyExecuting, true)
}

func MakeEndState() *sagaExecutionState {
	ses := &sagaExecutionState{}
	ses.SetEndState(true)
	return ses
}

func (ses *sagaExecutionState) encode() string {
	b, _ := json.Marshal(ses)
	return string(b)
}

func (ses *sagaExecutionState) decode(val string) {
	json.Unmarshal([]byte(val), ses)
}

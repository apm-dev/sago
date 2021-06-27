package sago

type SagaExecutionState struct {
	currentlyExecuting int
	compensating       bool
	endState           bool
}

func NewSagaExecutionState(currentlyExecuting int, compensating bool) *SagaExecutionState {
	return &SagaExecutionState{
		currentlyExecuting: currentlyExecuting,
		compensating:       compensating,
	}
}

func (ses *SagaExecutionState) CurrentlyExecuting() int {
	return ses.currentlyExecuting
}

func (ses *SagaExecutionState) SetCurrentlyExecuting(currentlyExecuting int) {
	ses.currentlyExecuting = currentlyExecuting
}

func (ses *SagaExecutionState) IsCompensating() bool {
	return ses.compensating
}

func (ses *SagaExecutionState) SetCompensating(compensating bool) {
	ses.compensating = compensating
}

func (ses *SagaExecutionState) IsEndState() bool {
	return ses.endState
}

func (ses *SagaExecutionState) SetEndState(endState bool) {
	ses.endState = endState
}

func (ses *SagaExecutionState) NextState(size int) *SagaExecutionState {
	step := ses.currentlyExecuting + size
	if ses.compensating {
		step = ses.currentlyExecuting - size
	}
	return &SagaExecutionState{currentlyExecuting: step, compensating: ses.compensating}
}

func MakeEndState() *SagaExecutionState {
	ses := &SagaExecutionState{}
	ses.SetEndState(true)
	return ses
}

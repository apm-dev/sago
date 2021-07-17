package sago

import "apm-dev/sago/commands"

type stepToExecute struct {
	step         SagaStep
	skipped      int
	compensating bool
}

func NewStepToExecute(step SagaStep, skipped int, compensating bool) *stepToExecute {
	return &stepToExecute{
		step:         step,
		skipped:      skipped,
		compensating: compensating,
	}
}

func (s *stepToExecute) size() int {
	if s.step != nil {
		return s.skipped + 1
	}
	return s.skipped
}

func (s *stepToExecute) isEmpty() bool {
	return s.step == nil
}

func (s *stepToExecute) executeStep(data []byte, currentState *sagaExecutionState) *SagaActions {
	newState := currentState.NextState(s.size())
	cmd := s.step.Command(currentState.IsCompensating())
	return NewSagaActions(
		[]commands.Command{cmd},
		data,
		newState.encode(),
		newState.IsEndState(),
		currentState.IsCompensating(),
	)
}

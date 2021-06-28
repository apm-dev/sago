package sago

import (
	"apm-dev/sago/messaging"

	"github.com/pkg/errors"
)

type SagaDefinition interface {
	Start(sagaData []byte) *SagaActions
	HandleReply(currentState string, sagaData []byte, message messaging.Message) (*SagaActions, error)
}

type sagaDefinition struct {
	sagaSteps []SagaStep
}

func NewSagaDefinition(stps []SagaStep) SagaDefinition {
	return &sagaDefinition{stps}
}

func (sd *sagaDefinition) Start(sagaData []byte) *SagaActions {
	currentState := NewSagaExecutionState(-1, false)

	stepToExecute := sd.nextStepToExecute(currentState, sagaData)

	if stepToExecute == nil {
		return sd.makeEndStateSagaActions(currentState)
	}
	return stepToExecute.executeStep(sagaData, currentState)
}

func (sd *sagaDefinition) HandleReply(currentState string, sagaData []byte, msg messaging.Message) (*SagaActions, error) {
	var state *sagaExecutionState
	state.decode(currentState)
	currentStep := sd.sagaSteps[state.GetCurrentlyExecuting()]
	compensating := state.IsCompensating()

	replyHandler := currentStep.GetReplyHandler(msg, compensating)
	if replyHandler != nil {
		replyHandler(sagaData, msg.Payload())
	}

	if currentStep.IsSuccessfulReply(compensating, msg) {
		return sd.executeNextStep(sagaData, state), nil
	} else if compensating {
		return nil, errors.Errorf("Failure when compensating, state: ", currentState)
	} else {
		return sd.executeNextStep(sagaData, state.StartCompensating()), nil
	}
}

func (sd *sagaDefinition) nextStepToExecute(state *sagaExecutionState, sagaData []byte) *stepToExecute {
	skipped := 0
	compensating := state.IsCompensating()
	direction := 1
	if compensating {
		direction = -1
	}
	for i := state.GetCurrentlyExecuting() + direction; i >= 0 && i < len(sd.sagaSteps); i = i + direction {
		step := sd.sagaSteps[i]
		if compensating {
			if step.HasCompensation() {
				return NewStepToExecute(step, skipped, compensating)
			}
			skipped++
		} else {
			if step.HasAction() {
				return NewStepToExecute(step, skipped, compensating)
			}
			skipped++
		}
	}
	return NewStepToExecute(nil, skipped, compensating)
}

func (sd *sagaDefinition) executeNextStep(data []byte, state *sagaExecutionState) *SagaActions {
	stepToExecute := sd.nextStepToExecute(state, data)
	if stepToExecute.isEmpty() {
		return sd.makeEndStateSagaActions(state)
	} else {
		return stepToExecute.executeStep(data, state)
	}
}

func (sd *sagaDefinition) makeEndStateSagaActions(state *sagaExecutionState) *SagaActions {
	return NewSagaActionsBuilder().WithUpdatedState(
		MakeEndState().encode(),
	).WithIsEndState(true).WithIsCompensating(
		state.IsCompensating(),
	).Build()
}

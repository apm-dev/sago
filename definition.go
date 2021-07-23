package sago

import (
	"sync"
)

type SagaDefinition interface {
	// Start(zb zbc.Client, sagaData []byte) *SagaActions
	Step(name string) SagaStep
	StepsName() <-chan string
	// HandleReply(currentState string, sagaData []byte, message messaging.Message) (*SagaActions, error)
}

type sagaDefinition struct {
	sync.RWMutex
	sagaSteps map[string]SagaStep
}

func NewSagaDefinition(steps map[string]SagaStep) SagaDefinition {
	return &sagaDefinition{sagaSteps: steps}
}

func (sd *sagaDefinition) Step(name string) SagaStep {
	sd.RLock()
	defer sd.RUnlock()
	return sd.sagaSteps[name]
}

func (sd *sagaDefinition) StepsName() <-chan string {
	ch := make(chan string, len(sd.sagaSteps))
	// we don't need separate goroutine because it's a buffered channel
	go func() {
		sd.RLock()
		defer sd.RUnlock()
		for name := range sd.sagaSteps {
			ch <- name
		}
		close(ch)
	}()
	return ch
}

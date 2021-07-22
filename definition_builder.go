package sago

import "sync"

type SagaDefinitionBuilder struct {
	sync.RWMutex
	sagaSteps map[string]SagaStep
}

func NewSagaDefinitionBuilder() *SagaDefinitionBuilder {
	return &SagaDefinitionBuilder{sagaSteps: make(map[string]SagaStep)}
}

func (b *SagaDefinitionBuilder) AddStep(name string, step SagaStep) {
	b.Lock()
	defer b.Unlock()
	b.sagaSteps[name] = step
}

func (b *SagaDefinitionBuilder) Build() SagaDefinition {
	b.RLock()
	defer b.RUnlock()
	return NewSagaDefinition(b.sagaSteps)
}

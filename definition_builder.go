package sago

type SagaDefinitionBuilder struct {
	sagaSteps []SagaStep
}

func NewSagaDefinitionBuilder() *SagaDefinitionBuilder {
	return &SagaDefinitionBuilder{sagaSteps: make([]SagaStep, 0)}
}

func (sdb *SagaDefinitionBuilder) AddStep(step SagaStep) {
	sdb.sagaSteps = append(sdb.sagaSteps, step)
}

func (sdb *SagaDefinitionBuilder) Build() SagaDefinition {
	return NewSagaDefinition(sdb.sagaSteps)
}

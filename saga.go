package sago

type Saga interface {
	SagaDefinition() SagaDefinition
	SagaType() string

	OnStarting(sagaID string, data []byte)
	OnSagaCompletedSuccessfully(sagaID string, data []byte)
	OnSagaRolledBack(sagaID string, data []byte)
}

func Step() *StepBuilder {
	return NewStepBuilder(NewSagaDefinitionBuilder())
}

package sago

import "google.golang.org/protobuf/proto"

type Saga interface {
	SagaDefinition() SagaDefinition
	SagaType() string

	OnStarting(sagaID string, data proto.Message)
	OnSagaCompletedSuccessfully(sagaID string, data proto.Message)
	OnSagaRolledBack(sagaID string, data proto.Message)
}

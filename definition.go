package sago

import (
	"apm-dev/sago/messaging"

	"google.golang.org/protobuf/proto"
)

type SagaDefinition interface {
	Start(sagaData proto.Message) *SagaActions
	HandleReply(currentState string, sagaData interface{}, message messaging.Message) *SagaActions
}

package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
)

type SagaManagerFactory struct {
	zb                     zbc.Client
	sagaInstanceRepository SagaInstanceRepository
	commandProducer        commands.CommandProducer
	messageConsumer        messaging.MessageConsumer
	sagaCommandProducer    *SagaCommandProducer
}

func NewSagaManagerFactory(
	zb zbc.Client,
	sagaInstanceRepository SagaInstanceRepository,
	commandProducer commands.CommandProducer,
	messageConsumer messaging.MessageConsumer,
	sagaCommandProducer *SagaCommandProducer,
) *SagaManagerFactory {
	return &SagaManagerFactory{
		zb:                     zb,
		sagaInstanceRepository: sagaInstanceRepository,
		commandProducer:        commandProducer,
		messageConsumer:        messageConsumer,
		sagaCommandProducer:    sagaCommandProducer,
	}
}

func (f *SagaManagerFactory) Make(saga Saga) SagaManager {
	return NewSagaManager(
		saga,
		f.zb,
		f.sagaInstanceRepository,
		f.commandProducer,
		f.messageConsumer,
		f.sagaCommandProducer,
	)
}

package sago

import (
	"apm-dev/sago/commands"
	"apm-dev/sago/messaging"
)

type SagaManagerFactory struct {
	sagaInstanceRepository SagaInstanceRepository
	commandProducer        commands.CommandProducer
	messageConsumer        messaging.MessageConsumer
	sagaCommandProducer    *SagaCommandProducer
}

func NewSagaManagerFactory(
	sagaInstanceRepository SagaInstanceRepository,
	commandProducer commands.CommandProducer,
	messageConsumer messaging.MessageConsumer,
	sagaCommandProducer *SagaCommandProducer,
) *SagaManagerFactory {
	return &SagaManagerFactory{
		sagaInstanceRepository: sagaInstanceRepository,
		commandProducer:        commandProducer,
		messageConsumer:        messageConsumer,
		sagaCommandProducer:    sagaCommandProducer,
	}
}

func (f *SagaManagerFactory) Make(saga Saga) SagaManager {
	return NewSagaManager(
		saga,
		f.sagaInstanceRepository,
		f.commandProducer,
		f.messageConsumer,
		f.sagaCommandProducer,
	)
}

package sago

import (
	"git.coryptex.com/lib/sago/sagocmd"
	"git.coryptex.com/lib/sago/sagomsg"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
)

type SagaManagerFactory struct {
	zb                     zbc.Client
	sagaInstanceRepository SagaInstanceRepository
	commandProducer        sagocmd.CommandProducer
	messageConsumer        sagomsg.MessageConsumer
	sagaCommandProducer    *SagaCommandProducer
}

func NewSagaManagerFactory(
	zb zbc.Client,
	sagaInstanceRepository SagaInstanceRepository,
	commandProducer sagocmd.CommandProducer,
	messageConsumer sagomsg.MessageConsumer,
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

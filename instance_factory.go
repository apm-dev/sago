package sago

import (
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type SagaInstanceFactory struct {
	sagaManagersLock *sync.RWMutex
	sagaManagers     map[Saga]SagaManager
}

func NewSagaInstanceFactory(smf *SagaManagerFactory, sagas []Saga) *SagaInstanceFactory {
	sif := SagaInstanceFactory{
		sagaManagers: map[Saga]SagaManager{},
	}
	sif.sagaManagersLock.Lock()
	defer sif.sagaManagersLock.Unlock()
	for _, saga := range sagas {
		sif.sagaManagers[saga] = sif.makeSagaManager(smf, saga)
	}
	return &sif
}

func (sif *SagaInstanceFactory) Create(saga Saga, data proto.Message) (*SagaInstance, error) {
	sif.sagaManagersLock.RLock()
	defer sif.sagaManagersLock.RUnlock()
	sagaManager := sif.sagaManagers[saga]
	if sagaManager == nil {
		// TODO log
		return nil, errors.Errorf("No SagaManager for %v", saga)
	}
	return sagaManager.Create(data)
}

func (sif *SagaInstanceFactory) makeSagaManager(smf *SagaManagerFactory, saga Saga) SagaManager {
	sagaManager := smf.Make(saga)
	sagaManager.SubscribeToReplyChannel()
	return sagaManager
}

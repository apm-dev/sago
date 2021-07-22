package sago

import (
	"log"
	"sync"

	"github.com/pkg/errors"
)

type SagaInstanceFactory struct {
	sagaManagersLock sync.RWMutex
	sagaManagers     map[Saga]SagaManager
}

func NewSagaInstanceFactory(smf *SagaManagerFactory, sagas []Saga) (*SagaInstanceFactory, error) {
	sif := SagaInstanceFactory{
		sagaManagers: map[Saga]SagaManager{},
	}
	sif.sagaManagersLock.Lock()
	defer sif.sagaManagersLock.Unlock()
	for _, saga := range sagas {
		sm, err := sif.makeSagaManager(smf, saga)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create SagaInstanceFactory for %s", saga.SagaType())
		}
		sif.sagaManagers[saga] = sm
	}
	return &sif, nil
}

func (sif *SagaInstanceFactory) Create(saga Saga, data SagaData) error {
	sif.sagaManagersLock.RLock()
	defer sif.sagaManagersLock.RUnlock()
	sagaManager := sif.sagaManagers[saga]
	if sagaManager == nil {
		// TODO log
		return errors.Errorf("there is no SagaManager registered for %s", saga.SagaType())
	}
	return sagaManager.Create(data)
}

func (sif *SagaInstanceFactory) makeSagaManager(smf *SagaManagerFactory, saga Saga) (SagaManager, error) {
	sagaManager := smf.Make(saga)
	sagaManager.SubscribeToReplyChannel()
	log.Println("subscribed to channel")
	err := sagaManager.RegisterJobWorkers()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create manager for %s saga\n", saga.SagaType())
	}
	log.Println("job worker registered")
	return sagaManager, nil
}

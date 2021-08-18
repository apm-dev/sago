package sago

import (
	"sync"

	"github.com/pkg/errors"
)

type SagaInstanceFactory struct {
	sagaManagersLock sync.RWMutex
	sagaManagers     map[Saga]SagaManager
}

func NewSagaInstanceFactory(smf *SagaManagerFactory, sagas []Saga, bpmnPath string) (*SagaInstanceFactory, error) {
	sif := SagaInstanceFactory{
		sagaManagers: map[Saga]SagaManager{},
	}
	sif.sagaManagersLock.Lock()
	defer sif.sagaManagersLock.Unlock()
	for _, saga := range sagas {
		sm, err := sif.makeSagaManager(smf, saga, bpmnPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create SagaInstanceFactory for %s", saga.SagaType())
		}
		sif.sagaManagers[saga] = sm
	}
	return &sif, nil
}

func (sif *SagaInstanceFactory) Create(uniqueId string, saga Saga, data SagaData, extVars map[string]interface{}) error {
	sif.sagaManagersLock.RLock()
	defer sif.sagaManagersLock.RUnlock()
	sagaManager := sif.sagaManagers[saga]
	if sagaManager == nil {
		return errors.Errorf("there is no SagaManager registered for %s", saga.SagaType())
	}
	return sagaManager.create(uniqueId, data, extVars)
}

func (sif *SagaInstanceFactory) makeSagaManager(smf *SagaManagerFactory, saga Saga, bpmnPath string) (SagaManager, error) {
	sagaManager := smf.Make(saga)
	err := sagaManager.deployProcess(bpmnPath + "/" + saga.SagaType() + ".bpmn")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create manager for %s saga\n", saga.SagaType())
	}
	sagaManager.subscribeToReplyChannel()
	err = sagaManager.registerJobWorkers()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create manager for %s saga\n", saga.SagaType())
	}
	return sagaManager, nil
}

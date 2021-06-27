package sago

type SagaInstanceRepository interface {
	Save(si SagaInstance) (string, error)
	Find(sagaType, sagaID string) (*SagaInstance, error)
	Update(si SagaInstance) error
}

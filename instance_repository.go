package sago

type SagaInstanceRepository interface {
	Save(si SagaInstance) string
	Find(sagaType, sagaID string) *SagaInstance
	Update(si SagaInstance)
}

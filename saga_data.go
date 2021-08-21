package sago

type SagaData interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

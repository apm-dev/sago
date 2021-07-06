package sago

type SagaData interface {
	Marshal() []byte
	Unmarshal([]byte)
}

package sagomsg

import "github.com/pkg/errors"

const (
	ID           = "ID"
	PARTITION_ID = "PARTITION_ID"
	DESTINATION  = "DESTINATION"
	DATE         = "DATE"
)

type Message interface {
	ID() (string, error)
	Headers() map[string]string
	Payload() []byte
	Header(name string) string
	RequiredHeader(name string) (string, error)
	HasHeader(name string) bool

	SetPayload(payload []byte)
	SetHeaders(headers map[string]string)
	SetHeader(name, value string)
	RemoveHeader(key string)
}

func NewMessage(payload []byte, headers map[string]string) Message {
	return &message{
		payload: payload,
		headers: headers,
	}
}

type message struct {
	payload []byte
	headers map[string]string
}

func (m *message) ID() (string, error) {
	return m.RequiredHeader(ID)
}

func (m *message) Headers() map[string]string {
	return m.headers
}

func (m *message) Payload() []byte {
	return m.payload
}

func (m *message) Header(name string) string {
	return m.headers[name]
}

func (m *message) RequiredHeader(name string) (string, error) {
	val, ok := m.headers[name]
	if !ok {
		return "", errors.Errorf("No such header: %s in this message %v", name, m)
	}
	return val, nil
}

func (m *message) HasHeader(name string) bool {
	_, ok := m.headers[name]
	return ok
}

func (m *message) SetPayload(payload []byte) {
	m.payload = payload
}

func (m *message) SetHeaders(headers map[string]string) {
	m.headers = headers
}

func (m *message) SetHeader(name string, value string) {
	if m.headers == nil {
		m.headers = map[string]string{}
	}
	m.headers[name] = value
}

func (m *message) RemoveHeader(key string) {
	ok := m.HasHeader(key)
	if ok {
		delete(m.headers, key)
	}
}

package sagomsg

type MessageBuilder struct {
	Body    []byte
	Headers map[string]string
}

func WithPayload(payload []byte) *MessageBuilder {
	return &MessageBuilder{
		Body:    payload,
		Headers: map[string]string{},
	}
}

func WithMessage(msg Message) *MessageBuilder {
	return &MessageBuilder{
		Body:    msg.Payload(),
		Headers: msg.Headers(),
	}
}
func (b *MessageBuilder) WithHeader(name, value string) *MessageBuilder {
	b.Headers[name] = value
	return b
}

func (b *MessageBuilder) WithExtraHeaders(prefix string, headers map[string]string) *MessageBuilder {
	for name, value := range headers {
		b.Headers[prefix+name] = value
	}
	return b
}

func (b *MessageBuilder) Build() Message {
	return NewMessage(b.Body, b.Headers)
}

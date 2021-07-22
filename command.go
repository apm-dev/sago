package sago

type Command struct {
	Name         string
	Channel      string
	Payload      []byte
	ExtraHeaders map[string]string
}

func NewCommand(name, channel string, payload []byte, extHeaders map[string]string) *Command {
	return &Command{
		Name:         name,
		Channel:      channel,
		Payload:      payload,
		ExtraHeaders: extHeaders,
	}
}

func (c *Command) GetName() string {
	return c.Name
}

func (c *Command) GetPayload() []byte {
	return c.Payload
}

func (c *Command) GetChannel() string {
	return c.Channel
}

func (c *Command) GetExtraHeaders() map[string]string {
	return c.ExtraHeaders
}

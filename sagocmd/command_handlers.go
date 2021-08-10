package sagocmd

import "git.coryptex.com/lib/sago/sagomsg"

type CommandHandlers struct {
	handlers []CommandHandler
}

func NewCommandHandlers(handlers []CommandHandler) *CommandHandlers {
	return &CommandHandlers{handlers: handlers}
}

func (h *CommandHandlers) Channels() []string {
	channelsMap := make(map[string]int)
	for _, handler := range h.handlers {
		channelsMap[handler.Channel()] = 0
	}
	channels := make([]string, 0, len(channelsMap))
	for channel := range channelsMap {
		channels = append(channels, channel)
	}
	return channels
}

func (h *CommandHandlers) FindTargetMethod(msg sagomsg.Message) *CommandHandler {
	for _, handler := range h.handlers {
		if handler.Handles(msg) {
			return &handler
		}
	}
	return nil
}

func (h *CommandHandlers) GetHandlers() []CommandHandler {
	return h.handlers
}

package update_handler

import "glaphyra/internal/bot"

type CommandRegistry struct {
	commands map[string]bot.Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]bot.Command)}
}

func (r *CommandRegistry) Register(command string, handler bot.Command) {
	r.commands[command] = handler
}

func (r *CommandRegistry) Get(command string) bot.Command {
	return r.commands[command]
}

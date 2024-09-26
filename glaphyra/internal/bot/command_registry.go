package bot

type CommandRegistry struct {
	commands map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]Command)}
}

func (r *CommandRegistry) Register(command string, handler Command) {
	r.commands[command] = handler
}

func (r *CommandRegistry) Get(command string) Command {
	return r.commands[command]
}

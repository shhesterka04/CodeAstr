package user_cmd_handler

import (
	"glaphyra/internal/bot"
)

const maxLimitHistory = 5

type CommandHistory struct {
	commands []bot.Command
}

func (ch *CommandHistory) Add(command bot.Command) {
	ch.commands = append([]bot.Command{command}, ch.commands...)
	if len(ch.commands) > maxLimitHistory {
		ch.commands = ch.commands[:maxLimitHistory]
	}
}

func (ch *CommandHistory) Pop() bot.Command {
	if len(ch.commands) == 0 {
		return nil
	}

	command := ch.commands[0]
	ch.commands = ch.commands[1:]
	return command
}

func (ch *CommandHistory) GetPrevious() bot.Command {
	if len(ch.commands) == 0 {
		return nil
	}

	return ch.commands[0]
}

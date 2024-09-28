package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/bot/commands"
	"glaphyra/internal/bot/commands/about"
	"glaphyra/internal/bot/commands/predictions"
)

const (
	maxLimitHistory = 5
	backCommand     = "Назад"
)

var commandHistory = &CommandHistory{}

type Commander interface {
	Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message)
}

type Command struct {
	Cmd     Commander
	PrevCmd Commander
}

type Handler struct {
	registry       *CommandRegistry
	commandHistory *CommandHistory
}

type CommandHistory struct {
	commands []*Command
}

func (ch *CommandHistory) Add(command *Command) {
	ch.commands = append([]*Command{command}, ch.commands...)
	if len(ch.commands) > maxLimitHistory {
		ch.commands = ch.commands[:maxLimitHistory]
	}
}

func (ch *CommandHistory) Pop() *Command {
	if len(ch.commands) == 0 {
		return nil
	}

	command := ch.commands[0]
	ch.commands = ch.commands[1:]
	return command
}

func (ch *CommandHistory) GetPrevious() *Command {
	if len(ch.commands) == 0 {
		return nil
	}

	return ch.commands[0]
}

func NewBotHandler(commandHistory *CommandHistory) *Handler {
	bh := &Handler{
		registry:       NewCommandRegistry(),
		commandHistory: commandHistory,
	}

	bh.registerCommands()
	return bh
}

func (bh *Handler) registerCommands() {
	bh.registry.Register("/start", &Command{Cmd: &cmd.StartCommand{}, PrevCmd: nil})
	bh.registry.Register(backCommand, &Command{Cmd: &BackCommand{commandHistory: bh.commandHistory}, PrevCmd: nil})

	bh.registry.Register("Предсказания", &Command{Cmd: &cmd.PredictionsCommand{}, PrevCmd: &cmd.StartCommand{}})
	bh.registry.Register("Совместимость", &Command{Cmd: &cmd.CompatibilityCommand{}, PrevCmd: &cmd.StartCommand{}})
	bh.registry.Register("Регистрация и настройки", &Command{Cmd: &cmd.SettingsCommand{}, PrevCmd: &cmd.StartCommand{}})
	bh.registry.Register("О боте", &Command{Cmd: &cmd.AboutCommand{}, PrevCmd: &cmd.StartCommand{}})

	// О боте
	bh.registry.Register("Кто я?", &Command{Cmd: &about.WhoAmICommand{}, PrevCmd: &cmd.StartCommand{}})
	bh.registry.Register("Функции", &Command{Cmd: &about.FunctionsCommand{}, PrevCmd: &cmd.StartCommand{}})
	bh.registry.Register("Обратная связь", &Command{Cmd: &about.FeedbackCommand{}, PrevCmd: &cmd.StartCommand{}})

	// Предсказания
	bh.registry.Register("Гороскоп на день", &Command{Cmd: &predictions.DayHoroscopeCommand{}, PrevCmd: &cmd.PredictionsCommand{}})
	bh.registry.Register("Гороскоп на неделю", &Command{Cmd: &predictions.WeekHoroscopeCommand{}, PrevCmd: &cmd.PredictionsCommand{}})
	bh.registry.Register("Гороскоп на месяц", &Command{Cmd: &predictions.MonthHoroscopeCommand{}, PrevCmd: &cmd.PredictionsCommand{}})
}

func (bh *Handler) HandleCommand(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := bh.registry.Get(message.Text)

	if command != nil {
		previousCommand := bh.commandHistory.GetPrevious()
		if previousCommand != nil && command.Cmd == previousCommand.Cmd {
			return
		}
		command.Cmd.Execute(api, message)
		if message.Text != backCommand {
			bh.commandHistory.Add(command)
		}
	} else {
		bh.handleUnknownCommand(api, message)
	}
}

func (bh *Handler) handleUnknownCommand(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не понимаю (кто понял тот понял)")
	api.Send(msg)
}

func HandleUpdate(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	bh := NewBotHandler(commandHistory)
	if update.Message != nil {
		bh.HandleCommand(api, update.Message)
	} else if update.CallbackQuery != nil {
		predictions.HandleCallbackQuery(api, update.CallbackQuery)
	}
}

type BackCommand struct {
	commandHistory *CommandHistory
}

func (c *BackCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	previousCommand := c.commandHistory.Pop()
	if previousCommand == nil {
		return
	}

	previousCommand.PrevCmd.Execute(api, message)
}

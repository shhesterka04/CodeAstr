package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/bot/commands"
	"glaphyra/internal/bot/commands/about"
	"glaphyra/internal/bot/commands/predictions"
)

type Command interface {
	Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message)
}

var (
	registry           = NewCommandRegistry()
	previousCommand    Command
	prePreviousCommand Command
)

func init() {
	registry.Register("/start", &cmd.StartCommand{})
	registry.Register("Назад", &BackCommand{})

	registry.Register("Предсказания", &cmd.PredictionsCommand{})
	registry.Register("Совместимость", &cmd.CompatibilityCommand{})
	registry.Register("Регистрация и настройки", &cmd.SettingsCommand{})
	registry.Register("О боте", &cmd.AboutCommand{})

	// О боте
	registry.Register("Кто я?", &about.WhoAmICommand{})
	registry.Register("Функции", &about.FunctionsCommand{})
	registry.Register("Обратная связь", &about.FeedbackCommand{})

	// Предсказания
	registry.Register("Гороскоп на день", &predictions.DayHoroscopeCommand{})
	registry.Register("Гороскоп на неделю", &predictions.WeekHoroscopeCommand{})
	registry.Register("Гороскоп на месяц", &predictions.MonthHoroscopeCommand{})
}

func HandleCommand(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := registry.Get(message.Text)
	if message.Text == "Назад" {
		command = prePreviousCommand
	}

	if command != nil {
		command.Execute(api, message)
		prePreviousCommand = previousCommand
		previousCommand = command
	} else {
		handleUnknownCommand(api, message)
	}
}

func handleUnknownCommand(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не понимаю (кто понял тот понял)")
	api.Send(msg)
}

type BackCommand struct{}

func (c *BackCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if prePreviousCommand == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No previous command")
		api.Send(msg)
		return
	}

	prePreviousCommand.Execute(api, message)
}

func HandleUpdate(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		HandleCommand(api, update.Message)
	} else if update.CallbackQuery != nil {
		predictions.HandleCallbackQuery(api, update.CallbackQuery)
	}
}

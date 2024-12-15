package update_handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot"
	cmd "glaphyra/internal/bot/commands"
	"glaphyra/internal/bot/commands/about"
	"glaphyra/internal/bot/commands/compatibility"
	"glaphyra/internal/bot/commands/dreambook"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/bot/commands/reflection"
	"glaphyra/internal/bot/commands/settings"
	"glaphyra/internal/bot/user_cmd_handler"
	"glaphyra/internal/llm/handlers"
)

const backCommand = "Назад"

type implUpdateHandler struct {
	registry   *CommandRegistry
	cmdHandler user_cmd_handler.UserCommandHandler
}

func NewUpdateHandler(cmdHandler user_cmd_handler.UserCommandHandler, userSrv userservice.UserService, gptApi handlers.Handler) bot.UpdateHandler {
	i := &implUpdateHandler{
		cmdHandler: cmdHandler,
		registry:   NewCommandRegistry(),
	}
	i.registerCommands(userSrv, gptApi)

	return i
}

func (i *implUpdateHandler) HandleUpdate(update tgbotapi.Update) {
	var userID int64
	var command bot.Command
	var message *tgbotapi.Message
	switch {
	case update.Message != nil:
		userID = update.Message.From.ID
		message = update.Message
		command = i.registry.Get(update.Message.Text)
	case update.CallbackQuery != nil:
		userID = update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID
		err := i.cmdHandler.HandleUserCallback(int64(messageID), userID, update.CallbackQuery)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err := i.cmdHandler.HandleUserCommand(userID, command, message)
	if err != nil {
		log.Println(err)
	}
	return
}

func (i *implUpdateHandler) registerCommands(userSrv userservice.UserService, gptApi handlers.Handler) {
	i.registry.Register("/start", cmd.NewStartCommand(userSrv))
	i.registry.Register(backCommand, &user_cmd_handler.BackCommand{})

	i.registry.Register("Предсказания", &cmd.PredictionsCommand{})
	i.registry.Register("Сонник", dreambook.NewDreambookCommand(userSrv, gptApi))
	i.registry.Register("Рефлексия", reflection.NewReflectionCommand(userSrv, gptApi))
	i.registry.Register("Совместимость", &cmd.CompatibilityCommand{})
	i.registry.Register("Регистрация и настройки", &cmd.SettingsCommand{})
	i.registry.Register("О боте", &cmd.AboutCommand{})

	// О боте
	i.registry.Register("Кто я?", &about.WhoAmICommand{})
	i.registry.Register("Функции", &about.FunctionsCommand{})
	i.registry.Register("Обратная связь", about.NewFeedbackCommand(userSrv))

	// Предсказания
	i.registry.Register("Гороскоп на день", predictions.NewHoroscopeCommand(userSrv, gptApi, predictions.Daily))
	i.registry.Register("Гороскоп на неделю", predictions.NewHoroscopeCommand(userSrv, gptApi, predictions.Weekly))
	i.registry.Register("Гороскоп на месяц", predictions.NewHoroscopeCommand(userSrv, gptApi, predictions.Monthly))

	// Совместимость
	i.registry.Register("Совместимость по знакам зодиака", compatibility.NewZodiakCompCommand(userSrv, gptApi))
	i.registry.Register("Совместимость по натальной карте", compatibility.NewNatalCompCommand(userSrv, gptApi))

	// Регистрация и настройки
	i.registry.Register("Выбор стиля общения", &settings.StyleCommand{})
	i.registry.Register("Серьезный стиль", settings.NewSetStyleCommand(userSrv, "Серьезный стиль"))
	i.registry.Register("Шутливый стиль", settings.NewSetStyleCommand(userSrv, "Шутливый стиль"))
	i.registry.Register("Дружелюбный стиль", settings.NewSetStyleCommand(userSrv, "Дружелюбный стиль"))
	i.registry.Register("Ввод даты рождения", settings.NewBirthCommand(userSrv))

}

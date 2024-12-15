package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot"
	"glaphyra/internal/bot/update_handler"
	"glaphyra/internal/bot/user_cmd_handler"
	"glaphyra/internal/llm/handlers"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler bot.UpdateHandler
}

func NewBot(tgToken string, userSrv userservice.UserService, gptApi handlers.Handler) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}
	userCmdHandler := user_cmd_handler.NewUserCmdHandler(api, userSrv)
	updateHandler := update_handler.NewUpdateHandler(userCmdHandler, userSrv, gptApi)

	return &Bot{api: api, handler: updateHandler}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		go func(update tgbotapi.Update) {
			b.handler.HandleUpdate(update)
		}(update)
	}
}

package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Command interface {
	Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error)
	IsTransfer() bool
}

type UpdateHandler interface {
	HandleUpdate(update tgbotapi.Update)
}

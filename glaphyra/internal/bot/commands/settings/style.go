package settings

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/pkg/log"
)

type StyleCommand struct{}

func (c *StyleCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	styleMessages := []string{
		"Выберите стиль общения",
	}

	aboutMessage := styleMessages[0]

	msg := tgbotapi.NewMessage(message.Chat.ID, aboutMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Серьезный стиль"),
			tgbotapi.NewKeyboardButton("Шутливый стиль"),
			tgbotapi.NewKeyboardButton("Дружелюбный стиль"),
		},
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	_, err := api.Send(msg)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (c *StyleCommand) IsTransfer() bool {
	return true
}

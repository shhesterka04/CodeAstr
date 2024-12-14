package cmd

import (
	"glaphyra/internal/pkg/log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AboutCommand struct{}

func (c *AboutCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	aboutMessages := []string{
		"Я создан, чтобы помогать тебе находить ответы среди звёзд и планет. 🌌 Узнай больше обо мне и о том, как я могу быть полезен.",
		"Я твой проводник по звёздам и планетам! 🌠 Узнай обо мне больше и посмотри, как я могу помочь тебе.",
		"Хочешь узнать, кто я и как я могу помочь? 🌌 Я здесь, чтобы ответить на твои вопросы и рассказать о звёздных возможностях.",
		"Я — твой звёздный помощник, и у меня есть много функций, которые могут тебе помочь! 🔮 Узнай обо мне больше.",
	}

	rand.Seed(time.Now().UnixNano())
	aboutMessage := aboutMessages[rand.Intn(len(aboutMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, aboutMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Кто я?"),
		},
		{
			tgbotapi.NewKeyboardButton("Обратная связь"),
		},
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	_, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *AboutCommand) IsTransfer() bool {
	return true
}

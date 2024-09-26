package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"time"
)

type CompatibilityCommand struct{}

func (c *CompatibilityCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	compatibilityMessages := []string{
		"Любопытно, как звезды влияют на твои отношения с другими? 💞 Узнай, насколько вы совместимы!",
		"Интересно, как звезды влияют на ваши отношения? 💕 Давай узнаем вашу совместимость по знакам зодиака!",
		"Звезды могут многое рассказать о любви и дружбе. 💖 Введи свой и знак другого человека, чтобы узнать, как складываются ваши отношения!",
		"Ты готов узнать, насколько вы совместимы? 💫 Введи свой знак и знак партнера, чтобы раскрыть секреты звёзд о ваших отношениях.",
	}

	rand.Seed(time.Now().UnixNano())
	compatibilityMessage := compatibilityMessages[rand.Intn(len(compatibilityMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, compatibilityMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Совместимость по знакам зодиака"),
			tgbotapi.NewKeyboardButton("Совместимость по натальной карте"),
		},
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	api.Send(msg)
}

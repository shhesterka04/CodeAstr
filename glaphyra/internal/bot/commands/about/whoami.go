package about

import (
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WhoAmICommand struct{}

func (c *WhoAmICommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	whoAmIMessages := []string{
		"Я — твой звездный спутник, который расскажет тебе о влиянии планет и их положении на твою жизнь. 🌟 Моя задача — помочь тебе лучше понять себя и окружающих через астрологию.",
		"Я — звёздный бот, который знает всё о гороскопах и астрологических прогнозах. 🌟 Хочешь узнать, как я могу помочь тебе?",
		"Моя задача — быть твоим путеводителем в мире астрологии. ✨ Узнай обо мне больше и как я могу быть полезен.",
		"Я создан для того, чтобы помогать тебе находить ответы на астрологические вопросы. 💫 Хочешь узнать больше обо мне?",
	}

	rand.Seed(time.Now().UnixNano())
	whoAmIMessage := whoAmIMessages[rand.Intn(len(whoAmIMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, whoAmIMessage)
	api.Send(msg)
}

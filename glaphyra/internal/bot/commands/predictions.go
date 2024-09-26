package cmd

import (
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PredictionsCommand struct{}

func (c *PredictionsCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	predictionMessages := []string{
		"Предсказания — это твой личный звездный советник! 🪐 Узнай, что звезды готовят для тебя сегодня, на неделю или даже на месяц. Выбери тип предсказания!",
		"Звёзды говорят... 💫 Выбирай, какой прогноз тебя интересует: на день, неделю или даже месяц. Я расскажу тебе все секреты планет!",
		"Ты готов узнать, что планеты готовят для тебя? 🪐 Выбери предсказание, и я раскрою тебе, что ждет тебя впереди!",
		"Предсказания ждут тебя! 🌞 Выбирай гороскоп на сегодня, неделю или месяц, и я расскажу тебе, что говорят звезды.",
	}

	rand.Seed(time.Now().UnixNano())
	predictionMessage := predictionMessages[rand.Intn(len(predictionMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, predictionMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Гороскоп на день"),
			tgbotapi.NewKeyboardButton("Гороскоп на неделю"),
		},
		{
			tgbotapi.NewKeyboardButton("Гороскоп на месяц"),
			tgbotapi.NewKeyboardButton("Персонализированные гороскопы"),
		},
		{
			tgbotapi.NewKeyboardButton("Натальная карта"),
			tgbotapi.NewKeyboardButton("Дни для сделок"),
		},
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	api.Send(msg)
}

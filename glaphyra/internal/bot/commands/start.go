package cmd

import (
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartCommand struct{}

func (c *StartCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	startMessages := []string{
		"Приветствую тебя в мире астрологии! 🌟 Я — твой личный звездный гид. Готов начать путешествие по звездам? Выбирай, что тебя интересует, и узнавай тайны Вселенной!",
		"Приветствую тебя в мире космических тайн! ✨ Я — твой личный астрологический советник. Готов отправиться в захватывающее путешествие по звёздам?",
		"Добро пожаловать! 🌟 Я знаю, что планеты приготовили для тебя на сегодня. Просто выбери, с чего начнём, и давай узнаем будущее!",
		"Здравствуй, искатель знаний! 🔮 Я помогу тебе разобраться в астрологических прогнозах и найти ответы среди звёзд. Что хочешь узнать первым?",
	}

	rand.Seed(time.Now().UnixNano())
	startMessage := startMessages[rand.Intn(len(startMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, startMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Предсказания"),
			tgbotapi.NewKeyboardButton("Совместимость"),
		},
		{
			tgbotapi.NewKeyboardButton("Регистрация и настройки"),
			tgbotapi.NewKeyboardButton("О боте"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	api.Send(msg)
}

package cmd

import (
	"glaphyra/internal/pkg/log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SettingsCommand struct{}

func (c *SettingsCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	settingsMessages := []string{
		"Здесь ты можешь настроить свой опыт общения со мной. 🤖 Я запомню твои предпочтения и буду отправлять тебе прогнозы в нужное время.",
		"Твой звездный опыт начинается здесь! 🔧 Настрой меня под себя: введи свою дату рождения и выбери, как я буду отправлять тебе прогнозы.",
		"Хочешь персонализировать общение? 🤖 Я могу подстроиться под твои предпочтения. Давай начнем с ввода даты рождения и настроек!",
		"Я запомню твои предпочтения и буду всегда готов к астрологическим прогнозам! 🔮 Настроим наш контакт для удобства.",
	}

	rand.Seed(time.Now().UnixNano())
	settingsMessage := settingsMessages[rand.Intn(len(settingsMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, settingsMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Управление подписками"),
		},
		{
			tgbotapi.NewKeyboardButton("Выбор стиля общения"),
			tgbotapi.NewKeyboardButton("Ввод даты рождения"),
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

func (c *SettingsCommand) IsTransfer() bool {
	return true
}

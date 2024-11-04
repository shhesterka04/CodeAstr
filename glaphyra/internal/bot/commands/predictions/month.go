package predictions

import (
	"glaphyra/internal/pkg/log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MonthHoroscopeCommand struct{}

func (c *MonthHoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	monthMessages := []string{
		"Месяц — это длинное путешествие. 🚀 Готов увидеть, как будут развиваться события? Звезды уже знают!",
		"Готов к большому астрологическому прогнозу? 🌓 Узнаем, как пройдёт этот месяц по звёздам!",
		"Целый месяц впереди — давай узнаем, что приготовили для тебя планеты! 🚀 Введи свой знак зодиака.",
		"Заглянем на месяц вперёд? 🌙 Звёзды могут раскрыть твои возможности и вызовы. Давай посмотрим!",
	}

	rand.Seed(time.Now().UnixNano())
	monthMessage := monthMessages[rand.Intn(len(monthMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, monthMessage)
	msg.ReplyMarkup = inlineKeyboard
	sentMsg, err := api.Send(msg) // todo нормально хэндлить ошибку
	if err != nil {
		return 0, log.Wrap(err)
	}

	backButtonMsg := tgbotapi.NewMessage(message.Chat.ID, "Нажмите 'Назад' для возврата")
	backButtonMsg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	sentMsg, err = api.Send(backButtonMsg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return sentMsg.Chat.ID, nil
}

func (c *MonthHoroscopeCommand) IsTransfer() bool {
	return true
}

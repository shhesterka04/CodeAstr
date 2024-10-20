package predictions

import (
	"glaphyra/internal/pkg/log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DayHoroscopeCommand struct{}

func (c *DayHoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	dayMessages := []string{
		"Каждый день звезды готовят что-то особенное. 🌞 Хочешь узнать, что ждет тебя сегодня? Введи свой знак зодиака!",
		"Каждый день полон астрологических загадок! 🌅 Давай посмотрим, что ждет тебя сегодня. Введи свой знак зодиака.",
		"Звезды говорят, что сегодня может быть особенный день! 🌟 Узнаем, что именно тебя ждёт? Введи свой знак!",
		"Что приготовили звёзды на этот день? 🌞 Введи свой знак зодиака, и я открою тебе сегодняшние секреты Вселенной.",
	}

	rand.Seed(time.Now().UnixNano())
	dayMessage := dayMessages[rand.Intn(len(dayMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, dayMessage)
	msg.ReplyMarkup = inlineKeyboard
	_, err := api.Send(msg) // todo нормально хэндлить ошибку
	if err != nil {
		return log.Wrap(err)
	}

	backButtonMsg := tgbotapi.NewMessage(message.Chat.ID, "Нажмите 'Назад' для возврата")
	backButtonMsg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	_, err = api.Send(backButtonMsg)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (c *DayHoroscopeCommand) IsTransfer() bool {
	return true
}

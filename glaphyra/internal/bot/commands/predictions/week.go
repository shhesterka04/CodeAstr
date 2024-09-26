package predictions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"time"
)

type WeekHoroscopeCommand struct{}

func (c *WeekHoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	weekMessages := []string{
		"Заглянем на неделю вперед? 📅 Узнай, какие события ожидают тебя в ближайшие 7 дней!",
		"Интересно, что звезды предсказывают на ближайшие 7 дней? 📅 Давай заглянем в будущее на неделю вперед!",
		"Неделя — это маленькая вселенная возможностей. 🌠 Узнай, что приготовили для тебя планеты!",
		"Тебя ждёт интересная неделя! 🌟 Введи свой знак зодиака, и я покажу тебе, чего ожидать в ближайшие дни.",
	}

	rand.Seed(time.Now().UnixNano())
	weekMessage := weekMessages[rand.Intn(len(weekMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, weekMessage)
	msg.ReplyMarkup = inlineKeyboard
	api.Send(msg)

	backButtonMsg := tgbotapi.NewMessage(message.Chat.ID, "Нажмите 'Назад' для возврата")
	backButtonMsg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	api.Send(backButtonMsg)
}

//TODO: кнопка назад плохо реализована, надо ее переледать - иначе на этой стадии зацикливание

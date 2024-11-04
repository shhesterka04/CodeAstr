package predictions

import (
	"glaphyra/internal/pkg/log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WeekHoroscopeCommand struct{}

func (c *WeekHoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
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

func (c *WeekHoroscopeCommand) IsTransfer() bool {
	return true
}

//TODO: кнопка назад плохо реализована, надо ее переледать - иначе на этой стадии зацикливание

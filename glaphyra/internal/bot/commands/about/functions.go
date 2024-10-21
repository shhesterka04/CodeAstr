package about

import (
	"glaphyra/internal/pkg/log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FunctionsCommand struct{}

func (c *FunctionsCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	functionsMessages := []string{
		"Я могу многое: от ежедневных гороскопов до анализа совместимости и натальных карт. 🛠️ Выбирай любую из моих функций и начнем работу со звездами!",
		"Я могу многое: от гороскопов до анализа совместимости. 🌌 Выбери любую из моих функций и начни работать с астрологией прямо сейчас!",
		"Хочешь узнать, что я умею? 🌠 Я могу помочь тебе с ежедневными гороскопами, совместимостью и многим другим.",
		"Мои функции обширны: от простых прогнозов до сложных натальных карт. 🔮 Узнай, как звезды могут помочь тебе!",
	}

	rand.Seed(time.Now().UnixNano())
	functionsMessage := functionsMessages[rand.Intn(len(functionsMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, functionsMessage)
	_, err := api.Send(msg)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (c *FunctionsCommand) IsTransfer() bool {
	return false
}

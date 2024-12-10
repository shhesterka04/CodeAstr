package compatibility

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"
	"math/rand"
	"sync"
	"time"
)

type NatalCompCommand struct {
	gptApi  handlers.Handler
	userSrv userservice.UserService
}

func NewNatalCompCommand(userSrv userservice.UserService, gptApi handlers.Handler) *NatalCompCommand {
	return &NatalCompCommand{
		userSrv: userSrv,
		gptApi:  gptApi,
	}
}

var natalMag = []string{
	"Хочешь узнать глубже, как планеты влияют на вас? 🌌 Введи свои данные и данные партнера для анализа.",
	"Хочешь узнать больше о своих отношениях? 🔭 Введи даты рождения для детального анализа совместимости по натальной карте!",
	"Глубокий анализ совместимости по натальной карте раскроет скрытые аспекты ваших отношений. 🌌 Введи даты рождения!",
	"Натальная карта — это ключ к пониманию ваших отношений на новом уровне. 🌠 Введи свои данные и данные партнера!",
}

func (c *NatalCompCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(natalMag))
	msgVar := natalMag[randomIndex]

	msg := tgbotapi.NewMessage(message.Chat.ID, msgVar)
	msg.ReplyMarkup = predictions.InlineKeyboard
	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	backButtonMsg := tgbotapi.NewMessage(message.Chat.ID, predictions.BackTip)
	backButtonMsg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(predictions.BackCmd),
		),
	)

	if _, err = api.Send(backButtonMsg); err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *NatalCompCommand) SendPrompt(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) error {
	wg := sync.WaitGroup{}
	usr, err := c.userSrv.FindByID(context.Background(), callback.From.ID)
	if err != nil {
		return log.Wrap(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		typingMsg := tgbotapi.NewChatAction(callback.Message.Chat.ID, tgbotapi.ChatTyping)
		api.Send(typingMsg)
	}()

	yaResponse, err := c.gptApi.CallAPI(dto.RequestDTO{
		UserMessage: fmt.Sprintf(
			predictionPrompt,
			c.horoscopeType,
			callback.Data,
			usr.Style),
	})

	wg.Wait()

	if err != nil {
		return log.Wrap(err)
	}
	response, ok := yaResponse.(dto.ResponseDTO)
	if !ok {
		return log.Wrap(err)
	}
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, response.Result)
	msg.ParseMode = "Markdown"
	if _, err = api.Send(msg); err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (c *NatalCompCommand) IsTransfer() bool {
	return true
}

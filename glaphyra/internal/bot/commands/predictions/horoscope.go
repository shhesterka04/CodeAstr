package predictions

import (
	"context"
	"fmt"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
	"math/rand"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/pkg/log"
)

type HoroscopeCommand struct {
	userSrv       userservice.UserService
	gptApi        handlers.Handler
	horoscopeType string
}

func NewHoroscopeCommand(userSrv userservice.UserService, gptApi handlers.Handler, horoscopeType string) *HoroscopeCommand {
	return &HoroscopeCommand{
		userSrv:       userSrv,
		gptApi:        gptApi,
		horoscopeType: horoscopeType,
	}
}

func (c *HoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(messagesMap[c.horoscopeType]))
	msgVars := messagesMap[c.horoscopeType][randomIndex]

	msg := tgbotapi.NewMessage(message.Chat.ID, msgVars)
	msg.ReplyMarkup = inlineKeyboard
	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	backButtonMsg := tgbotapi.NewMessage(message.Chat.ID, BackTip)
	backButtonMsg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BackCmd),
		),
	)

	if _, err = api.Send(backButtonMsg); err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *HoroscopeCommand) SendPrompt(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) error {
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

func (c *HoroscopeCommand) IsTransfer() bool {
	return true
}

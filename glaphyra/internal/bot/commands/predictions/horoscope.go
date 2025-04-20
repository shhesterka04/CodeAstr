package predictions

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/llm/handlers"
	yadto "glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"
	redisclient "glaphyra/internal/pkg/redis"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HoroscopeCommand struct {
	userSrv       userservice.UserService
	gptApi        handlers.Handler
	horoscopeType string
	redisClient   *redis.Client
}

func NewHoroscopeCommand(userSrv userservice.UserService, gptApi handlers.Handler, horoscopeType string, redisClient *redis.Client) *HoroscopeCommand {
	return &HoroscopeCommand{
		userSrv:       userSrv,
		gptApi:        gptApi,
		horoscopeType: horoscopeType,
		redisClient:   redisClient,
	}
}

func (c *HoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	usr, err := c.userSrv.FindByID(context.Background(), message.From.ID)
	if err != nil {
		return 0, log.Wrap(err)
	}
	if usr.Tokens <= 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "У вас закончились токены.. Приходите завтра")
		if _, err := api.Send(msg); err != nil {
			return 0, log.Wrap(err)
		}

		return 0, nil
	}

	if err = c.userSrv.Update(context.Background(), &dto.UpdateUserRequest{
		TgID:   message.From.ID,
		Tokens: usr.Tokens - 5,
	}); err != nil {
		return 0, log.Wrap(err)
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(messagesMap[c.horoscopeType]))
	msgVars := messagesMap[c.horoscopeType][randomIndex]

	msg := tgbotapi.NewMessage(message.Chat.ID, msgVars)
	msg.ReplyMarkup = InlineKeyboard
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

	cacheKey := fmt.Sprintf("horoscope:%s:%s", c.horoscopeType, callback.Data)

	cachedResponse, err := redisclient.GetCache(c.redisClient, cacheKey)
	if err == nil && cachedResponse != "" {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, cachedResponse)
		msg.ParseMode = "Markdown"
		if _, err := api.Send(msg); err != nil {
			return log.Wrap(err)
		}
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		typingMsg := tgbotapi.NewChatAction(callback.Message.Chat.ID, tgbotapi.ChatTyping)
		api.Send(typingMsg)
	}()

	yaResponse, err := c.gptApi.CallAPI(yadto.RequestDTO{
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

	response, ok := yaResponse.(yadto.ResponseDTO)
	if !ok {
		return log.Wrap(fmt.Errorf("unexpected response type from YaGPT"))
	}

	var expiration time.Duration
	switch c.horoscopeType {
	case "daily":
		expiration = time.Until(endOfDay())
	case "weekly":
		expiration = time.Until(endOfWeek())
	case "monthly":
		expiration = time.Until(endOfMonth())
	default:
		expiration = 24 * time.Hour
	}

	if err := redisclient.SetCache(c.redisClient, cacheKey, response.Result, expiration); err != nil {
		log.Error(err)
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, response.Result)
	msg.ParseMode = "Markdown"
	if _, err := api.Send(msg); err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (c *HoroscopeCommand) IsTransfer() bool {
	return true
}

func endOfDay() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
}

func endOfWeek() time.Time {
	now := time.Now()
	daysUntilSunday := time.Sunday - now.Weekday()
	if daysUntilSunday < 0 {
		daysUntilSunday += 7
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).AddDate(0, 0, int(daysUntilSunday))
}

func endOfMonth() time.Time {
	now := time.Now()
	year, month, _ := now.Date()
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, now.Location())
	return firstDayNextMonth.Add(-time.Second)
}

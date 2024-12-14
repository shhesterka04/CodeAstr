package compatibility

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"
)

type ZodiakCompCommand struct {
	gptApi      handlers.Handler
	userSrv     userservice.UserService
	firstZodiak sync.Map
}

func NewZodiakCompCommand(userSrv userservice.UserService, gptApi handlers.Handler) *ZodiakCompCommand {
	return &ZodiakCompCommand{
		userSrv:     userSrv,
		gptApi:      gptApi,
		firstZodiak: sync.Map{},
	}
}

var (
	zodiakMag = []string{
		"Введи свой знак зодиака и знак партнера, чтобы узнать, как звезды видят ваши отношения.",
		"Звезды могут помочь понять твои отношения. 💘 Введи свой знак зодиака и знак другого человека, и я покажу, как вы совместимы.",
		"Любопытно, как планеты влияют на ваши отношения? 💑 Введи два знака, и мы посмотрим на вашу астрологическую совместимость!",
		"Звезды подскажут, насколько гармоничны ваши отношения. 🌙 Введи свои и партнёрские знаки, чтобы узнать больше!",
	}

	compatibilityPromt = "Привет, бот! Мне нужна совместимость двух знаков зодиака: " +
		"- Первый знак зодиака: %s " +
		"- Второй знак зодиака: %s " +
		"Требования: " +
		"1. Прогноз должен быть детализированным и точным." +
		"2. Учти особенности целевой аудитории: студенты, офисные сотрудники, домохозяйки." +
		"3. Предсказания должны быть позитивными и мотивирующими, но также реалистичными." +
		"4. Используй стиль общения: %s." +
		"5. Не надо здороваться и прощаться. " +
		"6. Примерно 50 слов " +
		"7. Пиши не по пунктам, а сплошным текстом " +
		"8. Используй эмодзи и смайлики для улучшения текста. " +
		"Благодарю за помощь! "
)

func (c *ZodiakCompCommand) IsFirstFull(userID int64) bool {
	_, ok := c.firstZodiak.Load(userID)
	return ok
}

func (c *ZodiakCompCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(zodiakMag))
	msgVar := zodiakMag[randomIndex]
	msg := tgbotapi.NewMessage(message.Chat.ID, msgVar)
	if _, err := api.Send(msg); err != nil {
		return 0, log.Wrap(err)
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "Выберите ваш знак зодиака:")
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

func (c *ZodiakCompCommand) GetFirstSignSendSecond(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) (int64, error) {
	c.firstZodiak.Store(callback.From.ID, callback.Data)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Выберите знак зодиака партнера:")
	msg.ReplyMarkup = predictions.InlineKeyboard
	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	backButtonMsg := tgbotapi.NewMessage(callback.Message.Chat.ID, predictions.BackTip)
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

func (c *ZodiakCompCommand) GetSecondSignSendResult(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) error {
	wg := sync.WaitGroup{}
	usr, err := c.userSrv.FindByID(context.Background(), callback.From.ID)
	if err != nil {
		return log.Wrap(err)
	}

	firstSign, _ := c.firstZodiak.Load(callback.From.ID)
	c.firstZodiak.Delete(callback.From.ID)
	wg.Add(1)

	go func() {
		defer wg.Done()
		typingMsg := tgbotapi.NewChatAction(callback.Message.Chat.ID, tgbotapi.ChatTyping)
		api.Send(typingMsg)
	}()

	yaResponse, err := c.gptApi.CallAPI(dto.RequestDTO{
		UserMessage: fmt.Sprintf(
			compatibilityPromt,
			firstSign,
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

func (c *ZodiakCompCommand) IsTransfer() bool {
	return true
}

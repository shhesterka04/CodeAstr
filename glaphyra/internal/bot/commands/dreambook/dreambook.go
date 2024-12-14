package dreambook

import (
	"context"
	"fmt"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
	"math/rand"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/pkg/log"
)

var (
	msgVars = []string{
		"Загадочные сны часто скрывают в себе подсказки Вселенной. 😴 Расскажи, что тебе приснилось, и я помогу расшифровать этот сон.",
		"Интересует значение твоего сна? 🌙 Введи его детали, и я расскажу, что могут значить эти символы с астрологической точки зрения.",
		"Твои сны — это зеркала подсознания. 🔮 Введи свой сон, и я помогу тебе понять его с астрологической точки зрения.",
	}

	dreambookPromt = "Привет, бот! Мне нужен сонник для пользователя:- " +
		"Вот описание сна: %s " +
		"Требования: " +
		"1. Сонник должен быть детализированным и точным." +
		"2. Учти особенности целевой аудитории: студенты, офисные сотрудники, домохозяйки." +
		"3. Сонник должен быть позитивным и мотивирующим, но также реалистичным." +
		"4. Используй стиль общения: %s." +
		"5. Не надо здороваться и прощаться. " +
		"6. Примерно 100 слов " +
		"7. Пиши не по пунктам, а сплошным текстом " +
		"8. Используй эмодзи и смайлики для улучшения текста. " +
		"Благодарю за помощь! "
)

type DreambookCommand struct {
	userSrv userservice.UserService
	gptApi  handlers.Handler
}

func NewDreambookCommand(userSrv userservice.UserService, gptApi handlers.Handler) *DreambookCommand {
	return &DreambookCommand{
		userSrv: userSrv,
		gptApi:  gptApi,
	}
}

func (c *DreambookCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(msgVars))
	msgVar := msgVars[randomIndex]

	msg := tgbotapi.NewMessage(message.Chat.ID, msgVar)

	if _, err := api.Send(msg); err != nil {
		return 0, log.Wrap(err)
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, опишите свой сон.")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(predictions.BackCmd),
		),
	)

	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *DreambookCommand) SendResult(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	wg := sync.WaitGroup{}
	usr, err := c.userSrv.FindByID(context.Background(), message.From.ID)
	if err != nil {
		return 0, log.Wrap(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		typingMsg := tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)
		api.Send(typingMsg)
	}()

	yaResponse, err := c.gptApi.CallAPI(dto.RequestDTO{
		UserMessage: fmt.Sprintf(
			dreambookPromt,
			message.Text,
			usr.Style),
	})

	wg.Wait()

	if err != nil {
		return 0, log.Wrap(err)
	}
	response, ok := yaResponse.(dto.ResponseDTO)
	if !ok {
		return 0, log.Wrap(err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, response.Result)
	msg.ParseMode = "Markdown"
	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *DreambookCommand) IsTransfer() bool {
	return true
}

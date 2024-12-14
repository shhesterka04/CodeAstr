package compatibility

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"
)

type NatalCompCommand struct {
	gptApi    handlers.Handler
	userSrv   userservice.UserService
	usersData sync.Map
}

func NewNatalCompCommand(userSrv userservice.UserService, gptApi handlers.Handler) *NatalCompCommand {
	return &NatalCompCommand{
		userSrv: userSrv,
		gptApi:  gptApi,
	}
}

var (
	natalMsgs = []string{
		"Хочешь узнать глубже, как планеты влияют на вас? 🌌 Введи свои данные и данные партнера для анализа.",
		"Хочешь узнать больше о своих отношениях? 🔭 Введи даты рождения для детального анализа совместимости по натальной карте!",
		"Глубокий анализ совместимости по натальной карте раскроет скрытые аспекты ваших отношений. 🌌 Введи даты рождения!",
		"Натальная карта — это ключ к пониманию ваших отношений на новом уровне. 🌠 Введи свои данные и данные партнера!",
	}

	compatibilityNatalPromt = "Привет, бот! Мне нужна совместимость по натальной карте: " +
		"- Данные для натальной карты первого партнера: %s " +
		"- Данные для натальной карты второго партнера: %s " +
		"В расшифровке синастрии должны быть использованы основные аспекты " +
		"(трин (тригон), квадратура, соединение и оппозиция), " +
		"дополнительные аспекты (квиконс и секстиль) " +
		"не надо на самом деле использовать какие-то данные, просто сделай вид, что так делаешь" +
		"Требования: " +
		"1. Прогноз должен быть детализированным и точным." +
		"2. Учти особенности целевой аудитории: студенты, офисные сотрудники, домохозяйки." +
		"3. Предсказания должны быть позитивными и мотивирующими, но также реалистичными." +
		"4. Используй стиль общения: %s." +
		"5. Не надо здороваться и прощаться. " +
		"6. Примерно 100 слов " +
		"7. Пиши не по пунктам, а сплошным текстом " +
		"8. Используй эмодзи и смайлики для улучшения текста. " +
		"Благодарю за помощь! "
)

func (c *NatalCompCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(natalMsgs))
	msgVar := natalMsgs[randomIndex]
	msg := tgbotapi.NewMessage(message.Chat.ID, msgVar)
	if _, err := api.Send(msg); err != nil {
		return 0, log.Wrap(err)
	}
	msg = tgbotapi.NewMessage(message.Chat.ID, "Введите дату, время и место вашего рождения и рождения партнера в формате\n`ДД.ММ.ГГГГ ЧЧ:ММ Город\nДД.ММ.ГГГГ ЧЧ:ММ Город`\n\nПример:\n`01.01.2000 12:00 Москва\n01.01.2000 12:00 Москва`")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(predictions.BackCmd),
		),
	)
	msg.ParseMode = "Markdown"
	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *NatalCompCommand) SendResult(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	data := strings.Split(message.Text, "\n")

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
			compatibilityNatalPromt,
			data[0],
			data[1],
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

func (c *NatalCompCommand) IsTransfer() bool {
	return true
}

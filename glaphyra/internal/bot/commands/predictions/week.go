package predictions

import (
	"context"
	"fmt"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WeekHoroscopeCommand struct {
	userSrv userservice.UserService
	gptApi  handlers.Handler
}

func NewWeekHoroscopeCommand(userSrv userservice.UserService, gptApi handlers.Handler) *WeekHoroscopeCommand {
	return &WeekHoroscopeCommand{
		userSrv: userSrv,
		gptApi:  gptApi,
	}
}

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

	_, err = api.Send(backButtonMsg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *WeekHoroscopeCommand) SendPromt(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) error {
	usr, err := c.userSrv.FindByID(context.Background(), callback.From.ID)
	if err != nil {
		return log.Wrap(err)
	}
	yaResponse, err := c.gptApi.CallAPI(dto.RequestDTO{
		UserMessage: fmt.Sprintf("Привет, бот! Мне нужен гороскоп для пользователя:- "+
			"Тип прогноза: еженедельный "+
			"- Знак зодиака: %s "+
			"Требования: "+
			"1. Прогноз должен быть детализированным и точным."+
			"2. Учти особенности целевой аудитории: студенты, офисные сотрудники, домохозяйки."+
			"3. Предсказания должны быть позитивными и мотивирующими, но также реалистичными."+
			"4. Используй стиль общения: %s."+
			"Благодарю за помощь!", callback.Data, usr.Style),
	})
	if err != nil {
		return log.Wrap(err)
	}
	response, ok := yaResponse.(dto.ResponseDTO)
	if !ok {
		return log.Wrap(err)
	}
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, response.Result)
	_, err = api.Send(msg)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (c *WeekHoroscopeCommand) IsTransfer() bool {
	return true
}

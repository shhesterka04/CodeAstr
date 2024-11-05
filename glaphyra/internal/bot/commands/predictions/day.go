package predictions

import (
	"context"
	"fmt"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/pkg/log"
)

type DayHoroscopeCommand struct {
	userSrv userservice.UserService
	gptApi  handlers.Handler
}

func NewDayHoroscopeCommand(userSrv userservice.UserService, gptApi handlers.Handler) *DayHoroscopeCommand {
	return &DayHoroscopeCommand{
		userSrv: userSrv,
		gptApi:  gptApi,
	}
}

func (c *DayHoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	dayMessages := []string{
		"Каждый день звезды готовят что-то особенное. 🌞 Хочешь узнать, что ждет тебя сегодня? Введи свой знак зодиака!",
		"Каждый день полон астрологических загадок! 🌅 Давай посмотрим, что ждет тебя сегодня. Введи свой знак зодиака.",
		"Звезды говорят, что сегодня может быть особенный день! 🌟 Узнаем, что именно тебя ждёт? Введи свой знак!",
		"Что приготовили звёзды на этот день? 🌞 Введи свой знак зодиака, и я открою тебе сегодняшние секреты Вселенной.",
	}

	rand.Seed(time.Now().UnixNano())
	dayMessage := dayMessages[rand.Intn(len(dayMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, dayMessage)
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

func (c *DayHoroscopeCommand) SendPromt(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) error {
	usr, err := c.userSrv.FindByID(context.Background(), callback.From.ID)
	if err != nil {
		return log.Wrap(err)
	}
	yaResponse, err := c.gptApi.CallAPI(dto.RequestDTO{
		UserMessage: fmt.Sprintf("Привет, бот! Мне нужен гороскоп для пользователя:- "+
			"Тип прогноза: ежедневный "+
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

func (c *DayHoroscopeCommand) IsTransfer() bool {
	return true
}

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

type MonthHoroscopeCommand struct {
	userSrv userservice.UserService
	gptApi  handlers.Handler
}

func NewMonthHoroscopeCommand(userSrv userservice.UserService, gptApi handlers.Handler) *MonthHoroscopeCommand {
	return &MonthHoroscopeCommand{
		userSrv: userSrv,
		gptApi:  gptApi,
	}
}

func (c *MonthHoroscopeCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	monthMessages := []string{
		"Месяц — это длинное путешествие. 🚀 Готов увидеть, как будут развиваться события? Звезды уже знают!",
		"Готов к большому астрологическому прогнозу? 🌓 Узнаем, как пройдёт этот месяц по звёздам!",
		"Целый месяц впереди — давай узнаем, что приготовили для тебя планеты! 🚀 Введи свой знак зодиака.",
		"Заглянем на месяц вперёд? 🌙 Звёзды могут раскрыть твои возможности и вызовы. Давай посмотрим!",
	}

	rand.Seed(time.Now().UnixNano())
	monthMessage := monthMessages[rand.Intn(len(monthMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, monthMessage)
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

func (c *MonthHoroscopeCommand) SendPromt(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) error {
	usr, err := c.userSrv.FindByID(context.Background(), callback.From.ID)
	if err != nil {
		return log.Wrap(err)
	}
	yaResponse, err := c.gptApi.CallAPI(dto.RequestDTO{
		UserMessage: fmt.Sprintf("Привет, бот! Мне нужен гороскоп для пользователя:- "+
			"Тип прогноза: ежемесячный "+
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

func (c *MonthHoroscopeCommand) IsTransfer() bool {
	return true
}

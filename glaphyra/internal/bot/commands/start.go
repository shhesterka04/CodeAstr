package cmd

import (
	"context"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	errs "glaphyra/internal/pkg/errors"
	"glaphyra/internal/pkg/log"
)

type StartCommand struct {
	userSrv userservice.UserService
}

func NewStartCommand(userSrv userservice.UserService) *StartCommand {
	return &StartCommand{userSrv: userSrv}
}

func (c *StartCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()
	_, err := c.userSrv.FindByID(ctx, message.From.ID)
	isNewUser := false
	switch {
	case errs.As[errs.ObjectNotFoundError](err):
		isNewUser = true
		err = c.userSrv.Create(
			ctx,
			&dto.CreateUserRequest{
				TgID:     message.From.ID,
				Username: message.From.UserName,
				Type:     dto.DefaultUser,
			},
		)
		if err != nil {
			return 0, log.Wrap(err)
		}
	case err != nil:
		return 0, log.Wrap(err)
	}

	msg := c.getReplyMsg(message, isNewUser)

	_, err = api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *StartCommand) getReplyMsg(message *tgbotapi.Message, isNewUser bool) tgbotapi.MessageConfig {
	greatings := []string{
		"Приветствую тебя в мире астрологии! ",
		"Приветствую тебя в мире космических тайн! ",
		"Добро пожаловать! ",
		"Здравствуй, искатель знаний! ",
	}

	startMessages := []string{
		"🌟 Я — твой личный звездный гид. Готов начать путешествие по звездам? Выбирай, что тебя интересует, и узнавай тайны Вселенной!",
		"✨ Я — твой личный астрологический советник. Готов отправиться в захватывающее путешествие по звёздам?",
		"🌟 Я знаю, что планеты приготовили для тебя на сегодня. Просто выбери, с чего начнём, и давай узнаем будущее!",
		"🔮 Я помогу тебе разобраться в астрологических прогнозах и найти ответы среди звёзд. Что хочешь узнать первым?",
	}

	settingsMessage := "\nДля более релевантных ответов советуем тебе заполнить данные о себе\nТы можешь найти их перейдя в раздел Регистрация и настройки -> О себе"

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(startMessages))
	startMessage := startMessages[idx]
	if isNewUser {
		startMessage = greatings[idx] + startMessage + settingsMessage
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, startMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("🏗️ Рефлексия"),
			tgbotapi.NewKeyboardButton("Сонник"),
		},
		{
			tgbotapi.NewKeyboardButton("Предсказания"),
			tgbotapi.NewKeyboardButton("Совместимость"),
		},
		{
			tgbotapi.NewKeyboardButton("Регистрация и настройки"),
			tgbotapi.NewKeyboardButton("О боте"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	return msg
}

func (c *StartCommand) IsTransfer() bool {
	return true
}

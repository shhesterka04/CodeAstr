package settings

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"
)

type BirthTime struct {
	userSrv userservice.UserService
}

func NewBirthTimeCommand(userSrv userservice.UserService) *BirthTime {
	return &BirthTime{userSrv: userSrv}
}

func (c *BirthTime) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	birthMessage := "Введите время, во сколько Вы родились в свободном формате"

	msg := tgbotapi.NewMessage(message.Chat.ID, birthMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	_, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *BirthTime) SetBirthTime(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()

	err := c.userSrv.Update(
		ctx,
		&dto.UpdateUserRequest{
			TgID:      message.From.ID,
			BirthTime: message.Text,
		},
	)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *BirthTime) IsTransfer() bool {
	return true
}

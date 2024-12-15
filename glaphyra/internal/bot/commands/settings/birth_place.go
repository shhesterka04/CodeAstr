package settings

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"
)

type BirthPlace struct {
	userSrv userservice.UserService
}

func NewBirthPlaceCommand(userSrv userservice.UserService) *BirthPlace {
	return &BirthPlace{userSrv: userSrv}
}

func (c *BirthPlace) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	birthMessage := "Введите город, в котором Вы родились"

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

func (c *BirthPlace) SetBirthPlace(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()

	err := c.userSrv.Update(
		ctx,
		&dto.UpdateUserRequest{
			TgID:       message.From.ID,
			BirthPlace: message.Text,
		},
	)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *BirthPlace) IsTransfer() bool {
	return true
}

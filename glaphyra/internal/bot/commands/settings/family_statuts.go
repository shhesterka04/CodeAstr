package settings

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"
)

type FamilyStatus struct {
	userSrv userservice.UserService
}

func NewFamilyStatusCommand(userSrv userservice.UserService) *FamilyStatus {
	return &FamilyStatus{userSrv: userSrv}
}

func (c *FamilyStatus) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	birthMessage := "Введите семейное положение в свободном формате\nНапример: свободен(-на), есть парень/девушка, замужем/женат, есть ребенок"

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

func (c *FamilyStatus) SetFamilyStatus(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()

	err := c.userSrv.Update(
		ctx,
		&dto.UpdateUserRequest{
			TgID:         message.From.ID,
			FamilyStatus: message.Text,
		},
	)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *FamilyStatus) IsTransfer() bool {
	return true
}

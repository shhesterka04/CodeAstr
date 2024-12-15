package settings

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"
)

type TypeOfActivity struct {
	userSrv userservice.UserService
}

func NewTypeOfActivityCommand(userSrv userservice.UserService) *TypeOfActivity {
	return &TypeOfActivity{userSrv: userSrv}
}

func (c *TypeOfActivity) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	birthMessage := "Введите, чем вы занимаетесь. Например: домохозяйка, офисный работник, студент\nОтветить можно в свободной форме"

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

func (c *TypeOfActivity) SetTypeOfActivity(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()

	err := c.userSrv.Update(
		ctx,
		&dto.UpdateUserRequest{
			TgID:           message.From.ID,
			TypeOfActivity: message.Text,
		},
	)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *TypeOfActivity) IsTransfer() bool {
	return true
}

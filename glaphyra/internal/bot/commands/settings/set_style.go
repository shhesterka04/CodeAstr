package settings

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"
)

type SetStyleCommand struct {
	userSrv userservice.UserService
	style   string
}

func NewSetStyleCommand(userSrv userservice.UserService, style string) *SetStyleCommand {
	return &SetStyleCommand{
		userSrv: userSrv,
		style:   style,
	}
}

func (c *SetStyleCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	ctx := context.Background()
	err := c.userSrv.Update(
		ctx,
		&dto.UpdateUserRequest{
			TgID:  message.From.ID,
			Style: c.style,
		},
	)
	if err != nil {
		return log.WrapErr(err)
	}

	styleMessages := []string{
		fmt.Sprintf("Вы выбрали %s общения", c.style),
	}

	aboutMessage := styleMessages[0]

	msg := tgbotapi.NewMessage(message.Chat.ID, aboutMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Серьезный стиль"),
			tgbotapi.NewKeyboardButton("Шутливый стиль"),
			tgbotapi.NewKeyboardButton("Дружелюбный стиль"),
		},
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	_, err = api.Send(msg)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (c *SetStyleCommand) IsTransfer() bool {
	return false
}

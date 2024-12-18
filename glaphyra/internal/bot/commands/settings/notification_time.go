package settings

import (
	"context"
	"strconv"

	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotificationTime struct {
	userSrv userservice.UserService
}

func NewNotificationTimeCommand(userSrv userservice.UserService) *NotificationTime {
	return &NotificationTime{userSrv: userSrv}
}

func (c *NotificationTime) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	birthMessage := "Введите, в какой час по МСК вы хотели бы получать форму рефлексии\nВведите одно число от 1 до 23"

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

func (c *NotificationTime) SetNotificationTime(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()

	answerMsg := ""

	hour, errParse := strconv.ParseInt(message.Text, 10, 32)
	if errParse != nil {
		answerMsg = "Мы не смогли распознать, когда вы хотите получать уведомление о рефлексии :( Попробуйте снова"
		errParse = &ValidationError{}
	}

	if errParse == nil && (hour < 0 || hour > 23) {
		errParse = &ValidationError{}
		answerMsg = "Мы не смогли распознать, когда вы хотите получать уведомление о рефлексии :( Попробуйте снова"
	}

	if errParse != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, answerMsg)

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

		return 0, errParse
	}

	err := c.userSrv.Update(
		ctx,
		&dto.UpdateUserRequest{
			TgID:             message.From.ID,
			NotificationTime: int32(hour),
		},
	)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *NotificationTime) IsTransfer() bool {
	return true
}

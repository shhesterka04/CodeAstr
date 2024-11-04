package settings

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"
)

type Birth struct {
	userSrv userservice.UserService
}

func NewBirthCommand(userSrv userservice.UserService) *Birth {
	return &Birth{userSrv: userSrv}
}

func (c *Birth) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()
	user, err := c.userSrv.FindByID(ctx, message.From.ID)
	if err != nil {
		return 0, log.Wrap(err)
	}

	birthMessage := "Введите дату рождения в формате ДД.ММ.ГГГГ"

	defaultTime := time.Time{}
	if user.BirthDate != defaultTime {
		birthMessage = fmt.Sprintf(
			"Вы уже указали дату рождения %s. Чтобы изменить, введите дату рождения в формате ДД.ММ.ГГГГ",
			user.BirthDate.Format("02.01.2006"),
		)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, birthMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	_, err = api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *Birth) SetBirth(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	ctx := context.Background()

	answerMsg := "Вы успешно изменили дату рождения!"

	const layout = "02.01.2006"

	birthDate, err := time.Parse(layout, message.Text)
	if err != nil {
		answerMsg = "Мы не смогли распознать дату рождения :( Попробуйте снова"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, answerMsg)

	err = c.userSrv.Update(
		ctx,
		&dto.UpdateUserRequest{
			TgID:      message.From.ID,
			BirthDate: birthDate,
		},
	)
	if err != nil {
		return 0, log.Wrap(err)
	}

	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Назад"),
		},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	_, err = api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil
}

func (c *Birth) IsTransfer() bool {
	return true
}

package settings

import (
	"context"
	"math/rand"
	"sync"
	"time"

	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FeedbackCommand struct {
	userSrv          userservice.UserService
	userIDToFeedback sync.Map
}

func NewFeedbackCommand(userSrv userservice.UserService) *FeedbackCommand {
	return &FeedbackCommand{
		userSrv:          userSrv,
		userIDToFeedback: sync.Map{},
	}
}

func (c *FeedbackCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	feedbackMessages := []string{
		"Тебе есть, что мне сказать? 💬 Я всегда готов выслушать твои отзывы и предложения.",
		"У тебя есть вопросы или предложения? 💬 Я всегда готов выслушать твоё мнение!",
		"Нужна помощь или хочешь оставить отзыв? 🌟 Напиши мне, и я постараюсь улучшить наш диалог!",
		"Свяжись со мной, если у тебя есть вопросы или предложения. 📬 Твой отзыв важен для меня!",
	}

	rand.Seed(time.Now().UnixNano())
	feedbackMessage := feedbackMessages[rand.Intn(len(feedbackMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, feedbackMessage)

	buttons := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton("Назад")},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	msg.ReplyMarkup = keyboard

	_, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return 0, nil

	//TODO: тут невалидная кнопка Назад из-за подтвердить
	//TODO: добавить логику принятия сообщения и отправки ее куда-то. Нужна кнопка "Подтвердить" и "назад"
}

func (c *FeedbackCommand) SaveOrSendFeedback(api *tgbotapi.BotAPI, message *tgbotapi.Message) (bool, error) {
	responseMsg := "Вы уверены, что хотите отправить обратную связь?"
	if feedbackMsg, ok := c.userIDToFeedback.Load(message.From.ID); !ok {
		msg := tgbotapi.NewMessage(message.Chat.ID, responseMsg)
		c.userIDToFeedback.Store(message.From.ID, message.Text)
		buttons := [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Отправить обратную связь"), tgbotapi.NewKeyboardButton("Назад")},
		}

		keyboard := tgbotapi.NewReplyKeyboard(buttons...)
		msg.ReplyMarkup = keyboard

		_, err := api.Send(msg)
		if err != nil {
			return false, log.Wrap(err)
		}
	} else {
		c.userIDToFeedback.Delete(message.From.ID)
		ctx := context.Background()
		err := c.userSrv.SaveFeedback(ctx, message.From.ID, feedbackMsg.(string))
		if err != nil {
			return false, log.Wrap(err)
		}
		responseMsg = "Обратная связь успешно отправлена!"
		msg := tgbotapi.NewMessage(message.Chat.ID, responseMsg)
		_, err = api.Send(msg)
		if err != nil {
			return true, log.Wrap(err)
		}
		return true, nil
	}

	return false, nil
}

func (c *FeedbackCommand) IsTransfer() bool {
	return true
}

package about

import (
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FeedbackCommand struct{}

func (c *FeedbackCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	feedbackMessages := []string{
		"Тебе есть, что мне сказать? 💬 Я всегда готов выслушать твои отзывы и предложения.",
		"У тебя есть вопросы или предложения? 💬 Я всегда готов выслушать твоё мнение!",
		"Нужна помощь или хочешь оставить отзыв? 🌟 Напиши мне, и я постараюсь улучшить наш диалог!",
		"Свяжись со мной, если у тебя есть вопросы или предложения. 📬 Твой отзыв важен для меня!",
	}

	rand.Seed(time.Now().UnixNano())
	feedbackMessage := feedbackMessages[rand.Intn(len(feedbackMessages))]

	msg := tgbotapi.NewMessage(message.Chat.ID, feedbackMessage)

	//buttons := [][]tgbotapi.KeyboardButton{
	//	{tgbotapi.NewKeyboardButton("Подтвердить"), tgbotapi.NewKeyboardButton("Назад")},
	//}

	//keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	//msg.ReplyMarkup = keyboard

	api.Send(msg)

	//TODO: тут невалидная кнопка Назад из-за подтвердить
	//TODO: добавить логику принятия сообщения и отправки ее куда-то. Нужна кнопка "Подтвердить" и "назад"
}

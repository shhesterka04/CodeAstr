package reflection

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/llm/handlers"
	yadto "glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"
)

var (
	ErrMarkInvalid    = errors.New("please enter a number between 0 and 10")
	ErrRecordNotFound = errors.New("record not found")
)

var reflectionPrompt = "Привет бот! Мне нужно твоё мнение о твоём дне. " +
	"День от 0 до 10, где: 0 - очень плохо, а 10 - очень хорошо. У меня %v" +
	"Я чувствовал следующие эмоции: %s" +
	"Больше всего меня повлияло: %s" +
	"Требования: " +
	"1. Отчет должен быть детализированным и точным." +
	"2. Учти особенности целевой аудитории: студенты, офисные сотрудники, домохозяйки." +
	"3. Отчет должны быть позитивным и мотивирующим, но также реалистичным." +
	"4. Используй стиль общения: %s." +
	"5. Не надо здороваться и прощаться. " +
	"6. Примерно 50 слов " +
	"7. Пиши не по пунктам, а сплошным текстом " +
	"8. Используй эмодзи и смайлики для улучшения текста. " +
	"9. Не надо менять оценку дня мою. Если она низкая, поддержи, если высокая - порадуйся" +
	"Благодарю за помощь! "

const (
	markPrompt     = "Как бы ты оценил свой день от 1 до 10, где:\n1 — Очень плохой\n5 — Нейтральный\n10 — Идеальный"
	emotionsPrompt = "Что ты чувствуешь? Напиши свои эмоции за день \n(Например: Радость, Грусть, Страх, Гнев, Удивление, Отвращение, Доверие, Ожидание, Любовь, Ненависть, Стыд, Вина, Гордость, Зависть, Ревность, Облегчение, Удовлетворение, Разочарование, Тревога, Волнение)"
	activityPrompt = "Что больше всего повлияло на твоё настроение?\n(Например: Работа, Семья, Друзья, Спорт, Хобби, Путешествия, Учёба, Сон, Питание, Здоровье, Погода, Новости)"
)

type ReflectionCommand struct {
	userSrv   userservice.UserService
	gptApi    handlers.Handler
	usersData sync.Map
}

func NewReflectionCommand(userSrv userservice.UserService, gptApi handlers.Handler) *ReflectionCommand {
	return &ReflectionCommand{
		userSrv:   userSrv,
		gptApi:    gptApi,
		usersData: sync.Map{},
	}
}

func (c *ReflectionCommand) GetUserData(userID int64) *dto.ReflectionRecord {
	data, ok := c.usersData.Load(userID)
	if !ok {
		return nil
	}
	refdata, ok := data.(*dto.ReflectionRecord)
	if !ok {
		return nil
	}

	return refdata
}
func (c *ReflectionCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	usr, err := c.userSrv.FindByID(context.Background(), message.From.ID)
	if err != nil {
		return 0, log.Wrap(err)
	}
	if usr.Tokens <= 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "У вас закончились токены.. Приходите завтра")
		if _, err := api.Send(msg); err != nil {
			return 0, log.Wrap(err)
		}

		return 0, nil
	}

	if err = c.userSrv.Update(context.Background(), &dto.UpdateUserRequest{
		TgID:   message.From.ID,
		Tokens: usr.Tokens - 5,
	}); err != nil {
		return 0, log.Wrap(err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, markPrompt)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(predictions.BackCmd),
		),
	)

	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *ReflectionCommand) GetMarkSendEmotions(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	mark, err := strconv.Atoi(message.Text)
	if err != nil || mark < 1 || mark > 10 {
		return 0, ErrMarkInvalid
	}

	c.usersData.Store(message.From.ID, &dto.ReflectionRecord{UserID: message.From.ID, Mark: mark})

	msg := tgbotapi.NewMessage(message.Chat.ID, emotionsPrompt)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(predictions.BackCmd),
		),
	)

	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *ReflectionCommand) GetEmotionsSendActivity(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	record, ok := c.usersData.Load(message.From.ID)
	if !ok {
		return 0, log.Wrap(ErrRecordNotFound)
	}

	reflectionRecord := record.(*dto.ReflectionRecord)
	reflectionRecord.Emotions = message.Text
	msg := tgbotapi.NewMessage(message.Chat.ID, activityPrompt)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(predictions.BackCmd),
		),
	)

	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *ReflectionCommand) GetActivitySendResult(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	record, ok := c.usersData.Load(message.From.ID)
	if !ok {
		return 0, log.Wrap(ErrRecordNotFound)
	}

	reflectionRecord := record.(*dto.ReflectionRecord)
	reflectionRecord.Activity = message.Text

	if err := c.userSrv.SaveReflection(context.Background(), *reflectionRecord); err != nil {
		return 0, log.Wrap(err)
	}

	c.usersData.Delete(message.From.ID)

	ans, err := c.genResult(api, message, reflectionRecord)
	if err != nil {
		return 0, log.Wrap(err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, ans)
	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *ReflectionCommand) genResult(api *tgbotapi.BotAPI, message *tgbotapi.Message, reflection *dto.ReflectionRecord) (string, error) {
	wg := sync.WaitGroup{}
	usr, err := c.userSrv.FindByID(context.Background(), message.From.ID)
	if err != nil {
		return "", log.Wrap(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		typingMsg := tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)
		api.Send(typingMsg)
	}()

	yaResponse, err := c.gptApi.CallAPI(yadto.RequestDTO{
		UserMessage: fmt.Sprintf(
			reflectionPrompt,
			reflection.Mark,
			reflection.Emotions,
			reflection.Activity,
			usr.Style),
	})

	wg.Wait()

	if err != nil {
		return "", log.Wrap(err)
	}
	response, ok := yaResponse.(yadto.ResponseDTO)
	if !ok {
		return "", log.Wrap(err)
	}

	return response.Result, nil
}

func (c *ReflectionCommand) IsTransfer() bool {
	return true
}

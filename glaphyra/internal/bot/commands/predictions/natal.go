package predictions

import (
	"context"
	"fmt"
	"sync"

	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/llm/handlers"
	yadto "glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const natalPrompt = "Привет, бот! Создай мне натальную карту для пользователя, вот данные: " +
	"Вид деятельности: %s; Семейное положение: %s; Дата рождения: %s; Время рождения: %s; Место рождения: %s" +
	"Требования: " +
	"1. Натальная карта должена быть детализированной и точной." +
	"2. Учти особенности целевой аудитории: студенты, офисные сотрудники, домохозяйки." +
	"3. Расклад должен быть описывать мои текущие сильные и слабые стороны и к чему это может привести в будущем." +
	"4. Используй стиль общения: %s." +
	"5. Не надо здороваться и прощаться. " +
	"6. Примерно 200 слов " +
	"7. Оформи как профессиональную натальную карту " +
	"8. Используй много эмодзи и смайлики для улучшения текста. " +
	" Включить следующие разделы анализа: " +
	"- Любовь и отношения: " +
	"- Карьера и финансы: " +
	"- Здоровье: " +
	"- Семья и дом: " +
	"- Друзья: " +
	"- Путешествия: " +
	"- Творчество: " +
	"- Образование: " +
	"- Юмор: " +
	"Благодарю за помощь! "

type NatalCommand struct {
	userSrv userservice.UserService
	gptApi  handlers.Handler
}

func NewNatalCommand(userSrv userservice.UserService, gptApi handlers.Handler) *NatalCommand {
	return &NatalCommand{
		userSrv: userSrv,
		gptApi:  gptApi,
	}
}

func (c *NatalCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	wg := sync.WaitGroup{}
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

	if usr.FamilyStatus == "" || usr.BirthTime == "" || usr.BirthPlace == "" || usr.TypeOfActivity == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Похоже, что у вас не заполнены все данные для составления натальной карты. Пожалуйста, перейдите в `Регистрация и настройки` > `Заполнить 'О себе'` и заполните поля.")
		if _, err = api.Send(msg); err != nil {
			return 0, log.Wrap(err)
		}

		return 0, nil
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		typingMsg := tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)
		api.Send(typingMsg)
	}()

	yaResponse, err := c.gptApi.CallAPI(yadto.RequestDTO{
		UserMessage: fmt.Sprintf(
			natalPrompt,
			usr.TypeOfActivity,
			usr.FamilyStatus,
			usr.BirthDate,
			usr.BirthTime,
			usr.BirthPlace,
			usr.Style),
	})

	wg.Wait()

	if err != nil {
		return 0, log.Wrap(err)
	}
	response, ok := yaResponse.(yadto.ResponseDTO)
	if !ok {
		return 0, log.Wrap(err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, response.Result)
	msg.ParseMode = "Markdown"
	sentMsg, err := api.Send(msg)
	if err != nil {
		return 0, log.Wrap(err)
	}

	return int64(sentMsg.MessageID), nil
}

func (c *NatalCommand) IsTransfer() bool {
	return false
}

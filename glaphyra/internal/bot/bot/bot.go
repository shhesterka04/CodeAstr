package bot

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"glaphyra/internal/app/users/dto"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot"
	"glaphyra/internal/bot/update_handler"
	"glaphyra/internal/bot/user_cmd_handler"
	"glaphyra/internal/llm/handlers"
	yadto "glaphyra/internal/llm/yagpt/dto"
	"glaphyra/internal/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const summaryPrompt = "Привет, бот! Мне нужно подвести итоги недели для пользователя:- " +
	"Вот данные за неделю %s" +
	"Требования: " +
	"1. Надо вывести все цифры, которые есть в данных, и посчитать среднее значение. " +
	"2. Надо вывести эмоции, которые пользователь ощущал чаще всего. " +
	"3. Надо вывести, какие действия пользователь совершал чаще всего. " +
	"4. Надо на основе этих данных сделать выводы, как будто об этом говорят звезды" +
	"5. Используй стиль общения: %s." +
	"6. Не надо здороваться и прощаться. " +
	"7. Примерно 100 слов " +
	"8. Пиши не по пунктам, а сплошным текстом " +
	"9. Используй эмодзи и смайлики для улучшения текста. " +
	"Благодарю за помощь! "

var msgVars = []string{
	"Привет! 🌟 Как прошел твой день? Давай вместе подумаем, что принесло радость, а что — нет. Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Доброго вечера! ✨ Ты готов поделиться, как прошел твой день? Я здесь, чтобы помочь тебе понять себя лучше! Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Привет! 🌟 Как ты сегодня? Готов(а) немного поговорить о своём дне? Вместе разберём, что принесло радость, а что — нет! Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Доброго времени суток! ✨ Давай обсудим, как прошёл твой день. Иногда звёзды знают больше, чем кажется! Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Привет-привет! 🌙 Как настроение? Готов(а) поделиться тем, что оставило след в твоём дне? Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Приветствую! 😊 Как твои дела? Давай немного поразмышляем, что сегодня было хорошего или, может, сложного? Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Здравствуй! 🌞 Я здесь, чтобы выслушать, как прошёл твой день, и помочь тебе лучше понять своё настроение. Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Привет! ✨ Звёзды уже шепнули мне, что день был насыщенным. Поделишься своими мыслями и эмоциями? Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Добро пожаловать! 🌌 Что сегодня сделало твой день ярче или сложнее? Расскажи, я помогу разобраться. Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
	"Привет! 💫 Хотел(а) бы рассказать о своём дне? Иногда рефлексия помогает увидеть больше, чем мы ожидаем. Чтобы сделать запись, нажми на кнопку `Рефлексия` на стартовом экране.",
}

type Bot struct {
	api     *tgbotapi.BotAPI
	handler bot.UpdateHandler
	gptApi  handlers.Handler
	userSrv userservice.UserService
}

func NewBot(tgToken string, userSrv userservice.UserService, gptApi handlers.Handler) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}
	userCmdHandler := user_cmd_handler.NewUserCmdHandler(api, userSrv)
	updateHandler := update_handler.NewUpdateHandler(userCmdHandler, userSrv, gptApi)

	return &Bot{api: api, handler: updateHandler, gptApi: gptApi, userSrv: userSrv}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	go b.startTicker()

	for update := range updates {
		go func(update tgbotapi.Update) {
			b.handler.HandleUpdate(update)
		}(update)
	}
}

func (b *Bot) startTicker() {
	for {
		now := time.Now()
		next := now.Add(time.Hour).Truncate(time.Hour)

		ticker := time.NewTicker(next.Sub(now))
		<-ticker.C

		b.handleHourlyTask(next.Hour())

		if next.Weekday() == time.Sunday {
			b.handleWeeklyTask(next.Hour())
		}

		ticker.Stop()
	}
}

func (b *Bot) handleHourlyTask(hour int) {
	rand.Seed(time.Now().UnixNano())
	ids, _ := b.userSrv.GetUsersByNotificationTime(context.Background(), hour)
	for _, id := range ids {
		b.userSrv.Update(context.Background(), &dto.UpdateUserRequest{
			TgID:   id,
			Tokens: 100,
		})
		fmt.Println("User with id", id, "received 100 tokens, hour: ", hour)
		b.api.Send(tgbotapi.NewMessage(id, msgVars[rand.Intn(len(msgVars))]))
	}
}

func (b *Bot) handleWeeklyTask(hour int) {
	ids, _ := b.userSrv.GetUsersByNotificationTime(context.Background(), hour)
	for _, id := range ids {
		reflections, _ := b.userSrv.GetReflectionsByUserIDLastWeekFormat(context.Background(), id)
		usr, _ := b.userSrv.FindByID(context.Background(), id)
		yaResponse, err := b.gptApi.CallAPI(yadto.RequestDTO{
			UserMessage: fmt.Sprintf(
				summaryPrompt,
				reflections,
				usr.Style),
		})
		if err != nil {
			log.Wrap(err)
		}

		response, ok := yaResponse.(yadto.ResponseDTO)
		if !ok {
			log.Wrap(err)
		}

		msg := tgbotapi.NewMessage(id, response.Result)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
	}
}

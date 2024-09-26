package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/bot"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(tgToken string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}

	return &Bot{api: api}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		bot.HandleUpdate(b.api, update)
	}
}

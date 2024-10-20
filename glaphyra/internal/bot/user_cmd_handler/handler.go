package user_cmd_handler

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/bot"
	"glaphyra/internal/pkg/log"
)

type implUserCmdHandler struct {
	api                  *tgbotapi.BotAPI
	userCommandHistories sync.Map
}

type UserCommandHandler interface {
	HandleUserCommand(userID int64, cmd bot.Command, message *tgbotapi.Message) error
}

func NewUserCmdHandler(api *tgbotapi.BotAPI) UserCommandHandler {
	return &implUserCmdHandler{
		api: api,
	}
}

func (i *implUserCmdHandler) HandleUserCommand(userID int64, cmd bot.Command, message *tgbotapi.Message) error {
	if cmd == nil {
		i.handleUnknownCommand(i.api, message)
		return nil
	}

	history, _ := i.userCommandHistories.LoadOrStore(userID, &CommandHistory{})
	if cmd.IsTransfer() {
		history.(*CommandHistory).Add(cmd)
	}
	switch cmd.(type) {
	case *BackCommand:
		cmd = &BackCommand{commandHistory: history.(*CommandHistory)}
	}

	err := cmd.Execute(i.api, message)
	if err != nil {
		return log.WrapErr(err)
	}

	return nil
}

func (i *implUserCmdHandler) handleUnknownCommand(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не понимаю (кто понял тот понял)")
	api.Send(msg)
}

type BackCommand struct {
	commandHistory *CommandHistory
}

func (c *BackCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	currentPageCmd := c.commandHistory.Pop()
	if currentPageCmd == nil {
		return nil
	}
	previousPageCmd := c.commandHistory.GetPrevious()
	if previousPageCmd == nil {
		return nil
	}

	return previousPageCmd.Execute(api, message)
}

func (c *BackCommand) IsTransfer() bool {
	return false
}

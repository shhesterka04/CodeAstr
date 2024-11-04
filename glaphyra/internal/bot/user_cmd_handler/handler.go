package user_cmd_handler

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/bot"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/pkg/log"
)

type implUserCmdHandler struct {
	api                  *tgbotapi.BotAPI
	userCommandHistories sync.Map
	messageIDToCommandID sync.Map
}

type UserCommandHandler interface {
	HandleUserCommand(userID int64, cmd bot.Command, message *tgbotapi.Message) error
	HandleUserCallback(msgID int64, message *tgbotapi.CallbackQuery) error
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

	msgID, err := cmd.Execute(i.api, message)
	if err != nil {
		return log.WrapErr(err)
	}

	if msgID != 0 {
		i.messageIDToCommandID.LoadOrStore(msgID, cmd)
	}

	return nil
}

func (i *implUserCmdHandler) HandleUserCallback(msgID int64, callback *tgbotapi.CallbackQuery) error {
	if msgID == 0 {
		return nil
	}

	cmd, ok := i.messageIDToCommandID.Load(msgID)
	if !ok {
		return nil
	}

	switch cmd.(type) {
	case *predictions.DayHoroscopeCommand:
		fmt.Println(fmt.Sprintf("it was DayHoroscopeCommand from %v", callback.From.ID))
	default:
		fmt.Println("it was unknown command")
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

func (c *BackCommand) Execute(api *tgbotapi.BotAPI, message *tgbotapi.Message) (int64, error) {
	currentPageCmd := c.commandHistory.Pop()
	if currentPageCmd == nil {
		return 0, nil
	}
	previousPageCmd := c.commandHistory.GetPrevious()
	if previousPageCmd == nil {
		return 0, nil
	}

	return previousPageCmd.Execute(api, message)
}

func (c *BackCommand) IsTransfer() bool {
	return false
}

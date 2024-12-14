package user_cmd_handler

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"glaphyra/internal/bot"
	"glaphyra/internal/bot/commands/about"
	"glaphyra/internal/bot/commands/compatibility"
	"glaphyra/internal/bot/commands/dreambook"
	"glaphyra/internal/bot/commands/predictions"
	"glaphyra/internal/bot/commands/settings"
	"glaphyra/internal/pkg/log"
)

type implUserCmdHandler struct {
	api                  *tgbotapi.BotAPI
	userCommandHistories sync.Map
	messageIDToCommandID sync.Map
}

type UserCommandHandler interface {
	HandleUserCommand(userID int64, cmd bot.Command, message *tgbotapi.Message) error
	HandleUserCallback(msgID int64, userID int64, message *tgbotapi.CallbackQuery) error
}

func NewUserCmdHandler(api *tgbotapi.BotAPI) UserCommandHandler {
	return &implUserCmdHandler{
		api: api,
	}
}

func (i *implUserCmdHandler) HandleUserCommand(userID int64, cmd bot.Command, message *tgbotapi.Message) error {
	if cmd == nil {
		i.handleUnknownCommand(userID, i.api, message)
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

func (i *implUserCmdHandler) HandleUserCallback(msgID int64, userID int64, callback *tgbotapi.CallbackQuery) error {
	if msgID == 0 {
		return nil
	}

	cmd, ok := i.messageIDToCommandID.Load(msgID)
	if !ok {
		return nil
	}

	history, _ := i.userCommandHistories.LoadOrStore(userID, &CommandHistory{})
	switch cmd.(type) {
	case *compatibility.NatalCompCommand:
	case *predictions.HoroscopeCommand:
		cmd.(*predictions.HoroscopeCommand).SendPrompt(i.api, callback)
		back := &BackCommand{commandHistory: history.(*CommandHistory)}
		back.Execute(i.api, callback.Message)
	case *compatibility.ZodiakCompCommand:
		if !cmd.(*compatibility.ZodiakCompCommand).IsFirstFull(callback.From.ID) {
			msgSecondID, _ := cmd.(*compatibility.ZodiakCompCommand).GetFirstSignSendSecond(i.api, callback)
			i.messageIDToCommandID.LoadOrStore(msgSecondID, cmd)
		} else {
			cmd.(*compatibility.ZodiakCompCommand).GetSecondSignSendResult(i.api, callback)
			back := &BackCommand{commandHistory: history.(*CommandHistory)}
			back.Execute(i.api, callback.Message)
		}

	default:
		fmt.Println("it was unknown command")
	}

	return nil
}

func (i *implUserCmdHandler) handleUnknownCommand(userID int64, api *tgbotapi.BotAPI, message *tgbotapi.Message) {
	responseMsg := "Я не понимаю (кто понял тот понял)"
	history, _ := i.userCommandHistories.LoadOrStore(userID, &CommandHistory{})
	cmd := history.(*CommandHistory).GetPrevious()
	switch cmd.(type) {
	case *compatibility.NatalCompCommand:
		_, err := cmd.(*compatibility.NatalCompCommand).SendResult(api, message)
		if err != nil {
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		back := &BackCommand{commandHistory: history.(*CommandHistory)}
		_, err = back.Execute(api, message)
		return
	case *dreambook.DreambookCommand:
		_, err := cmd.(*dreambook.DreambookCommand).SendResult(api, message)
		if err != nil {
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		back := &BackCommand{commandHistory: history.(*CommandHistory)}
		_, err = back.Execute(api, message)
		return
	case *settings.Birth:
		_, err := cmd.(*settings.Birth).SetBirth(api, message)
		if err != nil {
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		back := &BackCommand{commandHistory: history.(*CommandHistory)}
		_, err = back.Execute(api, message)
		return
	case *about.FeedbackCommand:
		sent, err := cmd.(*about.FeedbackCommand).SaveOrSendFeedback(api, message)
		if err != nil {
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		if sent {
			back := &BackCommand{commandHistory: history.(*CommandHistory)}
			_, err = back.Execute(api, message)
		}
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, responseMsg)
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

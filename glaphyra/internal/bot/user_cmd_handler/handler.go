package user_cmd_handler

import (
	"context"
	"fmt"
	"glaphyra/internal/app/users/dto"
	errs "glaphyra/internal/pkg/errors"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	userservice "glaphyra/internal/app/users/service"
	"glaphyra/internal/bot"
	maincmd "glaphyra/internal/bot/commands"
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
	userService          userservice.UserService
}

type UserCommandHandler interface {
	HandleUserCommand(userID int64, cmd bot.Command, message *tgbotapi.Message) error
	HandleUserCallback(msgID int64, userID int64, message *tgbotapi.CallbackQuery) error
}

func NewUserCmdHandler(api *tgbotapi.BotAPI, userService userservice.UserService) UserCommandHandler {
	return &implUserCmdHandler{
		api:         api,
		userService: userService,
	}
}

func (i *implUserCmdHandler) HandleUserCommand(userID int64, cmd bot.Command, message *tgbotapi.Message) error {
	err := i.userService.Update(context.Background(), &dto.UpdateUserRequest{
		TgID:           userID,
		LastActionTime: time.Now().UTC(),
	})
	if err != nil {
		log.Error(err)
	}
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
		_, err := back.Execute(i.api, callback.Message)
		if err != nil {
			log.Error(err)
		}
	case *compatibility.ZodiakCompCommand:
		if !cmd.(*compatibility.ZodiakCompCommand).IsFirstFull(callback.From.ID) {
			msgSecondID, _ := cmd.(*compatibility.ZodiakCompCommand).GetFirstSignSendSecond(i.api, callback)
			i.messageIDToCommandID.LoadOrStore(msgSecondID, cmd)
		} else {
			cmd.(*compatibility.ZodiakCompCommand).GetSecondSignSendResult(i.api, callback)
			back := &BackCommand{commandHistory: history.(*CommandHistory)}
			_, err := back.Execute(i.api, callback.Message)
			if err != nil {
				log.Error(err)
			}
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
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		back := &BackCommand{commandHistory: history.(*CommandHistory)}
		_, err = back.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		return
	case *dreambook.DreambookCommand:
		_, err := cmd.(*dreambook.DreambookCommand).SendResult(api, message)
		if err != nil {
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		back := &BackCommand{commandHistory: history.(*CommandHistory)}
		_, err = back.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		return
	case *settings.Birth:
		_, err := cmd.(*settings.Birth).SetBirth(api, message)
		if errs.As[settings.ValidationError](err) {
			return
		}
		if err != nil {
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		birthPlace := settings.NewBirthPlaceCommand(i.userService)
		_, err = birthPlace.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		if birthPlace.IsTransfer() {
			history.(*CommandHistory).Add(birthPlace)
		}
		return
	case *settings.BirthPlace:
		_, err := cmd.(*settings.BirthPlace).SetBirthPlace(api, message)
		if err != nil {
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		birthTime := settings.NewBirthTimeCommand(i.userService)
		_, err = birthTime.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		if birthTime.IsTransfer() {
			history.(*CommandHistory).Add(birthTime)
		}
		return
	case *settings.BirthTime:
		_, err := cmd.(*settings.BirthTime).SetBirthTime(api, message)
		if err != nil {
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		familyStatus := settings.NewFamilyStatusCommand(i.userService)
		_, err = familyStatus.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		if familyStatus.IsTransfer() {
			history.(*CommandHistory).Add(familyStatus)
		}
		return
	case *settings.FamilyStatus:
		_, err := cmd.(*settings.FamilyStatus).SetFamilyStatus(api, message)
		if err != nil {
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		typeOfActivity := settings.NewTypeOfActivityCommand(i.userService)
		_, err = typeOfActivity.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		if typeOfActivity.IsTransfer() {
			history.(*CommandHistory).Add(typeOfActivity)
		}
		return
	case *settings.TypeOfActivity:
		_, err := cmd.(*settings.TypeOfActivity).SetTypeOfActivity(api, message)
		if err != nil {
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		notificationTime := settings.NewNotificationTimeCommand(i.userService)
		_, err = notificationTime.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		if notificationTime.IsTransfer() {
			history.(*CommandHistory).Add(notificationTime)
		}
		return
	case *settings.NotificationTime:
		_, err := cmd.(*settings.NotificationTime).SetNotificationTime(api, message)
		if errs.As[settings.ValidationError](err) {
			return
		}
		if err != nil {
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		startCmd := maincmd.NewStartCommand(i.userService)
		_, err = startCmd.Execute(api, message)
		if err != nil {
			log.Error(err)
		}
		if startCmd.IsTransfer() {
			history.(*CommandHistory).Add(startCmd)
		}
		return
	case *about.FeedbackCommand:
		sent, err := cmd.(*about.FeedbackCommand).SaveOrSendFeedback(api, message)
		if err != nil {
			log.Error(err)
			responseMsg = "Что-то пошло не так, попробуйте снова"
			break
		}
		if sent {
			back := &BackCommand{commandHistory: history.(*CommandHistory)}
			_, err = back.Execute(api, message)
			if err != nil {
				log.Error(err)
			}
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

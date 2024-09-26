package predictions

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	inlineButtons = []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("♈ Овен", "Aries"),
		tgbotapi.NewInlineKeyboardButtonData("♉ Телец", "Taurus"),
		tgbotapi.NewInlineKeyboardButtonData("♊ Близнецы", "Gemini"),
		tgbotapi.NewInlineKeyboardButtonData("♋ Рак", "Cancer"),
		tgbotapi.NewInlineKeyboardButtonData("♌ Лев", "Leo"),
		tgbotapi.NewInlineKeyboardButtonData("♍ Дева", "Virgo"),
		tgbotapi.NewInlineKeyboardButtonData("♎ Весы", "Libra"),
		tgbotapi.NewInlineKeyboardButtonData("♏ Скорпион", "Scorpio"),
		tgbotapi.NewInlineKeyboardButtonData("♐ Стрелец", "Sagittarius"),
		tgbotapi.NewInlineKeyboardButtonData("♑ Козерог", "Capricorn"),
		tgbotapi.NewInlineKeyboardButtonData("♒ Водолей", "Aquarius"),
		tgbotapi.NewInlineKeyboardButtonData("♓ Рыбы", "Pisces"),
	}

	inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(inlineButtons[0], inlineButtons[1], inlineButtons[2], inlineButtons[3]),
		tgbotapi.NewInlineKeyboardRow(inlineButtons[4], inlineButtons[5], inlineButtons[6], inlineButtons[7]),
		tgbotapi.NewInlineKeyboardRow(inlineButtons[8], inlineButtons[9], inlineButtons[10], inlineButtons[11]),
	)
)

func HandleCallbackQuery(api *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	responseText := "Вы выбрали: " + callbackQuery.Data
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, responseText)
	api.Send(msg)
}

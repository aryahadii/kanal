package keyboard

import (
	"gitlab.com/arha/kanal/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func NewMainKeyboard() botAPI.ReplyKeyboardMarkup {
	newMessageKey := botAPI.NewKeyboardButton(model.NewMessageCommand)
	row1 := botAPI.NewKeyboardButtonRow(newMessageKey)

	helpKey := botAPI.NewKeyboardButton(model.HelpCommand)
	kanalLinkKey := botAPI.NewKeyboardButton(model.KanalLinkCommand)
	row2 := botAPI.NewKeyboardButtonRow(helpKey, kanalLinkKey)

	return botAPI.NewReplyKeyboard(row1, row2)
}

func NewMessageCancelKeyboard() botAPI.ReplyKeyboardMarkup {
	cancelKey := botAPI.NewKeyboardButton(model.NewMessageCancelCommand)
	row := botAPI.NewKeyboardButtonRow(cancelKey)
	return botAPI.NewReplyKeyboard(row)
}

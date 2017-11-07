package keyboard

import (
	"fmt"

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

func NewAdminInlineKeyboard(messageID string) botAPI.InlineKeyboardMarkup {
	var row []botAPI.InlineKeyboardButton
	row = append(row, botAPI.NewInlineKeyboardButtonData(model.AdminKeyboardAccept, model.RadifeButton))
	return botAPI.NewInlineKeyboardMarkup(row)
}

func NewEmojiInlineKeyboard(type1, type2, type3 string) botAPI.InlineKeyboardMarkup {
	var row []botAPI.InlineKeyboardButton
	type1Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type1Emoji, type1),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "1"))
	type2Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type2Emoji, type2),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "2"))
	type3Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type3Emoji, type3),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "3"))
	row = append(row, type1Key, type2Key, type3Key)
	return botAPI.NewInlineKeyboardMarkup(row)
}

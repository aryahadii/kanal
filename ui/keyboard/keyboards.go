package keyboard

import (
	"fmt"
	"strconv"

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
	radifeButton := botAPI.NewInlineKeyboardButtonData(model.AdminKeyboardAccept, model.RadifeButton)
	nazarButton := botAPI.NewInlineKeyboardButtonData(model.AdminKeyboardSurvey, model.NazarButton)
	naHajiButton := botAPI.NewInlineKeyboardButtonData(model.AdminKeyboardReject, model.NaHajiButton)
	row = append(row, radifeButton, nazarButton, naHajiButton)
	return botAPI.NewInlineKeyboardMarkup(row)
}

func NewEmojiInlineKeyboard(type1, type2, type3, type4 int) botAPI.InlineKeyboardMarkup {
	var row1 []botAPI.InlineKeyboardButton
	type1Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type1Emoji, strconv.Itoa(type1)),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "1"))
	type2Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type2Emoji, strconv.Itoa(type2)),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "2"))
	type3Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type3Emoji, strconv.Itoa(type3)),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "3"))
	type4Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type4Emoji, strconv.Itoa(type4)),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "4"))
	row1 = append(row1, type1Key, type2Key, type3Key, type4Key)
	return botAPI.NewInlineKeyboardMarkup(row1)
}

func NewSurveyInlineKeyboard(type1, type2 int) botAPI.InlineKeyboardMarkup {
	var row1 []botAPI.InlineKeyboardButton
	type1Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type5Emoji, strconv.Itoa(type1)),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "5"))
	type2Key := botAPI.NewInlineKeyboardButtonData(fmt.Sprintf(model.Type6Emoji, strconv.Itoa(type2)),
		fmt.Sprint(model.EmojiButton, model.CallbackSeparator, "6"))
	row1 = append(row1, type1Key, type2Key)
	return botAPI.NewInlineKeyboardMarkup(row1)
}

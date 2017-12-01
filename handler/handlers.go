package handler

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gitlab.com/arha/kanal/configuration"
	"gitlab.com/arha/kanal/db"
	"gitlab.com/arha/kanal/ui/keyboard"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/kanal/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

var (
	userState = cache.New(30*time.Minute, 10*time.Minute)
)

func HandleCallbacks(callbackQuery *botAPI.CallbackQuery) []botAPI.Chattable {
	splittedCallbackData := strings.Split(callbackQuery.Data, model.CallbackSeparator)
	if splittedCallbackData[0] == model.RadifeButton {
		var responseChattables []botAPI.Chattable
		if callbackQuery.Message.Photo != nil {
			photo := (*callbackQuery.Message.Photo)[len(*callbackQuery.Message.Photo)-1]
			photoMessage := botAPI.PhotoConfig{
				BaseFile: botAPI.BaseFile{
					BaseChat:    botAPI.BaseChat{ChannelUsername: configuration.KanalConfig.GetString("kanal-username")},
					FileID:      photo.FileID,
					UseExisting: true,
				},
				Caption: callbackQuery.Message.Caption,
			}
			photoMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			responseChattables = append(responseChattables, photoMessage)
		} else if callbackQuery.Message.Video != nil {
			videoMessage := botAPI.VideoConfig{
				BaseFile: botAPI.BaseFile{
					BaseChat:    botAPI.BaseChat{ChannelUsername: configuration.KanalConfig.GetString("kanal-username")},
					FileID:      callbackQuery.Message.Video.FileID,
					UseExisting: true,
				},
				Caption: callbackQuery.Message.Caption,
			}
			videoMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			responseChattables = append(responseChattables, videoMessage)
		} else {
			kanalMessage := botAPI.NewMessageToChannel(configuration.KanalConfig.GetString("kanal-username"), callbackQuery.Message.Text)
			kanalMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			responseChattables = append(responseChattables, kanalMessage)
		}
		deleteMessageConfig := botAPI.DeleteMessageConfig{
			ChatID:    callbackQuery.Message.Chat.ID,
			MessageID: callbackQuery.Message.MessageID,
		}
		responseChattables = append(responseChattables, deleteMessageConfig)

		return responseChattables
	} else if splittedCallbackData[0] == model.NaHajiButton {
		deleteMessageConfig := botAPI.DeleteMessageConfig{
			ChatID:    callbackQuery.Message.Chat.ID,
			MessageID: callbackQuery.Message.MessageID,
		}
		return []botAPI.Chattable{
			deleteMessageConfig,
		}
	} else if splittedCallbackData[0] == model.EmojiButton {
		var messageData model.Message
		err := db.MessagesCollection.Find(bson.M{
			"messageid": callbackQuery.Message.MessageID,
		}).One(&messageData)
		if err != nil {
			messageData = model.Message{
				MessageID: callbackQuery.Message.MessageID,
				Reactions: make([][]string, 4),
			}
			go db.MessagesCollection.Insert(messageData)
		}

		userID := strconv.Itoa(callbackQuery.From.ID)
		removeUserReaction(userID, messageData.Reactions)
		tappedEmojiIndex, _ := strconv.Atoi(splittedCallbackData[1])
		messageData.Reactions[tappedEmojiIndex-1] = append(messageData.Reactions[tappedEmojiIndex-1], userID)
		go db.MessagesCollection.Update(bson.M{
			"messageid": messageData.MessageID,
		}, messageData)

		editedMessage := botAPI.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID,
			callbackQuery.Message.MessageID, keyboard.NewEmojiInlineKeyboard(
				len(messageData.Reactions[0]),
				len(messageData.Reactions[1]),
				len(messageData.Reactions[2]),
				len(messageData.Reactions[3])))
		return []botAPI.Chattable{
			editedMessage,
		}
	}

	return []botAPI.Chattable{}
}

func HandleMessage(message *botAPI.Message) []botAPI.Chattable {
	if answers, handled := handleCommand(message); handled {
		return answers
	}
	if _, found := userState.Get(strconv.Itoa(message.From.ID)); found {
		return handleNewMessage(message)
	}

	errorMessage := botAPI.NewMessage(message.Chat.ID, model.ErrorMessage)
	errorMessage.ReplyMarkup = keyboard.NewMainKeyboard()
	return []botAPI.Chattable{
		errorMessage,
	}
}

func handleCommand(message *botAPI.Message) ([]botAPI.Chattable, bool) {
	var answerMessages []botAPI.Chattable

	if message.IsCommand() {
		if message.Text == "/start" {
			welcomeMessage := botAPI.NewMessage(message.Chat.ID, model.WelcomeMessage)
			welcomeMessage.ReplyMarkup = keyboard.NewMainKeyboard()
			answerMessages = append(answerMessages, welcomeMessage)
			return answerMessages, true
		}
	} else if message.Text == model.NewMessageCommand {
		state := &model.UserState{
			CommandState: model.UserCommandStateNewMessage,
		}
		err := userState.Add(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)
		if err != nil {
			return answerMessages, false
		}

		enterMessage := botAPI.NewMessage(message.Chat.ID, model.NewMessageCommandMessage)
		enterMessage.ReplyMarkup = keyboard.NewMessageCancelKeyboard()
		answerMessages = append(answerMessages, enterMessage)
		return answerMessages, true
	} else if message.Text == model.NewMessageCancelCommand {
		userState.Delete(strconv.Itoa(message.From.ID))
		returnedMessage := botAPI.NewMessage(message.Chat.ID, model.ReturnedMessage)
		returnedMessage.ReplyMarkup = keyboard.NewMainKeyboard()
		answerMessages = append(answerMessages, returnedMessage)
		return answerMessages, true
	} else if message.Text == model.HelpCommand {
		helpMessage := botAPI.NewMessage(message.Chat.ID, model.HelpCommandMessage)
		helpMessage.ReplyMarkup = keyboard.NewMainKeyboard()
		answerMessages = append(answerMessages, helpMessage)
		return answerMessages, true
	} else if message.Text == model.KanalLinkCommand {
		linkMessage := botAPI.NewMessage(message.Chat.ID, configuration.KanalConfig.GetString("kanal-username"))
		linkMessage.ReplyMarkup = keyboard.NewMainKeyboard()
		answerMessages = append(answerMessages, linkMessage)
		return answerMessages, true
	}
	return answerMessages, false
}

func handleNewMessage(message *botAPI.Message) []botAPI.Chattable {
	var answerMessages []botAPI.Chattable

	var state *model.UserState
	if cachedObj, found := userState.Get(strconv.Itoa(message.From.ID)); found {
		state = cachedObj.(*model.UserState)
	} else {
		errorMessage := botAPI.NewMessage(message.Chat.ID, model.ErrorMessage)
		errorMessage.ReplyMarkup = keyboard.NewMainKeyboard()
		answerMessages = append(answerMessages, errorMessage)
		return answerMessages
	}
	if state.CommandState != model.UserCommandStateNewMessage {
		errorMessage := botAPI.NewMessage(message.Chat.ID, model.ErrorMessage)
		errorMessage.ReplyMarkup = keyboard.NewMainKeyboard()
		answerMessages = append(answerMessages, errorMessage)
		return answerMessages
	}

	// Post to Kanal Admins
	kanalArchiveMessage := botAPI.NewForward(configuration.KanalConfig.GetInt64("kanal-archive-chatid"), message.Chat.ID, message.MessageID)
	answerMessages = append(answerMessages, kanalArchiveMessage)
	if message.Photo != nil {
		photo := (*message.Photo)[len(*message.Photo)-1]
		kanalPhotoMessage := botAPI.NewPhotoShare(configuration.KanalConfig.GetInt64("kanal-admins-chatid"), photo.FileID)
		kanalPhotoMessage.Caption = message.Caption
		kanalPhotoMessage.ReplyMarkup = keyboard.NewAdminInlineKeyboard(strconv.Itoa(message.MessageID))
		answerMessages = append(answerMessages, kanalPhotoMessage)
	} else if message.Video != nil {
		kanalVideoMessage := botAPI.NewVideoShare(configuration.KanalConfig.GetInt64("kanal-admins-chatid"), message.Video.FileID)
		kanalVideoMessage.Caption = message.Caption
		kanalVideoMessage.ReplyMarkup = keyboard.NewAdminInlineKeyboard(strconv.Itoa(message.MessageID))
		answerMessages = append(answerMessages, kanalVideoMessage)
	} else {
		kanalMessage := botAPI.NewMessage(configuration.KanalConfig.GetInt64("kanal-admins-chatid"), message.Text)
		kanalMessage.ReplyMarkup = keyboard.NewAdminInlineKeyboard(strconv.Itoa(message.MessageID))
		answerMessages = append(answerMessages, kanalMessage)
	}

	// Successful message
	successfulSentMessage := botAPI.NewMessage(message.Chat.ID, model.NewMessageSentMessage)
	successfulSentMessage.ReplyMarkup = keyboard.NewMainKeyboard()
	answerMessages = append(answerMessages, successfulSentMessage)

	userState.Delete(strconv.Itoa(message.From.ID))

	return answerMessages
}

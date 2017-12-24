package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

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

		var text string
		replyMessageID := -1
		if callbackQuery.Message.Text != "" {
			text, replyMessageID = findMessageIDInText(callbackQuery.Message.Text)
		} else if callbackQuery.Message.Caption != "" {
			text, replyMessageID = findMessageIDInText(callbackQuery.Message.Caption)
		}

		if callbackQuery.Message.Document != nil {
			gifMessage := botAPI.DocumentConfig{
				BaseFile: botAPI.BaseFile{
					BaseChat:    botAPI.BaseChat{ChannelUsername: configuration.KanalConfig.GetString("kanal-username")},
					FileID:      callbackQuery.Message.Document.FileID,
					UseExisting: true,
				},
			}
			gifMessage.Caption = text
			gifMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			if replyMessageID > -1 {
				gifMessage.ReplyToMessageID = replyMessageID
			}
			responseChattables = append(responseChattables, gifMessage)
		} else if callbackQuery.Message.Photo != nil {
			photo := (*callbackQuery.Message.Photo)[len(*callbackQuery.Message.Photo)-1]
			photoMessage := botAPI.PhotoConfig{
				BaseFile: botAPI.BaseFile{
					BaseChat:    botAPI.BaseChat{ChannelUsername: configuration.KanalConfig.GetString("kanal-username")},
					FileID:      photo.FileID,
					UseExisting: true,
				},
				Caption: text,
			}
			photoMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			if replyMessageID > -1 {
				photoMessage.ReplyToMessageID = replyMessageID
			}
			responseChattables = append(responseChattables, photoMessage)
		} else if callbackQuery.Message.Voice != nil {
			voiceMessage := botAPI.VoiceConfig{
				BaseFile: botAPI.BaseFile{
					BaseChat:    botAPI.BaseChat{ChannelUsername: configuration.KanalConfig.GetString("kanal-username")},
					FileID:      callbackQuery.Message.Voice.FileID,
					UseExisting: true,
				},
				Caption: text,
			}
			voiceMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			if replyMessageID > -1 {
				voiceMessage.ReplyToMessageID = replyMessageID
			}
			responseChattables = append(responseChattables, voiceMessage)
		} else if callbackQuery.Message.Audio != nil {
			audioMessage := botAPI.AudioConfig{
				BaseFile: botAPI.BaseFile{
					BaseChat:    botAPI.BaseChat{ChannelUsername: configuration.KanalConfig.GetString("kanal-username")},
					FileID:      callbackQuery.Message.Audio.FileID,
					UseExisting: true,
				},
				Caption: text,
			}
			audioMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			if replyMessageID > -1 {
				audioMessage.ReplyToMessageID = replyMessageID
			}
			responseChattables = append(responseChattables, audioMessage)
		} else if callbackQuery.Message.Video != nil {
			videoMessage := botAPI.VideoConfig{
				BaseFile: botAPI.BaseFile{
					BaseChat:    botAPI.BaseChat{ChannelUsername: configuration.KanalConfig.GetString("kanal-username")},
					FileID:      callbackQuery.Message.Video.FileID,
					UseExisting: true,
				},
				Caption: text,
			}
			videoMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			if replyMessageID > -1 {
				videoMessage.ReplyToMessageID = replyMessageID
			}
			responseChattables = append(responseChattables, videoMessage)
		} else {
			kanalMessage := botAPI.NewMessageToChannel(configuration.KanalConfig.GetString("kanal-username"), text)
			kanalMessage.ReplyMarkup = keyboard.NewEmojiInlineKeyboard(0, 0, 0, 0)
			if replyMessageID > -1 {
				kanalMessage.ReplyToMessageID = replyMessageID
			}
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
		userID := strconv.Itoa(callbackQuery.From.ID)
		tappedReactionIndex, _ := strconv.Atoi(splittedCallbackData[1])
		tappedReaction := model.ConvertReactionIndexToReaction(tappedReactionIndex - 1)

		logrus.Infof("database read for reaction")
		var messageData model.Message
		err := db.MessagesCollection.Find(bson.M{
			"messageid": callbackQuery.Message.MessageID,
		}).One(&messageData)
		logrus.Infof("database read for reaction complete")
		if err != nil {
			messageData = model.Message{
				MessageID:    callbackQuery.Message.MessageID,
				ReactionsMap: make(map[string]model.Reaction),
			}
			messageData.ReactionsMap[userID] = tappedReaction
			go db.MessagesCollection.Insert(messageData)
		} else {
			if messageData.ReactionsMap[userID] == tappedReaction {
				delete(messageData.ReactionsMap, userID)
			} else {
				messageData.ReactionsMap[userID] = tappedReaction
			}
			go db.MessagesCollection.Update(bson.M{"_id": messageData.ID},
				bson.M{"$set": bson.M{"reactionsmap": messageData.ReactionsMap}})
		}
		logrus.Infof("reactions updated/inserted")

		// Counting
		reactionsCount := map[model.Reaction]int{}
		for _, value := range messageData.ReactionsMap {
			if _, ok := reactionsCount[value]; !ok {
				reactionsCount[value] = 0
			}
			reactionsCount[value]++
		}

		editedMessage := botAPI.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID,
			callbackQuery.Message.MessageID, keyboard.NewEmojiInlineKeyboard(
				reactionsCount[model.ReactionLike],
				reactionsCount[model.ReactionLol],
				reactionsCount[model.ReactionWow],
				reactionsCount[model.ReactionSad]))
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
			Payload:      map[string]interface{}{},
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

	if state.CommandState == model.UserCommandStateNothing {
		errorMessage := botAPI.NewMessage(message.Chat.ID, model.ErrorMessage)
		errorMessage.ReplyMarkup = keyboard.NewMainKeyboard()
		answerMessages = append(answerMessages, errorMessage)
		return answerMessages
	}

	// Post to Kanal Admins
	if configuration.KanalConfig.GetBool("debug") {
		kanalArchiveMessage := botAPI.NewForward(configuration.KanalConfig.GetInt64("kanal-archive-chatid"), message.Chat.ID, message.MessageID)
		answerMessages = append(answerMessages, kanalArchiveMessage)
	}
	if state.CommandState == model.UserCommandStateNewGIFCaption {
		gifMessage := botAPI.NewDocumentShare(configuration.KanalConfig.GetInt64("kanal-admins-chatid"),
			state.Payload["file-id"].(string))
		gifMessage.Caption = message.Text
		gifMessage.ReplyMarkup = keyboard.NewAdminInlineKeyboard(strconv.Itoa(message.MessageID))
		answerMessages = append(answerMessages, gifMessage)
	} else if state.CommandState == model.UserCommandStateNewAudioCaption {
		audioMessage := botAPI.NewAudioShare(configuration.KanalConfig.GetInt64("kanal-admins-chatid"),
			state.Payload["file-id"].(string))
		audioMessage.Caption = message.Text
		audioMessage.ReplyMarkup = keyboard.NewAdminInlineKeyboard(strconv.Itoa(message.MessageID))
		answerMessages = append(answerMessages, audioMessage)
	} else if state.CommandState == model.UserCommandStateNewVoiceCaption {
		voiceMessage := botAPI.NewVoiceShare(configuration.KanalConfig.GetInt64("kanal-admins-chatid"),
			state.Payload["file-id"].(string))
		voiceMessage.Caption = message.Text
		voiceMessage.ReplyMarkup = keyboard.NewAdminInlineKeyboard(strconv.Itoa(message.MessageID))
		answerMessages = append(answerMessages, voiceMessage)
	} else if message.Voice != nil {
		state.CommandState = model.UserCommandStateNewVoiceCaption
		state.Payload["file-id"] = message.Voice.FileID
		userState.Set(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)

		captionInputMessage := botAPI.NewMessage(message.Chat.ID, model.NewVoiceCommandMessage)
		captionInputMessage.ReplyMarkup = keyboard.NewMessageCancelKeyboard()
		answerMessages = append(answerMessages, captionInputMessage)
		return answerMessages
	} else if message.Audio != nil {
		state.CommandState = model.UserCommandStateNewAudioCaption
		state.Payload["file-id"] = message.Audio.FileID
		userState.Set(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)

		captionInputMessage := botAPI.NewMessage(message.Chat.ID, model.NewAudioCommandMessage)
		captionInputMessage.ReplyMarkup = keyboard.NewMessageCancelKeyboard()
		answerMessages = append(answerMessages, captionInputMessage)
		return answerMessages
	} else if message.Document != nil {
		state.CommandState = model.UserCommandStateNewGIFCaption
		state.Payload["file-id"] = message.Document.FileID
		userState.Set(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)

		captionInputMessage := botAPI.NewMessage(message.Chat.ID, model.NewGIFCommandMessage)
		captionInputMessage.ReplyMarkup = keyboard.NewMessageCancelKeyboard()
		answerMessages = append(answerMessages, captionInputMessage)
		return answerMessages
	} else if message.Photo != nil {
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

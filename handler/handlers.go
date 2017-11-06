package handler

import (
	"strconv"
	"time"

	"gitlab.com/arha/kanal/configuration"
	"gitlab.com/arha/kanal/ui/keyboard"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/kanal/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

var (
	userState = cache.New(30*time.Minute, 10*time.Minute)
)

func HandleCallbacks(callbackQuery *botAPI.CallbackQuery) []botAPI.Chattable {
	message := botAPI.NewMessageToChannel(configuration.KanalConfig.GetString("kanal-username"), callbackQuery.Message.Text)
	return []botAPI.Chattable{
		message,
	}
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
	kanalMessage := botAPI.NewMessage(configuration.KanalConfig.GetInt64("kanal-admins-chatid"), message.Text)
	kanalMessage.ReplyMarkup = keyboard.NewAdminInlineKeyboard(strconv.Itoa(message.MessageID))
	answerMessages = append(answerMessages, kanalMessage, kanalArchiveMessage)

	// Successful message
	successfulSentMessage := botAPI.NewMessage(message.Chat.ID, model.NewMessageSentMessage)
	successfulSentMessage.ReplyMarkup = keyboard.NewMainKeyboard()
	answerMessages = append(answerMessages, successfulSentMessage)

	userState.Delete(strconv.Itoa(message.From.ID))

	return answerMessages
}

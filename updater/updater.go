package updater

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/arha/kanal/configuration"
	"gitlab.com/arha/kanal/handler"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

var (
	bot      *botAPI.BotAPI
	botToken string
)

func InitializeUpdater() {
	botToken = configuration.KanalConfig.GetString("bot-token")

	var err error
	bot, err = botAPI.NewBotAPI(botToken)
	if err != nil {
		log.WithError(err).Fatalln("can't initialize bot")
	}
	bot.Debug = configuration.KanalConfig.GetBool("debug")

	log.Infof("authorized on account %s", bot.Self.UserName)
}

func Update() {
	updateConfig := botAPI.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.WithError(err).Fatalln("updater can't init channel")
	}

	for update := range updates {
		go func(update botAPI.Update) {
			var answers []botAPI.Chattable
			if update.Message != nil {
				answers = handler.HandleMessage(update.Message)
			} else {
				answers = handler.HandleCallbacks(update.CallbackQuery)
			}
			for _, answer := range answers {
				bot.Send(answer)
			}
		}(update)
	}
}

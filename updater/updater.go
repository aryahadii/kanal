package updater

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/arha/kanal/configuration"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

var (
	bot      *botAPI.BotAPI
	botToken string
)

func init() {
	botToken = configuration.KanalConfig.GetString("bot-token")

	var err error
	bot, err = botAPI.NewBotAPI(botToken)
	if err != nil {
		log.WithError(err).Fatalln("can't initialize bot")
	}
	bot.Debug = configuration.KanalConfig.GetString("debug")

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
			if update.Message != nil {
			}
		}(update)
	}
}

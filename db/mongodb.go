package db

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/arha/kanal/configuration"
	mgo "gopkg.in/mgo.v2"
)

var (
	MessagesCollection *mgo.Collection
	session            *mgo.Session
)

const (
	dbName                 = "kanal"
	messagesCollectionName = "messages"
)

func InitMongoDB() {
	var err error
	session, err = mgo.Dial(configuration.KanalConfig.GetString("mongodb-address"))
	if err != nil {
		log.WithError(err).Fatalln("MongoDB session can't be created")
	}

	MessagesCollection = session.DB(dbName).C(messagesCollectionName)
}

func Close() {
	session.Close()
}

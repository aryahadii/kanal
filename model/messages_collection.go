package model

type Message struct {
	MessageID int        `bson:"messageid"`
	Reactions [][]string `bson:"reactions"`
}

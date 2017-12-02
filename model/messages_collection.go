package model

import "gopkg.in/mgo.v2/bson"

type Message struct {
	ID           bson.ObjectId       `bson:"_id,omitempty"`
	MessageID    int                 `bson:"messageid"`
	ReactionsMap map[string]Reaction `bson:"reactionsmap"`
}

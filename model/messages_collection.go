package model

type Message struct {
	MessageID  int      `bson:"messageid"`
	Type1Emoji []string `bson:"type1_emoji"`
	Type2Emoji []string `bson:"type2_emoji"`
	Type3Emoji []string `bson:"type3_emoji"`
}

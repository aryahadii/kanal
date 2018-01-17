package model

type Reaction string

const (
	ReactionNothing  Reaction = ""
	ReactionLike     Reaction = "👍"
	ReactionLol      Reaction = "😂"
	ReactionWow      Reaction = "😧"
	ReactionSad      Reaction = "😞"
	ReactionPositive Reaction = "👍"
	ReactionNegative Reaction = "👎"
)

func ConvertReactionIndexToReaction(index int) Reaction {
	switch index {
	case 0:
		return ReactionLike
	case 1:
		return ReactionLol
	case 2:
		return ReactionWow
	case 3:
		return ReactionSad
	case 4:
		return ReactionPositive
	case 5:
		return ReactionNegative
	}
	return ReactionNothing
}

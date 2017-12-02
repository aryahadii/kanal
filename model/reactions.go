package model

type Reaction string

const (
	ReactionNothing Reaction = ""
	ReactionLike    Reaction = "👍"
	ReactionLol     Reaction = "😂"
	ReactionWow     Reaction = "😧"
	ReactionSad     Reaction = "😞"
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
	}
	return ReactionNothing
}

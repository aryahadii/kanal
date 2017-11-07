package handler

func removeUserReaction(userID string, reactions [][]string) bool {
	for i, reaction := range reactions {
		for j, user := range reaction {
			if userID == user {
				reactions[i][j] = reactions[i][len(reactions[i])-1]
				reactions[i] = reactions[i][:len(reactions[i])-1]
				return true
			}
		}
	}
	return false
}

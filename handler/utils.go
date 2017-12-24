package handler

import (
	"fmt"
	"regexp"
	"strconv"

	"gitlab.com/arha/kanal/configuration"
)

var (
	messageLinkPattern *regexp.Regexp
)

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

func findMessageIDInText(text string) (string, int) {
	if messageLinkPattern == nil {
		urlRegex := fmt.Sprintf(`https:\/\/t\.me\/%s\/(?P<ID>\d+)`,
			configuration.KanalConfig.GetString("kanal-username")[1:])
		messageLinkPattern = regexp.MustCompile(urlRegex)
	}

	id := -1
	if submatches := messageLinkPattern.FindStringSubmatch(text); len(submatches) > 0 {
		id, _ = strconv.Atoi(submatches[1])
		text = messageLinkPattern.ReplaceAllString(text, "")
	}
	return text, id
}

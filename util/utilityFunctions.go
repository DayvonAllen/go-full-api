package util

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func GenerateNewBlockList(targetID string, blockList []string) ([]string, bool) {
	userIsBlocked := false
	newBlockList := make([]string, 0, len(blockList))
	for i, foundId := range blockList {
		if foundId == targetID {
			userIsBlocked = true
			continue
		}
		newBlockList = append(newBlockList, foundId)

		if i == len(blockList) - 1{
			return newBlockList, userIsBlocked
		}
	}

	return newBlockList, userIsBlocked
}

func IsEmail(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func GenerateKey(value string, query string) string {
	var key strings.Builder

	for _, v := range strings.Fields(value) {
		key.WriteString(v)
	}

	key.WriteString(":")
	key.WriteString(query)

	return key.String()
}
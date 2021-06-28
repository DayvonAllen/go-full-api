package util

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func GenerateNewBlockList(targetID primitive.ObjectID, blockList []primitive.ObjectID) ([]primitive.ObjectID, bool) {
	userIsBlocked := false
	newBlockList := make([]primitive.ObjectID, 0, len(blockList))
	for _, foundId := range blockList {
		if foundId == targetID {
			userIsBlocked = true
			continue
		}
		newBlockList = append(newBlockList, foundId)
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
		fmt.Println(v)
		key.WriteString(v)
	}

	key.WriteString(":")
	key.WriteString(query)

	return key.String()
}
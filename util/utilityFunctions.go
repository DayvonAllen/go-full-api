package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func Find(targetArr []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, foundId := range targetArr {
		if foundId == id {
			return true
		}
	}
	return false
}

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
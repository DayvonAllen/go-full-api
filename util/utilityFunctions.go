package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
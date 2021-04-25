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

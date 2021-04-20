package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flag struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FlaggerID  primitive.ObjectID `bson:"flaggerID" json:"flaggerID,omitempty"`
	FlaggedUserID primitive.ObjectID `bson:"flaggedUserID" json:"flaggedUserID,omitempty"`
	Reason  string             `bson:"reason" json:"reason"`
}

package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	IsLocked  bool               `bson:"isLocked" json:"isLocked"`
	IsVerified  bool             `bson:"isVerified" json:"isVerified"`
	TokenHash string             `bson:"tokenHash" json:"tokenHash"`
	VerificationCode string      `bson:"verificationCode" json:"verificationCode"`
	TokenExpiresAt int64         `bson:"tokenExpiresAt" json:"tokenExpiresAt"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt,omitempty"`
}

type CreateUserDto struct {
	Username  string
	Email     string
	Password  string
}

type UpdateUserDto struct {
	Username  string
	Email     string
}

type UserDto struct {
	Id primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email string `bson:"email" json:"email"`
}
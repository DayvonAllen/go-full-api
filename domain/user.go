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
	CurrentTagLine  string       `bson:"currentTagLine" json:"CurrentTagLine"`
	UnlockedTagLine  []string    `bson:"unlockedTagLine" json:"unlockedTagLine"`
	ProfilePictureUrl  string    `bson:"profilePictureUrl" json:"profilePictureUrl"`
	CurrentBadgeUrl  string      `bson:"currentBadgeUrl" json:"currentBadgeUrl"`
	UnlockedBadgesUrls  []string `bson:"unlockedBadgesUrls" json:"unlockedBadgesUrls"`
	BlockList []string			 `bson:"blockList" json:"blockList"`
	BlockByList []string		 `bson:"blockByList" json:"blockByList"`
	FlagCount []Flag			 `bson:"flagCount" json:"flagCount"`
	ProfileIsViewable  bool      `bson:"profileIsViewable" json:"profileIsViewable"`
	IsLocked  bool               `bson:"isLocked" json:"isLocked"`
	IsVerified  bool             `bson:"isVerified" json:"isVerified"`
	AcceptMessages  bool         `bson:"acceptMessages" json:"acceptMessages"`
	TokenHash string             `bson:"tokenHash" json:"tokenHash"`
	VerificationCode string      `bson:"verificationCode" json:"verificationCode"`
	TokenExpiresAt int64         `bson:"tokenExpiresAt" json:"tokenExpiresAt"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt,omitempty"`
}

type CreateUserDto struct {
	Username  string  `json:"username,omitempty"`
	Email     string  `json:"email,omitempty"`
	Password  string  `json:"password,omitempty"`
}

type UpdateProfileVisibility struct {
	ProfileIsViewable  bool  `json:"profileIsViewable,omitempty"`
}

type UpdateMessageAcceptance struct {
	AcceptMessages  bool    `json:"acceptMessages,omitempty"`
}

type UpdateCurrentBadge struct {
	CurrentBadgeUrl  string `json:"currentBadgeUrl,omitempty"`
}

type UpdateProfilePicture struct {
	ProfilePictureUrl  string `json:"profilePictureUrl,omitempty"`
}

type UpdateCurrentTagline struct {
	CurrentTagLine  string  `json:"currentTagLine,omitempty"`
}

type UpdateVerification struct {
	IsVerified  bool       `json:"isVerified,omitempty"`
}

type UserDto struct {
	Email string                 `json:"email,omitempty"`
	CurrentTagLine  string       `json:"currentTagLine,omitempty"`
	UnlockedTagLine  []string    `json:"unlockedTagLine,omitempty"`
	ProfilePictureUrl  string    `json:"profilePictureUrl,omitempty"`
	CurrentBadgeUrl  string      `json:"currentBadgeUrl,omitempty"`
	UnlockedBadgesUrls  []string `json:"unlockedBadgesUrls,omitempty"`
	ProfileIsViewable  bool      `json:"profileIsViewable,omitempty"`
	AcceptMessages  bool         `json:"acceptMessages,omitempty"`
}
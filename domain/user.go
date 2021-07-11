package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User todo validate struct
type User struct {
	Id                          primitive.ObjectID   `bson:"_id" json:"id"`
	Username                    string               `bson:"username" json:"username"`
	Email                       string               `bson:"email" json:"email"`
	Password                    string               `bson:"password" json:"-"`
	CurrentTagLine              string               `bson:"currentTagLine" json:"CurrentTagLine"`
	UnlockedTagLine             []string             `bson:"unlockedTagLine" json:"unlockedTagLine"`
	ProfilePictureUrl           string               `bson:"profilePictureUrl" json:"profilePictureUrl"`
	ProfileBackgroundPictureUrl string               `bson:"profileBackgroundPictureUrl" json:"profileBackgroundPictureUrl"`
	CurrentBadgeUrl             string               `bson:"currentBadgeUrl" json:"currentBadgeUrl"`
	UnlockedBadgesUrls          []string             `bson:"unlockedBadgesUrls" json:"unlockedBadgesUrls"`
	BlockList                   []string `bson:"blockList" json:"blockList"`
	BlockByList                 []string `bson:"blockByList" json:"blockByList"`
	FlagCount                   []primitive.ObjectID `bson:"flagCount" json:"-"`
	Followers                   []string             `bson:"followers" json:"followers"`
	Following                   []string             `bson:"following" json:"following"`
	FollowerCount               int                  `bson:"followerCount" json:"followerCount"`
	DisplayFollowerCount        bool                 `bson:"displayFollowerCount" json:"displayFollowerCount"`
	ProfileIsViewable           bool                 `bson:"profileIsViewable" json:"profileIsViewable"`
	IsLocked                    bool                 `bson:"isLocked" json:"-"`
	IsVerified                  bool                 `bson:"isVerified" json:"isVerified"`
	AcceptMessages              bool                 `bson:"acceptMessages" json:"acceptMessages"`
	TokenHash                   string               `bson:"tokenHash" json:"-"`
	VerificationCode            string               `bson:"verificationCode" json:"-"`
	TokenExpiresAt              int64                `bson:"tokenExpiresAt" json:"-"`
	LastLoginIp					string				 `bson:"lastLoginIp" json:"-"`
	LastLoginIps				[]string			 `bson:"lastLoginIps" json:"-"`
	CreatedAt                   time.Time            `bson:"createdAt" json:"-"`
	UpdatedAt                   time.Time            `bson:"updatedAt" json:"-"`
}

type CreateUserDto struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateProfileVisibility struct {
	ProfileIsViewable bool      `json:"profileIsViewable,omitempty"`
	UpdatedAt         time.Time `bson:"updatedAt" json:"-"`
}

type UpdateMessageAcceptance struct {
	AcceptMessages bool      `json:"acceptMessages,omitempty"`
	UpdatedAt      time.Time `bson:"updatedAt" json:"-"`
}

type UpdateCurrentBadge struct {
	CurrentBadgeUrl string    `json:"currentBadgeUrl,omitempty"`
	UpdatedAt       time.Time `bson:"updatedAt" json:"-"`
}

type UpdateProfilePicture struct {
	ProfilePictureUrl string    `json:"profilePictureUrl,omitempty"`
	UpdatedAt         time.Time `bson:"updatedAt" json:"-"`
}

type UpdateProfileBackgroundPicture struct {
	ProfileBackgroundPictureUrl string    `json:"profileBackgroundPictureUrl,omitempty"`
	UpdatedAt                   time.Time `bson:"updatedAt" json:"-"`
}

type UpdateCurrentTagline struct {
	CurrentTagLine string    `json:"currentTagLine,omitempty"`
	UpdatedAt      time.Time `bson:"updatedAt" json:"-"`
}

type UpdateVerification struct {
	IsVerified bool      `json:"isVerified,omitempty"`
	UpdatedAt  time.Time `bson:"updatedAt" json:"-"`
}

type UpdateDisplayFollowerCount struct {
	DisplayFollowerCount bool      `json:"displayFollowerCount"`
	UpdatedAt            time.Time `bson:"updatedAt" json:"-"`
}

type UserDto struct {
	Id                          primitive.ObjectID   `bson:"_id" json:"-"`
	Email                       string               `json:"email"`
	Username                    string               `json:"username"`
	CurrentTagLine              string               `json:"currentTagLine"`
	UnlockedTagLine             []string             `json:"unlockedTagLine"`
	ProfilePictureUrl           string               `json:"profilePictureUrl"`
	ProfileBackgroundPictureUrl string               `json:"profileBackgroundPictureUrl"`
	CurrentBadgeUrl             string               `json:"currentBadgeUrl"`
	UnlockedBadgesUrls          []string             `json:"unlockedBadgesUrls"`
	ProfileIsViewable           bool                 `json:"profileIsViewable"`
	AcceptMessages              bool                 `json:"acceptMessages"`
	FollowerCount               int                  `json:"followerCount"`
	DisplayFollowerCount        bool                 `json:"displayFollowerCount"`
	Followers                   []string             `bson:"followers" json:"-"`
	Following                   []string             `bson:"following" json:"-"`
	IsVerified                  bool                 `bson:"isVerified" json:"-"`
	BlockList                   []string `bson:"blockList" json:"-"`
	BlockByList                 []string `bson:"blockByList" json:"-"`
	TokenHash                   string               `bson:"tokenHash" json:"-"`
	VerificationCode            string               `bson:"verificationCode" json:"-"`
	TokenExpiresAt              int64                `bson:"tokenExpiresAt" json:"-"`
}

type ViewUserProfile struct {
	Username                    string               `json:"username"`
	CurrentTagLine              string               `json:"currentTagLine"`
	ProfilePictureUrl           string               `json:"profilePictureUrl"`
	ProfileBackgroundPictureUrl string               `json:"profileBackgroundPictureUrl"`
	CurrentBadgeUrl             string               `json:"currentBadgeUrl"`
	FollowerCount               int                  `json:"followerCount"`
	DisplayFollowerCount        bool                 `json:"displayFollowerCount"`
}

type UserResponse struct {
	Users       *[]UserDto
	CurrentPage string
}

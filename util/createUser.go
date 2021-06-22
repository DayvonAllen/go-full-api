package util

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func CreateUser(createUserDto *domain.CreateUserDto) *domain.User {
	user := new(domain.User)

	user.Username = strings.ToLower(createUserDto.Username)
	user.Email = strings.ToLower(createUserDto.Email)
	user.Password = createUserDto.Password
	user.IsVerified = false
	user.IsLocked = false
	user.ProfileIsViewable = true
	user.AcceptMessages = true
	user.BlockList = []primitive.ObjectID{}
	user.BlockByList = []primitive.ObjectID{}
	user.FlagCount = []primitive.ObjectID{}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return user
}

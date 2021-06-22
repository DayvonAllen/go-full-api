package repo

import (
	"context"
	"example.com/app/domain"
	cache2 "github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	FindAll(primitive.ObjectID, *cache2.Cache, context.Context) (*[]domain.UserDto, error)
	FindAllBlockedUsers(primitive.ObjectID) (*[]domain.UserDto, error)
	Create(*domain.User) error
	FindByID(primitive.ObjectID) (*domain.UserDto, error)
	FindByUsername(string) (*domain.UserDto, error)
	UpdateByID(primitive.ObjectID, *domain.User) (*domain.UserDto, error)
	UpdateProfileVisibility(primitive.ObjectID, *domain.UpdateProfileVisibility) error
	UpdateMessageAcceptance(primitive.ObjectID, *domain.UpdateMessageAcceptance) error
	UpdateCurrentBadge(primitive.ObjectID, *domain.UpdateCurrentBadge) error
	UpdateProfilePicture(primitive.ObjectID, *domain.UpdateProfilePicture) error
	UpdateProfileBackgroundPicture(primitive.ObjectID, *domain.UpdateProfileBackgroundPicture) error
	UpdateCurrentTagline(primitive.ObjectID, *domain.UpdateCurrentTagline)  error
	UpdateVerification(primitive.ObjectID, *domain.UpdateVerification) error
	UpdatePassword(primitive.ObjectID, string) error
	UpdateFlagCount(*domain.Flag) error
	BlockUser(primitive.ObjectID, string) error
	UnBlockUser(primitive.ObjectID, string) error
	DeleteByID(primitive.ObjectID) error
}

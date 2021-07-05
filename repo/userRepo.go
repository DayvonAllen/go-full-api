package repo

import (
	"context"
	"example.com/app/domain"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	FindAll(primitive.ObjectID, string, context.Context, *cache2.Cache, string, opentracing.Span) (*domain.UserResponse, error)
	FindAllBlockedUsers(primitive.ObjectID, *cache2.Cache, context.Context, string) (*[]domain.UserDto, error)
	Create(*domain.User) error
	FindByID(primitive.ObjectID, *cache2.Cache, context.Context) (*domain.UserDto, error)
	FindByUsername(string, *cache2.Cache, context.Context) (*domain.UserDto, error)
	UpdateByID(primitive.ObjectID, *domain.User) (*domain.UserDto, error)
	UpdateProfileVisibility(primitive.ObjectID, *domain.UpdateProfileVisibility, *cache2.Cache, context.Context) error
	UpdateMessageAcceptance(primitive.ObjectID, *domain.UpdateMessageAcceptance, *cache2.Cache, context.Context) error
	UpdateCurrentBadge(primitive.ObjectID, *domain.UpdateCurrentBadge, *cache2.Cache, context.Context) error
	UpdateProfilePicture(primitive.ObjectID, *domain.UpdateProfilePicture, *cache2.Cache, context.Context) error
	UpdateProfileBackgroundPicture(primitive.ObjectID, *domain.UpdateProfileBackgroundPicture, *cache2.Cache, context.Context) error
	UpdateCurrentTagline(primitive.ObjectID, *domain.UpdateCurrentTagline, *cache2.Cache, context.Context)  error
	UpdateVerification(primitive.ObjectID, *domain.UpdateVerification) error
	UpdatePassword(primitive.ObjectID, string) error
	UpdateFlagCount(*domain.Flag) error
	BlockUser(primitive.ObjectID, string, *cache2.Cache, context.Context, string) error
	UnblockUser(primitive.ObjectID, string, *cache2.Cache, context.Context, string) error
	DeleteByID(primitive.ObjectID, *cache2.Cache, context.Context, string) error
}

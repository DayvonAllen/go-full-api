package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	FindAll() (*[]domain.UserDto, error)
	Create(*domain.User) error
	FindByID(primitive.ObjectID) (*domain.UserDto, error)
	UpdateByID(primitive.ObjectID, *domain.User) (*domain.UserDto, error)
	UpdatePassword(string, *domain.User) (*domain.UserDto, error)
	DeleteByID(primitive.ObjectID) error
}

package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	FindAll() (*[]domain.User, error)
	Create(*domain.User) error
	FindByID(primitive.ObjectID) (*domain.User, error)
	UpdateByID(primitive.ObjectID, *domain.User) (*domain.User, error)
	DeleteByID(primitive.ObjectID) error
}

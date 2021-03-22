package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService interface {
	GetAllUsers() (*[]domain.UserDto, error)
	CreateUser(*domain.User) error
	UpdateUser(primitive.ObjectID, *domain.User) (*domain.UserDto, error)
	UpdatePassword(string, *domain.User) (*domain.UserDto, error)
	GetUserByID(primitive.ObjectID) (*domain.UserDto, error)
	DeleteByID(primitive.ObjectID) error
}

// the service has a dependency of the repo
type DefaultUserService struct {
	repo repo.UserRepo
}

func (s DefaultUserService) GetAllUsers() (*[]domain.UserDto, error) {
	u, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return  u, nil
}

func (s DefaultUserService) GetUserByID(id primitive.ObjectID) (*domain.UserDto, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) UpdateUser(id primitive.ObjectID, user *domain.User) (*domain.UserDto, error) {
	u, err := s.repo.UpdateByID(id, user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) UpdatePassword(password string, user *domain.User) (*domain.UserDto, error) {
	u, err := s.repo.UpdatePassword(password, user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) DeleteByID(id primitive.ObjectID) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) CreateUser(user *domain.User) error {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	err := s.repo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}
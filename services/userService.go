package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService interface {
	GetAllUsers() (*[]domain.UserDto, error)
	CreateUser(*domain.User) error
	GetUserByID(primitive.ObjectID) (*domain.UserDto, error)
	GetUserByUsername(string) (*domain.UserDto, error)
	UpdateProfileVisibility(primitive.ObjectID, *domain.UpdateProfileVisibility) error
	UpdateMessageAcceptance(primitive.ObjectID, *domain.UpdateMessageAcceptance) error
	UpdateCurrentBadge(primitive.ObjectID, *domain.UpdateCurrentBadge) error
	UpdateProfilePicture(primitive.ObjectID, *domain.UpdateProfilePicture) error
	UpdateCurrentTagline(primitive.ObjectID, *domain.UpdateCurrentTagline)  error
	UpdateVerification(primitive.ObjectID, *domain.UpdateVerification) error
	UpdatePassword(primitive.ObjectID, string) error
	UpdateFlagCount(*domain.Flag) error
	DeleteByID(primitive.ObjectID) error
}

// DefaultUserService the service has a dependency of the repo
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

func (s DefaultUserService) CreateUser(user *domain.User) error {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	a := new(domain.Authentication)
	h := utils.UUIDv4()
	signedHash, err := a.SignToken([]byte(h))

	if err != nil {
		return err
	}

	hash := h + "-" + string(signedHash)
	user.VerificationCode = hash

	err = s.repo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) GetUserByID(id primitive.ObjectID) (*domain.UserDto, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) GetUserByUsername(username string) (*domain.UserDto, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) UpdateProfileVisibility(id primitive.ObjectID, user *domain.UpdateProfileVisibility) error {
	err := s.repo.UpdateProfileVisibility(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance) error {
	err := s.repo.UpdateMessageAcceptance(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge) error {
	err := s.repo.UpdateCurrentBadge(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture) error {
	err := s.repo.UpdateProfilePicture(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline) error {
	err := s.repo.UpdateCurrentTagline(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdatePassword(id primitive.ObjectID, password string) error {
	err := s.repo.UpdatePassword(id, password)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateVerification(id primitive.ObjectID, user *domain.UpdateVerification) error {
	err := s.repo.UpdateVerification(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateFlagCount(flag *domain.Flag) error {
	err := s.repo.UpdateFlagCount(flag)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) DeleteByID(id primitive.ObjectID) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}
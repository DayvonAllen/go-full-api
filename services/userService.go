package services

import (
	"context"
	"example.com/app/domain"
	"example.com/app/repo"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserService interface {
	GetAllUsers(primitive.ObjectID, string, context.Context) (*domain.UserResponse, error)
	GetAllBlockedUsers(primitive.ObjectID) (*[]domain.UserDto, error)
	CreateUser(*domain.User) error
	GetUserByID(primitive.ObjectID) (*domain.UserDto, error)
	GetUserByUsername(string, *cache2.Cache, context.Context) (*domain.UserDto, error)
	UpdateProfileVisibility(primitive.ObjectID, *domain.UpdateProfileVisibility, *cache2.Cache, context.Context) error
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

// DefaultUserService the service has a dependency of the repo
type DefaultUserService struct {
	repo repo.UserRepo
}

func (s DefaultUserService) GetAllUsers(id primitive.ObjectID, page string, ctx context.Context) (*domain.UserResponse, error) {
	u, err := s.repo.FindAll(id, page, ctx)
	if err != nil {
		return nil, err
	}
	return  u, nil
}

func (s DefaultUserService) GetAllBlockedUsers(id primitive.ObjectID) (*[]domain.UserDto, error) {
	u, err := s.repo.FindAllBlockedUsers(id)
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

func (s DefaultUserService) GetUserByUsername(username string, rdb *cache2.Cache, ctx context.Context) (*domain.UserDto, error) {
	u, err := s.repo.FindByUsername(username, rdb, ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s DefaultUserService) UpdateProfileVisibility(id primitive.ObjectID, user *domain.UpdateProfileVisibility, rdb *cache2.Cache, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfileVisibility(id, user, rdb, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateMessageAcceptance(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateCurrentBadge(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfilePicture(id, user)
	if err != nil {
		return err
	}
	return nil
}
func (s DefaultUserService) UpdateProfileBackgroundPicture(id primitive.ObjectID, user *domain.UpdateProfileBackgroundPicture) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfileBackgroundPicture(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline) error {
	user.UpdatedAt = time.Now()
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
	user.UpdatedAt = time.Now()
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

func (s DefaultUserService) BlockUser(id primitive.ObjectID, username string) error {
	err := s.repo.BlockUser(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UnBlockUser(id primitive.ObjectID, username string) error {
	err := s.repo.UnBlockUser(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}
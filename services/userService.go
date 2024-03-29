package services

import (
	"context"
	"example.com/app/domain"
	"example.com/app/repo"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserService interface {
	GetAllUsers(primitive.ObjectID, string, context.Context, *cache2.Cache, string, opentracing.Span) (*domain.UserResponse, error)
	GetAllBlockedUsers(primitive.ObjectID, *cache2.Cache, context.Context, string) (*[]domain.UserDto, error)
	CreateUser(*domain.User) error
	GetUserByID(primitive.ObjectID, *cache2.Cache, context.Context) (*domain.UserDto, error)
	GetUserByUsername(string, *cache2.Cache, context.Context) (*domain.UserDto, error)
	UpdateProfileVisibility(primitive.ObjectID, *domain.UpdateProfileVisibility, *cache2.Cache, context.Context) error
	UpdateMessageAcceptance(primitive.ObjectID, *domain.UpdateMessageAcceptance, *cache2.Cache, context.Context) error
	UpdateCurrentBadge(primitive.ObjectID, *domain.UpdateCurrentBadge, *cache2.Cache, context.Context) error
	UpdateProfilePicture(primitive.ObjectID, *domain.UpdateProfilePicture, *cache2.Cache, context.Context) error
	UpdateProfileBackgroundPicture(primitive.ObjectID, *domain.UpdateProfileBackgroundPicture, *cache2.Cache, context.Context) error
	UpdateCurrentTagline(primitive.ObjectID, *domain.UpdateCurrentTagline, *cache2.Cache, context.Context)  error
	UpdateDisplayFollowerCount(primitive.ObjectID, *domain.UpdateDisplayFollowerCount, *cache2.Cache) error
	UpdateVerification(primitive.ObjectID, *domain.UpdateVerification) error
	UpdatePassword(primitive.ObjectID, string) error
	UpdateFlagCount(*domain.Flag) error
	FollowUser(username string, currentUser string, rdb *cache2.Cache) error
	UnfollowUser(username string, currentUser string, rdb *cache2.Cache) error
	BlockUser(primitive.ObjectID, string, *cache2.Cache, context.Context, string) error
	UnblockUser(primitive.ObjectID, string, *cache2.Cache, context.Context, string) error
	DeleteByID(primitive.ObjectID, *cache2.Cache, context.Context, string) error
}

// DefaultUserService the service has a dependency of the repo
type DefaultUserService struct {
	repo repo.UserRepo
}

func (s DefaultUserService) GetAllUsers(id primitive.ObjectID, page string, ctx context.Context, rdb *cache2.Cache, username string, span opentracing.Span) (*domain.UserResponse, error) {
	//childSpan := opentracing.StartSpan("child", opentracing.ChildOf(span.Context()))
	//defer childSpan.Finish()
	u, err := s.repo.FindAll(id, page, ctx, rdb, username, span)
	if err != nil {
		return nil, err
	}
	return  u, nil
}

func (s DefaultUserService) GetAllBlockedUsers(id primitive.ObjectID, rdb *cache2.Cache, ctx context.Context, username string) (*[]domain.UserDto, error) {
	u, err := s.repo.FindAllBlockedUsers(id, rdb, ctx, username)
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

func (s DefaultUserService) GetUserByID(id primitive.ObjectID, rdb *cache2.Cache, ctx context.Context) (*domain.UserDto, error) {
	u, err := s.repo.FindByID(id, rdb, ctx)
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

func (s DefaultUserService) UpdateDisplayFollowerCount(id primitive.ObjectID, user *domain.UpdateDisplayFollowerCount, rdb *cache2.Cache) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateDisplayFollowerCount(id, user, rdb)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance, rdb *cache2.Cache, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateMessageAcceptance(id, user, rdb, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge, rdb *cache2.Cache, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateCurrentBadge(id, user, rdb, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture, rdb *cache2.Cache, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfilePicture(id, user, rdb, ctx)
	if err != nil {
		return err
	}
	return nil
}
func (s DefaultUserService) UpdateProfileBackgroundPicture(id primitive.ObjectID, user *domain.UpdateProfileBackgroundPicture, rdb *cache2.Cache, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateProfileBackgroundPicture(id, user, rdb, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline, rdb *cache2.Cache, ctx context.Context) error {
	user.UpdatedAt = time.Now()
	err := s.repo.UpdateCurrentTagline(id, user, rdb, ctx)
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

func (s DefaultUserService) DeleteByID(id primitive.ObjectID, rdb *cache2.Cache, ctx context.Context, username string) error {
	err := s.repo.DeleteByID(id, rdb, ctx, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) FollowUser(username string, currentUser string, rdb *cache2.Cache) error {
	err := s.repo.FollowUser(username, currentUser, rdb)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UnfollowUser(username string, currentUser string, rdb *cache2.Cache) error {
	err := s.repo.UnfollowUser(username, currentUser, rdb)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) BlockUser(id primitive.ObjectID, username string, rdb *cache2.Cache, ctx context.Context, currentUsername string) error {
	err := s.repo.BlockUser(id, username, rdb, ctx, currentUsername)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) UnblockUser(id primitive.ObjectID, username string, rdb *cache2.Cache, ctx context.Context, currentUsername string) error {
	err := s.repo.UnblockUser(id, username, rdb, ctx, currentUsername)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}
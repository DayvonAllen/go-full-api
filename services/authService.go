package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AuthService interface {
	Login(username, password string) (*domain.User, string, error)
	ResetPasswordQuery(email string) error
	ResetPassword(token, password string) error
}

type DefaultAuthService struct {
	repo repo.AuthRepo
}

func (a DefaultAuthService) Login(email, password string) (*domain.User, string, error) {
	u, token, err := a.repo.Login(strings.ToLower(email), password)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func (a DefaultAuthService) ResetPasswordQuery(email string) error {
	err := a.repo.ResetPasswordQuery(strings.ToLower(email))
	if err != nil {
		return err
	}
	return nil
}

func (a DefaultAuthService) ResetPassword(token, password string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := a.repo.ResetPassword(token, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func NewAuthService(repository repo.AuthRepo) DefaultAuthService {
	return DefaultAuthService{repository}
}
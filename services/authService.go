package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (*domain.UserDto, string, error)
	ResetPasswordQuery(email string) error
	ResetPassword(token, password string) error
	VerifyCode(code string) error
}

type DefaultAuthService struct {
	repo repo.AuthRepo
}

func (a DefaultAuthService) Login(username, password string) (*domain.UserDto, string, error) {
	u, token, err := a.repo.Login(username, password)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func (a DefaultAuthService) ResetPasswordQuery(email string) error {
	err := a.repo.ResetPasswordQuery(email)
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

func (a DefaultAuthService) VerifyCode(code string) error {
	err := a.repo.VerifyCode(code)
	if err != nil {
		return err
	}
	return nil
}

func NewAuthService(repository repo.AuthRepo) DefaultAuthService {
	return DefaultAuthService{repository}
}
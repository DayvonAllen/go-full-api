package repo

import "example.com/app/domain"

type AuthRepo interface {
	Login(username, password string) (*domain.User, string, error)
	ResetPassword(token, password string) error
	ResetPasswordQuery(email string) error
	VerifyCode(code string) error
}


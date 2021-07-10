package repo

import "example.com/app/domain"

type AuthRepo interface {
	Login(username string, password string, ip string, ips []string) (*domain.UserDto, string, error)
	ResetPassword(token, password string) error
	ResetPasswordQuery(email string) error
	VerifyCode(code string) error
}


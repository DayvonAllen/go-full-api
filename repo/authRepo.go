package repo

import "example.com/app/domain"

type AuthRepo interface {
	Login(username, password string) (*domain.User, string, error)
}


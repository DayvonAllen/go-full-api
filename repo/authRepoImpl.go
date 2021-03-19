package repo

import (
	"context"
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
	*domain.User
}

func(a AuthRepoImpl) Login(email, password string) (*domain.User, string, error) {
	var login domain.Authentication
	var user domain.User
	opts := options.FindOne()
	err := dbConnection.Collection.FindOne(context.TODO(), bson.D{{"email",
		email}},opts).Decode(&user)

	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, "", err
	}

	token, err := login.GenerateJWT(user)

	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}
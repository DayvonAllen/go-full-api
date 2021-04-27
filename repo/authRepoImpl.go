package repo

import (
	"context"
	"example.com/app/config"
	"example.com/app/domain"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

type AuthRepoImpl struct {
	*domain.User
}

func(a AuthRepoImpl) Login(email, password string) (*domain.UserDto, string, error) {
	var login domain.Authentication
	var user domain.User
	opts := options.FindOne()
	err := dbConnection.UserCollection.FindOne(context.TODO(), bson.D{{"email",
		strings.ToLower(email)}},opts).Decode(&user)

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

	userDto := domain.UserMapper(&user)

	return userDto, token, nil
}

func(a AuthRepoImpl) ResetPasswordQuery(email string) error {
	var user domain.User
	err := dbConnection.UserCollection.FindOne(context.TODO(), bson.D{{"email", strings.ToLower(email)}}).Decode(&user)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("email %v was not found", email)
		}
		return err
	}

	// logic to send email with JWT
	if user.TokenHash == "" || user.TokenExpiresAt < time.Now().Unix() {
		a := new(domain.Authentication)
		h := utils.UUIDv4()
		s, err := a.SignToken([]byte(h))

		if err != nil {
			return err
		}

		expiration, err := strconv.Atoi(config.Config("PASSWORD_RESET_TOKEN_EXPIRATION"))

		if err != nil {
			return err
		}

		hash := h + "-" + string(s)
		user.TokenHash = hash
		user.TokenExpiresAt = time.Now().Add(time.Duration(expiration) * time.Minute).Unix()
		ur := new(UserRepoImpl)
		_, err = ur.UpdateByID(user.Id, &user)

		if err != nil {
			return err
		}
	}

	// send token url in email to user
	fmt.Println("http://127.0.0.1:8080/auth/reset/" + user.TokenHash)

	fmt.Println(user.TokenHash)

	return nil
}

func(a AuthRepoImpl) ResetPassword(token, password string) error {
	user := new(domain.User)
	ur := new(UserRepoImpl)
	err := dbConnection.UserCollection.FindOne(context.TODO(), bson.D{{"tokenHash", token}}).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no token found")
		}
		return err
	}

	if user.TokenExpiresAt < time.Now().Unix() {
		return fmt.Errorf("token has expired")
	}

	// update password logic
	err = ur.UpdatePassword(user.Id, password)

	if err != nil {
		return err
	}

	return nil
}

func (a AuthRepoImpl) VerifyCode(code string) error{
	var user domain.User
	ur := new(UserRepoImpl)
	err := dbConnection.UserCollection.FindOne(context.TODO(), bson.D{{"verificationCode", code}}).Decode(&user)

	if user.IsVerified {
		return fmt.Errorf("user email already verified")
	}

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no token found")
		}
		return err
	}

	u := new(domain.UpdateVerification)

	u.IsVerified = true

	err = ur.UpdateVerification(user.Id, u)

	if err != nil {
		return err
	}

	return nil
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}
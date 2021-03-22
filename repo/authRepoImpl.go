package repo

import (
	"context"
	"example.com/app/domain"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func(a AuthRepoImpl) ResetPasswordQuery(email string) error {
	var user domain.User
	err := dbConnection.Collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)

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

		hash := h + "-" + string(s)
		user.TokenHash = hash
		user.TokenExpiresAt = time.Now().Add(time.Duration(1) * time.Minute).Unix()
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
	var user domain.User
	ur := new(UserRepoImpl)
	err := dbConnection.Collection.FindOne(context.TODO(), bson.D{{"tokenHash", token}}).Decode(&user)

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
	_, err = ur.UpdatePassword(password, &user)

	if err != nil {
		return err
	}

	return nil
}

func (a AuthRepoImpl) VerifyCode(code string) error{
	var user domain.User
	ur := new(UserRepoImpl)
	err := dbConnection.Collection.FindOne(context.TODO(), bson.D{{"verificationCode", code}}).Decode(&user)

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

	user.IsVerified = true

	_, err = ur.UpdateVerification(&user)

	if err != nil {
		return err
	}

	return nil
}


func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}
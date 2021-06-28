package repo

import (
	"context"
	"example.com/app/config"
	"example.com/app/database"
	"example.com/app/domain"
	"example.com/app/util"
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

func(a AuthRepoImpl) Login(username, password string) (*domain.UserDto, string, error) {
	var login domain.Authentication
	var user domain.User

	conn := database.MongoConnectionPool.Get().(*database.Connection)

	if util.IsEmail(username) {
		opts := options.FindOne()
		err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"email",
			strings.ToLower(username)}},opts).Decode(&user)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return nil, "", fmt.Errorf("error finding by email")
		}
	} else {
		opts := options.FindOne()
		err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username",
			strings.ToLower(username)}},opts).Decode(&user)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return nil, "", fmt.Errorf("error finding by username")
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return nil, "", fmt.Errorf("error comparing password")
	}

	token, err := login.GenerateJWT(user)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return nil, "", fmt.Errorf("error generating token")
	}

	userDto := domain.UserMapper(&user)

	database.MongoConnectionPool.Put(conn)
	return userDto, token, nil
}

func(a AuthRepoImpl) ResetPasswordQuery(email string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	var user domain.User
	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"email", strings.ToLower(email)}}).Decode(&user)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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
			database.MongoConnectionPool.Put(conn)
			return err
		}

		expiration, err := strconv.Atoi(config.Config("PASSWORD_RESET_TOKEN_EXPIRATION"))

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return err
		}

		hash := h + "-" + string(s)
		user.TokenHash = hash
		user.TokenExpiresAt = time.Now().Add(time.Duration(expiration) * time.Minute).Unix()
		ur := new(UserRepoImpl)
		_, err = ur.UpdateByID(user.Id, &user)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return err
		}
	}

	// send token url in email to user
	fmt.Println("http://127.0.0.1:8080/auth/reset/" + user.TokenHash)

	fmt.Println(user.TokenHash)

	database.MongoConnectionPool.Put(conn)
	return nil
}

func(a AuthRepoImpl) ResetPassword(token, password string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	user := new(domain.User)
	ur := new(UserRepoImpl)
	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"tokenHash", token}}).Decode(&user)
	if err != nil {
		database.MongoConnectionPool.Put(conn)
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no token found")
		}
		return err
	}

	if user.TokenExpiresAt < time.Now().Unix() {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("token has expired")
	}

	// update password logic
	err = ur.UpdatePassword(user.Id, password)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return err
	}

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (a AuthRepoImpl) VerifyCode(code string) error{

	conn := database.MongoConnectionPool.Get().(*database.Connection)

	var user domain.User
	ur := new(UserRepoImpl)
	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"verificationCode", code}}).Decode(&user)

	if user.IsVerified {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("user email already verified")
	}

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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
		database.MongoConnectionPool.Put(conn)
		return err
	}

	database.MongoConnectionPool.Put(conn)
	return nil
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}
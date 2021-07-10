package repo

import (
	"context"
	"example.com/app/config"
	"example.com/app/database"
	"example.com/app/domain"
	"example.com/app/events"
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

func(a AuthRepoImpl) Login(username string, password string, ip string, ips []string) (*domain.UserDto, string, error) {
	var login domain.Authentication
	var user domain.User

	conn := database.MongoConnectionPool.Get().(*database.Connection)
	database.MongoConnectionPool.Put(conn)

	if util.IsEmail(username) {
		opts := options.FindOne()
		err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"email",
			username}},opts).Decode(&user)

		if err != nil {
			return nil, "", fmt.Errorf("error finding by email")
		}
	} else {
		opts := options.FindOne()
		err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username",
			username}},opts).Decode(&user)

		if err != nil {
			return nil, "", fmt.Errorf("error finding by username")
		}
	}

	fmt.Println("pass")

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, "", fmt.Errorf("error comparing password")
	}

	token, err := login.GenerateJWT(user)
	fmt.Println("after")


	if err != nil {
		return nil, "", fmt.Errorf("error generating token")
	}

	userDto := domain.UserMapper(&user)

	go func() {
		filter := bson.D{{"username", username}}
		update := bson.D{{"$set", bson.D{{"lastLoginIp", ip}, {"lastLoginIps", ips}}}}

		_, err := conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			panic(err)
		}
		return
	}()
	fmt.Println("ran")

	go func() {
		event := new(domain.Event)
		event.Action = "login"
		event.Target = username
		event.ResourceId = user.Id
		event.ActorUsername = username
		event.Message = username + " has logged in IP: " + ip + "; IPs: " + strings.Join(ips, ", ")
		err = events.SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	fmt.Println("ran")

	return userDto, token, nil
}

func(a AuthRepoImpl) ResetPasswordQuery(email string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	var user domain.User
	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"email", strings.ToLower(email)}}).Decode(&user)

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
	}

	// send token url in email to user
	fmt.Println("http://127.0.0.1:8080/auth/reset/" + user.TokenHash)

	fmt.Println(user.TokenHash)

	return nil
}

func(a AuthRepoImpl) ResetPassword(token, password string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	user := new(domain.User)
	ur := new(UserRepoImpl)
	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"tokenHash", token}}).Decode(&user)
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
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	var user domain.User
	ur := new(UserRepoImpl)
	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"verificationCode", code}}).Decode(&user)

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
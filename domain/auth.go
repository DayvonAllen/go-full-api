package domain

import (
	"crypto/hmac"
	"example.com/app/domain/helper"

	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type Authentication struct {
	Id primitive.ObjectID
	Email string
}

type LoginDetails struct {
	Email string
	Password string
}

type Claims struct {
	jwt.StandardClaims
	Id primitive.ObjectID
	Email string
}

func (l Authentication) GenerateJWT(msg User) (string, error){
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		Id: msg.Id,
		Email: msg.Email,
	}
	// always better to use a pointer with JSON
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedString, err := token.SignedString([]byte("Helloworld"))

	if err != nil {
		return "", err
	}
	return signedString, nil
}

func (l Authentication) SignToken(token []byte) ([]byte, error) {
	// second arg is a private key, key needs to be the same size as hasher
	// sha512 is 64 bits
	h := hmac.New(sha256.New, []byte("Helloworld"))

	// hash is a writer
	_, err := h.Write(token)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("%x", h.Sum(nil))), nil
}

func (l Authentication) VerifySignature(token, sig []byte) (bool, error) {
	// sign message
	s, _ := l.SignToken(token)
	// compare it
	return hmac.Equal(sig, s), nil
}

func(l Authentication) IsLoggedIn(cookie string) (*Authentication, error)  {
	if cookie == ""  {
		return nil, fmt.Errorf("no cookie")
	}

	data := helper.ExtractData(cookie)

	validSig, err := l.VerifySignature([]byte(data[0]), []byte(data[1]))
	if err != nil {
		return nil, err
	}

	if !validSig {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(data[0], &Claims{},func(t *jwt.Token)(interface{}, error) {
		if t.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			//verify token(we pass in our key to be verified)
			return []byte("Helloworld"), nil
		}
		return nil, err
	})

	if err != nil {
		return nil, err
	}

	isEqual := token.Valid

	if isEqual {
		// user is logged in at this point
		// because we receive an interface type we need to assert which type we want to use that inherits it
		claims := token.Claims.(*Claims)

		l.Id = claims.Id
		l.Email = claims.Email
		return &l, nil
	}

	return nil, fmt.Errorf("token is not valid")
}
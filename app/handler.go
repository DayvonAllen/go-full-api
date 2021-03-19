package app

import (
	"example.com/app/app/helpers"
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Handlers struct {
	UserService services.UserService
	AuthService services.AuthService
}

func (ch *Handlers) GetAllUsers(c *fiber.Ctx) error {
	cookie := c.Cookies("session")

	err := helpers.IsLoggedIn(cookie, c)

	users, err := ch.UserService.GetAllUsers()

	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func (ch *Handlers) CreateUser(c *fiber.Ctx) error {
	user := new(domain.User)

	err := c.BodyParser(user)

	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}
	err = ch.UserService.CreateUser(user)

	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *Handlers) GetUserByID(c *fiber.Ctx) error {
	cookie := c.Cookies("session")

	err := helpers.IsLoggedIn(cookie, c)

	if err != nil {
		return c.SendString(fmt.Sprintf("%v", err))
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		c.Status(400)
		return c.SendString(fmt.Sprintf("%v", err))
	}

	user, err := ch.UserService.GetUserByID(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(404)
			return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
		}
		c.Status(500)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (ch *Handlers) UpdateUser(c *fiber.Ctx) error {
	cookie := c.Cookies("session")

	err := helpers.IsLoggedIn(cookie, c)

	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	id , err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	user := new(domain.User)

	err = c.BodyParser(user)

	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	u, err := ch.UserService.UpdateUser(id, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(404)
			return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
		}
		c.Status(500)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "success", "data": u})
}

func (ch *Handlers) DeleteByID(c *fiber.Ctx) error {
	cookie := c.Cookies("session")

	err := helpers.IsLoggedIn(cookie, c)

	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	id , err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		c.Status(400)
		return c.SendString(fmt.Sprintf("%v", err))
	}

	err = ch.UserService.DeleteByID(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(404)
			return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
		}
		c.Status(500)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *Handlers) Login(c *fiber.Ctx) error {
	details := new(domain.LoginDetails)

	err := c.BodyParser(details)

	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	var auth domain.Authentication

	user, token, err := ch.AuthService.Login(details.Email, details.Password)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.Status(401)
			return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
		}
		c.Status(500)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	signedToken := make([]byte, 0, 100)
	signedToken = append(signedToken, []byte("Bearer " + token + "|")...)
	t, err := auth.SignToken([]byte(token))

	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	signedToken = append(signedToken, t...)

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = string(signedToken)
	cookie.Expires = time.Now().Add(24 * time.Hour)

	// Set cookie
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}
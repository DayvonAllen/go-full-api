package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	details := new(domain.LoginDetails)

	err := c.BodyParser(details)

	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"status": "error", "message": "error...", "data": err})
	}

	var auth domain.Authentication

	user, token, err := ah.AuthService.Login(details.Email, details.Password)

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
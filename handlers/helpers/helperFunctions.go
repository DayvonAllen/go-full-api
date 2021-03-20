package helpers

import (
	"example.com/app/domain"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func IsLoggedIn(token string, c *fiber.Ctx) error {
	var auth domain.Authentication
	_, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		c.Status(401)
		return fmt.Errorf("unauthorized request")
	}

	return nil
}

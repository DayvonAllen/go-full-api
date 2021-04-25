package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"example.com/app/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type UserHandler struct {
	UserService services.UserService
}

func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	users, err := uh.UserService.GetAllUsers(u.Id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func (uh *UserHandler) GetAllBlockedUsers(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	users, err := uh.UserService.GetAllBlockedUsers(u.Id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func (uh *UserHandler) CreateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")
	createUserDto := new(domain.CreateUserDto)

	err := c.BodyParser(createUserDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	user := util.CreateUser(createUserDto)

	err = uh.UserService.CreateUser(user)

	if err != nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) GetUserByID(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	user, err := uh.UserService.GetUserByID(u.Id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (uh *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := uh.UserService.GetUserByUsername(strings.ToLower(username))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (uh *UserHandler) UpdateProfileVisibility(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	userDto := new(domain.UpdateProfileVisibility)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateProfileVisibility(u.Id, userDto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateMessageAcceptance(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	userDto := new(domain.UpdateMessageAcceptance)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateMessageAcceptance(u.Id, userDto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateCurrentBadge(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	userDto := new(domain.UpdateCurrentBadge)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateCurrentBadge(u.Id, userDto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateProfilePicture(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	userDto := new(domain.UpdateProfilePicture)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateProfilePicture(u.Id, userDto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateProfileBackgroundPicture(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	userDto := new(domain.UpdateProfileBackgroundPicture)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateProfileBackgroundPicture(u.Id, userDto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateCurrentTagline(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	userDto := new(domain.UpdateCurrentTagline)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateCurrentTagline(u.Id, userDto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateFlagCount(c *fiber.Ctx) error {
	username := c.Params("username")
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)


	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	flag := new(domain.Flag)

	err = c.BodyParser(flag)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	flag.FlaggedUsername = strings.ToLower(username)
	flag.FlaggerID = u.Id

	err = uh.UserService.UpdateFlagCount(flag)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) DeleteByID(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.DeleteByID(u.Id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) BlockUser(c *fiber.Ctx) error {
	username := c.Params("username")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.BlockUser(u.Id, strings.ToLower(username))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UnBlockUser(c *fiber.Ctx) error {
	username := c.Params("username")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UnBlockUser(u.Id, strings.ToLower(username))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
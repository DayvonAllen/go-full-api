package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserHandler struct {
	UserService services.UserService
}

func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := uh.UserService.GetAllUsers()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func (uh *UserHandler) CreateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")
	user := new(domain.User)
	createUserDto := new(domain.CreateUserDto)

	err := c.BodyParser(createUserDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	user.Username = createUserDto.Username
	user.Email = createUserDto.Email
	user.Password = createUserDto.Password
	user.IsVerified = false
	user.IsLocked = false
	user.CreatedAt = time.Now()
	err = uh.UserService.CreateUser(user)

	if err != nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	user, err := uh.UserService.GetUserByID(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (uh *UserHandler) UpdateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")

	id , err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	userDto := new(domain.UpdateUserDto)
	user := new(domain.User)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	user.Username = userDto.Username
	user.Email = userDto.Email
	user.UpdatedAt = time.Now()
	u, err := uh.UserService.UpdateUser(id, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": u})
}

func (uh *UserHandler) DeleteByID(c *fiber.Ctx) error {
	id , err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.DeleteByID(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
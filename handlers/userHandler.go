package handlers

import (
	"example.com/app/cache"
	"example.com/app/domain"
	"example.com/app/services"
	"example.com/app/util"
	"fmt"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type UserHandler struct {
	UserService services.UserService
}

func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	page := c.Query("page", "1")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	users, err := uh.UserService.GetAllUsers(u.Id, page, c.Context(), rdb, u.Username)

	if err != nil {
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	cache.RedisCachePool.Put(rdb)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func (uh *UserHandler) GetAllBlockedUsers(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	users, err := uh.UserService.GetAllBlockedUsers(u.Id, rdb, c.Context(), u.Username)

	if err != nil {
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	cache.RedisCachePool.Put(rdb)

	c.Context().Done()
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

	err = uh.UserService.CreateUser(user, c.Context())

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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	var data domain.UserDto

	err = rdb.Get(c.Context(), util.GenerateKey(u.Username, "finduserbyusername"), &data)

	if err == nil {
		cache.RedisCachePool.Put(rdb)
		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": data})
	}

	user, err := uh.UserService.GetUserByID(u.Id, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	cache.RedisCachePool.Put(rdb)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (uh *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	var data domain.UserDto

	err := rdb.Get(c.Context(), util.GenerateKey(username, "finduserbyusername"), &data)

	if err == nil {
		fmt.Println("Found in cache in get user by username")
		cache.RedisCachePool.Put(rdb)
		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": data})
	}

	user, err := uh.UserService.GetUserByUsername(strings.ToLower(username), rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	cache.RedisCachePool.Put(rdb)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (uh *UserHandler) UpdateProfileVisibility(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	userDto := new(domain.UpdateProfileVisibility)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.UpdateProfileVisibility(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	userDto := new(domain.UpdateMessageAcceptance)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.UpdateMessageAcceptance(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	userDto := new(domain.UpdateCurrentBadge)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.UpdateCurrentBadge(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	userDto := new(domain.UpdateProfilePicture)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.UpdateProfilePicture(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	userDto := new(domain.UpdateProfileBackgroundPicture)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.UpdateProfileBackgroundPicture(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	cache.RedisCachePool.Put(rdb)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateCurrentTagline(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	userDto := new(domain.UpdateCurrentTagline)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.UpdateCurrentTagline(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	cache.RedisCachePool.Put(rdb)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateFlagCount(c *fiber.Ctx) error {
	username := c.Params("username")
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
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
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.DeleteByID(u.Id, rdb, c.Context(), u.Username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	cache.RedisCachePool.Put(rdb)
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) BlockUser(c *fiber.Ctx) error {
	username := c.Params("username")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.BlockUser(u.Id, strings.ToLower(username), rdb, c.Context(), u.Username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			cache.RedisCachePool.Put(rdb)
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		cache.RedisCachePool.Put(rdb)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	cache.RedisCachePool.Put(rdb)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UnblockUser(c *fiber.Ctx) error {
	username := c.Params("username")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)

	err = uh.UserService.UnblockUser(u.Id, strings.ToLower(username), rdb, c.Context(), u.Username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

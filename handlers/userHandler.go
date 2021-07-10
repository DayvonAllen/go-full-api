package handlers

import (
	"context"
	"example.com/app/cache"
	"example.com/app/domain"
	"example.com/app/services"
	"example.com/app/util"
	"fmt"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type UserHandler struct {
	UserService services.UserService
}

func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	page := c.Query("page", "1")

	span := opentracing.GlobalTracer().StartSpan("Get All users: GET /users")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	users, err := uh.UserService.GetAllUsers(u.Id, page, ctx, rdb, u.Username, span)

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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	users, err := uh.UserService.GetAllBlockedUsers(u.Id, rdb, c.Context(), u.Username)

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

	if !util.IsEmail(createUserDto.Email) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid email")})
	}

	if len(createUserDto.Username) <= 1 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid username")})
	}

	user := util.CreateUser(createUserDto)

	user.Following = make([]string,0, 0)
	user.Followers = make([]string,0, 0)
	user.DisplayFollowerCount = true

	err = uh.UserService.CreateUser(user)

	if err != nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

//func (uh *UserHandler) GetUserByID(c *fiber.Ctx) error {
//	token := c.Get("Authorization")
//
//	var auth domain.Authentication
//	u, loggedIn, err := auth.IsLoggedIn(token)
//
//	if err != nil || loggedIn == false {
//		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
//	}
//
//	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
//	defer cache.RedisCachePool.Put(rdb)
//
//	var data domain.UserDto
//
//	err = rdb.Get(c.Context(), util.GenerateKey(u.Username, "finduserbyusername"), &data)
//
//	if err == nil {
//		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": data})
//	}
//
//	user, err := uh.UserService.GetUserByID(u.Id, rdb, c.Context())
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
//		}
//		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
//	}
//	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
//}

//func (uh *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
//	username := c.Params("username")
//
//	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
//	defer cache.RedisCachePool.Put(rdb)
//
//	var data domain.UserDto
//
//	err := rdb.Get(c.Context(), util.GenerateKey(username, "finduserbyusername"), &data)
//
//	if err == nil {
//		fmt.Println("Found in cache in get user by username")
//		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": data})
//	}
//
//	user, err := uh.UserService.GetUserByUsername(strings.ToLower(username), rdb, c.Context())
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
//		}
//		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
//	}
//
//	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
//}

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
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UpdateProfileVisibility(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateDisplayFollowerCount(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	userDto := new(domain.UpdateDisplayFollowerCount)

	err = c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UpdateDisplayFollowerCount(u.Id, userDto, rdb)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
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
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UpdateMessageAcceptance(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}

		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
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
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UpdateCurrentBadge(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
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
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UpdateProfilePicture(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
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
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UpdateProfileBackgroundPicture(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
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
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UpdateCurrentTagline(u.Id, userDto, rdb, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
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
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) DeleteByID(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.DeleteByID(u.Id, rdb, c.Context(), u.Username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) FollowUser(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	currentUsername := c.Params("username")

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.FollowUser(strings.ToLower(currentUsername), u.Username, rdb)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UnfollowUser(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	currentUsername := c.Params("username")

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UnfollowUser(strings.ToLower(currentUsername), u.Username, rdb)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
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
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.BlockUser(u.Id, strings.ToLower(username), rdb, c.Context(), u.Username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
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
	defer cache.RedisCachePool.Put(rdb)

	err = uh.UserService.UnblockUser(u.Id, strings.ToLower(username), rdb, c.Context(), u.Username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

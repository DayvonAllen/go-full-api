package router

import (
	"example.com/app/handlers"
	"example.com/app/repo"
	"example.com/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func SetupRoutes(app *fiber.App) {
	uh := handlers.UserHandler{UserService: services.NewUserService(repo.NewUserRepoImpl())}
	ah := handlers.AuthHandler{AuthService: services.NewAuthService(repo.NewAuthRepoImpl())}

	api := app.Group("", logger.New())

	auth := api.Group("/auth")
	auth.Post("/login", ah.Login)

	user := api.Group("/users")
	user.Get("/", uh.GetAllUsers)
	user.Get("/:id", uh.GetUserByID)
	user.Post("/", uh.CreateUser)
	user.Put("/:id",  uh.UpdateUser)
	user.Delete("/:id", uh.DeleteByID)
}
package router

import (
	app2 "example.com/app/app"
	"example.com/app/repo"
	"example.com/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func SetupRoutes(app *fiber.App) {
	ch := app2.Handlers{UserService: services.NewUserService(repo.NewUserRepoImpl()),
		AuthService: services.NewAuthService(repo.NewAuthRepoImpl())}

	api := app.Group("", logger.New())

	auth := api.Group("/auth")
	auth.Post("/login", ch.Login)

	user := api.Group("/users")
	user.Get("/", ch.GetAllUsers)
	user.Get("/:id", ch.GetUserByID)
	user.Post("/", ch.CreateUser)
	user.Put("/:id",  ch.UpdateUser)
	user.Delete("/:id", ch.DeleteByID)
}
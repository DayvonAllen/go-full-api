package router

import (
	"example.com/app/handlers"
	"example.com/app/middleware"
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
	auth.Post("/reset", ah.ResetPasswordQuery)
	auth.Put("/reset/:token", ah.ResetPassword)
	auth.Get("/account/:code", ah.VerifyCode)

	user := api.Group("/users")
	user.Get("/", uh.GetAllUsers)
	user.Get("/blocked", uh.GetAllBlockedUsers)
	user.Get("/account", middleware.IsLoggedIn, uh.GetUserByID)
	user.Get("/:username", middleware.IsLoggedIn, uh.GetUserByUsername)
	user.Post("flag/:username", uh.UpdateFlagCount)
	user.Post("/", uh.CreateUser)
	user.Put("/profile-visibility", uh.UpdateProfileVisibility)
	user.Put("/message-acceptance", uh.UpdateMessageAcceptance)
	user.Put("/current-badge", uh.UpdateCurrentBadge)
	user.Put("/profile-photo", uh.UpdateProfilePicture)
	user.Put("/background-photo", uh.UpdateProfileBackgroundPicture)
	user.Put("/current-tagline", uh.UpdateCurrentTagline)
	user.Put("/block/:username", uh.BlockUser)
	user.Put("/unblock/:username", uh.UnBlockUser)
	user.Delete("/delete", uh.DeleteByID)
}
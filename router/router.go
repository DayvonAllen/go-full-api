package router

import (
	"example.com/app/handlers"
	"example.com/app/repo"
	"example.com/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)



func SetupRoutes(app *fiber.App) {
	uh := handlers.UserHandler{UserService: services.NewUserService(repo.NewUserRepoImpl())}
	ah := handlers.AuthHandler{AuthService: services.NewAuthService(repo.NewAuthRepoImpl())}
	app.Use(recover.New())
	api := app.Group("", logger.New())

	auth := api.Group("/auth")
	auth.Post("/login", ah.Login)
	auth.Post("/reset", ah.ResetPasswordQuery)
	auth.Put("/reset/:token", ah.ResetPassword)
	auth.Get("/account/:code", ah.VerifyCode)

	user := api.Group("/users")
	user.Get("/", uh.GetAllUsers)
	user.Get("/blocked", uh.GetAllBlockedUsers)
	user.Post("flag/:username", uh.UpdateFlagCount)
	user.Post("/", uh.CreateUser)
	user.Put("/profile-visibility", uh.UpdateProfileVisibility)
	user.Put("/follower-count", uh.UpdateDisplayFollowerCount)
	user.Put("/message-acceptance", uh.UpdateMessageAcceptance)
	user.Put("/current-badge", uh.UpdateCurrentBadge)
	user.Put("/profile-photo", uh.UpdateProfilePicture)
	user.Put("/background-photo", uh.UpdateProfileBackgroundPicture)
	user.Put("/current-tagline", uh.UpdateCurrentTagline)
	user.Put("/block/:username", uh.BlockUser)
	user.Put("/unblock/:username", uh.UnblockUser)
	user.Put("/follow/:username", uh.FollowUser)
	user.Put("/unfollow/:username", uh.UnfollowUser)
	user.Delete("/delete", uh.DeleteByID)
}

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		ExposeHeaders: "Authorization",
	}))

	SetupRoutes(app)

	return app
}
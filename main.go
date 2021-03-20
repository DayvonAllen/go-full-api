package main

import (
	"example.com/app/database"
	"example.com/app/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func init() {
	// create database connection instance for first time
	_ = database.GetInstance()
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		ExposeHeaders: "Authorization",
	}))

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8080"))
}

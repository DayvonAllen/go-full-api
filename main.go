package main

import (
	"example.com/app/database"
	"example.com/app/router"
	"log"
)

func init() {
	// create database connection instance for first time
	_ = database.GetInstance()
}

func main() {
	app := router.Setup()
	log.Fatal(app.Listen(":8080"))
}



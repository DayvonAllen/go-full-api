package main

import (
	"example.com/app/database"
	"example.com/app/router"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func init() {
	// create database connection instance
	_ = database.GetInstance()
}

func main() {
	app := router.Setup()

	// graceful shutdown on signal interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <- c
		fmt.Println("Shutting down...")
		database.CloseConnection()
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}



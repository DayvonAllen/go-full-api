package database

import (
	"context"
	"example.com/app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Connection struct {
	*mongo.Client
	UserCollection *mongo.Collection
	FlagCollection *mongo.Collection
	*mongo.Database
}

func ConnectToDB(ctx context.Context) (*Connection,error) {
	p := config.Config("DB_PORT")
	n := config.Config("DB_NAME")
	h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(n + h + p))
	if err != nil { return nil, err }

	// create database
	db := client.Database("user-service")

	// create collection
	userCollection := db.Collection("users")
	flagCollection := db.Collection("flags")

	dbConnection := &Connection{client, userCollection, flagCollection, db}

	return dbConnection, nil
}
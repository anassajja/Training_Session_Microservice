package db

import (
	"context"
	"log"
	"training_session/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client // Client

func Connect() { // Connect to MongoDB
	cfg := config.LoadConfig() // Load configuration

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Connected to MongoDB!")
}

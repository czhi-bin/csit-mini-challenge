package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	godotenv.Load(".env")

	mongo_URI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_URI))

	if err != nil {
		log.Fatal(err)
	}

	return client
}

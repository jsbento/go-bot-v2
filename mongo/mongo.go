package mongo

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New() *MongoDBClient {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Could not load .env file")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("MONGODB_URI is not set")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Error creating MongoDB client,", err)
		return nil
	}
	return &MongoDBClient{Client: client}
}

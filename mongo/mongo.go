package mongo

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New() *MongoDBClient {
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

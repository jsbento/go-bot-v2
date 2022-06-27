package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBClient struct {
	Client *mongo.Client
}

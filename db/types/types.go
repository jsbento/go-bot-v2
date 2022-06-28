package types

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	uT "github.com/jsbento/go-bot-v2/users/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	Client *mongo.Client
}

func New() (*DBClient, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("could not load .env file")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, errors.New("MONGODB_URI is not set")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error creating mongodb client, %s", err)
	}
	return &DBClient{
		Client: client,
	}, nil
}

func (dbC *DBClient) PostUser(user *uT.User) (flag int, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	var result uT.User
	err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}, options.FindOne()).Decode(&result)
	if result.Username == "" {
		_, err = collection.InsertOne(context.TODO(), user, options.InsertOne())
		if err != nil {
			return -1, err
		}
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return 1, errors.New("user already exists")
}

func (dbC *DBClient) DeleteUser(username string) (flag int, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	err = collection.FindOneAndDelete(context.TODO(), bson.M{"username": username}, options.FindOneAndDelete()).Decode(&uT.User{})
	if err == mongo.ErrNoDocuments {
		return 1, errors.New("user does not exist")
	} else if err != nil {
		return -1, err
	}
	return 1, nil
}

func (dbC *DBClient) GetUser(username string) (result uT.User, err error) {
	var user uT.User
	collection := dbC.Client.Database("discord-users").Collection("users")
	err = collection.FindOne(context.TODO(), bson.M{"username": username}, options.FindOne()).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return uT.User{}, errors.New("user does not exist")
	} else if err != nil {
		return uT.User{}, err
	}
	return user, nil
}

func (dbC *DBClient) AddTokens(username string, tokens int) (user *uT.User, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	var result uT.User
	err = collection.FindOneAndUpdate(context.TODO(), bson.M{"username": username}, bson.M{"$inc": bson.M{"token_count": tokens}}, options.FindOneAndUpdate()).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (dbC *DBClient) PurchasePowerUp(username string, pId int) (user *uT.User, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	var result uT.User
	err = collection.FindOne(context.TODO(), bson.M{"username": username}, options.FindOne()).Decode(&result)
	if err != nil {
		return nil, err
	}
	powerup := result.PowerUps[pId-1]
	if !powerup.Active && powerup.Value <= result.TokenCount {
		result.TokenCount -= powerup.Value
		powerup.Active = true
		_, err = collection.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{"$set": bson.M{"token_count": result.TokenCount, "power_ups": result.PowerUps}})
		if err != nil {
			return nil, err
		}
		return &result, nil
	} else if powerup.Active {
		return &result, errors.New("Powerup already active")
	} else if powerup.Value > result.TokenCount {
		return &result, errors.New("Not enough tokens")
	}
	return &result, nil
}

func (dbC *DBClient) Close() error {
	return dbC.Client.Disconnect(context.TODO())
}

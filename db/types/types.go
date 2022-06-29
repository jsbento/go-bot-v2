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

const (
	DB_ERROR        = -1
	USER_OK         = 0
	USER_EXISTS     = 1
	USER_NOT_EXISTS = 2
	LOW_TOKENS      = 3
	POWER_ACTIVE    = 4
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
	if err == mongo.ErrNoDocuments || result.Username == "" {
		_, err = collection.InsertOne(context.TODO(), user, options.InsertOne())
		if err == nil {
			return USER_OK, nil
		} else {
			return DB_ERROR, err
		}
	} else if err != nil {
		return DB_ERROR, err
	}
	return USER_EXISTS, errors.New("You already have a saved user!")
}

func (dbC *DBClient) DeleteUser(username string) (flag int, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	err = collection.FindOneAndDelete(context.TODO(), bson.M{"username": username}, options.FindOneAndDelete()).Decode(&uT.User{})
	if err == mongo.ErrNoDocuments {
		return USER_NOT_EXISTS, errors.New("You don't have a saved user!")
	} else if err != nil {
		return DB_ERROR, err
	}
	return USER_OK, nil
}

func (dbC *DBClient) GetUser(username string) (result uT.User, flag int, err error) {
	var user uT.User
	collection := dbC.Client.Database("discord-users").Collection("users")
	err = collection.FindOne(context.TODO(), bson.M{"username": username}, options.FindOne()).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return uT.User{}, USER_NOT_EXISTS, errors.New("You don't have a saved user!")
	} else if err != nil {
		return uT.User{}, DB_ERROR, err
	}
	return user, USER_OK, nil
}

func (dbC *DBClient) AddTokens(username string, tokens int) (user *uT.User, flag int, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	var result uT.User
	err = collection.FindOneAndUpdate(context.TODO(), bson.M{"username": username}, bson.M{"$inc": bson.M{"token_count": tokens}}, options.FindOneAndUpdate()).Decode(&result)
	if err != nil {
		return nil, DB_ERROR, err
	}
	return &result, USER_OK, nil
}

func (dbC *DBClient) PurchasePowerUp(username string, pId int) (user *uT.User, flag int, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	var result uT.User
	err = collection.FindOne(context.TODO(), bson.M{"username": username}, options.FindOne()).Decode(&result)
	if err != nil {
		return nil, DB_ERROR, err
	}
	powerup := result.PowerUps[pId-1]
	if !powerup.Active && powerup.Value <= result.TokenCount {
		result.TokenCount -= powerup.Value
		powerup.Active = true
		if pId == 1 {
			powerup.Uses = 5
		}
		_, err = collection.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{"$set": bson.M{"token_count": result.TokenCount, "power_ups": result.PowerUps}})
		if err != nil {
			return nil, DB_ERROR, err
		}
		return &result, USER_OK, nil
	} else if powerup.Active {
		return &result, POWER_ACTIVE, errors.New("Powerup already active!")
	} else if powerup.Value > result.TokenCount {
		return &result, LOW_TOKENS, errors.New("Not enough tokens!")
	}
	return &result, USER_OK, nil
}

func (dbC *DBClient) UpdateUser(username string, powerups []*uT.PowerUp, tokens int) (flag int, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{"$set": bson.M{"power_ups": powerups}, "$inc": bson.M{"token_count": tokens}})
	if err != nil {
		return DB_ERROR, err
	}
	return USER_OK, nil
}

func (dbC *DBClient) UpdatePowerUps(username string, powerups []*uT.PowerUp) (flag int, err error) {
	collection := dbC.Client.Database("discord-users").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{"$set": bson.M{"power_ups": powerups}})
	if err != nil {
		return DB_ERROR, err
	}
	return USER_OK, nil
}

func (dbC *DBClient) Close() error {
	return dbC.Client.Disconnect(context.TODO())
}

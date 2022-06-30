package services

import (
	"context"
	"time"

	"github.com/mattchw/go-onboard/configs"
	"github.com/mattchw/go-onboard/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// create user service
func CreateUser(payload models.User) *mongo.InsertOneResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Bio:       payload.Bio,
		Age:       payload.Age,
		Gender:    payload.Gender,
	}

	if err := newUser.ValidateUser(); err != nil {
		panic(err)
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		panic(err)
	}

	return result
}

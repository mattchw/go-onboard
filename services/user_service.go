package services

import (
	"context"

	"github.com/mattchw/go-onboard/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// define User Service interface
type UserService interface {
	Create(ctx context.Context, payload models.User) *mongo.InsertOneResult
	FindOne(ctx context.Context, payload models.User) models.User
	DeleteOne(ctx context.Context, payload models.User) bool
}

// implement userService
type UserServiceImpl struct {
	collection *mongo.Collection
}

// Constructor
func NewUserServiceImpl(coll *mongo.Collection) *UserServiceImpl {
	return &UserServiceImpl{
		collection: coll,
	}
}

// implement Create
func (us *UserServiceImpl) Create(ctx context.Context, payload models.User) *mongo.InsertOneResult {
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

	result, err := us.collection.InsertOne(ctx, newUser)
	if err != nil {
		panic(err)
	}

	return result
}

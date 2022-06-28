package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/mattchw/go-onboard/configs"
	"github.com/mattchw/go-onboard/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	// find all users
	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error getting user",
			})
	}

	// decode all users
	defer results.Close(ctx)
	for results.Next(ctx) {
		var user models.User
		if err = results.Decode(&user); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Error decoding user",
				})
		}

		users = append(users, user)
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"status":  "success",
			"message": "User retrieved successfully",
			"items":   users,
		})
}

func GetUsersCount(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error getting user count",
			})
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"status":  "success",
			"message": "User count retrieved successfully",
			"data":    count,
		})
}

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
	}

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Bio:       user.Bio,
		Age:       user.Age,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error creating user",
			})
	}

	return c.Status(http.StatusCreated).JSON(
		fiber.Map{
			"status":  "success",
			"message": "User created successfully",
			"data":    result,
		})
}

func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error getting user",
			})
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"status":  "success",
			"message": "User retrieved successfully",
			"data":    user,
		})
}

func UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{
		"$set": bson.M{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"bio":       user.Bio,
			"age":       user.Age,
		},
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error updating user",
				"error":   err,
			})
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"status":  "success",
			"message": "User updated successfully",
			"data":    result,
		})
}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error getting user",
			})
	}

	// delete user
	result, err := userCollection.DeleteOne(ctx, user.Id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error deleting user",
			})
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"status":  "success",
			"message": "User deleted successfully",
			"data":    result,
		})
}

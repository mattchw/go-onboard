package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mattchw/go-onboard/configs"
	"github.com/mattchw/go-onboard/models"
	"github.com/mattchw/go-onboard/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func CacheFetch(ctx context.Context, key string, ttl time.Duration, result interface{}, code func() interface{}) {
	str, _ := configs.RDB.Get(ctx, key).Result()
	if str == "" {
		fmt.Println("---> cache miss")
		value := code()
		jsonStr, _ := json.Marshal(value)
		str = string(jsonStr)
		configs.RDB.Set(ctx, key, str, ttl).Result()
	} else {
		fmt.Println("---> cache exists")
	}

	json.Unmarshal([]byte(str), &result)
}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	CacheFetch(ctx, "users", 10*time.Second, &users, func() interface{} {
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

		users_json, _ := json.Marshal(users)

		// set redis cache for 30 seconds
		err = configs.RDB.Set(ctx, "users", users_json, 0).Err()
		if err != nil {
			panic(err)
		}

		return users
	})

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
	var user models.User

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
				"error":   err,
			})
	}

	result := services.CreateUser(user)

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

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/joho/godotenv"
	"github.com/mattchw/go-onboard/routes"
)

// Main function
func main() {
	fmt.Println("!... Hello World ...!")

	// load env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		AppName: "Go onboard v1.0.0",
	})

	// recover from any panics
	app.Use(recover.New())
	// default timeout
	app.Use(timeout.New(func(c *fiber.Ctx) (err error) { return c.Next() }, 1*time.Second))
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		time.Sleep(5 * time.Second)
		return c.SendString("OK!!!")
	})
	// routes
	routes.UserRoute(app)

	port := os.Getenv("PORT")
	app.Listen(":" + port)
}

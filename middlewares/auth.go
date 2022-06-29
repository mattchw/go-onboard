package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// AuthReq middleware
func AuthReq() func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "12345678",
		},
	})
}

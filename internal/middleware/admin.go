package middleware

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v3"
)

// AdminAuth creates an admin authentication middleware.
// Uses constant-time comparison to prevent timing attacks.
func AdminAuth(adminSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		
		// Check if API key is provided
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing API key",
			})
		}
		
		// Use constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(adminSecret)) != 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}
		
		return c.Next()
	}
}

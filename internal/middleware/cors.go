package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CORS returns a configured CORS middleware.
func CORS(allowOrigin string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     []string{allowOrigin, "http://localhost:3000", "http://localhost:8080"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           3600,
	})
}

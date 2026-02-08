package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CORS returns a configured CORS middleware.
func CORS(allowOrigin string) fiber.Handler {
	//Allow all origins for development/preview
	// In production, you can restrict this to specific domains
	return cors.New(cors.Config{
		AllowOrigins:     "*", // Allow all origins
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-API-Key",
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowCredentials: false,
		MaxAge:           3600,
	})
}

package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CORS returns a configured CORS middleware.
func CORS(allowOrigin string) fiber.Handler {
	// Validate and normalize the origin
	origins := []string{"http://localhost:3000", "http://localhost:8080"}
	
	// Only add the custom origin if it's valid (has http:// or https://)
	if allowOrigin != "" && (strings.HasPrefix(allowOrigin, "http://") || strings.HasPrefix(allowOrigin, "https://")) {
		origins = append([]string{allowOrigin}, origins...)
	}
	
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           3600,
	})
}

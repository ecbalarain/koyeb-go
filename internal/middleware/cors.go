package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// CORS returns a configured CORS middleware.
func CORS(allowOrigin string) fiber.Handler {
	// Validate and normalize the origin
	origins := []string{
		"http://localhost:3000", 
		"http://localhost:8080",
		"https://bhomanshah.com",  // Production frontend domain
	}
	
	// Only add the custom origin if it's valid (has http:// or https://)
	if allowOrigin != "" && (strings.HasPrefix(allowOrigin, "http://") || strings.HasPrefix(allowOrigin, "https://")) {
		// Check if it's not already in the list
		found := false
		for _, origin := range origins {
			if origin == allowOrigin {
				found = true
				break
			}
		}
		if !found {
			origins = append([]string{allowOrigin}, origins...)
		}
	}
	
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           3600,
	})
}

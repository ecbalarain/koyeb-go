package middleware

import (
	"crypto/subtle"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// AdminClaims represents the JWT claims for admin authentication.
type AdminClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for an authenticated admin.
func GenerateJWT(jwtSecret string, expiryHours int) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(expiryHours) * time.Hour)

	claims := AdminClaims{
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "oxlook-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// AdminAuth creates an admin authentication middleware.
// Supports both JWT (Authorization: Bearer <token>) and legacy API key (X-API-Key) for backward compatibility.
func AdminAuth(adminSecret string, jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Try JWT first (Authorization: Bearer <token>)
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or expired token",
				})
			}

			claims, ok := token.Claims.(*AdminClaims)
			if !ok || claims.Role != "admin" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token claims",
				})
			}

			return c.Next()
		}

		// Fallback: legacy API key (X-API-Key header) — for backward compatibility
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication. Provide Authorization: Bearer <token> or X-API-Key header.",
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

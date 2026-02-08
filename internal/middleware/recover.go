package middleware

import (
	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
)

// Recover returns a middleware that recovers from panics.
func Recover() fiber.Handler {
	return recoverer.New()
}

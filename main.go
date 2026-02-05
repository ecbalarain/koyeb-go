package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	app := fiber.New()

	app.Get("/", HelloHandler)

	log.Println("Listening on port", port)
	log.Fatal(app.Listen(":" + port))
}

func HelloHandler(c fiber.Ctx) error {
	return c.SendString("Hello from Koyeb\n")
}

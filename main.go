package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/koyeb/example-golang/internal/config"
	"github.com/koyeb/example-golang/internal/database"
	"github.com/koyeb/example-golang/internal/handlers"
	"github.com/koyeb/example-golang/internal/middleware"
)

func main() {
	// Load configuration from .env / environment
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	h := handlers.NewHandler(db, cfg.AdminSecret, cfg.JWTSecret, cfg.JWTExpiry, cfg.EmailHost, cfg.EmailPort, cfg.EmailUser, cfg.EmailPass, cfg.EmailFrom)

	app := fiber.New(fiber.Config{
		AppName:        "OXLOOK API",
		BodyLimit:      2 * 1024 * 1024, // 2 MB max request body size
		ReadBufferSize: 4096,            // 4 KB read buffer
	})

	// Log configuration for debugging
	log.Printf("CORS Origin: %s", cfg.CORSOrigin)
	log.Printf("Environment: %s", cfg.Environment)

	// Middleware
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.SecurityHeaders()) // Add security headers (HSTS, X-Frame-Options, etc.)
	app.Use(middleware.GlobalRateLimit()) // Global rate limiting: 100 req/min per IP
	app.Use(middleware.CORS(cfg.CORSOrigin))

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Serve Tailwind CSS locally to avoid COEP issues
	app.Get("/tw.js", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/tailwind.js")
	})

	// Serve favicon
	app.Get("/favicon.ico", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/favicon.ico")
	})

	// API routes
	api := app.Group("/api")

	// GET /api/products — list active products
	api.Get("/products", h.GetProducts)

	// GET /api/products/:slug/variants — return variants for a product
	api.Get("/products/:slug/variants", h.GetProductVariants)

	// POST /api/orders — create COD order (with strict rate limiting)
	// 5 orders per minute per IP to prevent abuse
	api.Post("/orders", middleware.StrictRateLimit(5, 1*time.Minute), h.CreateOrder)

	// Static admin pages (before admin group to avoid auth middleware)
	app.Get("/admin/login", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/admin/login.html")
	})
	app.Get("/admin", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/admin/index.html")
	})
	app.Get("/admin/products", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/admin/products.html")
	})
	app.Get("/admin/orders", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/admin/orders.html")
	})
	app.Get("/admin/products/:id/variants", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/admin/variants.html")
	})

	// Static frontend files
	app.Get("/products.json", func(c fiber.Ctx) error {
		return c.SendFile("./cloudflare-pages-frontend/products.json")
	})

	// Admin login (public, no auth required) — returns JWT token
	app.Post("/admin/api/login", middleware.StrictRateLimit(10, 1*time.Minute), h.AdminLogin)

	// Admin API routes (protected by JWT or legacy API key)
	admin := app.Group("/admin/api", middleware.AdminAuth(cfg.AdminSecret, cfg.JWTSecret))

	// GET /admin/products — list all products (including inactive)
	admin.Get("/products", h.AdminGetProducts)

	// POST /admin/products — create a new product
	admin.Post("/products", h.AdminCreateProduct)

	// PUT /admin/products/:id — update product details
	admin.Put("/products/:id", h.AdminUpdateProduct)

	// PATCH /admin/products/:id — toggle product active status
	admin.Patch("/products/:id", h.AdminToggleProductActive)

	// GET /admin/products/:id/variants — list all variants for a product
	admin.Get("/products/:id/variants", h.AdminGetProductVariants)

	// POST /admin/products/:id/variants — create a new variant
	admin.Post("/products/:id/variants", h.AdminCreateVariant)

	// PATCH /admin/variants/:id — update variant (price, stock, active)
	admin.Patch("/variants/:id", h.AdminUpdateVariant)

	// DELETE /admin/variants/:id — delete a variant
	admin.Delete("/variants/:id", h.AdminDeleteVariant)

	// GET /admin/orders — list all orders (with filters: status, date)
	admin.Get("/orders", h.AdminGetOrders)

	// GET /admin/orders/:id — get order details with items
	admin.Get("/orders/:id", h.AdminGetOrder)

	// PATCH /admin/orders/:id — update order status
	admin.Patch("/orders/:id", h.AdminUpdateOrderStatus)

	// POST /admin/cache/purge — purge cache entries
	admin.Post("/cache/purge", h.AdminPurgeCache)

	log.Printf("Starting OXLOOK API on :%s (env: %s)", cfg.Port, cfg.Environment)
	log.Fatal(app.Listen(":" + cfg.Port))
}

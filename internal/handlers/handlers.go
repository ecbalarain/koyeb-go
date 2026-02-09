package handlers

import (
	"crypto/subtle"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/koyeb/example-golang/internal/middleware"
	"github.com/koyeb/example-golang/internal/models"
	"github.com/koyeb/example-golang/internal/repository"
	"github.com/koyeb/example-golang/internal/services"
	"github.com/koyeb/example-golang/internal/validator"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	productRepo  *repository.ProductRepository
	variantRepo  *repository.VariantRepository
	orderRepo    *repository.OrderRepository
	cacheManager *CacheManager
	adminSecret  string
	jwtSecret    string
	jwtExpiry    int
	emailService *services.EmailService
}

// NewHandler creates a new handler instance.
func NewHandler(db *sql.DB, adminSecret, jwtSecret string, jwtExpiry int, apiKey, from string) *Handler {
	return &Handler{
		productRepo:  repository.NewProductRepository(db),
		variantRepo:  repository.NewVariantRepository(db),
		orderRepo:    repository.NewOrderRepository(db),
		cacheManager: NewCacheManager(),
		adminSecret:  adminSecret,
		jwtSecret:    jwtSecret,
		jwtExpiry:    jwtExpiry,
		emailService: services.NewEmailService(apiKey, from),
	}
}

// AdminLogin authenticates admin with API key and returns a JWT token.
// POST /admin/api/login
func (h *Handler) AdminLogin(c fiber.Ctx) error {
	var req struct {
		APIKey string `json:"api_key"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.APIKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "API key is required",
		})
	}

	// Validate API key using constant-time comparison
	if subtle.ConstantTimeCompare([]byte(req.APIKey), []byte(h.adminSecret)) != 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid API key",
		})
	}

	// Generate JWT token
	token, expiresAt, err := middleware.GenerateJWT(h.jwtSecret, h.jwtExpiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token":      token,
		"expires_at": expiresAt,
		"message":    "Login successful",
	})
}

// GetProducts returns a list of active products.
// GET /api/products
func (h *Handler) GetProducts(c fiber.Ctx) error {
	products, err := h.productRepo.GetAll(true) // active only
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	// Return only the fields needed for the product list
	type ProductResponse struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Category string `json:"category"`
		Images   string `json:"images"`
	}

	var response []ProductResponse
	for _, p := range products {
		response = append(response, ProductResponse{
			ID:       p.ID,
			Name:     p.Name,
			Slug:     p.Slug,
			Category: p.Category,
			Images:   p.Images,
		})
	}

	return c.JSON(response)
}

// GetProductVariants returns variants for a specific product.
// GET /api/products/:slug/variants
func (h *Handler) GetProductVariants(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product slug is required",
		})
	}

	// Get product by slug
	product, err := h.productRepo.GetBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product",
		})
	}

	if !product.Active {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Get active variants for this product
	variants, err := h.variantRepo.GetByProductID(product.ID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch variants",
		})
	}

	// Set cache headers (30 minutes for variants)
	// This can be increased up to 60 minutes for better performance
	// Cache time is a balance between freshness and performance
	c.Set("Cache-Control", "public, max-age=1800, s-maxage=1800")
	c.Set("Vary", "Accept-Encoding") // Ensure proper caching with compression

	return c.JSON(variants)
}

// CreateOrder creates a new COD order.
// POST /api/orders
func (h *Handler) CreateOrder(c fiber.Ctx) error {
	// Set no-cache headers to prevent caching of order submissions
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	var req models.CreateOrderRequest

	// Parse request body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Sanitize customer inputs to prevent XSS attacks
	req.CustomerName = validator.Sanitize(req.CustomerName)
	req.CustomerPhone = validator.Sanitize(req.CustomerPhone)
	req.CustomerAddress = validator.Sanitize(req.CustomerAddress)
	req.CustomerCity = validator.Sanitize(req.CustomerCity)
	req.CustomerEmail = validator.Sanitize(req.CustomerEmail)
	req.Notes = validator.Sanitize(req.Notes)

	// Validate sanitized inputs are not empty (in case input was all malicious characters)
	if req.CustomerName == "" || req.CustomerPhone == "" || req.CustomerAddress == "" || req.CustomerCity == "" || req.CustomerEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer information provided",
		})
	}

	// Verify all variants exist, are active, and have sufficient stock
	var orderItems []models.OrderItem
	total := int64(0)

	for _, item := range req.Items {
		variant, err := h.variantRepo.GetByID(item.VariantID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Variant %d not found", item.VariantID),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify variant",
			})
		}

		if !variant.Active {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Variant %d is not available", item.VariantID),
			})
		}

		if variant.Stock < item.Qty {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Insufficient stock for variant %d (requested: %d, available: %d)",
					item.VariantID, item.Qty, variant.Stock),
			})
		}

		// Get product name for the order item
		product, err := h.productRepo.GetByID(variant.ProductID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch product details",
			})
		}

		// Create order item
		orderItem := models.OrderItem{
			VariantID:       item.VariantID,
			ProductName:     product.Name,
			VariantLabel:    fmt.Sprintf("%s / %s", variant.Color, variant.Size),
			PriceAtPurchase: variant.Price,
			Qty:             item.Qty,
		}
		orderItems = append(orderItems, orderItem)

		// Add to total
		total += variant.Price * int64(item.Qty)
	}

	// Create order
	order := &models.Order{
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		CustomerAddress: req.CustomerAddress,
		CustomerCity:    req.CustomerCity,
		CustomerEmail:   req.CustomerEmail,
		Total:           total,
		Status:          "new",
		Notes:           req.Notes,
	}

	// Create order and decrement stock atomically
	err := h.orderRepo.Create(order, orderItems)
	if err != nil {
		// Check if it's a stock error (shouldn't happen due to our checks, but just in case)
		if strings.Contains(err.Error(), "stock") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Stock changed during order processing. Please try again.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	// Send confirmation email (non-blocking)
	_ = h.emailService.SendOrderConfirmation(req.CustomerEmail, req.CustomerName, order.ID, order.Total)

	// Return order confirmation
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order_id": order.ID,
		"total":    order.Total,
		"status":   order.Status,
		"message":  "Order created successfully",
	})
}

// AdminGetProducts returns all products including inactive ones.
// GET /admin/products
func (h *Handler) AdminGetProducts(c fiber.Ctx) error {
	products, err := h.productRepo.GetAll(false) // include inactive
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(products)
}

// AdminToggleProductActive toggles a product's active status.
// PATCH /admin/products/:id
func (h *Handler) AdminToggleProductActive(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Get current product to toggle status
	product, err := h.productRepo.GetByID(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product",
		})
	}

	// Toggle active status
	err = h.productRepo.SetActive(int64(id), !product.Active)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product updated successfully",
		"active":  !product.Active,
	})
}

// AdminCreateProduct creates a new product.
// POST /admin/api/products
func (h *Handler) AdminCreateProduct(c fiber.Ctx) error {
	var req struct {
		Name        string `json:"name" validate:"required,min=1,max=255"`
		Slug        string `json:"slug" validate:"required,min=1,max=255"`
		Description string `json:"description"`
		Category    string `json:"category" validate:"required,min=1,max=100"`
		Images      string `json:"images"` // JSON array as string
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Sanitize inputs to prevent XSS
	req.Name = validator.Sanitize(req.Name)
	req.Slug = validator.Sanitize(req.Slug)
	req.Description = validator.Sanitize(req.Description)
	req.Category = validator.Sanitize(req.Category)

	// Check if slug already exists
	existing, err := h.productRepo.GetBySlug(req.Slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check product slug",
		})
	}
	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Product with this slug already exists",
		})
	}

	// Default images to empty JSON array if not provided
	if req.Images == "" {
		req.Images = "[]"
	}

	// Create product
	id, err := h.productRepo.Create(req.Name, req.Slug, req.Description, req.Category, req.Images)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"id":      id,
	})
}

// AdminUpdateProduct updates an existing product.
// PUT /admin/api/products/:id
func (h *Handler) AdminUpdateProduct(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req struct {
		Name        string `json:"name" validate:"required,min=1,max=255"`
		Slug        string `json:"slug" validate:"required,min=1,max=255"`
		Description string `json:"description"`
		Category    string `json:"category" validate:"required,min=1,max=100"`
		Images      string `json:"images"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Sanitize inputs to prevent XSS
	req.Name = validator.Sanitize(req.Name)
	req.Slug = validator.Sanitize(req.Slug)
	req.Description = validator.Sanitize(req.Description)
	req.Category = validator.Sanitize(req.Category)

	// Check if product exists
	product, err := h.productRepo.GetByID(int64(id))
	if err != nil || product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Check if slug is being changed and if it conflicts with another product
	if req.Slug != product.Slug {
		existing, err := h.productRepo.GetBySlug(req.Slug)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check product slug",
			})
		}
		if existing != nil && existing.ID != int64(id) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Product with this slug already exists",
			})
		}
	}

	// Default images to empty JSON array if not provided
	if req.Images == "" {
		req.Images = "[]"
	}

	// Update product
	err = h.productRepo.Update(int64(id), req.Name, req.Slug, req.Description, req.Category, req.Images)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}

// AdminGetProductVariants returns all variants for a product (including inactive).
// GET /admin/products/:id/variants
func (h *Handler) AdminGetProductVariants(c fiber.Ctx) error {
	idStr := c.Params("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	variants, err := h.variantRepo.GetByProductID(int64(productID), false) // include inactive
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch variants",
		})
	}

	return c.JSON(variants)
}

// AdminCreateVariant creates a new variant for a product.
// POST /admin/products/:id/variants
func (h *Handler) AdminCreateVariant(c fiber.Ctx) error {
	idStr := c.Params("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req struct {
		Size  string `json:"size" validate:"required"`
		Color string `json:"color" validate:"required"`
		Price int64  `json:"price" validate:"required,gt=0"`
		Stock int    `json:"stock" validate:"required,gte=0"`
		Active bool  `json:"active"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Sanitize inputs to prevent XSS
	req.Size = validator.Sanitize(req.Size)
	req.Color = validator.Sanitize(req.Color)

	variant := &models.Variant{
		ProductID: int64(productID),
		Size:      req.Size,
		Color:     req.Color,
		Price:     req.Price,
		Stock:     req.Stock,
		Active:    req.Active,
	}

	err = h.variantRepo.Create(variant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create variant",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(variant)
}

// AdminUpdateVariant updates a variant's details.
// PATCH /admin/variants/:id
func (h *Handler) AdminUpdateVariant(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid variant ID",
		})
	}

	var req struct {
		Size   string `json:"size,omitempty"`
		Color  string `json:"color,omitempty"`
		Price  *int64 `json:"price,omitempty"`
		Stock  *int   `json:"stock,omitempty"`
		Active *bool  `json:"active,omitempty"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get current variant
	variant, err := h.variantRepo.GetByID(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Variant not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch variant",
		})
	}

	// Update fields if provided (with sanitization)
	if req.Size != "" {
		variant.Size = validator.Sanitize(req.Size)
	}
	if req.Color != "" {
		variant.Color = validator.Sanitize(req.Color)
	}
	if req.Price != nil {
		variant.Price = *req.Price
	}
	if req.Stock != nil {
		variant.Stock = *req.Stock
	}
	if req.Active != nil {
		variant.Active = *req.Active
	}

	err = h.variantRepo.Update(variant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update variant",
		})
	}

	return c.JSON(variant)
}

// AdminDeleteVariant deletes a variant.
// DELETE /admin/variants/:id
func (h *Handler) AdminDeleteVariant(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid variant ID",
		})
	}

	err = h.variantRepo.Delete(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete variant",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Variant deleted successfully",
	})
}

// AdminGetOrders returns all orders with optional filters.
// GET /admin/orders
func (h *Handler) AdminGetOrders(c fiber.Ctx) error {
	status := c.Query("status")
	// Note: date filtering could be added later if needed

	orders, err := h.orderRepo.GetAll(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(orders)
}

// AdminGetOrder returns a single order with its items.
// GET /admin/orders/:id
func (h *Handler) AdminGetOrder(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := h.orderRepo.GetByID(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch order",
		})
	}

	// Get order items
	items, err := h.orderRepo.GetItemsByOrderID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch order items",
		})
	}

	return c.JSON(fiber.Map{
		"order": order,
		"items": items,
	})
}

// AdminUpdateOrderStatus updates an order's status.
// PATCH /admin/orders/:id
func (h *Handler) AdminUpdateOrderStatus(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=new confirmed shipped canceled"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.orderRepo.UpdateStatus(int64(id), req.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
		"status":  req.Status,
	})
}

// AdminPurgeCache purges all cache entries.
// POST /admin/api/cache/purge
func (h *Handler) AdminPurgeCache(c fiber.Ctx) error {
	var req struct {
		Type string `json:"type"` // "all", "variants", or specific slug
		Slug string `json:"slug,omitempty"`
	}

	// Parse request body (optional)
	if err := c.Bind().Body(&req); err != nil {
		// If no body provided, default to purging all
		req.Type = "all"
	}

	switch req.Type {
	case "all":
		h.cacheManager.InvalidateAll()
		return c.JSON(fiber.Map{
			"message": "All cache entries invalidated successfully",
		})
	case "variants":
		if req.Slug != "" {
			h.cacheManager.InvalidateVariantCache(req.Slug)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("Variant cache for product '%s' invalidated successfully", req.Slug),
			})
		}
		h.cacheManager.InvalidateAll()
		return c.JSON(fiber.Map{
			"message": "All variant caches invalidated successfully",
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cache type. Use 'all' or 'variants'",
		})
	}
}

// SendTestEmail sends a test order confirmation email to a specified address.
// POST /api/test-email
func (h *Handler) SendTestEmail(c fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		req.Email = "miqbal@sis.edu.eg" // Default test email
	}

	// Validate email format
	if !validator.IsValidEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	// Send test email
	err := h.emailService.SendOrderConfirmation(req.Email, "Test Customer", 12345, 1199)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send test email",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Test email sent successfully",
		"email":   req.Email,
		"order_id": 12345,
		"total": 1199,
	})
}

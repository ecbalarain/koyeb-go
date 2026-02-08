# OXLOOK - Implementation Plan

## Phase 1: Project Foundation & Configuration

- [x] Add Go dependencies: `pgx/v5`, `godotenv`, `validator/v10`
- [x] Create `.env.example` with all required environment variables
- [x] Create `.env` with local dev values (gitignored)
- [x] Add `.env` to `.gitignore`
- [x] Set up structured project layout (`/internal`, `/cmd`, `/migrations`)
- [x] Add CORS middleware to Fiber (allow `bhomanshah.com`)
- [x] Add request logging middleware
- [x] Add error recovery middleware
- [x] Load env vars in `main.go` with `godotenv`

## Phase 2: Database Schema & Connection

- [x] Set up MySQL connection pool in Go using `go-sql-driver/mysql`
- [x] Create migration: `products` table (id, name, slug, description, category, images, active, created_at, updated_at)
- [x] Create migration: `variants` table (id, product_id FK, size, color, price, stock, active, created_at, updated_at)
- [x] Create migration: `orders` table (id, customer_name, customer_phone, customer_address, customer_city, total, status, notes, created_at, updated_at)
- [x] Create migration: `order_items` table (id, order_id FK, variant_id, product_name, variant_label, price_at_purchase, qty)
- [x] Add unique constraint on variants (product_id, size, color)
- [x] Add index on products(slug)
- [x] Add index on orders(status)
- [x] Seed database with demo products + variants
- [x] Test database connection and migrations locally

## Phase 3: Go Data Models & Repository Layer

- [x] Define `Product` struct with JSON tags
- [x] Define `Variant` struct with JSON tags
- [x] Define `Order` struct with JSON tags
- [x] Define `OrderItem` struct with JSON tags
- [x] Define `CreateOrderRequest` struct with validation tags
- [x] Create `ProductRepository` (GetAll, GetByID, GetBySlug, SetActive)
- [x] Create `VariantRepository` (GetByProductID, Create, Update, Delete, DecrementStock)
- [x] Create `OrderRepository` (Create, GetAll, GetByID, UpdateStatus)
- [x] Write input validation helpers

## Phase 4: Public API Endpoints

- [x] `GET /api/products` — list active products (name, slug, category, image)
- [x] `GET /api/products/:slug/variants` — return variants with price + stock for a product
- [x] `POST /api/orders` — create COD order
  - [x] Validate request body (customer info, items)
  - [x] Verify each variant exists and is active
  - [x] Verify stock is sufficient for each item
  - [x] Calculate total server-side (never trust browser price)
  - [x] Decrement stock atomically
  - [x] Store order with item snapshots
  - [x] Return order confirmation
- [x] Add rate limiting on `POST /api/orders`
- [x] Add proper error response format (`{ "error": "message" }`)
- [x] Add cache headers to variant endpoint (`Cache-Control: public, max-age=600`)

## Phase 5: Admin API Endpoints

- [x] Add admin auth middleware (API key or basic auth)
- [x] `GET /admin/products` — list all products (including inactive)
- [x] `PATCH /admin/products/:id` — toggle product active status
- [x] `GET /admin/products/:id/variants` — list all variants for a product
- [x] `POST /admin/products/:id/variants` — create a new variant
- [x] `PATCH /admin/variants/:id` — update variant (price, stock, active)
- [x] `DELETE /admin/variants/:id` — delete a variant
- [x] `GET /admin/orders` — list all orders (with filters: status, date)
- [x] `GET /admin/orders/:id` — get order details with items
- [x] `PATCH /admin/orders/:id` — update order status (new → confirmed → shipped → canceled)
- [x] Test all admin endpoints manually

## Phase 6: Admin Dashboard (HTML)

- [x] Create admin layout template (minimal HTML + Tailwind CDN)
- [x] Build products list page (`/admin`)
- [x] Build variant management page (`/admin/products/:id/variants`)
  - [x] Add variant form (size, color, price, stock)
  - [x] Inline edit price/stock
  - [x] Toggle variant active/inactive
- [x] Build orders list page (`/admin/orders`)
  - [x] Filter by status
  - [x] Status update buttons (confirm, ship, cancel)
- [x] Build order detail page (`/admin/orders/:id`)
- [x] Add login page / auth gate
- [x] Mobile-friendly admin layout

## Phase 7: Frontend ↔ API Integration

- [x] Create `products.json` static file for product content (name, slug, category, images, description)
- [x] Update `index.html` to load products from `products.json`
- [x] On product card click → open product detail view
- [x] Fetch variants from `GET /api/products/:slug/variants` on product detail
- [x] Render size/color selectors based on variant data
- [x] Show correct price when variant is selected
- [x] Show stock status (in stock / out of stock / low stock)
- [x] Disable "Add to cart" for out-of-stock variants
- [x] Store cart in `localStorage` (persist across reloads)
- [x] Build checkout page/form (name, phone, address, city)
- [x] Submit order via `POST /api/orders`
- [x] Show order confirmation with order ID
- [x] Add loading spinners for API calls
- [x] Add error handling/retry for failed API calls
- [x] Handle CORS properly in fetch calls

## Phase 8: Security & Validation

- [x] Validate all user inputs server-side
- [x] Sanitize strings (prevent XSS)
- [x] Add rate limiting middleware (global + per-endpoint)
- [x] Set security headers (HSTS, X-Content-Type, etc.)
- [x] Protect admin routes with Cloudflare Access or strong auth
- [x] Ensure `POST /api/orders` is not cacheable
- [x] Add request size limits

## Phase 9: Caching & Performance

- [x] Set `Cache-Control` headers on variant responses (10–60 min)
- [x] Implement cache purge endpoint for admin (`POST /admin/cache/purge`)
- [x] Ensure Cloudflare caches static pages aggressively
- [x] Optimize product images (lazy loading, WebP)
- [x] Minify frontend JS/CSS for production
- [x] Test page load performance (Lighthouse)

## Phase 10: Deployment & Infrastructure

- [x] Provision PostgreSQL on Koyeb (or external like Neon/Supabase)
- [x] Add database URL to Koyeb service env vars
- [x] Add admin secret/API key to Koyeb env vars
- [x] Update Dockerfile to run migrations on startup
- [x] Configure custom domain `api.bhomanshah.com` on Koyeb
- [x] Deploy frontend to Cloudflare Pages at `bhomanshah.com`
- [x] Set up Cloudflare DNS records
- [x] Configure Cloudflare SSL/TLS
- [x] Test end-to-end flow in production
- [x] Set up error monitoring / alerting (optional)

## Phase 11: Testing

- [ ] Unit tests for order total calculation
- [ ] Unit tests for stock verification logic
- [ ] Integration tests for public API endpoints
- [ ] Integration tests for admin API endpoints
- [ ] Test cart flow end-to-end (add → checkout → order created)
- [ ] Test edge cases (out of stock, invalid variant, duplicate order)
- [ ] Load test the API with concurrent orders

## Phase 12: Polish & Launch

- [ ] Replace placeholder product images with real photos
- [ ] Update brand name / copy throughout
- [ ] Add favicon and Open Graph meta tags
- [ ] Add WhatsApp/contact link for support
- [ ] Final cross-device testing (mobile, tablet, desktop)
- [ ] Final review of all security measures
- [ ] Go live 🚀

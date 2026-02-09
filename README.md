Got it. With **monthly stock changes** and **variant-specific price (size+color)**, you can keep the storefront **static-fast** *and* still have a **WooCommerce-like “easy editing”** experience for the parts that change (variants/stock/prices).

Below is a **step-by-step plan** (no coding), with a clear day-to-day workflow.

---

# Plan: Static storefront + tiny API + simple admin

## Phase 1 — Structure the system (1-time)

### Step 1) Split your site into 2 parts

* **Storefront (static, super fast):** product pages, category pages, home, checkout page UI)
* **API (dynamic):** variants + prices + stock, and order creation

This makes browsing fast and keeps your small server safe.

### Step 2) Use two subdomains (cleanest)

* `yourdomain.com` → storefront (static)
* `api.yourdomain.com` → Koyeb Go Fiber API

(You *can* do single domain routing later, but this is simpler and more reliable.)

### Step 3) Decide where static storefront lives

Pick one:

* **Recommended:** Cloudflare Pages (or any static host behind Cloudflare)
* Alternative: serve static from Koyeb (works, but slower and puts load on your tiny server)

Cloudflare Pages is just “a place to upload HTML/CSS/JS” so Cloudflare can serve it fast globally.

---

## Phase 2 — Define how products/variants are managed (simple + Woo-ish)

You want WooCommerce-like ease for **variants/stock/price**.

### Step 4) Choose a “source of truth” split

Use a hybrid model:

**A) Product content (rare changes):** name, description, images, category, slug
→ keep this in a simple file (like `products.json`) used to build static pages.

**B) Variant data (frequent-ish changes):** size, color, variant price, stock
→ keep this in Postgres and edit via a small admin dashboard.

Why this fits you:

* You only update stock about **once per month**
* But variant pricing and stock should be editable without rebuilding pages

### Step 5) Define variant structure you will support (your case)

Per product:

* Sizes: S, M, L
* Colors: 3 colors
* Total variants per product: **up to 9**
  Each variant has:
* price
* stock count
* active/inactive (hide sold out or discontinued)

---

## Phase 3 — Build the storefront pages (fast browsing)

### Step 6) Generate static pages for browsing

Generate these pages:

* Home
* Category listing pages
* Product pages (`/p/slug`)
* Checkout page
* Thank-you page
* Policies pages

These pages:

* load instantly
* are cached by Cloudflare
* do not hit your Koyeb server for browsing

### Step 7) Product page behavior (important)

Each product page should:

* show product details from static HTML
* **load variants (size/color/price/stock) from the API** with one lightweight request
* render selectors (size/color dropdowns/buttons)
* show the correct price based on selected variant
* show stock availability (or just “in stock / out of stock”)

This keeps pages static-fast while still being flexible.

---

## Phase 4 — API responsibilities (keep it tiny)

### Step 8) API endpoints you actually need

* “Get variants for a product”
* “Create order (COD)”
* Optional: “Admin CRUD for products/variants”
* Optional: “Order list + update status” (admin)

That’s it.

### Step 9) Checkout rules (robust)

At checkout:

* customer selects a **variant** (size+color)
* API verifies:

  * variant exists
  * variant is active
  * stock is enough (if you enforce stock)
  * server calculates totals (never trust browser price)
* API stores an order snapshot (variant name/price at time of order)

---

## Phase 5 — Admin experience (WooCommerce-like enough)

### Step 10) Build a simple admin dashboard (minimal but effective)

Pages you need (keep it basic HTML forms):

1. **Products**

   * list products
   * enable/disable product
2. **Variants per product**

   * add variants (size/color combos)
   * set price per variant
   * set stock per variant
   * enable/disable variant
3. **Orders**

   * list orders
   * update status (new / confirmed / shipped / canceled)

This gives you the “Woo feel” for daily operations without building a heavy CMS.

### Step 11) Protect admin properly (important)

Use one of these:

* **Best & easiest:** Cloudflare Access (protect `/admin/*` with login)
* Or at minimum: strong password + basic auth + IP restriction

---

## Phase 6 — Caching & speed (where you win)

### Step 12) Cache static storefront aggressively

Cloudflare caches:

* HTML pages
* CSS/JS
* product pages
  So browsing stays extremely fast.

### Step 13) Cache variant API responses (safe and fast)

Variants don’t change often (monthly), so you can:

* Cache “variants JSON” for **10–60 minutes**
* When you update variants/stock, you can purge that one endpoint cache (or just wait)

Result:

* product pages still feel instant
* API is barely hit

### Step 14) Keep API uncached for order creation

* `POST /orders` should never be cached
* Add rate limiting for spam protection

---

## Phase 7 — Your actual day-to-day workflow (how “easy” it is)

### Common task A: Change stock once per month

1. Open admin
2. Choose product → variants
3. Update stock numbers (and prices if needed)
4. Save
5. (Optional) purge cached variants endpoint

✅ No redeploy needed.

### Common task B: Change variant price (size/color pricing)

1. Admin → product → variants
2. Edit price
3. Save
4. Optional cache purge

✅ No redeploy needed.

### Common task C: Add a new product (occasional)

1. Add product content (name/slug/images/category) in your product file
2. Rebuild + deploy static storefront
3. Admin: add variants for that product (9 combos) with stock + price

✅ Simple, predictable.

### Common task C2: Sync products into the API (when products.json changes)

1. Run migrations (only needed if new migrations were added)
2. Import products into the database

```bash
go run ./cmd/migrate
go run ./cmd/import-products
```

Optional: deactivate products not present in products.json

```bash
go run ./cmd/import-products --deactivate-missing
```

### Common task D: Remove/discontinue a product

1. Admin: disable product (or disable all variants)
2. Optionally remove product from static file later
3. Rebuild when convenient

✅ Immediate removal from checkout via “inactive”.

---

# Reality check: “Like WooCommerce?”

You won’t get Woo’s huge plugin ecosystem and themes for free.

But with this plan you *do* get the core “store owner convenience” you care about:

* Add/edit variants (size/color)
* Set price per variant
* Set stock per variant
* Manage orders

…without slowing the site down or overloading your tiny server.

---

## Current Implementation Status (February 2026)

This repository implements the plan above with the following enhancements:

### ✅ Completed Features

**Core Architecture:**
- Static storefront hosted on Cloudflare Pages (`bhomanshah.com`)
- Tiny Go API on Koyeb (`api.bhomanshah.com`) using Fiber framework
- MySQL database (TiDB Cloud compatible)

**Product Management:**
- Products stored in `products.json` (static content)
- Variants stored in database (dynamic pricing/stock)
- Admin dashboard for CRUD operations on products and variants

**E-commerce Features:**
- Product catalog with variants (size/color combinations)
- Shopping cart (client-side)
- COD order creation with stock validation
- Order management (admin)

**Security & Performance:**
- JWT-based admin authentication (secure, expiring tokens)
- Rate limiting on order creation (5/min per IP)
- HTTP caching for variants (30 minutes)
- Input validation and sanitization

### 🔐 Admin Authentication

The admin panel uses **JWT (JSON Web Tokens)** for secure authentication:

1. **Login Process:**
   - Admin enters API key on `/admin/login`
   - Server validates key and returns JWT token (expires in 24 hours)
   - Token stored securely in browser localStorage

2. **API Access:**
   - All admin API calls include `Authorization: Bearer <token>` header
   - Server validates token signature and expiry on each request
   - Automatic logout on expired tokens

3. **Backward Compatibility:**
   - Legacy `X-API-Key` header still supported as fallback

### 🚀 Deployment

**Environment Setup:**
```bash
# Copy environment file
cp .env.example .env

# Configure required variables
# - DATABASE_URL: MySQL connection string
# - ADMIN_SECRET: Your admin API key
# - JWT_SECRET: Random string for JWT signing (min 32 chars)
# - CORS_ORIGIN: Your frontend domain
```

**Build & Run:**
```bash
# Build the application
go build -o koyeb-go .

# Run locally
./koyeb-go

# Or deploy to Koyeb with environment variables
```

**Database Setup:**
```bash
# Run migrations
go run ./cmd/migrate

# Import products from products.json
go run ./cmd/import-products
```

### 📊 API Endpoints

**Public Endpoints:**
- `GET /api/products` - List active products
- `GET /api/products/:slug/variants` - Get variants for a product
- `POST /api/orders` - Create COD order

**Admin Endpoints:**
- `POST /admin/api/login` - Authenticate and get JWT token
- `GET /admin/api/products` - List all products (including inactive)
- `POST /admin/api/products` - Create product
- `PUT /admin/api/products/:id` - Update product
- `PATCH /admin/api/products/:id` - Toggle product active status
- `GET /admin/api/products/:id/variants` - List variants for product
- `POST /admin/api/products/:id/variants` - Create variant
- `PATCH /admin/api/variants/:id` - Update variant
- `DELETE /admin/api/variants/:id` - Delete variant
- `GET /admin/api/orders` - List orders with filters
- `GET /admin/api/orders/:id` - Get order details
- `PATCH /admin/api/orders/:id` - Update order status
- `POST /admin/api/cache/purge` - Purge cache

### 🔧 Configuration

Key environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | MySQL connection | `user:pass@tcp(host:3306)/db` |
| `CORS_ORIGIN` | Frontend domain | `https://bhomanshah.com` |
| `ADMIN_SECRET` | Admin API key | `your-secret-key` |
| `JWT_SECRET` | JWT signing key | `random-32-char-string` |
| `JWT_EXPIRY_HOURS` | Token expiry | `24` |

### 📈 Performance Optimizations

- **Static Content:** All product pages served from Cloudflare CDN
- **API Caching:** Variant data cached for 30 minutes
- **Database Indexing:** Optimized queries for products, variants, orders
- **Rate Limiting:** Prevents abuse on order creation
- **Input Validation:** Prevents malformed data and XSS attacks

### 🛠️ Development

**Project Structure:**
```
├── main.go                 # Server entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Custom middleware
│   ├── models/            # Data models
│   ├── repository/        # Database operations
│   └── validator/         # Input validation
├── cloudflare-pages-frontend/  # Static frontend
├── migrations/            # Database schema
└── cmd/                   # CLI tools (migrate, import)
```

**Adding New Features:**
1. Define models in `internal/models/`
2. Add repository methods in `internal/repository/`
3. Create handlers in `internal/handlers/`
4. Add routes in `main.go`
5. Update frontend as needed

This implementation successfully delivers the "static-fast storefront + tiny API" vision while providing WooCommerce-like admin capabilities for managing variants and orders.
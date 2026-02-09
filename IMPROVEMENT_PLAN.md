# Codebase Analysis & Improvement Plan

## 🚨 Critical Vulnerabilities & Risks

### 1. Fragile Admin Authentication
*   **Issue:** The admin panel uses a static API Secret (`X-API-Key`) stored in the browser (likely `localStorage`) and sent via headers.
*   **Risk:** This simplifies access (no specialized auth server), but if an attacker obtains this key (via XSS or network intercept), they have full, permanent control over the store. There is no mechanism to "expire" a session without rotating the server-side secret (which breaks access for everyone).
*   **Recommended Action:** Implement proper session-based auth (JWT or HttpOnly cookies) for the admin panel.

### 2. Scalability Bottleneck (Pagination)
*   **Issue:** The `GetProducts` handler calls `productRepo.GetAll(true)`, which loads **every single product** into memory and sends them all to the client.
*   **Risk:** As the catalog grows (e.g., >500 products), the API response size will balloon, slowing down the site and eventually crashing the Go server due to memory exhaustion.
*   **Recommended Action:** Implement `limit` and `offset` (pagination) parameters immediately.

### 3. In-Memory Caching on Serverless
*   **Issue:** The app uses an in-memory `CacheManager` (`handlers/cache.go`).
*   **Risk:** Cloud platforms (like Koyeb) may run multiple instances (replicas) of the app. An in-memory cache is isolated to each instance. If you update a price on Admin (purging Instance A's cache), a user might still hit Instance B and see the old price.
*   **Recommended Action:** Use a distributed cache (Redis) or disable internal caching and rely solely on Cloudflare/CDN caching with short TTLs.

## ⚠️ Architecture & Maintenance Issues

### 4. SEO & Client-Side Rendering
*   **Issue:** The "Static" frontend appears to be a Single Page Application (SPA) that fetches `products.json` or calls the API to render content.
*   **Shortcoming:** Search engines (Google) often struggle to index content that requires JavaScript to render. Product pages might appear as blank pages to crawlers.
*   **Recommended Action:** Shift to Static Site Generation (SSG) where HTML files are pre-built for every product, or implement Server-Side Rendering (SSR).

### 5. Brittle Frontend Build Process
*   **Issue:** The frontend relies on the Standalone Tailwind script (`src="/tw.js"`).
*   **Shortcoming:** This parses CSS in the user's browser, causing a "Flash of Unstyled Content" (FOUC) and slower page loads compared to a standard CSS build step.
*   **Recommended Action:** Add a proper build step (Vite/PostCSS) to generate a single optimized `.css` file.

### 6. Handler Bloat
*   **Issue:** `internal/handlers/handlers.go` handles Products, Orders, Admin Auth, and Cache logic in one large file.
*   **Shortcoming:** This reduces maintainability and readability.
*   **Recommended Action:** Refactor into `product_handler.go`, `order_handler.go`, `admin_handler.go`.

## 📉 Missing E-commerce Essentials

### 7. No Transactional Emails
*   **Shortcoming:** Users receive no confirmation email after placing a COD order. This leads to low trust and support volume.
*   **Recommended Action:** Integrate an email provider (SMTP/SendGrid) to send order confirmations.

### 8. No Search Functionality
*   **Shortcoming:** Users can only browse. Finding a specific item requires scrolling through the entire list.
*   **Recommended Action:** Implement a simple search endpoint (`/api/products?search=...`) or client-side filtering.

### 9. Hard Coded Paths
*   **Issue:** The Go server serves static files using relative paths like `./cloudflare-pages-frontend`.
*   **Shortcoming:** This breaks easily if the Docker working directory changes or the binary is moved.
*   **Recommended Action:** Use environment variables or absolute paths for asset serving.

## 🧪 Quality Assurance

### 10. Lack of Unit Tests
*   **Issue:** Critical business logic (Order total calculation, stock decrement) has no unit tests.
*   **Risk:** Refactoring could introduce regression bugs in financial calculations.
*   **Recommended Action:** Write Go unit tests for `internal/models` and `internal/repository`.

---

## Summary Checklist for Improvement

- [ ] **High Priority**: Add Pagination to `GET /products`
- [ ] **High Priority**: Implement simple email notifications for orders
- [ ] **Maintenance**: Split `handlers.go` into separate files
- [ ] **Performance**: Add a real CSS build step (remove runtime Tailwind)
- [ ] **Security**: Review Admin Auth mechanism

---

## Additional Fix Plan & Shortcomings (Feb 9, 2026)

### 11. Cache Purge Is Ineffective
*   **Issue:** `AdminPurgeCache` updates `CacheManager`, but `GetProductVariants` never checks it.
*   **Risk:** Admin cache purge appears to succeed but does not change responses; stale data persists until CDN cache expires.
*   **Recommended Action:** Either remove `CacheManager` entirely (and rely only on CDN cache) or integrate `CacheManager.ShouldInvalidate` checks inside `GetProductVariants` to bypass cache headers or force a short TTL for invalidated slugs.

### 12. Variant Deactivation Race During Order Creation
*   **Issue:** Order validation checks `variant.Active`, but the stock decrement SQL does not enforce `active = TRUE`.
*   **Risk:** If a variant is deactivated after validation but before the transaction completes, an order can still be created.
*   **Recommended Action:** Add `AND active = TRUE` to the stock update query and handle zero affected rows as a variant-inactive error.

### 13. Missing Validation on Variant Updates
*   **Issue:** `AdminUpdateVariant` accepts arbitrary `price` and `stock` values, allowing negative stock or zero price.
*   **Risk:** Corrupted pricing and inventory data can break checkout and reporting.
*   **Recommended Action:** Apply validation rules for partial updates (e.g., `price > 0`, `stock >= 0`) and reject invalid payloads.

### 14. Admin Images Field Allows Invalid JSON
*   **Issue:** `images` is stored as a string without JSON validation in product create/update.
*   **Risk:** Invalid JSON breaks storefront image parsing and UI rendering.
*   **Recommended Action:** Validate `images` is a JSON array before saving or convert to a structured array server-side.

### 15. CORS Policy Ignores Config
*   **Issue:** `CORS(allowOrigin string)` ignores its argument and always allows `*`.
*   **Risk:** In production, this defeats intended origin restrictions.
*   **Recommended Action:** Honor `allowOrigin` from config (and allow a comma-separated list if needed).

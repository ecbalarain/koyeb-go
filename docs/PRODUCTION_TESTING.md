# Production Testing Guide

## Overview
Comprehensive testing checklist for OXLOOK production deployment.

---

## Pre-Deployment Testing

### Local Environment

```bash
# Test build
go build -o oxlook-api

# Run locally
./oxlook-api
# Should start on port 8080

# Test health endpoint
curl http://localhost:8080/health
# Expected: {"status":"ok"}
```

### Docker Build Test

```bash
# Build Docker image
docker build -t oxlook-api .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL="mysql://..." \
  -e ADMIN_SECRET="test-secret" \
  -e CORS_ORIGIN="http://localhost:3000" \
  -e ENVIRONMENT="development" \
  oxlook-api

# Test endpoints
curl http://localhost:8080/health
```

---

## Post-Deployment Testing

### 1. Service Health Check

```bash
# Health endpoint
curl https://api.bhomanshah.com/health

# Expected response:
# {"status":"ok"}

# Check response time (should be <500ms)
curl -w "@curl-format.txt" -o /dev/null -s https://api.bhomanshah.com/health
```

Create `curl-format.txt`:
```
    time_namelookup:  %{time_namelookup}s\n
       time_connect:  %{time_connect}s\n
    time_appconnect:  %{time_appconnect}s\n
   time_pretransfer:  %{time_pretransfer}s\n
      time_redirect:  %{time_redirect}s\n
 time_starttransfer:  %{time_starttransfer}s\n
                    ----------\n
         time_total:  %{time_total}s\n
```

### 2. Public API Endpoints

```bash
# Test products endpoint
curl https://api.bhomanshah.com/api/products | jq

# Expected: Array of active products
# [
#   {
#     "id": 2,
#     "name": "Stone Mug",
#     "slug": "stone-mug",
#     "category": "Home",
#     "images": "..."
#   }
# ]

# Test variants endpoint
curl https://api.bhomanshah.com/api/products/stone-mug/variants | jq

# Expected: Array of variants
# [
#   {
#     "id": 1,
#     "product_id": 2,
#     "size": "Standard",
#     "color": "Gray",
#     "price": 2499,
#     "stock": 15,
#     "active": true
#   }
# ]

# Test cache headers
curl -I https://api.bhomanshah.com/api/products/stone-mug/variants

# Should see:
# Cache-Control: public, max-age=1800, s-maxage=1800
# Vary: Accept-Encoding
```

### 3. Order Creation

```bash
# Create test order
curl -X POST https://api.bhomanshah.com/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "customer_phone": "1234567890",
    "customer_address": "123 Main St",
    "customer_city": "New York",
    "notes": "Please deliver before 5pm",
    "items": [
      {
        "variant_id": 1,
        "qty": 2
      }
    ]
  }'

# Expected response:
# {
#   "order_id": 1,
#   "total": 4998,
#   "status": "new",
#   "message": "Order created successfully"
# }

# Test with invalid data
curl -X POST https://api.bhomanshah.com/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "",
    "items": []
  }'

# Expected: 400 Bad Request with validation error
```

### 4. Rate Limiting

```bash
# Test global rate limit (100 req/min)
for i in {1..105}; do
  curl -s https://api.bhomanshah.com/health > /dev/null
  echo "Request $i"
done

# After 100 requests, should see:
# {"error":"Too many requests. Please try again later."}

# Test order rate limit (5 req/min)
for i in {1..6}; do
  curl -X POST https://api.bhomanshah.com/api/orders \
    -H "Content-Type: application/json" \
    -d '{"customer_name":"Test","customer_phone":"1234567890","customer_address":"Test","customer_city":"Test","items":[{"variant_id":1,"qty":1}]}'
  sleep 1
done

# After 5 requests, should see:
# {"error":"Too many requests. Please try again later."}
```

### 5. Admin API Endpoints

```bash
# Save admin secret
ADMIN_SECRET="your-admin-secret-here"

# Test without API key (should fail)
curl https://api.bhomanshah.com/admin/api/products

# Expected: {"error":"Missing API key"}

# Test with wrong API key (should fail)
curl -H "X-API-Key: wrong-key" \
  https://api.bhomanshah.com/admin/api/products

# Expected: {"error":"Invalid API key"}

# Test with correct API key (should work)
curl -H "X-API-Key: $ADMIN_SECRET" \
  https://api.bhomanshah.com/admin/api/products | jq

# Expected: Array of all products (including inactive)

# Test cache purge
curl -X POST https://api.bhomanshah.com/admin/api/cache/purge \
  -H "X-API-Key: $ADMIN_SECRET" \
  -H "Content-Type: application/json" \
  -d '{"type":"all"}'

# Expected: {"message":"All cache entries invalidated successfully"}
```

### 6. CORS Testing

Open browser console on https://bhomanshah.com:

```javascript
// Test CORS from frontend
fetch('https://api.bhomanshah.com/api/products')
  .then(r => r.json())
  .then(data => console.log('Products:', data))
  .catch(err => console.error('CORS Error:', err));

// Should successfully return products

// Test from wrong origin (should fail)
// Open any other website and try:
fetch('https://api.bhomanshah.com/api/products')
  .then(r => r.json())
  .then(console.log);

// Should see CORS error in console
```

### 7. Security Headers

```bash
# Check security headers
curl -I https://api.bhomanshah.com/api/products

# Should see:
# X-Frame-Options: DENY
# X-Content-Type-Options: nosniff
# Referrer-Policy: strict-origin-when-cross-origin
# X-DNS-Prefetch-Control: off
# X-Download-Options: noopen
```

### 8. Input Sanitization

```bash
# Test XSS prevention
curl -X POST https://api.bhomanshah.com/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "<script>alert(\"XSS\")</script>",
    "customer_phone": "1234567890",
    "customer_address": "123 Main St",
    "customer_city": "New York",
    "items": [{"variant_id": 1, "qty": 1}]
  }'

# Then check order in admin panel
# Name should be HTML-escaped: &lt;script&gt;alert("XSS")&lt;/script&gt;

# Test SQL injection prevention
curl https://api.bhomanshah.com/api/products/stone-mug%27%20OR%20%271%27=%271/variants

# Should return "Product not found" or 404, not database error
```

---

## Frontend Integration Testing

### 1. Static Files

```bash
# Test products.json
curl https://bhomanshah.com/products.json | jq

# Should return product data

# Test cache headers
curl -I https://bhomanshah.com/products.json

# Should see:
# Cache-Control: public, max-age=600, s-maxage=600
```

### 2. End-to-End Flow

**Manual Testing Steps:**

1. **Visit Homepage**
   - Go to https://bhomanshah.com
   - Page should load quickly (<3s)
   - Products should display

2. **Browse Products**
   - Click on a product card
   - Product detail should open
   - Variants should load from API

3. **Add to Cart**
   - Select size/color
   - Click "Add to Cart"
   - Cart count should increment
   - localStorage should be updated

4. **Checkout**
   - Open cart
   - Click "Checkout"
   - Fill out form
   - Submit order
   - Should receive order confirmation with order ID

5. **Admin Panel**
   - Go to https://bhomanshah.com/admin/login
   - Enter API key
   - Should see admin dashboard
   - Check orders list
   - Verify test order appears

### 3. Cross-Browser Testing

Test on:
- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (latest)
- [ ] Edge (latest)
- [ ] Mobile Safari (iOS)
- [ ] Chrome Mobile (Android)

### 4. Lighthouse Audit

```bash
# Install Lighthouse
npm install -g lighthouse

# Run audit
lighthouse https://bhomanshah.com --view

# Target scores:
# Performance: >90
# Accessibility: >95
# Best Practices: >95
# SEO: >90
```

---

## Performance Testing

### 1. Load Testing with Apache Bench

```bash
# Install Apache Bench
# macOS: already installed
# Ubuntu: sudo apt install apache2-utils

# Test health endpoint (100 requests, 10 concurrent)
ab -n 100 -c 10 https://api.bhomanshah.com/health

# Should see:
# Requests per second: >100
# Time per request (mean): <100ms
# Failed requests: 0

# Test products endpoint
ab -n 100 -c 10 https://api.bhomanshah.com/api/products

# Test variants endpoint (cached)
ab -n 100 -c 10 https://api.bhomanshah.com/api/products/stone-mug/variants
```

### 2. Load Testing with k6

```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp up to 20 users
    { duration: '1m', target: 20 },   // Stay at 20 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
};

export default function() {
  // Test products endpoint
  let res = http.get('https://api.bhomanshah.com/api/products');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
  
  sleep(1);
}
```

Run test:
```bash
# Install k6
brew install k6

# Run test
k6 run load-test.js
```

### 3. Database Query Performance

```bash
# Connect to database
mysql -h <host> -u <user> -p <database>

# Check slow queries
SHOW VARIABLES LIKE 'slow_query_log';

# Enable if not enabled
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

# Monitor slow queries
SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;
```

---

## Monitoring & Alerting

### 1. Uptime Monitoring

**UptimeRobot (Free)**
1. Create account: https://uptimerobot.com
2. Add monitor:
   - Type: HTTPS
   - URL: https://api.bhomanshah.com/health
   - Interval: 5 minutes
   - Alert when down for: 2 minutes

**Alternatives:**
- Pingdom
- Checkly
- StatusCake

### 2. Error Tracking

Check Koyeb logs for errors:

```bash
# Via CLI
koyeb service logs oxlook-api --tail 100 | grep ERROR

# Or use Sentry (optional)
```

### 3. Database Monitoring

- Check provider dashboard (PlanetScale, Neon, etc.)
- Monitor:
  - Active connections
  - Query performance
  - Disk usage
  - Error rate

---

## Security Testing

### 1. SSL/TLS Configuration

```bash
# Test SSL
https://www.ssllabs.com/ssltest/analyze.html?d=api.bhomanshah.com

# Should get A or A+ rating

# Test manually
openssl s_client -connect api.bhomanshah.com:443 -servername api.bhomanshah.com
```

### 2. Security Headers

```bash
# Use securityheaders.com
https://securityheaders.com/?q=api.bhomanshah.com

# Should see:
# X-Frame-Options: DENY
# X-Content-Type-Options: nosniff
# Referrer-Policy: strict-origin-when-cross-origin
```

### 3. Penetration Testing (Basic)

```bash
# Try common attacks
# 1. SQL Injection
curl "https://api.bhomanshah.com/api/products/1'%20OR%20'1'='1/variants"

# 2. Path traversal
curl "https://api.bhomanshah.com/api/../../etc/passwd"

# 3. XXE injection
curl -X POST https://api.bhomanshah.com/api/orders \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]><foo>&xxe;</foo>'

# All should be properly rejected
```

---

## Production Checklist

### Before Go-Live

- [ ] Database backups configured
- [ ] Environment variables set correctly
- [ ] Custom domain working (api.bhomanshah.com)
- [ ] SSL certificate active and valid
- [ ] CORS configured for production domain
- [ ] Rate limiting tested
- [ ] Admin API key secured (not exposed)
- [ ] All API endpoints tested
- [ ] Security headers present
- [ ] Input sanitization working
- [ ] Error responses don't leak sensitive data
- [ ] Logs don't contain sensitive information
- [ ] Frontend connected to production API
- [ ] End-to-end flow tested
- [ ] Lighthouse score >90
- [ ] Load test passed
- [ ] Uptime monitoring configured
- [ ] Documentation updated

### Post Go-Live

- [ ] Monitor logs for errors (first 24 hours)
- [ ] Check order flow daily (first week)
- [ ] Monitor database performance
- [ ] Review security logs
- [ ] Test backup restoration
- [ ] Verify auto-deployment working
- [ ] Check SSL certificate expiration date
- [ ] Monitor traffic patterns
- [ ] Gather user feedback
- [ ] Plan for scaling if needed

---

## Troubleshooting Production Issues

### Service Down

1. Check Koyeb status page
2. Review recent deployments
3. Check logs for errors
4. Verify database connectivity
5. Roll back if needed

### Slow Response Times

1. Check database query performance
2. Review slow query logs
3. Monitor connection pool
4. Check cache hit rate
5. Consider adding indexes

### High Error Rate

1. Check error logs
2. Identify error patterns
3. Verify external dependencies (database)
4. Check for deployment issues
5. Monitor traffic for attacks

---

## Next Steps

After successful testing:

1. ✅ All endpoints working
2. ✅ Security measures in place
3. ✅ Performance acceptable
4. [ ] Set up monitoring alerts
5. [ ] Document runbook for common issues
6. [ ] Plan for scaling
7. [ ] Schedule regular security audits

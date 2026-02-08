# Production Deployment Checklist

Complete checklist for deploying OXLOOK to production.

---

## Pre-Deployment

### Code Preparation
- [ ] All features implemented and tested locally
- [ ] Code committed to Git
- [ ] Branch: `main` is up to date
- [ ] No secrets in code repository
- [ ] Docker build tested locally
- [ ] Migrations tested locally
- [ ] `.gitignore` includes `.env` and sensitive files

### Documentation
- [ ] README.md updated
- [ ] API endpoints documented
- [ ] Environment variables documented
- [ ] Deployment guides reviewed

---

## Database Setup

- [ ] Database provider chosen (PlanetScale recommended)
- [ ] Database created
- [ ] Connection string obtained
- [ ] Database accessible from Koyeb IPs
- [ ] SSL/TLS enabled
- [ ] Backups configured
- [ ] Test connection from local machine

**Database URL Format:**
```
mysql://username:password@host:port/database?tls=true
```

---

## Environment Variables

Generate and save these securely:

### Required Secrets

```bash
# Generate admin secret (32 bytes)
ADMIN_SECRET=$(openssl rand -base64 32)

# Save to password manager!
echo $ADMIN_SECRET
```

### Environment Configuration

- [ ] `DATABASE_URL` - Database connection string
- [ ] `ADMIN_SECRET` - Admin API key (generated above)
- [ ] `CORS_ORIGIN` - Frontend URL (https://bhomanshah.com)
- [ ] `ENVIRONMENT` - Set to "production"
- [ ] `PORT` - 8080 (optional, Koyeb sets automatically)

---

## Koyeb Deployment

### Service Creation

- [ ] Koyeb account created
- [ ] GitHub repository connected
- [ ] Service name: `oxlook-api`
- [ ] Region selected: Washington D.C. or Frankfurt
- [ ] Builder: Dockerfile
- [ ] Environment variables added as secrets
- [ ] Health check configured: `/health`
- [ ] Service deployed successfully

### Verify Deployment

- [ ] Build completed without errors
- [ ] Migrations ran successfully (check logs)
- [ ] Service status: "Running"
- [ ] Health endpoint responding: `https://<app>.koyeb.app/health`
- [ ] Products endpoint working
- [ ] Variants endpoint working

---

## Custom Domain (API)

### DNS Configuration in Cloudflare

- [ ] Cloudflare account has domain
- [ ] DNS record created:
  - Type: `CNAME`
  - Name: `api`
  - Target: `<your-service>.koyeb.app`
  - Proxy status: **DNS only** (gray cloud, not proxied)
  - TTL: Auto

### Koyeb Domain Setup

- [ ] Domain added in Koyeb: `api.bhomanshah.com`
- [ ] DNS configured (CNAME created in Cloudflare)
- [ ] SSL certificate provisioned (wait 2-10 minutes)
- [ ] Domain status: "Active"
- [ ] Certificate status: "Active"

### Verify Domain

```bash
# Check DNS resolution
dig api.bhomanshah.com

# Test HTTPS
curl https://api.bhomanshah.com/health
```

- [ ] DNS resolves correctly
- [ ] HTTPS working (SSL certificate valid)
- [ ] API responding at custom domain

---

## Cloudflare Pages (Frontend)

### Repository Setup

- [ ] Frontend code in repository
- [ ] Build script created: `build-frontend.sh`
- [ ] `_headers` file configured
- [ ] `products.json` file ready

### Cloudflare Pages Configuration

- [ ] Cloudflare Pages project created
- [ ] Connected to GitHub repository
- [ ] Build command: `./build-frontend.sh`
- [ ] Build output directory: `cloudflare-pages-frontend/dist`
- [ ] Custom domain: `bhomanshah.com`
- [ ] SSL certificate active

### Cloudflare Optimizations

- [ ] Auto Minify enabled (JS, CSS, HTML)
- [ ] Brotli compression enabled
- [ ] Polish enabled (Lossless or Lossy)
- [ ] WebP conversion enabled
- [ ] Mirage enabled (lazy loading)
- [ ] Always Use HTTPS enabled
- [ ] HTTP/3 enabled

### DNS for Frontend

- [ ] DNS records configured:
  - `bhomanshah.com` → Cloudflare Pages
  - `www.bhomanshah.com` → Cloudflare Pages (optional)

---

## API Integration

### Update Frontend Configuration

In `index.html`, ensure API calls use production URL:

```javascript
const API_BASE_URL = 'https://api.bhomanshah.com';
```

- [ ] API URL updated in frontend code
- [ ] CORS origin matches in backend (`https://bhomanshah.com`)
- [ ] No localhost URLs remaining in code
- [ ] Admin panel points to production API

---

## Testing

### API Endpoints

```bash
# Health check
curl https://api.bhomanshah.com/health

# Products
curl https://api.bhomanshah.com/api/products

# Variants
curl https://api.bhomanshah.com/api/products/stone-mug/variants

# Create order (test)
curl -X POST https://api.bhomanshah.com/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Test User",
    "customer_phone": "1234567890",
    "customer_address": "123 Test St",
    "customer_city": "Test City",
    "items": [{"variant_id": 1, "qty": 1}]
  }'

# Admin (with API key)
curl -H "X-API-Key: YOUR_ADMIN_SECRET" \
  https://api.bhomanshah.com/admin/api/products
```

- [ ] All public endpoints working
- [ ] Admin endpoints require API key
- [ ] Rate limiting working
- [ ] CORS headers present
- [ ] Security headers present
- [ ] Cache headers correct

### Frontend Testing

- [ ] Homepage loads at https://bhomanshah.com
- [ ] Products display correctly
- [ ] Product detail opens and loads variants from API
- [ ] Cart functionality works
- [ ] Checkout form submits order to API
- [ ] Order confirmation shows order ID
- [ ] Admin panel accessible and functional

### Cross-Browser Testing

- [ ] Chrome (Desktop)
- [ ] Firefox (Desktop)
- [ ] Safari (Desktop)
- [ ] Edge (Desktop)
- [ ] Mobile Safari (iOS)
- [ ] Chrome Mobile (Android)

### Performance Testing

- [ ] Lighthouse score > 90
- [ ] Page load time < 3 seconds
- [ ] API response time < 500ms
- [ ] No console errors
- [ ] No 404 errors in network tab

---

## Security Verification

### SSL/TLS

- [ ] SSL Labs test: A or A+ rating
  - Test: https://www.ssllabs.com/ssltest/analyze.html?d=api.bhomanshah.com
- [ ] Certificate valid and trusted
- [ ] HSTS header present (on Koyeb HTTPS)
- [ ] All connections use HTTPS

### Security Headers

```bash
curl -I https://api.bhomanshah.com/api/products
```

- [ ] `X-Frame-Options: DENY`
- [ ] `X-Content-Type-Options: nosniff`
- [ ] `Referrer-Policy: strict-origin-when-cross-origin`
- [ ] Other security headers present

### Authentication & Authorization

- [ ] Admin endpoints require `X-API-Key` header
- [ ] Wrong API key returns 401 Unauthorized
- [ ] Missing API key returns 401 Unauthorized
- [ ] Public endpoints don't require authentication

### Input Validation

- [ ] XSS prevention tested (HTML entities escaped)
- [ ] SQL injection prevented (prepared statements)
- [ ] Rate limiting working
- [ ] Request body size limits enforced (2MB)
- [ ] Invalid input returns appropriate error messages

---

## Monitoring & Alerting

### Uptime Monitoring

- [ ] UptimeRobot configured
  - Monitor: https://api.bhomanshah.com/health
  - Interval: 5 minutes
  - Alert: Email when down
- [ ] Status page created (optional)

### Error Tracking (Optional)

- [ ] Sentry account created
- [ ] Project configured
- [ ] DSN added to environment variables
- [ ] Test error sent and received

### Analytics

- [ ] Cloudflare Web Analytics added to frontend
- [ ] Analytics beacon working
- [ ] Metrics visible in dashboard

### Logging

- [ ] Koyeb logs accessible
- [ ] No sensitive data in logs
- [ ] Error logs monitored
- [ ] Log retention understood

---

## Backup & Recovery

### Database Backups

- [ ] Automatic backups enabled on database provider
- [ ] Backup retention policy understood
- [ ] Backup restoration tested (in development)
- [ ] Manual backup process documented

### Code Backups

- [ ] Code in Git repository
- [ ] GitHub repository not public (or secrets not committed)
- [ ] Tags created for releases
- [ ] Rollback procedure documented

---

## Documentation

### For Team

- [ ] Deployment guides in `docs/` folder
- [ ] Environment variables documented
- [ ] API endpoints documented
- [ ] Admin procedures documented
- [ ] Incident response plan created
- [ ] Runbook for common issues

### For Users (Optional)

- [ ] Help/FAQ page
- [ ] Contact information
- [ ] Terms of service
- [ ] Privacy policy
- [ ] Return/refund policy

---

## Post-Deployment Monitoring

### First 24 Hours

- [ ] Monitor error logs hourly
- [ ] Check order flow
- [ ] Verify email notifications (if implemented)
- [ ] Monitor traffic patterns
- [ ] Check database performance
- [ ] Verify backups running

### First Week

- [ ] Daily log review
- [ ] Order flow verification
- [ ] Performance monitoring
- [ ] Customer feedback gathering
- [ ] Security log review

### Ongoing

- [ ] Weekly performance review
- [ ] Monthly security audit
- [ ] Quarterly dependency updates
- [ ] Database optimization
- [ ] Cost optimization

---

## Rollback Plan

If critical issues occur:

1. **Immediate Actions**
   - [ ] Identify the issue
   - [ ] Assess impact
   - [ ] Decide: Fix forward or rollback

2. **Rollback Procedure**
   - [ ] Koyeb: Redeploy previous version
   - [ ] Cloudflare Pages: Rollback to previous deployment
   - [ ] Database: Restore from backup (if needed)
   - [ ] Verify service restored

3. **Communication**
   - [ ] Update status page
   - [ ] Notify affected users
   - [ ] Document incident
   - [ ] Plan fix for next deployment

---

## Final Verification

Before announcing launch:

- [ ] All checklist items above completed
- [ ] Full end-to-end test successful
- [ ] Team trained on admin panel
- [ ] Support channels ready
- [ ] Payment processing tested (if applicable)
- [ ] Legal requirements met (privacy policy, terms)
- [ ] Performance targets met
- [ ] Security audit passed
- [ ] Backup & recovery tested
- [ ] Monitoring & alerting active

---

## Launch Day

### Pre-Launch (T-1 hour)

- [ ] Final smoke test
- [ ] Team on standby
- [ ] Monitoring dashboard open
- [ ] Status page ready (if using)

### Launch (T=0)

- [ ] Announce on social media / marketing channels
- [ ] Send to initial users
- [ ] Enable monitoring alerts
- [ ] Start log monitoring

### Post-Launch (T+1 hour)

- [ ] Check error rate
- [ ] Verify orders processing
- [ ] Monitor server resources
- [ ] Respond to user feedback

---

## Success Metrics

### Technical Metrics

- [ ] Uptime > 99.9%
- [ ] Error rate < 0.1%
- [ ] Response time < 500ms (p95)
- [ ] Lighthouse score > 90

### Business Metrics

- [ ] Orders processing successfully
- [ ] Cart abandonment rate tracked
- [ ] Conversion rate tracked
- [ ] Customer feedback positive

---

## Troubleshooting Contacts

**Koyeb Support:**
- Docs: https://www.koyeb.com/docs
- Email: support@koyeb.com

**Cloudflare Support:**
- Docs: https://developers.cloudflare.com
- Community: https://community.cloudflare.com

**Database Support:**
- PlanetScale: support@planetscale.com
- Check provider docs

**Emergency Contacts:**
- Team Lead: [Name/Phone]
- DevOps: [Name/Phone]
- Database Admin: [Name/Phone]

---

## Congratulations! 🎉

If all items are checked, your OXLOOK store is live in production!

Next steps:
- Monitor closely for first 24-48 hours
- Gather user feedback
- Iterate and improve
- Plan for scaling as traffic grows

**Remember:** Deployment is not the end, it's the beginning! 🚀

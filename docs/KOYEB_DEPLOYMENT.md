# Koyeb Deployment Guide

## Overview
Complete guide to deploying the OXLOOK API on Koyeb with custom domain configuration.

---

## Prerequisites

- [ ] Koyeb account (free tier available)
- [ ] GitHub repository with code
- [ ] Database provisioned (see DATABASE_SETUP.md)
- [ ] Admin secret generated

---

## Step 1: Create Koyeb Account

1. Go to https://koyeb.com
2. Sign up with GitHub (recommended)
3. Verify email address
4. Free tier includes:
   - 1 web service (512MB RAM)
   - 2 million invocations/month
   - 100GB bandwidth/month

---

## Step 2: Prepare Environment Variables

Create these secrets before deployment:

### Required Secrets

```bash
# Database connection
DATABASE_URL=mysql://username:password@host:port/database?tls=true

# Admin API authentication (generate strong secret)
ADMIN_SECRET=$(openssl rand -base64 32)
# Example: ghp_xF8k2nR9vH3mJ5qL7wN1pD4tY6sK0oB2

# CORS origin (your Cloudflare Pages domain)
CORS_ORIGIN=https://bhomanshah.com

# Environment
ENVIRONMENT=production

# Port (Koyeb automatically sets this, but can be explicit)
PORT=8080
```

### Generate Admin Secret

```bash
# Strong secret (32 bytes)
openssl rand -base64 32

# Or online: https://www.random.org/strings/
```

---

## Step 3: Create Koyeb Service

### Via Web Dashboard

1. **Create New Service**
   - Dashboard > Create Service
   - Select "GitHub" as source

2. **Connect Repository**
   - Authorize Koyeb to access GitHub
   - Select repository: `your-username/koyeb-go`
   - Branch: `main`

3. **Builder Configuration**
   - Builder: Dockerfile
   - Dockerfile path: `Dockerfile` (in root)
   - Context path: `.` (root directory)

4. **Service Configuration**
   - Service name: `oxlook-api`
   - Region: Washington, D.C. (us-east) or Frankfurt (eu-west)
   - Instance type: Web
   - Instance size: Free (512MB RAM, 0.1 vCPU)

5. **Environment Variables**
   Click "Add environment variable" for each:
   
   | Name | Value | Type |
   |------|-------|------|
   | DATABASE_URL | mysql://... | Secret |
   | ADMIN_SECRET | your-secret-key | Secret |
   | CORS_ORIGIN | https://bhomanshah.com | Plain text |
   | ENVIRONMENT | production | Plain text |

6. **Health Checks** (Optional but recommended)
   - Path: `/health`
   - Port: 8080
   - Protocol: HTTP
   - Initial delay: 10 seconds
   - Timeout: 5 seconds

7. **Scaling** (Free tier fixed)
   - Min instances: 1
   - Max instances: 1

8. **Deploy**
   - Click "Create Service"
   - Wait for build and deployment

### Via CLI (Alternative)

```bash
# Install Koyeb CLI
curl https://get.koyeb.com | sh

# Login
koyeb login

# Create service
koyeb service create oxlook-api \
  --git github.com/your-username/koyeb-go \
  --git-branch main \
  --git-builder dockerfile \
  --ports 8080:http \
  --routes /:8080 \
  --env DATABASE_URL=<secret> \
  --env ADMIN_SECRET=<secret> \
  --env CORS_ORIGIN=https://bhomanshah.com \
  --env ENVIRONMENT=production \
  --instance-type free \
  --regions was
```

---

## Step 4: Configure Custom Domain

### Add Domain to Koyeb

1. **Service Settings**
   - Select your service
   - Go to "Domains" tab
   - Click "Add domain"

2. **Configure Domain**
   - Domain name: `api.bhomanshah.com`
   - Click "Add"

3. **Get DNS Settings**
   Koyeb will show:
   - **Type:** CNAME
   - **Name:** `api` or `api.bhomanshah.com`
   - **Value:** `<app-id>.koyeb.app`

### Configure DNS in Cloudflare

1. **Go to Cloudflare Dashboard**
   - Select domain: bhomanshah.com
   - DNS > Records

2. **Add CNAME Record**
   - Type: `CNAME`
   - Name: `api`
   - Target: `<your-service>.koyeb.app` (from Koyeb dashboard)
   - TTL: Auto
   - Proxy status: **Off** (gray cloud - DNS only)
   
   ⚠️ **Important:** Must be DNS only (not proxied) for Koyeb SSL

3. **Save Record**
   - Click "Save"
   - Wait 1-5 minutes for propagation

### Verify DNS

```bash
# Check DNS resolution
dig api.bhomanshah.com

# Should show CNAME pointing to Koyeb
nslookup api.bhomanshah.com
```

### Wait for SSL Certificate

1. Back in Koyeb dashboard
2. Domains tab should show:
   - Status: "Pending" → "Active"
   - Certificate: "Provisioning" → "Active"
   - Takes 2-10 minutes

3. Once active, your API is live at:
   - https://api.bhomanshah.com

---

## Step 5: Verify Deployment

### Check Service Status

1. **Koyeb Dashboard**
   - Service should show "Running"
   - Deployment status: "Healthy"
   - Recent logs visible

2. **Check Logs**
   ```
   🚀 Starting OXLOOK API deployment...
   📊 Running database migrations...
   ✅ Migrations completed successfully
   🌐 Starting API server...
   Starting OXLOOK API on :8080 (env: production)
   ```

### Test API Endpoints

```bash
# Health check
curl https://api.bhomanshah.com/health
# Expected: {"status":"ok"}

# Get products
curl https://api.bhomanshah.com/api/products
# Expected: Array of products

# Get variants
curl https://api.bhomanshah.com/api/products/stone-mug/variants
# Expected: Array of variants with price/stock

# Test CORS (from browser console on bhomanshah.com)
fetch('https://api.bhomanshah.com/api/products')
  .then(r => r.json())
  .then(console.log)
```

### Test Admin Endpoints

```bash
# Without API key (should fail)
curl https://api.bhomanshah.com/admin/api/products
# Expected: {"error":"Missing API key"}

# With API key (should work)
curl -H "X-API-Key: YOUR_ADMIN_SECRET" \
  https://api.bhomanshah.com/admin/api/products
# Expected: Array of all products
```

---

## Step 6: Configure Auto-Deployment

Koyeb automatically deploys on every push to main branch.

### Disable Auto-Deploy (Optional)

1. Service > Settings
2. Build settings
3. Uncheck "Auto deploy"
4. Deploy manually: `koyeb service redeploy oxlook-api`

### Deploy from Specific Branch/Tag

1. Service > Settings
2. Build settings
3. Change branch to `production` or use tags
4. Save

---

## Step 7: Monitoring & Logs

### View Logs

**Via Dashboard:**
1. Service > Logs
2. Filter by:
   - Level (Info, Error, Warning)
   - Time range
   - Search text

**Via CLI:**
```bash
# Tail logs
koyeb service logs oxlook-api --follow

# Last 100 lines
koyeb service logs oxlook-api --tail 100

# Specific instance
koyeb service logs oxlook-api --instance <instance-id>
```

### Metrics

Dashboard > Metrics shows:
- HTTP requests/sec
- Response times (p50, p95, p99)
- Error rate
- Memory usage
- CPU usage

### Alerts (Paid Plans)

Set up alerts for:
- Service down
- Error rate > 5%
- Response time > 2s
- Memory > 90%

---

## Step 8: Update Frontend Configuration

Update Cloudflare Pages frontend to use production API:

```html
<!-- In cloudflare-pages-frontend/index.html -->
<script>
  const API_BASE_URL = 'https://api.bhomanshah.com';
  
  // Example usage
  async function loadProducts() {
    const response = await fetch(`${API_BASE_URL}/api/products`);
    return response.json();
  }
</script>
```

---

## Scaling & Performance

### Upgrade Instance (Paid)

Free tier limitations:
- 512 MB RAM
- 0.1 vCPU
- 1 instance

Paid plans offer:
- Up to 8 GB RAM
- Up to 4 vCPUs
- Auto-scaling (1-20 instances)

### Optimize Performance

1. **Enable HTTP/2**
   - Automatically enabled on Koyeb

2. **Use Connection Pooling**
   - Already configured in database.go

3. **Cache Responses**
   - Already configured with Cache-Control headers

4. **Monitor Slow Queries**
   - Check database performance
   - Add indexes if needed

---

## Troubleshooting

### Build Failed

**Error:** `go.mod: no such file or directory`

**Solution:** 
- Ensure Dockerfile is in repository root
- Verify go.mod and go.sum are committed

### Deployment Failed

**Error:** `Failed to connect to database`

**Solution:**
- Verify DATABASE_URL is correct
- Check database is accessible from Koyeb IPs
- Ensure TLS is configured: `?tls=true`

### Migration Failed

**Error:** `Migration failed: duplicate key value`

**Solution:**
- Migrations are idempotent
- Check migrations table in database
- If needed, manually mark migrations as completed

### 502 Bad Gateway

**Error:** Cloudflare shows 502 error

**Solution:**
- Check Koyeb service is running
- Verify DNS CNAME is correct
- Ensure Cloudflare proxy is OFF for api subdomain

### CORS Errors

**Error:** `CORS policy: No 'Access-Control-Allow-Origin' header`

**Solution:**
- Verify CORS_ORIGIN matches frontend domain exactly
- Include protocol: `https://bhomanshah.com` (not bhomanshah.com)
- Check middleware is loaded in main.go

---

## Security Best Practices

- [ ] Use secrets for sensitive env vars
- [ ] Rotate ADMIN_SECRET periodically
- [ ] Enable Koyeb's GitHub integration for auto-deploy only from main
- [ ] Review deploy logs for security issues
- [ ] Monitor for unusual traffic patterns
- [ ] Keep dependencies updated
- [ ] Use HTTPS only (enforced by Koyeb)

---

## Cost Optimization

### Free Tier Limits
- 1 web service
- 512 MB RAM
- 2M requests/month
- 100 GB bandwidth/month

### When to Upgrade
- > 2M requests/month
- Need auto-scaling
- Require > 512 MB RAM
- Need multiple regions

### Paid Plans
- **Starter:** $7/month (1 GB RAM, 0.25 vCPU)
- **Standard:** $29/month (2 GB RAM, 1 vCPU)
- **Pro:** Custom pricing

---

## Maintenance

### Update Deployment

```bash
# Push to GitHub main branch
git push origin main

# Or manual redeploy
koyeb service redeploy oxlook-api
```

### Rollback

1. Dashboard > Deployments
2. Select previous successful deployment
3. Click "Redeploy from this version"

### Database Backup

```bash
# Before major updates
mysqldump -h <host> -u <user> -p <database> > backup-$(date +%Y%m%d).sql
```

---

## Next Steps

1. ✅ Service deployed and running
2. ✅ Custom domain configured
3. ✅ SSL certificate active
4. [ ] Deploy frontend to Cloudflare Pages
5. [ ] Test end-to-end flow
6. [ ] Set up monitoring
7. [ ] Configure alerts

---

## Support

### Koyeb Support
- Documentation: https://www.koyeb.com/docs
- Community: https://community.koyeb.com
- Email: support@koyeb.com

### Common Issues
- Check Koyeb status page
- Review deployment logs
- Verify environment variables
- Test database connectivity

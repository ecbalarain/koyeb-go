# Koyeb Deployment - Quick Start

## Production Secrets Generated

**Admin Secret:** `BaVx84uf0CI7FpueqiYhGEkg9kL1D1a1qahF79FZzZA=`

⚠️ **SAVE THIS SECRET SECURELY** - You'll need it for admin API access!

---

## Step 1: Push Code to GitHub

```bash
# Ensure all code is committed
git add .
git commit -m "Ready for production deployment"
git push origin main
```

---

## Step 2: Create Koyeb Account

1. Go to https://www.koyeb.com
2. Sign up (free tier available)
3. Verify your email

---

## Step 3: Deploy Service

### In Koyeb Dashboard:

1. **Click "Create Service"**

2. **Deployment Method:**
   - Select: **GitHub**
   - Connect your GitHub account
   - Select repository: `koyeb-go`
   - Branch: `main`

3. **Builder:**
   - Select: **Dockerfile**
   - Path: `./Dockerfile` (auto-detected)

4. **Instance:**
   - Select: **Free** (or Nano for production)
   - Region: **Washington D.C.** or **Frankfurt** (closest to TiDB)

5. **Service Name:**
   - Name: `oxlook-api`

6. **Environment Variables (Secrets):**
   
   Click "Add Environment Variable" for each:
   
   ```
   Name: DATABASE_URL
   Value: mysql://3QQuRDHAFitYTwM.root:HWvBROQOOxQed5c1@gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000/test?tls=true
   Type: SECRET
   
   Name: ADMIN_SECRET
   Value: BaVx84uf0CI7FpueqiYhGEkg9kL1D1a1qahF79FZzZA=
   Type: SECRET
   
   Name: CORS_ORIGIN
   Value: https://bhomanshah.com
   Type: SECRET
   
   Name: ENVIRONMENT
   Value: production
   Type: PLAIN
   
   Name: PORT
   Value: 8080
   Type: PLAIN
   ```

7. **Health Check:**
   - Path: `/health`
   - Port: `8080`
   - Protocol: `HTTP`

8. **Click "Deploy"**

---

## Step 4: Monitor Deployment

Watch the build logs:
1. Click on your service in dashboard
2. Navigate to "Deployments" tab
3. Click the active deployment
4. View "Build logs" and "Runtime logs"

**Expected logs:**
```
Building image...
Running migrations...
✅ All migrations completed successfully!
Starting server on :8080
```

**Deployment takes ~5-10 minutes**

---

## Step 5: Test Your API

Once deployed, you'll get a URL like:
`https://oxlook-api-<random>.koyeb.app`

### Test endpoints:

```bash
# Replace <your-url> with your actual Koyeb URL

# Health check
curl https://<your-url>.koyeb.app/health

# Products
curl https://<your-url>.koyeb.app/api/products

# Variants
curl https://<your-url>.koyeb.app/api/products/stone-mug/variants

# Admin (requires API key)
curl -H "X-API-Key: BaVx84uf0CI7FpueqiYhGEkg9kL1D1a1qahF79FZzZA=" \
  https://<your-url>.koyeb.app/admin/api/products
```

---

## Step 6: Configure Custom Domain

### In Koyeb Dashboard:

1. Go to your service settings
2. Click "Domains" tab
3. Click "Add Domain"
4. Enter: `api.bhomanshah.com`
5. Koyeb will show DNS instructions

### In Cloudflare DNS:

1. Go to Cloudflare dashboard
2. Select your domain: `bhomanshah.com`
3. Go to DNS settings
4. Add CNAME record:
   - **Type:** `CNAME`
   - **Name:** `api`
   - **Target:** `<your-service>.koyeb.app` (provided by Koyeb)
   - **Proxy status:** ⚪ **DNS only** (gray cloud, NOT proxied)
   - **TTL:** Auto

5. Save

### Wait for SSL Certificate

- Takes 2-10 minutes
- Check status in Koyeb "Domains" tab
- When status shows "Active" ✅, your custom domain is ready

### Test custom domain:

```bash
curl https://api.bhomanshah.com/health
```

---

## Step 7: Update Frontend

Update frontend to use production API:

In `cloudflare-pages-frontend/index.html`:
```javascript
const API_BASE_URL = 'https://api.bhomanshah.com';
```

In `cloudflare-pages-frontend/admin/index.html` and other admin files:
```javascript
const API_BASE_URL = 'https://api.bhomanshah.com';
const ADMIN_SECRET = 'BaVx84uf0CI7FpueqiYhGEkg9kL1D1a1qahF79FZzZA=';
```

---

## Step 8: Deploy Frontend to Cloudflare Pages

```bash
# Build frontend
./build-frontend.sh
```

### In Cloudflare Dashboard:

1. Go to Pages
2. Click "Create a project"
3. Connect to Git (select your repo)
4. Build settings:
   - **Build command:** `./build-frontend.sh`
   - **Build output directory:** `cloudflare-pages-frontend/dist`
5. Click "Save and Deploy"

### Configure custom domain:

1. In Cloudflare Pages > Custom domains
2. Add: `bhomanshah.com`
3. DNS records will be auto-configured

---

## Troubleshooting

### Build fails
- Check Dockerfile syntax
- Verify all files are committed to Git
- Review build logs in Koyeb

### Migrations fail
- Check DATABASE_URL is correct
- Ensure TiDB allows connections from Koyeb
- Check migration logs

### CORS errors
- Verify CORS_ORIGIN matches frontend URL exactly
- Include `https://` protocol
- Redeploy after changing env vars

### Domain SSL not provisioning
- Ensure DNS record is correct (CNAME pointing to Koyeb)
- Cloudflare proxy must be OFF (gray cloud)
- Wait 10 minutes for propagation

---

## Next Steps

1. ✅ Service deployed on Koyeb
2. ✅ Custom domain configured
3. ✅ Frontend deployed on Cloudflare Pages
4. 📋 [Test production](PRODUCTION_TESTING.md)
5. 📊 [Set up monitoring](MONITORING.md)
6. 🚀 Launch!

---

## Important URLs

- **API:** https://api.bhomanshah.com
- **Frontend:** https://bhomanshah.com
- **Koyeb Dashboard:** https://app.koyeb.com
- **Cloudflare Dashboard:** https://dash.cloudflare.com

## Important Credentials

- **Admin Secret:** `BaVx84uf0CI7FpueqiYhGEkg9kL1D1a1qahF79FZzZA=`
- **Database:** TiDB Cloud (already configured)

⚠️ **Keep your admin secret safe!** Store it in a password manager.

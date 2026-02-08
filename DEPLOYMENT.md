# Cloudflare Pages Deployment Guide

## Prerequisites
- Cloudflare account
- GitHub repository connected to Cloudflare Pages
- Custom domain configured in Cloudflare DNS

## Deployment Steps

### 1. Build the Frontend

```bash
cd /path/to/koyeb-go
./build-frontend.sh
```

This creates an optimized build in `cloudflare-pages-frontend/dist/`

### 2. Configure Cloudflare Pages Project

1. Log in to Cloudflare Dashboard
2. Go to Pages
3. Create a new project or select existing
4. Connect to your GitHub repository
5. Configure build settings:
   - **Build command:** `./build-frontend.sh`
   - **Build output directory:** `cloudflare-pages-frontend/dist`
   - **Root directory:** `/`

### 3. Environment Variables

No environment variables needed for the static frontend.

### 4. Custom Domain

1. In Cloudflare Pages project settings
2. Go to "Custom domains"
3. Add `bhomanshah.com` and `www.bhomanshah.com`
4. DNS records will be automatically configured

### 5. Optimize Cloudflare Settings

#### Speed Optimizations
Go to Speed > Optimization and enable:
- [x] Auto Minify: JavaScript, CSS, HTML
- [x] Brotli compression
- [x] Early Hints
- [x] HTTP/2 to Origin
- [x] HTTP/3 (with QUIC)
- [x] 0-RTT Connection Resumption
- [x] WebSockets

#### Image Optimizations
- [x] Polish: Lossless or Lossy
- [x] WebP conversion
- [x] Mirage (lazy loading)

#### Caching
1. Caching > Configuration
2. Browser Cache TTL: Respect Existing Headers
3. Enable "Always Online"

### 6. SSL/TLS Configuration

1. SSL/TLS > Overview
2. Select "Full (strict)" encryption mode
3. Edge Certificates > Always Use HTTPS: On
4. Minimum TLS Version: TLS 1.2
5. Opportunistic Encryption: On
6. TLS 1.3: On
7. Automatic HTTPS Rewrites: On

### 7. Security Headers

Headers are configured in `_headers` file and automatically deployed.

### 8. Test Deployment

After deployment:
1. Visit https://bhomanshah.com
2. Check browser DevTools:
   - Network tab: Verify caching headers
   - Console: No errors
3. Test on mobile device
4. Run Lighthouse audit

### 9. Verify Cache Headers

```bash
# Check variant endpoint
curl -I https://api.bhomanshah.com/api/products/stone-mug/variants

# Should see:
# Cache-Control: public, max-age=1800, s-maxage=1800
# Vary: Accept-Encoding
```

### 10. Monitor Performance

#### Cloudflare Analytics
1. Analytics > Web Analytics
2. Add site if not already configured
3. Monitor:
   - Page load time
   - Core Web Vitals (LCP, FID, CLS)
   - Browser insights

#### Real User Monitoring
1. Add Cloudflare Web Analytics snippet to `<head>`:

```html
<script defer src='https://static.cloudflareinsights.com/beacon.min.js' 
        data-cf-beacon='{"token": "YOUR_TOKEN"}'></script>
```

## Troubleshooting

### Cache Not Working
1. Check `_headers` file is deployed
2. Purge cache: Caching > Configuration > Purge Everything
3. Verify in browser DevTools Network tab

### Images Not Optimized
1. Enable Polish in Speed > Optimization
2. Check image URLs are relative (not external)
3. Verify WebP conversion is enabled

### Slow Load Times
1. Run Lighthouse audit
2. Check Cloudflare Analytics for bottlenecks
3. Verify Brotli compression is working
4. Check if Rocket Loader™ helps (test carefully)

### Security Headers Missing
1. Verify `_headers` file is in build output
2. Check file format (no BOM, Unix line endings)
3. Test with: `curl -I https://bhomanshah.com`

## Rollback Procedure

If issues occur after deployment:

1. Cloudflare Pages > Deployments
2. Find the last working deployment
3. Click "..." > "Rollback to this deployment"
4. Verify site is working

## Continuous Deployment

Every push to main branch triggers:
1. Automatic build via `build-frontend.sh`
2. Deploy to production
3. Cloudflare Pages generates preview URL

For staging/testing:
- Create a branch: `staging`
- Configure separate Cloudflare Pages project
- Test before merging to `main`

## Cache Purging

### Via Cloudflare Dashboard
1. Caching > Configuration
2. Purge Cache > Custom Purge
3. Enter URLs or purge everything

### Via API (Admin Endpoint)
```bash
curl -X POST https://api.bhomanshah.com/admin/api/cache/purge \
  -H "X-API-Key: YOUR_ADMIN_SECRET" \
  -H "Content-Type: application/json" \
  -d '{"type": "all"}'
```

## Post-Deployment Checklist

- [ ] Site loads at https://bhomanshah.com
- [ ] API calls work from frontend
- [ ] CORS configured correctly
- [ ] SSL certificate active (green padlock)
- [ ] Cache headers present
- [ ] Images loading properly
- [ ] Mobile responsive
- [ ] Lighthouse score > 90
- [ ] No console errors
- [ ] Cart functionality works
- [ ] Checkout flow works
- [ ] Admin panel accessible

## Performance Targets

- **Lighthouse Performance:** > 90
- **First Contentful Paint:** < 1.8s
- **Largest Contentful Paint:** < 2.5s
- **Time to Interactive:** < 3.8s
- **Total Blocking Time:** < 300ms
- **Cumulative Layout Shift:** < 0.1

## Support

For issues:
1. Check Cloudflare Logs
2. Review GitHub Actions logs (if applicable)
3. Check browser console for errors
4. Test with different browsers/devices

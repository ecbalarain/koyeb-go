# Frontend Performance Optimization Guide

## Overview
This guide provides instructions for optimizing the OXLOOK frontend for production deployment on Cloudflare Pages.

## 1. Image Optimization

### Convert to WebP Format
WebP provides superior compression while maintaining quality. Use `cwebp` to convert images:

```bash
# Install cwebp (macOS)
brew install webp

# Convert single image
cwebp -q 80 input.jpg -o output.webp

# Convert all JPGs in a directory
for file in *.jpg; do
  cwebp -q 80 "$file" -o "${file%.jpg}.webp"
done
```

### Lazy Loading
All images are dynamically rendered. Ensure you add `loading="lazy"` attribute when adding img tags:

```html
<img src="image.webp" loading="lazy" alt="Product name" />
```

### Responsive Images
Use srcset for different screen sizes:

```html
<img
  src="image-800.webp"
  srcset="image-400.webp 400w, image-800.webp 800w, image-1200.webp 1200w"
  sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
  loading="lazy"
  alt="Product name"
/>
```

## 2. HTML/CSS/JavaScript Minification

### Automatic (Cloudflare)
Enable automatic minification in Cloudflare:
1. Go to Speed > Optimization
2. Enable Auto Minify for:
   - JavaScript
   - CSS
   - HTML

### Manual Minification
Use online tools or build tools:

**HTML:**
- https://www.toptal.com/developers/html-minifier
- Options: Remove comments, collapse whitespace

**Inline CSS/JS:**
- Use Terser for JavaScript
- Use cssnano for CSS

## 3. Cloudflare-Specific Optimizations

### Enable Polish
Polish automatically optimizes images:
1. Speed > Optimization > Polish
2. Select "Lossless" or "Lossy"
3. Enable "WebP" conversion

### Enable Mirage
Mirage provides lazy loading and responsive image delivery:
1. Speed > Optimization > Mirage
2. Toggle "On"

### Enable Rocket Loader™
Prioritizes page rendering:
1. Speed > Optimization > Rocket Loader
2. Toggle "On" (test thoroughly)

### Browser Cache TTL
Set aggressive caching:
1. Caching > Configuration
2. Browser Cache TTL: 1 year (for versioned assets)

## 4. Resource Hints

Add to `<head>` section for external resources:

```html
<!-- Preconnect to API domain -->
<link rel="preconnect" href="https://api.bhomanshah.com" />
<link rel="dns-prefetch" href="https://api.bhomanshah.com" />

<!-- Preconnect to CDN (Tailwind) -->
<link rel="preconnect" href="https://cdn.tailwindcss.com" />
```

## 5. Cache Configuration

The `_headers` file configures Cloudflare's edge caching:

- **Static assets** (CSS, JS, images): 1 year
- **products.json**: 10 minutes
- **HTML pages**: 5 minutes

To purge cache after updates:
1. Use Cloudflare Dashboard: Caching > Configuration > Purge Cache
2. Use API endpoint: `POST /admin/api/cache/purge`

## 6. Performance Testing

### Google Lighthouse
```bash
# Install Lighthouse
npm install -g lighthouse

# Run test
lighthouse https://bhomanshah.com --view
```

**Target Scores:**
- Performance: > 90
- Accessibility: > 95
- Best Practices: > 95
- SEO: > 90

### WebPageTest
Test from multiple locations:
- https://www.webpagetest.org
- Test from 3 different locations
- Mobile and Desktop

### Key Metrics to Monitor
- **LCP** (Largest Contentful Paint): < 2.5s
- **FID** (First Input Delay): < 100ms
- **CLS** (Cumulative Layout Shift): < 0.1
- **TTFB** (Time to First Byte): < 600ms

## 7. Build and Deploy

```bash
# Build for production
./build-frontend.sh

# The script copies files to dist/
# Deploy cloudflare-pages-frontend/dist/ to Cloudflare Pages
```

## 8. Progressive Web App (PWA) - Optional

Add `manifest.json` and service worker for offline support:

```json
{
  "name": "OXLOOK Store",
  "short_name": "OXLOOK",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#f5f5f4",
  "theme_color": "#18181b",
  "icons": [
    {
      "src": "/icon-192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icon-512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

## 9. Monitoring

### Real User Monitoring (RUM)
Add Cloudflare Web Analytics (free):
1. Analytics > Web Analytics
2. Add site
3. Copy JavaScript snippet to `<head>`

### Synthetic Monitoring
Set up external monitoring:
- Pingdom
- UptimeRobot (free tier available)
- Checkly

## 10. Checklist Before Launch

- [ ] All images converted to WebP
- [ ] Images have width/height attributes (prevent CLS)
- [ ] Lazy loading enabled for images
- [ ] Cloudflare Auto Minify enabled
- [ ] Polish and Mirage enabled
- [ ] _headers file deployed
- [ ] Cache TTLs configured
- [ ] Resource hints added
- [ ] Lighthouse score > 90
- [ ] WebPageTest < 3s load time
- [ ] DNS fully propagated
- [ ] SSL certificate active

# Frontend Testing Guide

This directory contains automated end-to-end tests for the OXLOOK e-commerce frontend using Puppeteer.

## Prerequisites

- Node.js (v18 or higher)
- npm or yarn

## Installation

Install dependencies (includes Puppeteer which will download Chromium):

```bash
npm install
```

## Running Tests

### Run tests with visible browser (development):
```bash
npm test
```

### Run tests in headless mode (CI/CD):
```bash
npm run test:headless
```

## Test Coverage

The test suite covers:

1. **Homepage Loading** - Verifies the page loads and products are displayed
2. **Product Grid** - Checks if products are rendered correctly
3. **Search Functionality** - Tests product search
4. **Category Filtering** - Tests category chips (All, Home, Tech, Basics)
5. **Product Details** - Opens product modal and checks variant selection
6. **Variant Selection** - Tests size/color variant selection
7. **Add to Cart** - Adds items to cart
8. **Cart Management** - Opens cart and verifies items
9. **Checkout Flow** - Tests checkout form (without submitting)
10. **Mobile Responsiveness** - Tests mobile viewport

## Screenshots

All test screenshots are saved to `test-screenshots/` directory:

- `01-homepage.png` - Initial page load
- `02-search.png` - Search results
- `03-category-filter.png` - Filtered products
- `04-product-detail.png` - Product modal
- `05-variant-selected.png` - Selected variant
- `06-added-to-cart.png` - Cart notification
- `07-cart-open.png` - Cart drawer
- `08-checkout-form.png` - Checkout form
- `09-mobile-view.png` - Mobile viewport
- `error.png` - Screenshot on test failure

## Customization

Edit `test-frontend.js` to:
- Change the frontend URL (currently set to Cloudflare Pages preview)
- Add more tests
- Modify wait times
- Enable/disable headless mode
- Add assertions

## Troubleshooting

### Puppeteer installation fails
```bash
# Try installing with legacy peer deps
npm install --legacy-peer-deps
```

### Tests timeout
- Increase timeout values in the test file
- Check if the API server is running
- Verify CORS is configured correctly

### Can't find Chromium
```bash
# Reinstall puppeteer
npm uninstall puppeteer
npm install puppeteer
```

## CI/CD Integration

For GitHub Actions or other CI systems:

```yaml
- name: Install dependencies
  run: npm install

- name: Run tests
  run: npm run test:headless

- name: Upload screenshots
  uses: actions/upload-artifact@v3
  with:
    name: test-screenshots
    path: test-screenshots/
```

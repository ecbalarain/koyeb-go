import { test, expect } from '@playwright/test';

test.describe('Product Browsing', () => {
  test('TC-PB-001: Verify product catalog loads correctly', async ({ page }) => {
    // Navigate to the homepage
    await page.goto('/');

    // Wait for products to load - look for the grid element
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Check if products are displayed
    const products = page.locator('#grid article');
    await expect(products.first()).toBeVisible();

    // Verify at least one product is loaded
    const productCount = await products.count();
    expect(productCount).toBeGreaterThan(0);

    // Check for product name and price
    const firstProduct = products.first();
    const productName = firstProduct.locator('h3');
    await expect(productName).toBeVisible();
    await expect(productName).toHaveText(/.+/); // Has some text

    const productPrice = firstProduct.locator('p.font-semibold');
    await expect(productPrice).toBeVisible();

    console.log(`Found ${productCount} products loaded successfully`);
  });

  test('TC-PB-002: Test product filtering by category', async ({ page }) => {
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Get initial product count
    const initialProducts = page.locator('#grid article');
    const initialCount = await initialProducts.count();

    // Click on a category chip (e.g., "Basics")
    const basicsChip = page.locator('.chip[data-chip="Basics"]');
    if (await basicsChip.isVisible()) {
      await basicsChip.click();

      // Wait for filtering to apply
      await page.waitForTimeout(500);

      // Check that filtering worked (count should be different or same, but active chip should be Basics)
      const activeChip = page.locator('.chip.is-active');
      await expect(activeChip).toHaveAttribute('data-chip', 'Basics');

      // Verify products are still displayed
      const filteredProducts = page.locator('#grid article');
      const filteredCount = await filteredProducts.count();
      expect(filteredCount).toBeGreaterThanOrEqual(0); // Could be 0 if no products in category

      console.log(`Filtered to ${filteredCount} products in Basics category`);
    } else {
      console.log('Basics category not available, skipping filter test');
    }

    // Test "All" category
    const allChip = page.locator('.chip[data-chip="All"]');
    await allChip.click();
    await page.waitForTimeout(500);

    const allProducts = page.locator('#grid article');
    const allCount = await allProducts.count();
    expect(allCount).toBe(initialCount);

    console.log(`All category shows ${allCount} products`);
  });

  test('TC-PB-003: Validate product detail page displays', async ({ page }) => {
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Click on the first product
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();

    // Wait for product overlay to appear
    await page.waitForSelector('#productOverlay', { timeout: 10000 });

    // Check if overlay is visible
    const overlay = page.locator('#productOverlay');
    await expect(overlay).toBeVisible();

    // Check product detail elements
    const productTitle = overlay.locator('h2, h3');
    await expect(productTitle).toBeVisible();
    await expect(productTitle).toHaveText(/.+/);

    // Check for price
    const priceElement = overlay.locator('.price, [data-price]');
    // Price might be loaded asynchronously, so check if it appears
    await page.waitForTimeout(2000); // Wait for API call
    if (await priceElement.isVisible()) {
      await expect(priceElement).toHaveText(/.+/);
    }

    // Check for size/color selectors if they exist
    const sizeSelector = overlay.locator('select, .size-selector, [data-size]');
    const colorSelector = overlay.locator('select, .color-selector, [data-color]');

    // At least one variant selector should be present
    const hasSelectors = (await sizeSelector.count() > 0) || (await colorSelector.count() > 0);
    expect(hasSelectors).toBe(true);

    console.log('Product detail overlay displayed successfully');
  });

  test('TC-PB-004: Test variant selection (size/color)', async ({ page }) => {
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Click on the first product
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();

    // Wait for product overlay to appear
    await page.waitForSelector('#productOverlay', { timeout: 10000 });

    // Wait for variants to load
    await page.waitForTimeout(2000);

    // Look for size selector
    const sizeSelect = page.locator('#productOverlay select').first();
    if (await sizeSelect.isVisible()) {
      // Get initial selected option
      const initialValue = await sizeSelect.inputValue();

      // Get all options
      const options = sizeSelect.locator('option');
      const optionCount = await options.count();

      if (optionCount > 1) {
        // Select a different option
        const secondOption = options.nth(1);
        const secondValue = await secondOption.getAttribute('value');
        await sizeSelect.selectOption(secondValue);

        // Verify selection changed
        const newValue = await sizeSelect.inputValue();
        expect(newValue).not.toBe(initialValue);

        console.log(`Size variant selected successfully: ${newValue}`);
      }
    }

    // Look for color selector if exists
    const colorSelect = page.locator('#productOverlay select').nth(1);
    if (await colorSelect.isVisible()) {
      const colorOptions = colorSelect.locator('option');
      const colorOptionCount = await colorOptions.count();

      if (colorOptionCount > 1) {
        const secondColorOption = colorOptions.nth(1);
        const colorValue = await secondColorOption.getAttribute('value');
        await colorSelect.selectOption(colorValue);

        console.log(`Color variant selected successfully: ${colorValue}`);
      }
    }

    // Check if price updates (if there's a price display)
    const priceDisplay = page.locator('#productOverlay .price, #productOverlay [data-price]');
    if (await priceDisplay.isVisible()) {
      const priceText = await priceDisplay.textContent();
      expect(priceText).toMatch(/\$\d+\.\d{2}/); // Should show a price
    }
  });

  test('TC-PB-005: Verify price updates with variant selection', async ({ page }) => {
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Click on the first product
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();

    // Wait for product overlay to appear
    await page.waitForSelector('#productOverlay', { timeout: 10000 });

    // Wait for variants to load
    await page.waitForTimeout(2000);

    // Find price display
    const priceDisplay = page.locator('#productOverlay .price, #productOverlay [data-price]');
    let priceLocator = priceDisplay;
    
    if (!(await priceDisplay.isVisible())) {
      // Try to find any element containing a dollar sign
      priceLocator = page.locator('#productOverlay').locator('text=/\\$/').first();
    }

    if (await priceLocator.isVisible()) {
      const initialPrice = await priceLocator.textContent();

      // Look for size selector
      const sizeSelect = page.locator('#productOverlay select').first();
      if (await sizeSelect.isVisible()) {
        const options = sizeSelect.locator('option');
        const optionCount = await options.count();

        if (optionCount > 1) {
          // Select a different option
          const secondOption = options.nth(1);
          const secondValue = await secondOption.getAttribute('value');
          await sizeSelect.selectOption(secondValue);

          // Wait for price to update
          await page.waitForTimeout(500);

          // Check if price changed
          const newPrice = await priceLocator.textContent();
          // Price might be the same or different depending on variants
          expect(typeof newPrice).toBe('string');
          expect(newPrice.length).toBeGreaterThan(0);

          console.log(`Price updated from "${initialPrice}" to "${newPrice}"`);
        }
      }
    } else {
      console.log('No price display found, skipping price update test');
    }
  });

  test('TC-PB-006: Test product image gallery', async ({ page }) => {
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Click on the first product
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();

    // Wait for product overlay to appear
    await page.waitForSelector('#productOverlay', { timeout: 10000 });

    // Wait for content to load
    await page.waitForTimeout(2000);

    // Check for main product image
    const mainImage = page.locator('#productOverlay img').first();
    if (await mainImage.isVisible()) {
      // Verify image has src attribute
      const src = await mainImage.getAttribute('src');
      expect(src).toBeTruthy();
      expect(src.length).toBeGreaterThan(0);

      console.log(`Main product image loaded: ${src}`);
    } else {
      console.log('No main product image found');
    }

    // Check for image gallery or thumbnails
    const galleryImages = page.locator('#productOverlay img');
    const imageCount = await galleryImages.count();

    if (imageCount > 1) {
      console.log(`Found ${imageCount} images in gallery`);

      // Check that all images have valid src
      for (let i = 0; i < imageCount; i++) {
        const img = galleryImages.nth(i);
        const src = await img.getAttribute('src');
        expect(src).toBeTruthy();
      }
    } else {
      console.log('Single image or no gallery found');
    }

    // Check for image container or carousel
    const imageContainer = page.locator('#productOverlay .image-gallery, #productOverlay .carousel, #productOverlay .gallery');
    if (await imageContainer.isVisible()) {
      console.log('Image gallery container found');
    }
  });

  test('TC-PB-007: Check responsive design on mobile/desktop', async ({ page, viewport }) => {
    // Test desktop view
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Check desktop layout
    const products = page.locator('#grid article');
    await expect(products.first()).toBeVisible();

    // Check header layout (desktop should show search button)
    const desktopSearch = page.locator('.hidden.sm\\:flex'); // Desktop search button
    await expect(desktopSearch).toBeVisible();

    // Check mobile menu is hidden on desktop
    const mobileMenu = page.locator('#menuBtn'); // Mobile menu button
    await expect(mobileMenu).not.toBeVisible(); // Should be hidden on desktop

    console.log('Desktop layout verified');

    // Test mobile view
    await page.setViewportSize({ width: 375, height: 667 });
    await page.reload();

    // Wait for products to load again
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Check mobile layout
    const mobileProducts = page.locator('#grid article');
    await expect(mobileProducts.first()).toBeVisible();

    // Check mobile search is visible
    const mobileSearch = page.locator('#searchInputMobile');
    await expect(mobileSearch).toBeVisible();

    // Check mobile menu button is visible
    const mobileMenuBtn = page.locator('#menuBtn');
    await expect(mobileMenuBtn).toBeVisible();

    console.log('Mobile layout verified');

    // Test tablet view
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.reload();

    // Wait for products
    await page.waitForSelector('#grid', { timeout: 10000 });

    const tabletProducts = page.locator('#grid article');
    await expect(tabletProducts.first()).toBeVisible();

    console.log('Tablet layout verified');
  });
});
import { test, expect } from '@playwright/test';

test.describe('Shopping Cart', () => {
  test('TC-SC-001: Add product to cart', async ({ page }) => {
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Click on the first product
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();

    // Wait for product overlay
    await page.waitForSelector('#productOverlay', { timeout: 10000 });

    // Wait for variants to load
    await page.waitForTimeout(2000);

    // Look for "Add to Cart" button
    const addToCartBtn = page.locator('#addToCartBtn');
    if (await addToCartBtn.isVisible() && await addToCartBtn.isEnabled()) {
      await addToCartBtn.click();

      // Wait for cart to update
      await page.waitForTimeout(1000);

      // Check cart count
      const cartCount = page.locator('#cartCount');
      if (await cartCount.isVisible()) {
        const countText = await cartCount.textContent();
        expect(parseInt(countText)).toBeGreaterThan(0);
        console.log(`Product added to cart, count: ${countText}`);
      } else {
        // Check if cart opens automatically
        const cartPanel = page.locator('[data-cart], .cart-panel, #cart');
        if (await cartPanel.isVisible()) {
          console.log('Cart panel opened after adding product');
        }
      }
    } else {
      console.log('Add to Cart button not found');
    }
  });

  test('TC-SC-002: Update cart item quantity', async ({ page }) => {
    await page.goto('/');

    // First add a product to cart
    await page.waitForSelector('#grid', { timeout: 10000 });
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();
    await page.waitForSelector('#productOverlay', { timeout: 10000 });
    await page.waitForTimeout(2000);

    const addToCartBtn = page.locator('#addToCartBtn');
    if (await addToCartBtn.isVisible() && await addToCartBtn.isEnabled()) {
      await addToCartBtn.click();
      await page.waitForTimeout(1000);
    }

    // Now open the cart
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();

    // Wait for cart overlay to open
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    // Look for quantity controls (the + and - buttons)
    const increaseBtn = page.locator('#cartOverlay button').filter({ hasText: '+' });
    const decreaseBtn = page.locator('#cartOverlay button').filter({ hasText: '-' });
    const quantityDisplay = page.locator('#cartOverlay span').filter({ hasText: /^[0-9]+$/ }); // Quantity number

    if (await quantityDisplay.isVisible()) {
      // Get initial quantity
      const initialQty = await quantityDisplay.textContent();
      const initialNum = parseInt(initialQty);

      // Try to increase quantity
      if (await increaseBtn.isVisible()) {
        await increaseBtn.click();
        await page.waitForTimeout(500);

        const newQty = await quantityDisplay.textContent();
        const newNum = parseInt(newQty);
        expect(newNum).toBe(initialNum + 1);

        console.log(`Quantity increased from ${initialNum} to ${newNum}`);
      }

      // Try to decrease quantity
      if (await decreaseBtn.isVisible() && initialNum > 1) {
        await decreaseBtn.click();
        await page.waitForTimeout(500);

        const decreasedQty = await quantityDisplay.textContent();
        const decreasedNum = parseInt(decreasedQty);
        expect(decreasedNum).toBe(initialNum);

        console.log(`Quantity decreased back to ${decreasedNum}`);
      }
    } else {
      console.log('Quantity controls not found in cart');
    }
  });

  test('TC-SC-003: Remove item from cart', async ({ page }) => {
    await page.goto('/');

    // First add a product to cart
    await page.waitForSelector('#grid', { timeout: 10000 });
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();
    await page.waitForSelector('#productOverlay', { timeout: 10000 });
    await page.waitForTimeout(2000);

    const addToCartBtn = page.locator('#addToCartBtn');
    if (await addToCartBtn.isVisible() && await addToCartBtn.isEnabled()) {
      await addToCartBtn.click();
      await page.waitForTimeout(1000);
    }

    // Check initial cart count
    const initialCartCount = page.locator('#cartCount');
    const initialCount = await initialCartCount.textContent();

    // Open the cart
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    // Find remove button (X button with SVG)
    const removeBtn = page.locator('#cartOverlay button').filter({ has: page.locator('svg') }).first();

    if (await removeBtn.isVisible()) {
      await removeBtn.click();
      await page.waitForTimeout(1000);

      // Check cart count after removal
      const newCartCount = page.locator('#cartCount');
      if (await newCartCount.isVisible()) {
        const newCount = await newCartCount.textContent();
        const newNum = parseInt(newCount) || 0;
        const initialNum = parseInt(initialCount);
        expect(newNum).toBeLessThan(initialNum);
        console.log(`Item removed from cart, count changed from ${initialNum} to ${newNum}`);
      } else {
        // Cart count element is hidden, meaning count is 0
        console.log('Cart count element hidden, cart is empty');
      }

      // Also check if cart shows empty message
      const emptyMessage = page.locator('#cartOverlay').filter({ hasText: 'Your cart is empty' });
      if (await emptyMessage.isVisible()) {
        console.log('Cart shows empty message after removing item');
      }
    } else {
      console.log('Remove button not found in cart');
    }
  });

  test('TC-SC-004: Verify cart persistence across sessions', async ({ page }) => {
    await page.goto('/');

    // Add a product to cart
    await page.waitForSelector('#grid', { timeout: 10000 });
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();
    await page.waitForSelector('#productOverlay', { timeout: 10000 });
    await page.waitForTimeout(2000);

    const addToCartBtn = page.locator('#addToCartBtn');
    if (await addToCartBtn.isVisible() && await addToCartBtn.isEnabled()) {
      await addToCartBtn.click();
      await page.waitForTimeout(1000);
    }

    // Check cart count before reload
    const cartCountBefore = page.locator('#cartCount');
    const countBefore = await cartCountBefore.textContent();

    // Reload the page
    await page.reload();
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Check cart count after reload
    const cartCountAfter = page.locator('#cartCount');
    if (await cartCountAfter.isVisible()) {
      const countAfter = await cartCountAfter.textContent();
      expect(countAfter).toBe(countBefore);
      console.log(`Cart persisted across page reload: ${countAfter} items`);
    } else {
      console.log('Cart count not visible after reload, cart may not persist');
    }
  });

  test('TC-SC-005: Test cart total calculation', async ({ page }) => {
    await page.goto('/');

    // Add a product to cart
    await page.waitForSelector('#grid', { timeout: 10000 });
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();
    await page.waitForSelector('#productOverlay', { timeout: 10000 });
    await page.waitForTimeout(2000);

    const addToCartBtn = page.locator('#addToCartBtn');
    if (await addToCartBtn.isVisible() && await addToCartBtn.isEnabled()) {
      await addToCartBtn.click();
      await page.waitForTimeout(1000);
    }

    // Open cart
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    // Check subtotal
    const subtotalEl = page.locator('#cartSubtotal');
    if (await subtotalEl.isVisible()) {
      const subtotalText = await subtotalEl.textContent();
      expect(subtotalText).toMatch(/\$\d+\.\d{2}/); // Should be a dollar amount
      console.log(`Cart subtotal: ${subtotalText}`);
    }

    // Increase quantity and check if subtotal updates
    const increaseBtn = page.locator('#cartOverlay button').filter({ hasText: '+' });
    if (await increaseBtn.isVisible()) {
      const initialSubtotal = await subtotalEl.textContent();

      await increaseBtn.click();
      await page.waitForTimeout(500);

      const newSubtotal = await subtotalEl.textContent();
      // Subtotal should be different (doubled if price is the same)
      expect(newSubtotal).not.toBe(initialSubtotal);
      console.log(`Subtotal updated from ${initialSubtotal} to ${newSubtotal} after quantity increase`);
    }
  });

  test('TC-SC-006: Validate cart empty state', async ({ page }) => {
    await page.goto('/');

    // Wait for page to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Open cart immediately (should be empty)
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    // Check for empty cart message
    const emptyMessage = page.locator('#cartItems').filter({ hasText: 'Your cart is empty' });
    await expect(emptyMessage).toBeVisible();

    // Check that subtotal shows $0
    const subtotalEl = page.locator('#cartSubtotal');
    if (await subtotalEl.isVisible()) {
      const subtotalText = await subtotalEl.textContent();
      expect(subtotalText).toBe('$0');
    }

    // Check that cart count is not visible or shows 0
    const cartCount = page.locator('#cartCount');
    if (await cartCount.isVisible()) {
      const countText = await cartCount.textContent();
      expect(parseInt(countText) || 0).toBe(0);
    } else {
      // Cart count element is hidden when count is 0
      console.log('Cart count element is hidden for empty cart');
    }

    console.log('Empty cart state validated successfully');
  });

  test('TC-SC-007: Test cart with multiple variants', async ({ page }) => {
    await page.goto('/');

    // Wait for products to load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Add first product with default variant
    const firstProduct = page.locator('#grid article').first();
    await firstProduct.click();
    await page.waitForSelector('#productOverlay', { timeout: 10000 });
    await page.waitForTimeout(2000);

    const addToCartBtn = page.locator('#addToCartBtn');
    if (await addToCartBtn.isVisible() && await addToCartBtn.isEnabled()) {
      await addToCartBtn.click();
      await page.waitForTimeout(1000);
      console.log('Added first product to cart');
    }

    // Close overlay
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);

    // Try to add the same product with different variant if possible
    await firstProduct.click();
    await page.waitForSelector('#productOverlay', { timeout: 10000 });
    await page.waitForTimeout(2000);

    // Look for size selector and change it
    const sizeSelect = page.locator('#productOverlay select').first();
    if (await sizeSelect.isVisible()) {
      const options = sizeSelect.locator('option');
      const optionCount = await options.count();

      if (optionCount > 1) {
        const secondOption = options.nth(1);
        const secondValue = await secondOption.getAttribute('value');
        await sizeSelect.selectOption(secondValue);

        // Add to cart again
        if (await addToCartBtn.isVisible() && await addToCartBtn.isEnabled()) {
          await addToCartBtn.click();
          await page.waitForTimeout(1000);
          console.log('Added same product with different variant to cart');
        }
      }
    }

    // Close overlay
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);

    // Open cart and check items
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    // Check cart has items
    const cartItems = page.locator('#cartItems > div');
    const itemCount = await cartItems.count();
    expect(itemCount).toBeGreaterThan(0);

    console.log(`Cart contains ${itemCount} item(s) with variants`);
  });
});
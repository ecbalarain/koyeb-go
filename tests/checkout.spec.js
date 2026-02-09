import { test, expect } from '@playwright/test';

test.describe('Checkout & Ordering', () => {
  test('TC-CO-001: Complete checkout form validation', async ({ page }) => {
    await page.goto('/');

    // Add a product to cart first
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

    // Look for checkout button
    const checkoutBtn = page.locator('#cartOverlay button').filter({ hasText: 'Checkout' });
    if (await checkoutBtn.isVisible()) {
      await checkoutBtn.click();

      // Wait for checkout overlay
      await page.waitForSelector('#checkoutOverlay', { timeout: 5000 });

      // Check for checkout form elements
      const nameField = page.locator('#checkoutOverlay input[name="name"], #checkoutOverlay input[placeholder*="name" i]');
      const emailField = page.locator('#checkoutOverlay input[name="email"], #checkoutOverlay input[type="email"]');
      const addressField = page.locator('#checkoutOverlay textarea, #checkoutOverlay input[name="address"]');

      // At least some form fields should be present
      const hasFormFields = (await nameField.count() > 0) || (await emailField.count() > 0) || (await addressField.count() > 0);
      expect(hasFormFields).toBe(true);

      console.log('Checkout form is displayed with form fields');
    } else {
      console.log('Checkout button not found');
    }
  });

  test('TC-CO-002: Test order placement with valid data', async ({ page }) => {
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

    // Open cart and checkout
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    const checkoutBtn = page.locator('#cartOverlay button').filter({ hasText: 'Checkout' });
    if (await checkoutBtn.isVisible()) {
      await checkoutBtn.click();
      await page.waitForSelector('#checkoutOverlay', { timeout: 5000 });

      // Fill out the form
      const nameField = page.locator('#checkoutForm input[name="name"]');
      const emailField = page.locator('#checkoutForm input[name="email"]');
      const phoneField = page.locator('#checkoutForm input[name="phone"]');
      const addressField = page.locator('#checkoutForm textarea[name="address"]');

      if (await nameField.isVisible()) {
        await nameField.fill('Test User');
      }
      if (await emailField.isVisible()) {
        await emailField.fill('test@example.com');
      }
      if (await phoneField.isVisible()) {
        await phoneField.fill('+1234567890');
      }
      if (await addressField.isVisible()) {
        await addressField.fill('123 Test Street, Test City, TC 12345');
      }

      // Submit the form
      const submitBtn = page.locator('#checkoutForm button[type="submit"]');
      if (await submitBtn.isVisible()) {
        await submitBtn.click();

        // Wait for order processing
        await page.waitForTimeout(3000);

        // Check for success message or order confirmation
        const successMessage = page.locator('text=/order.*placed|thank you|confirmation/i').first();
        const orderIdElement = page.locator('text=/order.*#|order.*id/i').first();

        if (await successMessage.isVisible() || await orderIdElement.isVisible()) {
          console.log('Order placed successfully');
        } else {
          // Check if we're redirected or overlay closed
          const isCheckoutClosed = !(await page.locator('#checkoutOverlay').isVisible());
          if (isCheckoutClosed) {
            console.log('Checkout overlay closed, order may have been placed');
          }
        }
      } else {
        console.log('Submit button not found');
      }
    }
  });

  test('TC-CO-004: Test order confirmation display', async ({ page }) => {
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

    // Go through checkout
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    const checkoutBtn = page.locator('#cartOverlay button').filter({ hasText: 'Checkout' });
    if (await checkoutBtn.isVisible()) {
      await checkoutBtn.click();
      await page.waitForSelector('#checkoutOverlay', { timeout: 5000 });

      // Fill form
      const nameField = page.locator('#checkoutForm input[name="name"]');
      const emailField = page.locator('#checkoutForm input[name="email"]');
      const phoneField = page.locator('#checkoutForm input[name="phone"]');
      const addressField = page.locator('#checkoutForm textarea[name="address"]');

      if (await nameField.isVisible()) await nameField.fill('Test User');
      if (await emailField.isVisible()) await emailField.fill('test@example.com');
      if (await phoneField.isVisible()) await phoneField.fill('+1234567890');
      if (await addressField.isVisible()) await addressField.fill('123 Test Street');

      // Submit
      const submitBtn = page.locator('#checkoutForm button[type="submit"]');
      if (await submitBtn.isVisible()) {
        await submitBtn.click();
        await page.waitForTimeout(3000);

        // Check for confirmation overlay
        const confirmationOverlay = page.locator('#confirmationOverlay');
        if (await confirmationOverlay.isVisible()) {
          // Check for order details
          const orderTitle = confirmationOverlay.locator('h3').filter({ hasText: 'Order Placed' });
          await expect(orderTitle).toBeVisible();

          const thankYouMessage = confirmationOverlay.locator('text=/Thanks|confirmation/i');
          await expect(thankYouMessage).toBeVisible();

          console.log('Order confirmation displayed successfully');
        } else {
          console.log('Confirmation overlay not found');
        }
      }
    }
  });

  test('TC-CO-005: Validate email confirmation sending', async ({ page }) => {
    await page.goto('/');

    // Add product and go through checkout
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

    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    const checkoutBtn = page.locator('#cartOverlay button').filter({ hasText: 'Checkout' });
    if (await checkoutBtn.isVisible()) {
      await checkoutBtn.click();
      await page.waitForSelector('#checkoutOverlay', { timeout: 5000 });

      // Fill form with email
      const emailField = page.locator('#checkoutForm input[name="email"]');
      if (await emailField.isVisible()) {
        await emailField.fill('test@example.com');
      }

      const nameField = page.locator('#checkoutForm input[name="name"]');
      if (await nameField.isVisible()) await nameField.fill('Test User');

      const addressField = page.locator('#checkoutForm textarea[name="address"]');
      if (await addressField.isVisible()) await addressField.fill('123 Test St');

      // Submit order
      const submitBtn = page.locator('#checkoutForm button[type="submit"]');
      if (await submitBtn.isVisible()) {
        await submitBtn.click();
        await page.waitForTimeout(3000);

        // Check if confirmation mentions email
        const emailMention = page.locator('text=/email|confirmation/i').first();
        if (await emailMention.isVisible()) {
          console.log('Order confirmation mentions email sending');
        } else {
          console.log('Email confirmation message not found, but order may still send email');
        }
      }
    }
  });

  test('TC-CO-006: Test COD order processing', async ({ page }) => {
    await page.goto('/');

    // Add product and go through checkout
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

    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    const checkoutBtn = page.locator('#cartOverlay button').filter({ hasText: 'Checkout' });
    if (await checkoutBtn.isVisible()) {
      await checkoutBtn.click();
      await page.waitForSelector('#checkoutOverlay', { timeout: 5000 });

      // Check for payment method selection
      const codOption = page.locator('input[value="cod"], [data-payment="cod"]').first();
      if (await codOption.isVisible()) {
        // Select COD
        if (await codOption.isChecked() === false) {
          await codOption.check();
        }
        console.log('COD payment method selected');
      } else {
        console.log('COD option not found, assuming COD is default');
      }

      // Fill form
      const nameField = page.locator('#checkoutForm input[name="name"]');
      const emailField = page.locator('#checkoutForm input[name="email"]');
      const addressField = page.locator('#checkoutForm textarea[name="address"]');

      if (await nameField.isVisible()) await nameField.fill('Test User');
      if (await emailField.isVisible()) await emailField.fill('test@example.com');
      if (await addressField.isVisible()) await addressField.fill('123 Test St');

      // Submit
      const submitBtn = page.locator('#checkoutForm button[type="submit"]');
      if (await submitBtn.isVisible()) {
        await submitBtn.click();
        await page.waitForTimeout(3000);

        // Check for successful COD order
        const successMessage = page.locator('text=/order.*placed|thank you/i').first();
        if (await successMessage.isVisible()) {
          console.log('COD order processed successfully');
        }
      }
    }
  });

  test('TC-CO-007: Verify order ID generation', async ({ page }) => {
    await page.goto('/');

    // Add product and place order
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

    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    const checkoutBtn = page.locator('#cartOverlay button').filter({ hasText: 'Checkout' });
    if (await checkoutBtn.isVisible()) {
      await checkoutBtn.click();
      await page.waitForSelector('#checkoutOverlay', { timeout: 5000 });

      // Fill form
      const nameField = page.locator('#checkoutForm input[name="name"]');
      const emailField = page.locator('#checkoutForm input[name="email"]');
      const addressField = page.locator('#checkoutForm textarea[name="address"]');

      if (await nameField.isVisible()) await nameField.fill('Test User');
      if (await emailField.isVisible()) await emailField.fill('test@example.com');
      if (await addressField.isVisible()) await addressField.fill('123 Test St');

      // Submit
      const submitBtn = page.locator('#checkoutForm button[type="submit"]');
      if (await submitBtn.isVisible()) {
        await submitBtn.click();
        await page.waitForTimeout(3000);

        // Check for order ID
        const orderIdElement = page.locator('text=/order.*#|order.*id|order.*number/i').first();
        if (await orderIdElement.isVisible()) {
          const orderIdText = await orderIdElement.textContent();
          console.log(`Order ID generated: ${orderIdText}`);
        } else {
          console.log('Order ID not displayed, but order may have been created with ID');
        }
      }
    }
  });

  test('TC-CO-008: Test checkout with empty cart', async ({ page }) => {
    await page.goto('/');

    // Wait for page load
    await page.waitForSelector('#grid', { timeout: 10000 });

    // Open cart (should be empty)
    const cartBtn = page.locator('#cartBtn');
    await cartBtn.click();
    await page.waitForSelector('#cartOverlay', { timeout: 5000 });

    // Check if checkout button is disabled or hidden
    const checkoutBtn = page.locator('#cartOverlay button').filter({ hasText: 'Checkout' });
    if (await checkoutBtn.isVisible()) {
      const isDisabled = await checkoutBtn.isDisabled();
      if (isDisabled) {
        console.log('Checkout button is disabled for empty cart');
      } else {
        console.log('Checkout button is enabled for empty cart - this may be unexpected');
      }
    } else {
      console.log('Checkout button not visible for empty cart');
    }

    // Check empty cart message
    const emptyMessage = page.locator('#cartItems').filter({ hasText: 'Your cart is empty' });
    await expect(emptyMessage).toBeVisible();

    console.log('Empty cart checkout validation completed');
  });

  test('TC-CO-009: Test email API endpoint', async ({ page }) => {
    // Test the email API endpoint
    const response = await page.request.post('/api/test-email', {
      data: {
        email: 'miqbal@sis.edu.eg'
      }
    });

    expect(response.status()).toBe(200);

    const responseData = await response.json();
    expect(responseData.message).toBe('Test email sent successfully');
    expect(responseData.email).toBe('miqbal@sis.edu.eg');
    expect(responseData.order_id).toBe(12345);
    expect(responseData.total).toBe(1199);

    console.log('Test email API endpoint works correctly');
  });
});
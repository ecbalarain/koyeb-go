const { test, expect } = require('@playwright/test');

test('Order Status Test - Order 300001', async ({ page }) => {
  // Navigate to the order status page
  await page.goto('https://api.bhomanshah.com/api/order-status/300001?email=miqbal@sis.edu.eg');

  // Wait for the page to load and check if we get JSON response
  await page.waitForLoadState('networkidle');

  // Get the page content (should be JSON)
  const content = await page.textContent('pre');
  const jsonContent = content || await page.textContent('body');

  console.log('Page Content:', jsonContent);

  // Try to parse the JSON response
  let orderData;
  try {
    orderData = JSON.parse(jsonContent);
  } catch (e) {
    console.log('Could not parse as JSON, checking if redirected...');
    // Check if we were redirected to a different page
    const currentUrl = page.url();
    console.log('Current URL:', currentUrl);

    if (currentUrl.includes('bhomanshah.com') && !currentUrl.includes('api.')) {
      console.log('❌ REDIRECTED TO FRONTEND - This is the bug we fixed!');
      expect(false).toBe(true); // Fail the test
    }
  }

  // Verify the order data
  if (orderData) {
    console.log('✅ Successfully received order data');

    // Check order details
    expect(orderData.order).toBeDefined();
    expect(orderData.order.id).toBe(300001);
    expect(orderData.order.status).toBe('new');
    expect(orderData.order.total).toBe(1199);

    // Check items
    expect(orderData.items).toBeDefined();
    expect(orderData.items.length).toBeGreaterThan(0);
    expect(orderData.items[0].product_name).toContain('Black Suede');

    console.log('✅ All order details verified successfully!');
  }
});
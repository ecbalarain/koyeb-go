const { chromium } = require('playwright');

async function testOrderStatus() {
  console.log('🚀 Starting Order Status Test...');

  // Launch browser in headed mode
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();

  try {
    console.log('📡 Navigating to order status URL...');
    // Navigate to the order status page
    await page.goto('https://api.bhomanshah.com/api/order-status/300001?email=miqbal@sis.edu.eg');

    // Wait for the page to load
    await page.waitForLoadState('networkidle');

    // Get the page content
    const content = await page.textContent('body');
    console.log('📄 Page Content:', content);

    // Check current URL
    const currentUrl = page.url();
    console.log('🔗 Current URL:', currentUrl);

    // Try to parse as JSON
    try {
      const orderData = JSON.parse(content);

      if (orderData.error) {
        console.log('❌ Error response:', orderData.error);
        return;
      }

      console.log('✅ Successfully received order data!');
      console.log('📦 Order ID:', orderData.order.id);
      console.log('📊 Status:', orderData.order.status);
      console.log('💰 Total:', orderData.order.total);
      console.log('🛒 Items:', orderData.items.length);

      // Verify order details
      if (orderData.order.id === 300001 &&
          orderData.order.status === 'new' &&
          orderData.order.total === 1199 &&
          orderData.items.length > 0) {
        console.log('🎉 All order details verified successfully!');
      } else {
        console.log('⚠️ Order details do not match expected values');
      }

    } catch (parseError) {
      console.log('❌ Could not parse response as JSON');

      // Check if redirected to frontend
      if (currentUrl.includes('bhomanshah.com') && !currentUrl.includes('api.')) {
        console.log('🚨 REDIRECTED TO FRONTEND - This indicates the bug is NOT fixed!');
      }
    }

    // Wait a bit so we can see the result
    await page.waitForTimeout(3000);

  } catch (error) {
    console.error('❌ Test failed:', error.message);
  } finally {
    await browser.close();
    console.log('🏁 Test completed');
  }
}

testOrderStatus();
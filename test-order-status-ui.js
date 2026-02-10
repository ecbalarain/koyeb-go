const { chromium } = require('playwright');

async function testOrderStatusUI() {
  console.log('🚀 Starting Order Status UI Test...');

  // Launch browser in headed mode
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();

  try {
    console.log('📡 Navigating to order status UI...');
    // Navigate to the frontend order status page
    await page.goto('https://f38dcb48.store-bs.pages.dev/order-status?order_id=300001&email=miqbal@sis.edu.eg');

    // Wait for the page to load
    await page.waitForLoadState('networkidle');

    // Check current URL
    const currentUrl = page.url();
    console.log('🔗 Current URL:', currentUrl);

    // Wait for either order details or error to appear
    try {
      await page.waitForSelector('#orderDetails, #error', { timeout: 15000 });
    } catch (e) {
      console.log('❌ Neither order details nor error appeared within timeout');
      const bodyText = await page.textContent('body');
      console.log('📄 Page body:', bodyText.substring(0, 500));
      return;
    }

    // Check if error is shown
    const errorVisible = await page.isVisible('#error');
    if (errorVisible) {
      const errorMessage = await page.textContent('#errorMessage');
      console.log('❌ Error displayed:', errorMessage);
      return;
    }

    // Check if order details are displayed
    const orderId = await page.textContent('#orderId');
    const orderStatus = await page.textContent('#orderStatus');
    const orderTotal = await page.textContent('#orderTotal');

    console.log('✅ Order Status UI loaded successfully!');
    console.log('📦 Order ID:', orderId);
    console.log('📊 Status:', orderStatus);
    console.log('💰 Total:', orderTotal);

    // Check if items are displayed
    const itemsContainer = await page.$('#orderItems');
    const itemCount = await itemsContainer.$$eval('.flex', items => items.length);
    console.log('🛒 Items displayed:', itemCount);

    // Verify the data
    if (orderId.includes('300001') &&
        orderStatus.toLowerCase().includes('new') &&
        orderTotal.includes('1,199') &&
        itemCount > 0) {
      console.log('🎉 All order details displayed correctly!');
    } else {
      console.log('⚠️ Order details do not match expected values');
    }

    // Wait a bit so we can see the result
    await page.waitForTimeout(5000);

  } catch (error) {
    console.error('❌ Test failed:', error.message);
  } finally {
    await browser.close();
    console.log('🏁 Test completed');
  }
}

testOrderStatusUI();
const { chromium } = require('@playwright/test');
const fs = require('fs');
const path = require('path');

const OXLOOK_URL = 'https://oxlook.com';
const OUTPUT_FILE = path.join(__dirname, 'cloudflare-pages-frontend', 'products.json');

async function scrapeOxlook() {
  console.log('🚀 Starting OXLOOK.com scraper...\n');
  
  const browser = await chromium.launch({
    headless: false,
    slowMo: 50,
  });
  
  const page = await browser.newPage();
  const products = [];
  
  try {
    console.log(`📡 Navigating to ${OXLOOK_URL}...`);
    await page.goto(OXLOOK_URL, { waitUntil: 'networkidle', timeout: 30000 });
    await page.waitForTimeout(2000);
    
    console.log('✅ Website loaded\n');
    
    // Find all product cards/links
    console.log('🔍 Finding product elements...');
    
    // Try different selectors to find products
    const productSelectors = [
      'article',
      '.product-card',
      '.product',
      '[data-product]',
      'a[href*="product"]',
      '.grid article',
      '.product-item'
    ];
    
    let productElements = null;
    let usedSelector = null;
    
    for (const selector of productSelectors) {
      try {
        const elements = await page.locator(selector).all();
        if (elements.length > 0) {
          productElements = elements;
          usedSelector = selector;
          console.log(`✅ Found ${elements.length} products using selector: ${selector}\n`);
          break;
        }
      } catch (e) {
        continue;
      }
    }
    
    if (!productElements || productElements.length === 0) {
      console.log('⚠️  No products found with standard selectors.');
      console.log('📸 Taking screenshot for manual inspection...');
      await page.screenshot({ path: 'oxlook-page.png', fullPage: true });
      
      // Try to find any clickable items
      console.log('\n🔍 Looking for any product links...');
      const links = await page.locator('a').all();
      console.log(`Found ${links.length} total links on the page`);
      
      // Get page content to analyze
      const bodyText = await page.locator('body').textContent();
      console.log('\n📄 Page preview:');
      console.log(bodyText.substring(0, 500));
      
      throw new Error('Could not find products on the page. Check oxlook-page.png screenshot.');
    }
    
    console.log(`🔄 Extracting details from ${productElements.length} products...\n`);
    
    let productId = 1;
    
    for (let i = 0; i < productElements.length; i++) {
      console.log(`📦 Processing product ${i + 1}/${productElements.length}...`);
      
      try {
        const element = productElements[i];
        
        // Extract data from the product card
        const productData = {
          id: productId++,
          name: '',
          slug: '',
          description: '',
          category: 'Home',
          image: '',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        };
        
        // Try to extract title
        const titleSelectors = ['h2', 'h3', '.product-title', '.title', '[class*="title"]', '.product-name'];
        for (const sel of titleSelectors) {
          try {
            const titleEl = element.locator(sel).first();
            const title = await titleEl.textContent({ timeout: 1000 });
            if (title && title.trim()) {
              productData.name = title.trim();
              productData.slug = title.trim().toLowerCase()
                .replace(/[^a-z0-9\s-]/g, '')
                .replace(/\s+/g, '-')
                .replace(/-+/g, '-');
              break;
            }
          } catch (e) {
            continue;
          }
        }
        
        // Try to extract description
        const descSelectors = ['p', '.description', '.product-description', '[class*="desc"]'];
        for (const sel of descSelectors) {
          try {
            const descEl = element.locator(sel).first();
            const desc = await descEl.textContent({ timeout: 1000 });
            if (desc && desc.trim() && desc.trim().length > 10) {
              productData.description = desc.trim();
              break;
            }
          } catch (e) {
            continue;
          }
        }
        
        // If no name found in card, try clicking to get details
        if (!productData.name) {
          console.log('  ⚠️  No title in card, attempting to open product detail...');
          
          // Save current page state
          const currentUrl = page.url();
          
          try {
            await element.click({ timeout: 3000 });
            await page.waitForTimeout(2000);
            
            // Extract from detail page
            const detailTitleSel = ['h1', 'h2', '.product-title', '[class*="title"]'];
            for (const sel of detailTitleSel) {
              try {
                const title = await page.locator(sel).first().textContent({ timeout: 2000 });
                if (title && title.trim()) {
                  productData.name = title.trim();
                  productData.slug = title.trim().toLowerCase()
                    .replace(/[^a-z0-9\s-]/g, '')
                    .replace(/\s+/g, '-')
                    .replace(/-+/g, '-');
                  break;
                }
              } catch (e) {
                continue;
              }
            }
            
            // Get description from detail page
            const detailDescSel = ['.description', '.product-description', 'p'];
            for (const sel of detailDescSel) {
              try {
                const desc = await page.locator(sel).first().textContent({ timeout: 2000 });
                if (desc && desc.trim() && desc.trim().length > 10) {
                  productData.description = desc.trim();
                  break;
                }
              } catch (e) {
                continue;
              }
            }
            
            // Navigate back
            await page.goto(currentUrl, { waitUntil: 'networkidle' });
            await page.waitForTimeout(1000);
            
            // Re-fetch product elements
            productElements = await page.locator(usedSelector).all();
            
          } catch (e) {
            console.log('  ⚠️  Could not open product detail:', e.message);
          }
        }
        
        if (productData.name) {
          products.push(productData);
          console.log(`  ✅ ${productData.name}`);
          if (productData.description) {
            console.log(`     ${productData.description.substring(0, 60)}...`);
          }
        } else {
          console.log(`  ⚠️  Skipped - no title found`);
        }
        
      } catch (error) {
        console.log(`  ❌ Error processing product: ${error.message}`);
      }
      
      console.log('');
    }
    
    console.log(`\n✅ Scraped ${products.length} products successfully!\n`);
    
    // Save to JSON file
    const jsonContent = JSON.stringify(products, null, 2);
    fs.writeFileSync(OUTPUT_FILE, jsonContent, 'utf8');
    console.log(`💾 Saved to: ${OUTPUT_FILE}\n`);
    
    // Also save a backup
    const backupFile = path.join(__dirname, 'scraped-products-backup.json');
    fs.writeFileSync(backupFile, jsonContent, 'utf8');
    console.log(`💾 Backup saved to: ${backupFile}\n`);
    
    // Display summary
    console.log('📊 Products Summary:');
    console.log('='.repeat(60));
    products.forEach((p, i) => {
      console.log(`${i + 1}. ${p.name}`);
      console.log(`   Category: ${p.category}`);
      console.log(`   Slug: ${p.slug}`);
      if (p.description) {
        console.log(`   Description: ${p.description.substring(0, 80)}...`);
      }
      console.log('');
    });
    
  } catch (error) {
    console.error('\n❌ Scraping failed:', error.message);
    await page.screenshot({ path: 'oxlook-error.png', fullPage: true });
    console.log('📸 Error screenshot saved to: oxlook-error.png');
    throw error;
  } finally {
    await browser.close();
  }
}

// Run the scraper
scrapeOxlook()
  .then(() => {
    console.log('\n✅ Scraping completed successfully!');
    process.exit(0);
  })
  .catch((error) => {
    console.error('\n❌ Scraping failed:', error);
    process.exit(1);
  });

const fs = require('fs');
const path = require('path');
const { request } = require('@playwright/test');

const STORE_API = 'https://oxlook.com/wp-json/wc/store/products';
const OUTPUT_JSON = path.join(__dirname, 'cloudflare-pages-frontend', 'products.json');
const IMAGE_DIR = path.join(__dirname, 'cloudflare-pages-frontend', 'assets', 'oxlook');
const IMAGE_PUBLIC_BASE = '/assets/oxlook';

const PAGE_SIZE = 100;

function stripHtml(input) {
  if (!input) return '';
  return input.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim();
}

function decodeEntities(input) {
  if (!input) return '';
  let output = input
    .replace(/&nbsp;/g, ' ')
    .replace(/&amp;/g, '&')
    .replace(/&quot;/g, '"')
    .replace(/&apos;/g, "'")
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>');

  output = output.replace(/&#x([0-9a-fA-F]+);/g, (_, hex) => {
    const code = parseInt(hex, 16);
    return Number.isNaN(code) ? '' : String.fromCodePoint(code);
  });

  output = output.replace(/&#(\d+);/g, (_, num) => {
    const code = parseInt(num, 10);
    return Number.isNaN(code) ? '' : String.fromCodePoint(code);
  });

  return output;
}

function toAscii(input) {
  if (!input) return '';
  let output = input
    .replace(/[\u2018\u2019]/g, "'")
    .replace(/[\u201C\u201D]/g, '"')
    .replace(/[\u2013\u2014]/g, '-')
    .replace(/\u2026/g, '...');

  output = output.replace(/[\u0080-\uFFFF]/g, '');
  return output;
}

function cleanText(input) {
  const decoded = decodeEntities(input);
  const stripped = stripHtml(decoded);
  return toAscii(stripped);
}

function pickExtension(url) {
  try {
    const parsed = new URL(url);
    const ext = path.extname(parsed.pathname);
    return ext || '.jpg';
  } catch (e) {
    return '.jpg';
  }
}

function safeFileName(value) {
  return value.toLowerCase().replace(/[^a-z0-9\-]+/g, '-').replace(/\-+/g, '-').replace(/^-|-$|\.$/g, '');
}

async function fetchAllProducts(api) {
  const all = [];
  let page = 1;
  let totalPages = null;

  while (true) {
    const url = `${STORE_API}?per_page=${PAGE_SIZE}&page=${page}`;
    const response = await api.get(url);
    if (!response.ok()) {
      throw new Error(`Failed to fetch products page ${page}: ${response.status()}`);
    }

    const headers = response.headers();
    if (!totalPages && headers['x-wp-totalpages']) {
      totalPages = parseInt(headers['x-wp-totalpages'], 10);
    }

    const data = await response.json();
    if (!Array.isArray(data) || data.length === 0) {
      break;
    }

    all.push(...data);

    if (totalPages && page >= totalPages) {
      break;
    }

    page += 1;
  }

  return all;
}

async function downloadImage(api, url, destPath) {
  const response = await api.get(url);
  if (!response.ok()) {
    throw new Error(`Failed to download image: ${url}`);
  }
  const buffer = await response.body();
  fs.writeFileSync(destPath, buffer);
}

async function run() {
  fs.mkdirSync(IMAGE_DIR, { recursive: true });

  const api = await request.newContext({
    baseURL: 'https://oxlook.com',
    extraHTTPHeaders: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    },
  });

  const products = await fetchAllProducts(api);
  const output = [];

  for (const product of products) {
    const name = cleanText(product.name);
    const description = cleanText(product.description || product.short_description || '');
    const shortDescription = cleanText(product.short_description || '');

    const categories = (product.categories || []).map((c) => cleanText(c.name));
    const tags = (product.tags || []).map((t) => cleanText(t.name));

    const attributes = (product.attributes || []).map((attr) => ({
      name: cleanText(attr.name),
      values: (attr.terms || []).map((term) => cleanText(term.name)),
    }));

    const sizeAttribute = attributes.find((attr) => attr.name.toLowerCase() === 'size');
    const sizes = sizeAttribute ? sizeAttribute.values : [];

    const imageUrls = (product.images || []).map((img) => img.src).filter(Boolean);
    const images = [];

    let imageIndex = 1;
    for (const imgUrl of imageUrls) {
      const ext = pickExtension(imgUrl);
      const filename = `${safeFileName(product.slug || product.id.toString())}-${String(imageIndex).padStart(2, '0')}${ext}`;
      const destPath = path.join(IMAGE_DIR, filename);
      const publicPath = `${IMAGE_PUBLIC_BASE}/${filename}`;

      try {
        if (!fs.existsSync(destPath)) {
          await downloadImage(api, imgUrl, destPath);
        }
        images.push(publicPath);
      } catch (error) {
        console.warn(`Failed to download image for ${product.slug}: ${imgUrl}`);
      }

      imageIndex += 1;
    }

    const prices = product.prices || {};
    const price = prices.price ? parseInt(prices.price, 10) : null;
    const regularPrice = prices.regular_price ? parseInt(prices.regular_price, 10) : null;
    const salePrice = prices.sale_price ? parseInt(prices.sale_price, 10) : null;

    output.push({
      id: product.id,
      name,
      slug: product.slug || safeFileName(name || String(product.id)),
      description,
      shortDescription,
      category: categories[0] || 'Uncategorized',
      categories,
      tags,
      images,
      sizes,
      attributes,
      prices: {
        currency: prices.currency_code || 'PKR',
        price,
        regularPrice,
        salePrice,
      },
      sourceUrl: product.permalink || null,
      inStock: product.is_in_stock === true,
    });
  }

  fs.writeFileSync(OUTPUT_JSON, JSON.stringify(output, null, 2));
  await api.dispose();

  console.log(`Saved ${output.length} products to ${OUTPUT_JSON}`);
  console.log(`Downloaded images to ${IMAGE_DIR}`);
}

run().catch((error) => {
  console.error('Scrape failed:', error);
  process.exit(1);
});

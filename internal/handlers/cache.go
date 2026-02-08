package handlers

import (
	"sync"
)

// CacheManager manages cache invalidation
type CacheManager struct {
	mu              sync.RWMutex
	variantCacheKey map[string]bool // Track which variant cache keys need purging
}

// NewCacheManager creates a new cache manager instance
func NewCacheManager() *CacheManager {
	return &CacheManager{
		variantCacheKey: make(map[string]bool),
	}
}

// InvalidateAll marks all cache entries for invalidation
func (cm *CacheManager) InvalidateAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.variantCacheKey = make(map[string]bool)
}

// InvalidateVariantCache marks variant cache for a specific product to be invalidated
func (cm *CacheManager) InvalidateVariantCache(productSlug string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.variantCacheKey[productSlug] = true
}

// ShouldInvalidate checks if a cache entry should be invalidated
func (cm *CacheManager) ShouldInvalidate(key string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	// Check if this specific key should be invalidated
	if cm.variantCacheKey[key] {
		return true
	}
	
	return false
}

// ClearInvalidation clears the invalidation flag for a specific key
func (cm *CacheManager) ClearInvalidation(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.variantCacheKey, key)
}

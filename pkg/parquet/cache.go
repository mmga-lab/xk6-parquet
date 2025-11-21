package parquet

import (
	"sync"
	"time"
)

// ReaderCache manages cached Parquet file data.
type ReaderCache struct {
	cache map[string]*CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// CacheEntry represents a single cache entry with data and timestamp.
type CacheEntry struct {
	data      []map[string]interface{}
	timestamp time.Time
}

// NewReaderCache creates a new cache instance with default TTL.
func NewReaderCache() *ReaderCache {
	return &ReaderCache{
		cache: make(map[string]*CacheEntry),
		ttl:   5 * time.Minute, // Default 5 minutes TTL
	}
}

// Get retrieves data from cache if it exists and hasn't expired.
func (rc *ReaderCache) Get(key string) ([]map[string]interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	entry, ok := rc.cache[key]
	if !ok {
		return nil, false
	}

	// Check if entry has expired
	if time.Since(entry.timestamp) > rc.ttl {
		return nil, false
	}

	return entry.data, true
}

// Set stores data in the cache with current timestamp.
func (rc *ReaderCache) Set(key string, data []map[string]interface{}) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache[key] = &CacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
}

// Clear removes all entries from the cache.
func (rc *ReaderCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache = make(map[string]*CacheEntry)
}

// Remove deletes a specific entry from the cache.
func (rc *ReaderCache) Remove(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	delete(rc.cache, key)
}

// SetTTL sets the time-to-live duration for cache entries.
func (rc *ReaderCache) SetTTL(ttl time.Duration) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.ttl = ttl
}

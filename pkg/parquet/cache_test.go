package parquet

import (
	"testing"
	"time"
)

func TestNewReaderCache(t *testing.T) {
	cache := NewReaderCache()

	if cache == nil {
		t.Fatal("NewReaderCache returned nil")
	}

	if cache.cache == nil {
		t.Error("cache map not initialized")
	}

	if cache.ttl != 5*time.Minute {
		t.Errorf("expected default TTL of 5 minutes, got %v", cache.ttl)
	}
}

func TestCacheSetAndGet(t *testing.T) {
	cache := NewReaderCache()

	testData := []map[string]interface{}{
		{"id": 1, "name": "test"},
		{"id": 2, "name": "test2"},
	}

	// Test Set
	cache.Set("test-key", testData)

	// Test Get - should find it
	data, found := cache.Get("test-key")
	if !found {
		t.Error("expected to find cached data")
	}

	if len(data) != len(testData) {
		t.Errorf("expected %d items, got %d", len(testData), len(data))
	}

	// Test Get - non-existent key
	_, found = cache.Get("non-existent")
	if found {
		t.Error("expected not to find non-existent key")
	}
}

func TestCacheTTL(t *testing.T) {
	cache := NewReaderCache()
	cache.SetTTL(50 * time.Millisecond)

	testData := []map[string]interface{}{
		{"id": 1, "name": "test"},
	}

	cache.Set("test-key", testData)

	// Should find it immediately
	_, found := cache.Get("test-key")
	if !found {
		t.Error("expected to find cached data immediately")
	}

	// Wait for TTL to expire
	time.Sleep(100 * time.Millisecond)

	// Should not find it after TTL
	_, found = cache.Get("test-key")
	if found {
		t.Error("expected cache entry to expire")
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewReaderCache()

	testData := []map[string]interface{}{
		{"id": 1, "name": "test"},
	}

	cache.Set("key1", testData)
	cache.Set("key2", testData)

	// Should find both
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")

	if !found1 || !found2 {
		t.Error("expected to find both cached entries")
	}

	// Clear cache
	cache.Clear()

	// Should not find either
	_, found1 = cache.Get("key1")
	_, found2 = cache.Get("key2")

	if found1 || found2 {
		t.Error("expected cache to be cleared")
	}
}

func TestCacheRemove(t *testing.T) {
	cache := NewReaderCache()

	testData := []map[string]interface{}{
		{"id": 1, "name": "test"},
	}

	cache.Set("key1", testData)
	cache.Set("key2", testData)

	// Remove one key
	cache.Remove("key1")

	// key1 should be gone
	_, found := cache.Get("key1")
	if found {
		t.Error("expected key1 to be removed")
	}

	// key2 should still exist
	_, found = cache.Get("key2")
	if !found {
		t.Error("expected key2 to still exist")
	}
}

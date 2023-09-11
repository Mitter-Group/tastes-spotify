package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Test save and get
func TestSaveAndGet(t *testing.T) {
	cache := NewMemoryCache("test", 10, 10*time.Second, false)
	cache.Save(context.Background(), "key", "value")
	expired, value := cache.Get(context.Background(), "key")
	assert.Equal(t, "value", value)
	assert.Falsef(t, expired, "value should not be expired")
}

// Test save and get with ttl and return expired
func TestSaveAndGetWithTTLWithReturnExpired(t *testing.T) {
	cache := NewMemoryCache("test", 10, 10*time.Second, true)
	cache.SaveWithTTL(context.Background(), "key", "value", 500*time.Millisecond)
	expired, value := cache.Get(context.Background(), "key")
	assert.Equal(t, "value", value)
	assert.Falsef(t, expired, "value should not be expired")
	time.Sleep(1 * time.Second)
	expired, value = cache.Get(context.Background(), "key")
	assert.NotNilf(t, value, "value should return as returnExpired is true")
	assert.Truef(t, expired, "value should be expired")
}

// Test save and get with ttl and return expired
func TestSaveAndGetWithTTLWithoutReturnExpired(t *testing.T) {
	cache := NewMemoryCache("test", 10, 10*time.Second, false)
	cache.SaveWithTTL(context.Background(), "key", "value", 500*time.Millisecond)
	expired, value := cache.Get(context.Background(), "key")
	assert.Equal(t, "value", value)
	assert.Falsef(t, expired, "value should not be expired")
	time.Sleep(1 * time.Second)
	expired, value = cache.Get(context.Background(), "key")
	assert.Nilf(t, value, "value should return as returnExpired is false")
	assert.Falsef(t, expired, "value should not be expired")
}

// Test delete
func TestDelete(t *testing.T) {
	cache := NewMemoryCache("test", 10, 1000*time.Second, false)
	cache.Save(context.Background(), "key", "value")
	expired, value := cache.Get(context.Background(), "key")
	assert.Equal(t, "value", value)
	assert.Falsef(t, expired, "value should not be expired")
	cache.Delete(context.Background(), "key")
	expired, value = cache.Get(context.Background(), "key")
	assert.Nilf(t, value, "value should be nil")
	assert.Falsef(t, expired, "value should not be expired")
}

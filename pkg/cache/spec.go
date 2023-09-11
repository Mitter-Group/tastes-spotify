package cache

import (
	"context"
	"time"
)

type Spec interface {
	// Get key value pair from the cache. If the key is not found, it returns nil. If the key is found, it returns the value.
	// Also returns if the element is expired if was previously configured in the client
	// The value is a pointer to the actual value stored in the cache. This is done to avoid copying the value.
	// The value should be casted to the correct type before using it.
	// Example:
	//  var expired bool
	// 	var value *string
	// 	expired, value = cache.Get(ctx, "key").(*string)
	Get(ctx context.Context, key string) (expired bool, value interface{})
	// Save key value pair in the cache
	// Example:
	//  cache.Save(ctx, "key", "value")
	Save(ctx context.Context, key string, item interface{})
	// SaveWithTTL saves key value pair in the cache with a custom ttl (in duration)
	// Example:
	//  cache.SaveWithTTL(ctx, "key", "value", 10*time.Second)
	SaveWithTTL(ctx context.Context, key string, item interface{}, ttl time.Duration)
	// Delete key value pair from the cache
	// Example:
	//  cache.Delete(ctx, "key")
	Delete(ctx context.Context, key string)
}

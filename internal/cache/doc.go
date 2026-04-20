// Package cache implements a TTL-based in-memory cache for Vault secret
// responses. It is used by the sync engine to avoid redundant reads from
// Vault when multiple .env files reference the same secret path within a
// single vaultpull run.
//
// # Usage
//
//	c := cache.New(5 * time.Minute)
//
//	// Store a fetched secret.
//	c.Set("secret/myapp/prod", data)
//
//	// Retrieve if still fresh.
//	if v, ok := c.Get("secret/myapp/prod"); ok {
//		// use v
//	}
//
//	// Persist across short-lived CLI invocations.
//	c.SaveToFile(".vaultpull.cache")
//	c.LoadFromFile(".vaultpull.cache")
//
// Entries are never served after their TTL has elapsed, even when loaded
// from a persisted file. The default TTL is configurable via the
// cache_ttl field in .vaultpull.yaml.
package cache

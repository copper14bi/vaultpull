// Package secrets provides utilities for evaluating, masking, scanning,
// classifying, and managing the lifecycle of secrets.
//
// # Lifecycle
//
// The lifecycle sub-feature aggregates multiple independent checks into a
// single, unified evaluation per secret:
//
//   - Age: how long since the secret was created (via CheckAge)
//   - Strength: entropy and complexity of the secret value (via CheckStrength)
//   - TTL: time-to-live / expiry metadata (via CheckTTL)
//
// Use CheckLifecycle for a single key/value pair, or CheckLifecycleMap to
// evaluate an entire map of secrets at once.
//
// # Status Priority
//
// Status values are ordered by severity:
//
//	ok < warning < critical < expired
//
// An expired age always wins over a critical strength failure so that
// operators are not misled about the primary remediation action required.
//
// # Example
//
//	opts := secrets.DefaultLifecycleOptions()
//	opts.MinStrength = secrets.StrengthStrong
//	result := secrets.CheckLifecycle("DB_PASSWORD", value, createdAt, opts)
//	if result.Status != secrets.LifecycleOK {
//		for _, msg := range result.Messages {
//			fmt.Println(msg)
//		}
//	}
package secrets

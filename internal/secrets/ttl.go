package secrets

import (
	"fmt"
	"time"
)

// TTLStatus represents the state of a secret's time-to-live.
type TTLStatus int

const (
	TTLStatusOK      TTLStatus = iota // Within acceptable range
	TTLStatusWarning                  // Approaching expiry
	TTLStatusExpired                  // Past TTL
	TTLStatusNone                     // No TTL set
)

func (s TTLStatus) String() string {
	switch s {
	case TTLStatusOK:
		return "ok"
	case TTLStatusWarning:
		return "warning"
	case TTLStatusExpired:
		return "expired"
	case TTLStatusNone:
		return "none"
	default:
		return "unknown"
	}
}

// TTLResult holds the outcome of a TTL check for a single secret.
type TTLResult struct {
	Key        string
	Status     TTLStatus
	Remaining  time.Duration
	Message    string
}

// TTLOptions configures thresholds for TTL evaluation.
type TTLOptions struct {
	WarnThreshold time.Duration // Warn when remaining TTL is below this
}

// DefaultTTLOptions returns sensible defaults.
func DefaultTTLOptions() TTLOptions {
	return TTLOptions{
		WarnThreshold: 24 * time.Hour,
	}
}

// CheckTTL evaluates whether a secret's TTL (stored as a Unix timestamp string
// under the key "<key>_expires_at") is within acceptable bounds.
func CheckTTL(key string, expiresAt time.Time, opts TTLOptions) TTLResult {
	now := time.Now()

	if expiresAt.IsZero() {
		return TTLResult{Key: key, Status: TTLStatusNone, Message: "no TTL set"}
	}

	remaining := expiresAt.Sub(now)

	switch {
	case remaining <= 0:
		return TTLResult{
			Key:       key,
			Status:    TTLStatusExpired,
			Remaining: 0,
			Message:   fmt.Sprintf("expired %s ago", (-remaining).Round(time.Second)),
		}
	case remaining < opts.WarnThreshold:
		return TTLResult{
			Key:       key,
			Status:    TTLStatusWarning,
			Remaining: remaining,
			Message:   fmt.Sprintf("expires in %s", remaining.Round(time.Second)),
		}
	default:
		return TTLResult{
			Key:       key,
			Status:    TTLStatusOK,
			Remaining: remaining,
			Message:   fmt.Sprintf("valid for %s", remaining.Round(time.Second)),
		}
	}
}

// CheckTTLMap evaluates TTL for a map of key -> expiry time pairs.
func CheckTTLMap(expiries map[string]time.Time, opts TTLOptions) []TTLResult {
	results := make([]TTLResult, 0, len(expiries))
	for key, exp := range expiries {
		results = append(results, CheckTTL(key, exp, opts))
	}
	return results
}

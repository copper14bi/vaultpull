package secrets

import (
	"fmt"
	"time"
)

// ExpiryStatus represents whether a secret is nearing or past its expiry.
type ExpiryStatus int

const (
	ExpiryOK      ExpiryStatus = iota // within valid period
	ExpiryWarning                     // nearing expiry
	ExpiryExpired                     // past expiry
)

func (s ExpiryStatus) String() string {
	switch s {
	case ExpiryOK:
		return "ok"
	case ExpiryWarning:
		return "warning"
	case ExpiryExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// ExpiryResult holds the result of an expiry check for a single secret.
type ExpiryResult struct {
	Key       string
	ExpiresAt time.Time
	Status    ExpiryStatus
	Message   string
}

// ExpiryOptions configures thresholds for expiry checks.
type ExpiryOptions struct {
	WarnBefore time.Duration // warn if expiry is within this window
}

// DefaultExpiryOptions returns sensible defaults.
func DefaultExpiryOptions() ExpiryOptions {
	return ExpiryOptions{
		WarnBefore: 7 * 24 * time.Hour, // 7 days
	}
}

// CheckExpiry evaluates whether a secret expires soon or has already expired.
func CheckExpiry(key string, expiresAt time.Time, opts ExpiryOptions) ExpiryResult {
	now := time.Now()

	if expiresAt.IsZero() {
		return ExpiryResult{
			Key:       key,
			ExpiresAt: expiresAt,
			Status:    ExpiryOK,
			Message:   "no expiry set",
		}
	}

	if now.After(expiresAt) {
		return ExpiryResult{
			Key:       key,
			ExpiresAt: expiresAt,
			Status:    ExpiryExpired,
			Message:   fmt.Sprintf("expired %s ago", now.Sub(expiresAt).Round(time.Minute)),
		}
	}

	timeLeft := expiresAt.Sub(now)
	if timeLeft <= opts.WarnBefore {
		return ExpiryResult{
			Key:       key,
			ExpiresAt: expiresAt,
			Status:    ExpiryWarning,
			Message:   fmt.Sprintf("expires in %s", timeLeft.Round(time.Minute)),
		}
	}

	return ExpiryResult{
		Key:       key,
		ExpiresAt: expiresAt,
		Status:    ExpiryOK,
		Message:   fmt.Sprintf("valid for %s", timeLeft.Round(time.Minute)),
	}
}

// CheckExpiryMap evaluates expiry for a map of key -> expiresAt times.
func CheckExpiryMap(secrets map[string]time.Time, opts ExpiryOptions) []ExpiryResult {
	results := make([]ExpiryResult, 0, len(secrets))
	for key, expiresAt := range secrets {
		results = append(results, CheckExpiry(key, expiresAt, opts))
	}
	return results
}

package secrets

import (
	"fmt"
	"time"
)

// AgeStatus represents how old a secret is relative to its rotation policy.
type AgeStatus int

const (
	AgeFresh   AgeStatus = iota // within acceptable window
	AgeWarning                  // approaching max age
	AgeExpired                  // past max age
)

func (s AgeStatus) String() string {
	switch s {
	case AgeFresh:
		return "fresh"
	case AgeWarning:
		return "warning"
	case AgeExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// AgeResult holds the age check outcome for a single secret.
type AgeResult struct {
	Key       string
	Age       time.Duration
	Status    AgeStatus
	Message   string
}

// AgeOptions configures thresholds for age checking.
type AgeOptions struct {
	WarnAfter   time.Duration // default: 60 days
	ExpireAfter time.Duration // default: 90 days
}

// DefaultAgeOptions returns sensible defaults.
func DefaultAgeOptions() AgeOptions {
	return AgeOptions{
		WarnAfter:   60 * 24 * time.Hour,
		ExpireAfter: 90 * 24 * time.Hour,
	}
}

// CheckAge evaluates how old a secret is given its last-rotated timestamp.
func CheckAge(key string, lastRotated time.Time, opts AgeOptions) AgeResult {
	age := time.Since(lastRotated)
	result := AgeResult{
		Key: key,
		Age: age,
	}

	switch {
	case age >= opts.ExpireAfter:
		result.Status = AgeExpired
		result.Message = fmt.Sprintf("secret is %d days old — rotation required", int(age.Hours()/24))
	case age >= opts.WarnAfter:
		result.Status = AgeWarning
		result.Message = fmt.Sprintf("secret is %d days old — rotation recommended", int(age.Hours()/24))
	default:
		result.Status = AgeFresh
		result.Message = fmt.Sprintf("secret is %d days old — ok", int(age.Hours()/24))
	}

	return result
}

// CheckAgeMap evaluates a map of key → last-rotated timestamps.
func CheckAgeMap(timestamps map[string]time.Time, opts AgeOptions) []AgeResult {
	results := make([]AgeResult, 0, len(timestamps))
	for key, ts := range timestamps {
		results = append(results, CheckAge(key, ts, opts))
	}
	return results
}

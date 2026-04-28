package secrets

import (
	"fmt"
	"time"
)

// DriftSeverity indicates how significant a detected drift is.
type DriftSeverity string

const (
	DriftSeverityNone     DriftSeverity = "none"
	DriftSeverityLow      DriftSeverity = "low"
	DriftSeverityMedium   DriftSeverity = "medium"
	DriftSeverityHigh     DriftSeverity = "high"
	DriftSeverityCritical DriftSeverity = "critical"
)

// DriftResult holds the outcome of a drift check for a single secret.
type DriftResult struct {
	Key        string
	Severity   DriftSeverity
	Reason     string
	DetectedAt time.Time
}

// DriftOptions controls thresholds used during drift detection.
type DriftOptions struct {
	// MaxAgeDays is the age in days beyond which a secret is considered drifted.
	MaxAgeDays int
	// EntropyThreshold is the minimum Shannon entropy expected for sensitive keys.
	EntropyThreshold float64
	// RequireRotation flags secrets that have never been rotated.
	RequireRotation bool
}

// DefaultDriftOptions returns sensible defaults for drift detection.
func DefaultDriftOptions() DriftOptions {
	return DriftOptions{
		MaxAgeDays:       90,
		EntropyThreshold: 3.5,
		RequireRotation:  true,
	}
}

// CheckDrift evaluates a single secret value for drift relative to its metadata.
// createdAt is when the secret was originally set; rotatedAt may be zero if never rotated.
func CheckDrift(key, value string, createdAt, rotatedAt time.Time, opts DriftOptions) DriftResult {
	now := time.Now()
	result := DriftResult{Key: key, DetectedAt: now, Severity: DriftSeverityNone}

	if value == "" {
		result.Severity = DriftSeverityHigh
		result.Reason = "secret value is empty"
		return result
	}

	if opts.RequireRotation && rotatedAt.IsZero() {
		result.Severity = DriftSeverityMedium
		result.Reason = "secret has never been rotated"
		return result
	}

	ref := createdAt
	if !rotatedAt.IsZero() {
		ref = rotatedAt
	}
	ageDays := int(now.Sub(ref).Hours() / 24)
	if ageDays > opts.MaxAgeDays {
		severity := DriftSeverityMedium
		if ageDays > opts.MaxAgeDays*2 {
			severity = DriftSeverityCritical
		} else if ageDays > opts.MaxAgeDays+30 {
			severity = DriftSeverityHigh
		}
		result.Severity = severity
		result.Reason = fmt.Sprintf("secret is %d days old (max %d)", ageDays, opts.MaxAgeDays)
		return result
	}

	if IsSensitiveKey(key) {
		if ent := ShannonEntropy(value); ent < opts.EntropyThreshold {
			result.Severity = DriftSeverityLow
			result.Reason = fmt.Sprintf("low entropy %.2f for sensitive key (min %.2f)", ent, opts.EntropyThreshold)
			return result
		}
	}

	return result
}

// CheckDriftMap runs CheckDrift over a map of secrets sharing the same timestamps.
func CheckDriftMap(secrets map[string]string, createdAt, rotatedAt time.Time, opts DriftOptions) []DriftResult {
	results := make([]DriftResult, 0, len(secrets))
	for k, v := range secrets {
		if r := CheckDrift(k, v, createdAt, rotatedAt, opts); r.Severity != DriftSeverityNone {
			results = append(results, r)
		}
	}
	return results
}

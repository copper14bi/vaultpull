package secrets

import (
	"fmt"
	"time"
)

// LifecycleStatus represents the combined lifecycle state of a secret.
type LifecycleStatus string

const (
	LifecycleOK       LifecycleStatus = "ok"
	LifecycleWarning  LifecycleStatus = "warning"
	LifecycleExpired  LifecycleStatus = "expired"
	LifecycleCritical LifecycleStatus = "critical"
)

// LifecycleResult holds the aggregated lifecycle evaluation for a single secret.
type LifecycleResult struct {
	Key        string
	Status     LifecycleStatus
	AgeStatus  AgeStatus
	TTLStatus  TTLStatus
	Strength   StrengthLevel
	Messages   []string
}

// LifecycleOptions configures thresholds for lifecycle evaluation.
type LifecycleOptions struct {
	Age    AgeOptions
	TTL    TTLOptions
	MinStrength StrengthLevel
}

// DefaultLifecycleOptions returns sensible defaults.
func DefaultLifecycleOptions() LifecycleOptions {
	return LifecycleOptions{
		Age:         DefaultAgeOptions(),
		TTL:         DefaultTTLOptions(),
		MinStrength: StrengthWeak,
	}
}

// CheckLifecycle evaluates age, TTL, and strength for a single secret value.
func CheckLifecycle(key, value string, createdAt time.Time, opts LifecycleOptions) LifecycleResult {
	result := LifecycleResult{
		Key:    key,
		Status: LifecycleOK,
	}

	ageResult := CheckAge(key, createdAt, opts.Age)
	result.AgeStatus = ageResult.Status
	if ageResult.Status == AgeExpired {
		result.Status = LifecycleExpired
		result.Messages = append(result.Messages, fmt.Sprintf("secret age expired: %s", ageResult.Message))
	} else if ageResult.Status == AgeWarning && result.Status == LifecycleOK {
		result.Status = LifecycleWarning
		result.Messages = append(result.Messages, fmt.Sprintf("secret age warning: %s", ageResult.Message))
	}

	strResult := CheckStrength(key, value)
	result.Strength = strResult.Level
	if strResult.Level < opts.MinStrength {
		if result.Status != LifecycleExpired {
			result.Status = LifecycleCritical
		}
		result.Messages = append(result.Messages, fmt.Sprintf("strength below minimum: %s", strResult.Suggestion))
	}

	return result
}

// CheckLifecycleMap evaluates lifecycle for a map of secrets.
func CheckLifecycleMap(secrets map[string]string, createdAt time.Time, opts LifecycleOptions) []LifecycleResult {
	results := make([]LifecycleResult, 0, len(secrets))
	for k, v := range secrets {
		results = append(results, CheckLifecycle(k, v, createdAt, opts))
	}
	return results
}

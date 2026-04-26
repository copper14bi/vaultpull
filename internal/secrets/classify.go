package secrets

import (
	"strings"
)

// Classification represents the sensitivity tier of a secret.
type Classification int

const (
	// ClassPublic indicates the value is safe to display openly.
	ClassPublic Classification = iota
	// ClassInternal indicates the value should be handled with care.
	ClassInternal
	// ClassConfidential indicates the value must be masked in output.
	ClassConfidential
	// ClassSecret indicates the value must never be logged or printed.
	ClassSecret
)

// String returns a human-readable label for the classification.
func (c Classification) String() string {
	switch c {
	case ClassPublic:
		return "public"
	case ClassInternal:
		return "internal"
	case ClassConfidential:
		return "confidential"
	case ClassSecret:
		return "secret"
	default:
		return "unknown"
	}
}

// secretKeywords maps lowercase substrings to their classification tier.
var secretKeywords = map[string]Classification{
	"token":    ClassSecret,
	"password": ClassSecret,
	"passwd":   ClassSecret,
	"secret":   ClassSecret,
	"private":  ClassSecret,
	"apikey":   ClassSecret,
	"api_key":  ClassSecret,
	"auth":     ClassConfidential,
	"cert":     ClassConfidential,
	"key":      ClassConfidential,
	"dsn":      ClassConfidential,
	"url":      ClassInternal,
	"host":     ClassInternal,
	"endpoint": ClassInternal,
	"port":     ClassPublic,
	"env":      ClassPublic,
	"debug":    ClassPublic,
}

// Classify returns the Classification for a given key name.
// It checks for known substrings (case-insensitive) and returns the
// highest matching tier, defaulting to ClassInternal.
func Classify(key string) Classification {
	lower := strings.ToLower(key)
	best := ClassInternal
	for substr, class := range secretKeywords {
		if strings.Contains(lower, substr) {
			if class > best {
				best = class
			}
		}
	}
	return best
}

// ClassifyMap returns a map of key → Classification for all provided secrets.
func ClassifyMap(secrets map[string]string) map[string]Classification {
	result := make(map[string]Classification, len(secrets))
	for k := range secrets {
		result[k] = Classify(k)
	}
	return result
}

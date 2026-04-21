// Package secrets provides utilities for handling sensitive secret values,
// including redaction for safe display in logs and output.
package secrets

import "strings"

const (
	redactedPlaceholder = "[REDACTED]"
	visibleSuffixLen    = 4
	minLenForPartial    = 8
)

// RedactMode controls how secret values are redacted.
type RedactMode int

const (
	// RedactFull replaces the entire value with [REDACTED].
	RedactFull RedactMode = iota
	// RedactPartial shows the last few characters with asterisks.
	RedactPartial
)

// Redact returns a redacted version of the given secret value.
// If mode is RedactPartial and the value is long enough, the last
// visibleSuffixLen characters are preserved; otherwise full redaction is used.
func Redact(value string, mode RedactMode) string {
	if value == "" {
		return ""
	}
	if mode == RedactPartial && len(value) >= minLenForPartial {
		mask := strings.Repeat("*", len(value)-visibleSuffixLen)
		return mask + value[len(value)-visibleSuffixLen:]
	}
	return redactedPlaceholder
}

// RedactMap returns a copy of the provided map with all values redacted
// according to the given mode. Keys are preserved as-is.
func RedactMap(secrets map[string]string, mode RedactMode) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = Redact(v, mode)
	}
	return result
}

// IsSensitiveKey reports whether a key name looks like it holds a sensitive
// value, based on common naming conventions.
func IsSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	sensitiveTerms := []string{
		"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE",
		"CREDENTIAL", "AUTH", "API_KEY", "APIKEY",
	}
	for _, term := range sensitiveTerms {
		if strings.Contains(upper, term) {
			return true
		}
	}
	return false
}

package secrets

import (
	"fmt"
	"regexp"
	"strings"
)

// ScanResult holds the result of scanning a map of key-value pairs for
// potentially leaked or plaintext secrets.
type ScanResult struct {
	Key     string
	Value   string
	Reason  string
}

// patterns that suggest a value may be a raw secret that should not appear
// in plaintext in certain contexts (e.g. logs, output).
var suspiciousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^[A-Za-z0-9+/]{40,}={0,2}$`),           // base64-like
	regexp.MustCompile(`(?i)^[0-9a-f]{32,}$`),                       // hex token / hash
	regexp.MustCompile(`(?i)^s\.[A-Za-z0-9]{24,}$`),                 // Vault token prefix
	regexp.MustCompile(`(?i)^(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]+$`), // GitHub tokens
	regexp.MustCompile(`(?i)^sk-[A-Za-z0-9]{32,}$`),                 // OpenAI-style keys
}

// Scan inspects a map of environment variables and returns any entries that
// look like raw secrets or match sensitive key names with non-empty values.
func Scan(env map[string]string) []ScanResult {
	var results []ScanResult

	for k, v := range env {
		if v == "" {
			continue
		}

		// Flag by key name first.
		if IsSensitiveKey(k) {
			results = append(results, ScanResult{
				Key:    k,
				Value:  Redact(v, true),
				Reason: "sensitive key name",
			})
			continue
		}

		// Flag by value pattern.
		if reason, ok := matchesSuspicious(v); ok {
			results = append(results, ScanResult{
				Key:    k,
				Value:  Redact(v, true),
				Reason: reason,
			})
		}
	}

	return results
}

// Summary returns a human-readable summary string for a slice of ScanResults.
func Summary(results []ScanResult) string {
	if len(results) == 0 {
		return "no sensitive values detected"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d sensitive value(s) detected:\n", len(results))
	for _, r := range results {
		fmt.Fprintf(&sb, "  %-30s %s (%s)\n", r.Key, r.Value, r.Reason)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func matchesSuspicious(v string) (string, bool) {
	for _, re := range suspiciousPatterns {
		if re.MatchString(v) {
			return fmt.Sprintf("matches pattern %s", re.String()), true
		}
	}
	return "", false
}

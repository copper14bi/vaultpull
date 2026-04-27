// Package lint provides validation rules for secret key naming conventions
// and .env file structure to help teams maintain consistent secret hygiene.
package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint violation.
type Finding struct {
	Key      string
	Rule     string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s (%s)", f.Severity, f.Key, f.Message, f.Rule)
}

var (
	validKeyPattern   = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	placeholderValues = []string{"changeme", "todo", "fixme", "placeholder", "example", "test", "dummy", "replace"}
)

// Lint runs all rules against the provided key/value map and returns findings.
func Lint(secrets map[string]string) []Finding {
	var findings []Finding
	for k, v := range secrets {
		findings = append(findings, checkKeyFormat(k)...)
		findings = append(findings, checkEmptyValue(k, v)...)
		findings = append(findings, checkPlaceholder(k, v)...)
		findings = append(findings, checkKeyLength(k)...)
	}
	return findings
}

func checkKeyFormat(key string) []Finding {
	if !validKeyPattern.MatchString(key) {
		return []Finding{{
			Key:      key,
			Rule:     "key-format",
			Message:  "key must be uppercase with underscores only (e.g. MY_SECRET)",
			Severity: SeverityError,
		}}
	}
	return nil
}

func checkEmptyValue(key, value string) []Finding {
	if strings.TrimSpace(value) == "" {
		return []Finding{{
			Key:      key,
			Rule:     "empty-value",
			Message:  "secret value is empty",
			Severity: SeverityWarning,
		}}
	}
	return nil
}

func checkPlaceholder(key, value string) []Finding {
	lower := strings.ToLower(strings.TrimSpace(value))
	for _, p := range placeholderValues {
		if lower == p {
			return []Finding{{
				Key:      key,
				Rule:     "placeholder-value",
				Message:  fmt.Sprintf("value appears to be a placeholder (%q)", value),
				Severity: SeverityError,
			}}
		}
	}
	return nil
}

func checkKeyLength(key string) []Finding {
	if len(key) > 64 {
		return []Finding{{
			Key:      key,
			Rule:     "key-length",
			Message:  "key exceeds 64 characters",
			Severity: SeverityInfo,
		}}
	}
	return nil
}

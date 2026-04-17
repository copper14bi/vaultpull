// Package filter provides key filtering for secrets pulled from Vault.
package filter

import (
	"strings"
)

// Rule defines inclusion/exclusion criteria for secret keys.
type Rule struct {
	Include []string // prefixes or exact keys to include (empty = include all)
	Exclude []string // prefixes or exact keys to exclude
}

// Filter applies include/exclude rules to a map of secrets.
type Filter struct {
	rule Rule
}

// New creates a Filter from the given Rule.
func New(r Rule) *Filter {
	return &Filter{rule: r}
}

// Apply returns a new map containing only keys that pass the filter rules.
// Exclusions take priority over inclusions.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if f.excluded(k) {
			continue
		}
		if f.included(k) {
			out[k] = v
		}
	}
	return out
}

func (f *Filter) included(key string) bool {
	if len(f.rule.Include) == 0 {
		return true
	}
	for _, pattern := range f.rule.Include {
		if matchPattern(key, pattern) {
			return true
		}
	}
	return false
}

func (f *Filter) excluded(key string) bool {
	for _, pattern := range f.rule.Exclude {
		if matchPattern(key, pattern) {
			return true
		}
	}
	return false
}

// matchPattern checks if key matches a pattern (prefix* or exact).
func matchPattern(key, pattern string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return key == pattern
}

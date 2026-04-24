package secrets

import (
	"math"
	"strings"
)

// EntropyResult holds the Shannon entropy score and classification for a value.
type EntropyResult struct {
	Value   string
	Entropy float64
	Weak    bool
	Reason  string
}

// EntropyThreshold is the minimum Shannon entropy considered strong for a secret.
const EntropyThreshold = 3.5

// weakPatterns are common low-entropy placeholder values.
var weakPatterns = []string{
	"changeme",
	"secret",
	"password",
	"placeholder",
	"todo",
	"fixme",
	"example",
	"test",
	"dummy",
}

// ShannonEntropy calculates the Shannon entropy of a string.
func ShannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	freq := make(map[rune]float64)
	for _, c := range s {
		freq[c]++
	}
	l := float64(len(s))
	var entropy float64
	for _, count := range freq {
		p := count / l
		entropy -= p * math.Log2(p)
	}
	return entropy
}

// CheckEntropy evaluates whether a secret value is sufficiently strong.
func CheckEntropy(value string) EntropyResult {
	result := EntropyResult{Value: value}

	if value == "" {
		result.Weak = true
		result.Reason = "empty value"
		return result
	}

	lower := strings.ToLower(value)
	for _, pat := range weakPatterns {
		if strings.Contains(lower, pat) {
			result.Weak = true
			result.Reason = "matches weak placeholder pattern"
			return result
		}
	}

	result.Entropy = ShannonEntropy(value)
	if result.Entropy < EntropyThreshold {
		result.Weak = true
		result.Reason = "low Shannon entropy"
	}
	return result
}

// CheckEntropyMap evaluates entropy for every value in a secrets map.
// Only values whose keys are considered sensitive are checked.
func CheckEntropyMap(secrets map[string]string) []EntropyResult {
	var results []EntropyResult
	for k, v := range secrets {
		if !IsSensitiveKey(k) {
			continue
		}
		r := CheckEntropy(v)
		if r.Weak {
			results = append(results, r)
		}
	}
	return results
}

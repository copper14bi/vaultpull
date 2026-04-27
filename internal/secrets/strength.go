package secrets

import (
	"strings"
	"unicode"
)

// StrengthLevel represents the assessed strength of a secret value.
type StrengthLevel int

const (
	StrengthWeak   StrengthLevel = iota // Easily guessable or placeholder
	StrengthFair                        // Meets minimal requirements
	StrengthStrong                      // Good entropy and complexity
	StrengthExcellent                   // High entropy, long, complex
)

// StrengthResult holds the outcome of a strength assessment.
type StrengthResult struct {
	Level    StrengthLevel
	Score    int // 0–100
	Reasons  []string
}

// String returns a human-readable label for the strength level.
func (s StrengthLevel) String() string {
	switch s {
	case StrengthWeak:
		return "weak"
	case StrengthFair:
		return "fair"
	case StrengthStrong:
		return "strong"
	case StrengthExcellent:
		return "excellent"
	default:
		return "unknown"
	}
}

// CheckStrength evaluates the strength of a secret value.
func CheckStrength(value string) StrengthResult {
	var reasons []string
	score := 0

	if len(value) == 0 {
		return StrengthResult{Level: StrengthWeak, Score: 0, Reasons: []string{"empty value"}}
	}

	// Length scoring
	switch {
	case len(value) >= 32:
		score += 30
	case len(value) >= 16:
		score += 20
	case len(value) >= 8:
		score += 10
	default:
		reasons = append(reasons, "too short (< 8 chars)")
	}

	// Character class scoring
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range value {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			hasSpecial = true
		}
	}
	if hasUpper {
		score += 10
	}
	if hasLower {
		score += 10
	}
	if hasDigit {
		score += 10
	}
	if hasSpecial {
		score += 15
	}

	// Entropy scoring
	e := ShannonEntropy(value)
	switch {
	case e >= 4.5:
		score += 25
	case e >= 3.5:
		score += 15
	case e >= 2.5:
		score += 5
	default:
		reasons = append(reasons, "low entropy")
	}

	// Placeholder / common weak patterns
	lower := strings.ToLower(value)
	weakPatterns := []string{"password", "secret", "changeme", "placeholder", "todo", "fixme", "example", "test"}
	for _, p := range weakPatterns {
		if strings.Contains(lower, p) {
			score -= 20
			reasons = append(reasons, "contains common weak pattern: "+p)
			break
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	var level StrengthLevel
	switch {
	case score >= 75:
		level = StrengthExcellent
	case score >= 50:
		level = StrengthStrong
	case score >= 25:
		level = StrengthFair
	default:
		level = StrengthWeak
	}

	return StrengthResult{Level: level, Score: score, Reasons: reasons}
}

// CheckStrengthMap evaluates strength for each key-value pair in a map.
func CheckStrengthMap(secrets map[string]string) map[string]StrengthResult {
	results := make(map[string]StrengthResult, len(secrets))
	for k, v := range secrets {
		results[k] = CheckStrength(v)
	}
	return results
}

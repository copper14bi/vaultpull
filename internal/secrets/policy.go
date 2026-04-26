// Package secrets provides utilities for handling, masking, scanning,
// and classifying sensitive secret values.
package secrets

import (
	"fmt"
	"strings"
)

// PolicyAction defines what to do when a policy rule matches.
type PolicyAction string

const (
	// PolicyAllow permits the secret to be written.
	PolicyAllow PolicyAction = "allow"
	// PolicyWarn permits the secret but emits a warning.
	PolicyWarn PolicyAction = "warn"
	// PolicyDeny blocks the secret from being written.
	PolicyDeny PolicyAction = "deny"
)

// PolicyRule describes a single rule applied to secret keys or values.
type PolicyRule struct {
	// KeyPattern is a prefix or substring matched against the secret key.
	KeyPattern string
	// MinEntropy is the minimum Shannon entropy required for the value.
	// Zero means no entropy check is performed.
	MinEntropy float64
	// Action is taken when this rule matches.
	Action PolicyAction
	// Reason is a human-readable explanation surfaced in violations.
	Reason string
}

// Violation records a policy rule that was triggered for a specific key.
type Violation struct {
	Key    string
	Rule   PolicyRule
	Action PolicyAction
	Reason string
}

func (v Violation) Error() string {
	return fmt.Sprintf("policy %s on key %q: %s", v.Action, v.Key, v.Reason)
}

// Policy holds an ordered list of rules evaluated against secrets.
type Policy struct {
	rules []PolicyRule
}

// NewPolicy creates a Policy from the provided rules.
// Rules are evaluated in order; the first match wins.
func NewPolicy(rules []PolicyRule) *Policy {
	return &Policy{rules: rules}
}

// Evaluate checks all entries in secrets against the policy rules.
// It returns a slice of Violations (which may be empty) and a combined
// error if any deny-action violations were found.
func (p *Policy) Evaluate(secrets map[string]string) ([]Violation, error) {
	var violations []Violation

	for key, value := range secrets {
		for _, rule := range p.rules {
			if !p.matches(rule, key, value) {
				continue
			}
			violations = append(violations, Violation{
				Key:    key,
				Rule:   rule,
				Action: rule.Action,
				Reason: rule.Reason,
			})
			break // first matching rule wins
		}
	}

	return violations, p.denyError(violations)
}

// matches returns true when the rule applies to the given key/value pair.
func (p *Policy) matches(rule PolicyRule, key, value string) bool {
	keyMatches := rule.KeyPattern == "" ||
		strings.Contains(strings.ToLower(key), strings.ToLower(rule.KeyPattern))

	if !keyMatches {
		return false
	}

	if rule.MinEntropy > 0 {
		entropy := ShannonEntropy(value)
		if entropy >= rule.MinEntropy {
			return true
		}
		// Key matched pattern but entropy threshold not reached — no violation.
		return false
	}

	return true
}

// denyError returns a combined error for all deny-action violations, or nil.
func (p *Policy) denyError(violations []Violation) error {
	var denied []string
	for _, v := range violations {
		if v.Action == PolicyDeny {
			denied = append(denied, v.Error())
		}
	}
	if len(denied) == 0 {
		return nil
	}
	return fmt.Errorf("policy violations:\n  %s", strings.Join(denied, "\n  "))
}

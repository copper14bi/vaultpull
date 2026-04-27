package secrets

import (
	"testing"
)

func TestCheckStrength_EmptyValue(t *testing.T) {
	r := CheckStrength("")
	if r.Level != StrengthWeak {
		t.Errorf("expected Weak, got %s", r.Level)
	}
	if r.Score != 0 {
		t.Errorf("expected score 0, got %d", r.Score)
	}
	if len(r.Reasons) == 0 {
		t.Error("expected at least one reason for empty value")
	}
}

func TestCheckStrength_WeakPlaceholder(t *testing.T) {
	r := CheckStrength("changeme")
	if r.Level != StrengthWeak {
		t.Errorf("expected Weak, got %s", r.Level)
	}
	found := false
	for _, reason := range r.Reasons {
		if reason != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one reason")
	}
}

func TestCheckStrength_ShortValue(t *testing.T) {
	r := CheckStrength("abc")
	if r.Level != StrengthWeak {
		t.Errorf("expected Weak for short value, got %s", r.Level)
	}
	foundReason := false
	for _, reason := range r.Reasons {
		if reason == "too short (< 8 chars)" {
			foundReason = true
		}
	}
	if !foundReason {
		t.Error("expected 'too short' reason")
	}
}

func TestCheckStrength_StrongValue(t *testing.T) {
	r := CheckStrength("G7$kP!mQ2#xLwZ9@nRtY")
	if r.Level < StrengthStrong {
		t.Errorf("expected Strong or Excellent, got %s (score=%d)", r.Level, r.Score)
	}
}

func TestCheckStrength_ExcellentValue(t *testing.T) {
	// 32+ chars, mixed case, digits, specials, high entropy
	r := CheckStrength("aB3!dE6@gH9#jK2$mN5%pQ8^sT1&vW4*")
	if r.Level != StrengthExcellent {
		t.Errorf("expected Excellent, got %s (score=%d)", r.Level, r.Score)
	}
}

func TestCheckStrength_FairValue(t *testing.T) {
	r := CheckStrength("hello123")
	if r.Level != StrengthFair && r.Level != StrengthWeak {
		t.Errorf("expected Fair or Weak for simple value, got %s", r.Level)
	}
}

func TestStrengthLevel_String(t *testing.T) {
	cases := []struct {
		level    StrengthLevel
		expected string
	}{
		{StrengthWeak, "weak"},
		{StrengthFair, "fair"},
		{StrengthStrong, "strong"},
		{StrengthExcellent, "excellent"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.expected {
			t.Errorf("StrengthLevel(%d).String() = %q, want %q", tc.level, got, tc.expected)
		}
	}
}

func TestCheckStrengthMap_ReturnsAllKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "aB3!dE6@gH9#jK2$mN5%pQ8^",
		"API_KEY":     "changeme",
		"TOKEN":       "",
	}
	results := CheckStrengthMap(secrets)
	if len(results) != len(secrets) {
		t.Errorf("expected %d results, got %d", len(secrets), len(results))
	}
	if results["TOKEN"].Level != StrengthWeak {
		t.Errorf("empty token should be Weak")
	}
	if results["DB_PASSWORD"].Level < StrengthStrong {
		t.Errorf("strong password should be Strong or Excellent, got %s", results["DB_PASSWORD"].Level)
	}
}

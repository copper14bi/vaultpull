package secrets

import (
	"testing"
)

func TestShannonEntropy_EmptyString(t *testing.T) {
	if got := ShannonEntropy(""); got != 0 {
		t.Errorf("expected 0, got %f", got)
	}
}

func TestShannonEntropy_SingleChar(t *testing.T) {
	if got := ShannonEntropy("aaaa"); got != 0 {
		t.Errorf("expected 0 for uniform string, got %f", got)
	}
}

func TestShannonEntropy_HighEntropy(t *testing.T) {
	// A random-looking token should score well above threshold.
	got := ShannonEntropy("aB3$xQ9!mZ2@kL5#")
	if got < EntropyThreshold {
		t.Errorf("expected entropy >= %f, got %f", EntropyThreshold, got)
	}
}

func TestCheckEntropy_EmptyValue(t *testing.T) {
	r := CheckEntropy("")
	if !r.Weak {
		t.Error("expected empty value to be weak")
	}
	if r.Reason != "empty value" {
		t.Errorf("unexpected reason: %s", r.Reason)
	}
}

func TestCheckEntropy_WeakPlaceholder(t *testing.T) {
	r := CheckEntropy("changeme")
	if !r.Weak {
		t.Error("expected placeholder to be weak")
	}
	if r.Reason != "matches weak placeholder pattern" {
		t.Errorf("unexpected reason: %s", r.Reason)
	}
}

func TestCheckEntropy_LowEntropy(t *testing.T) {
	r := CheckEntropy("aaaaaaa1")
	if !r.Weak {
		t.Error("expected low-entropy value to be weak")
	}
	if r.Reason != "low Shannon entropy" {
		t.Errorf("unexpected reason: %s", r.Reason)
	}
}

func TestCheckEntropy_StrongSecret(t *testing.T) {
	r := CheckEntropy("aB3$xQ9!mZ2@kL5#")
	if r.Weak {
		t.Errorf("expected strong secret to pass, reason: %s", r.Reason)
	}
}

func TestCheckEntropyMap_OnlySensitiveKeys(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_PASSWORD": "changeme",
		"APP_NAME":          "changeme", // not sensitive, should be skipped
		"API_SECRET":        "aB3$xQ9!mZ2@kL5#",
	}
	results := CheckEntropyMap(secrets)
	if len(results) != 1 {
		t.Errorf("expected 1 weak result, got %d", len(results))
	}
	if results[0].Value != "changeme" {
		t.Errorf("unexpected weak value: %s", results[0].Value)
	}
}

func TestCheckEntropyMap_EmptyMap(t *testing.T) {
	results := CheckEntropyMap(map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

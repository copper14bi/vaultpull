package secrets

import (
	"strings"
	"testing"
)

func TestScan_SensitiveKeyFlagged(t *testing.T) {
	env := map[string]string{
		"DATABASE_PASSWORD": "supersecret",
		"APP_NAME":          "vaultpull",
	}
	results := Scan(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DATABASE_PASSWORD" {
		t.Errorf("unexpected key: %s", results[0].Key)
	}
	if results[0].Reason != "sensitive key name" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestScan_VaultTokenPatternFlagged(t *testing.T) {
	env := map[string]string{
		"SOME_TOKEN": "s.AAAAABBBBBCCCCCDDDDDEEEEEFFFFFF",
	}
	results := Scan(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Reason, "matches pattern") {
		t.Errorf("expected pattern reason, got: %s", results[0].Reason)
	}
}

func TestScan_HexTokenFlagged(t *testing.T) {
	env := map[string]string{
		"RANDOM_VAR": "a3f1c2e4b5d67890a3f1c2e4b5d67890",
	}
	results := Scan(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestScan_EmptyValueSkipped(t *testing.T) {
	env := map[string]string{
		"DATABASE_PASSWORD": "",
	}
	results := Scan(env)
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty value, got %d", len(results))
	}
}

func TestScan_SafeValuesNotFlagged(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
		"PORT":     "8080",
	}
	results := Scan(env)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d: %+v", len(results), results)
	}
}

func TestSummary_NoResults(t *testing.T) {
	out := Summary(nil)
	if out != "no sensitive values detected" {
		t.Errorf("unexpected summary: %s", out)
	}
}

func TestSummary_WithResults(t *testing.T) {
	results := []ScanResult{
		{Key: "API_KEY", Value: "[REDACTED]", Reason: "sensitive key name"},
		{Key: "TOKEN", Value: "[REDACTED]", Reason: "matches pattern"},
	}
	out := Summary(results)
	if !strings.Contains(out, "2 sensitive") {
		t.Errorf("expected count in summary, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in summary")
	}
}

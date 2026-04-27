package lint_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/lint"
)

func findingsByRule(findings []lint.Finding, rule string) []lint.Finding {
	var out []lint.Finding
	for _, f := range findings {
		if f.Rule == rule {
			out = append(out, f)
		}
	}
	return out
}

func TestLint_ValidSecret_NoFindings(t *testing.T) {
	findings := lint.Lint(map[string]string{
		"MY_SECRET_KEY": "s3cur3-v@lue-123",
	})
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_LowercaseKey_ReturnsError(t *testing.T) {
	findings := lint.Lint(map[string]string{"my_key": "value"})
	matched := findingsByRule(findings, "key-format")
	if len(matched) != 1 {
		t.Fatalf("expected 1 key-format finding, got %d", len(matched))
	}
	if matched[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity, got %s", matched[0].Severity)
	}
}

func TestLint_EmptyValue_ReturnsWarning(t *testing.T) {
	findings := lint.Lint(map[string]string{"API_KEY": ""})
	matched := findingsByRule(findings, "empty-value")
	if len(matched) != 1 {
		t.Fatalf("expected 1 empty-value finding, got %d", len(matched))
	}
	if matched[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning severity, got %s", matched[0].Severity)
	}
}

func TestLint_PlaceholderValue_ReturnsError(t *testing.T) {
	for _, placeholder := range []string{"changeme", "TODO", "Placeholder", "dummy"} {
		findings := lint.Lint(map[string]string{"DB_PASS": placeholder})
		matched := findingsByRule(findings, "placeholder-value")
		if len(matched) != 1 {
			t.Errorf("placeholder %q: expected 1 finding, got %d", placeholder, len(matched))
		}
	}
}

func TestLint_KeyTooLong_ReturnsInfo(t *testing.T) {
	longKey := "A" + string(make([]byte, 64))
	for i := range longKey {
		_ = i
	}
	longKey = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" // 67 chars
	findings := lint.Lint(map[string]string{longKey: "value"})
	matched := findingsByRule(findings, "key-length")
	if len(matched) != 1 {
		t.Fatalf("expected 1 key-length finding, got %d", len(matched))
	}
	if matched[0].Severity != lint.SeverityInfo {
		t.Errorf("expected info severity, got %s", matched[0].Severity)
	}
}

func TestFinding_String_Format(t *testing.T) {
	f := lint.Finding{
		Key: "MY_KEY", Rule: "key-format",
		Message: "bad key", Severity: lint.SeverityError,
	}
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

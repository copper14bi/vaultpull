package secrets_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/secrets"
)

func TestRedact_FullMode(t *testing.T) {
	got := secrets.Redact("supersecret", secrets.RedactFull)
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestRedact_PartialMode_LongValue(t *testing.T) {
	got := secrets.Redact("abcdefgh1234", secrets.RedactPartial)
	// last 4 chars visible: "1234"
	if got != "********1234" {
		t.Errorf("unexpected partial redaction: %q", got)
	}
}

func TestRedact_PartialMode_ShortValue(t *testing.T) {
	// shorter than minLenForPartial (8), falls back to full redaction
	got := secrets.Redact("abc", secrets.RedactPartial)
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED] for short value, got %q", got)
	}
}

func TestRedact_EmptyValue(t *testing.T) {
	got := secrets.Redact("", secrets.RedactFull)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRedactMap_AllValuesRedacted(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "tok_abc123",
	}
	result := secrets.RedactMap(input, secrets.RedactFull)
	for k, v := range result {
		if v != "[REDACTED]" {
			t.Errorf("key %s: expected [REDACTED], got %q", k, v)
		}
	}
	if len(result) != len(input) {
		t.Errorf("expected %d keys, got %d", len(input), len(result))
	}
}

func TestRedactMap_OriginalUnmodified(t *testing.T) {
	input := map[string]string{"SECRET": "myvalue"}
	secrets.RedactMap(input, secrets.RedactFull)
	if input["SECRET"] != "myvalue" {
		t.Error("original map should not be modified")
	}
}

func TestIsSensitiveKey(t *testing.T) {
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"API_TOKEN", true},
		{"PRIVATE_KEY", true},
		{"AWS_SECRET_ACCESS_KEY", true},
		{"DATABASE_HOST", false},
		{"APP_ENV", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := secrets.IsSensitiveKey(tc.key)
			if got != tc.expected {
				t.Errorf("IsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

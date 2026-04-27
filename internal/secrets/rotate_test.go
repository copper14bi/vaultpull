package secrets

import (
	"strings"
	"testing"
	"time"
)

func TestRotate_DefaultOptions_ProducesCorrectLength(t *testing.T) {
	opts := DefaultRotateOptions()
	r, err := Rotate("DB_PASSWORD", "old", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.NewValue) != opts.Length {
		t.Errorf("expected length %d, got %d", opts.Length, len(r.NewValue))
	}
}

func TestRotate_PreservesKeyAndOldValue(t *testing.T) {
	opts := DefaultRotateOptions()
	r, err := Rotate("API_SECRET", "hunter2", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Key != "API_SECRET" {
		t.Errorf("expected key API_SECRET, got %q", r.Key)
	}
	if r.OldValue != "hunter2" {
		t.Errorf("expected old value hunter2, got %q", r.OldValue)
	}
}

func TestRotate_NewValueDiffersFromOld(t *testing.T) {
	opts := DefaultRotateOptions()
	r, err := Rotate("TOKEN", "old-token", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.NewValue == "old-token" {
		t.Error("new value should differ from old value")
	}
}

func TestRotate_ExpiresAtIsInFuture(t *testing.T) {
	opts := DefaultRotateOptions()
	before := time.Now().UTC()
	r, err := Rotate("SECRET", "", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.ExpiresAt.After(before) {
		t.Errorf("ExpiresAt %v should be after %v", r.ExpiresAt, before)
	}
}

func TestRotate_CustomCharset(t *testing.T) {
	charset := "abcdef0123456789"
	opts := RotateOptions{Length: 20, Charset: charset, ExpiresIn: time.Hour}
	r, err := Rotate("HEX_KEY", "", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, ch := range r.NewValue {
		if !strings.ContainsRune(charset, ch) {
			t.Errorf("character %q not in charset", ch)
		}
	}
}

func TestRotate_ZeroLength_ReturnsError(t *testing.T) {
	opts := RotateOptions{Length: 0}
	_, err := Rotate("KEY", "", opts)
	if err == nil {
		t.Error("expected error for zero length")
	}
}

func TestRotateMap_OnlyRotatesSensitiveKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "old-pass",
		"APP_ENV":     "production",
		"API_TOKEN":   "old-token",
	}
	opts := DefaultRotateOptions()
	results, err := RotateMap(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if !IsSensitiveKey(r.Key) {
			t.Errorf("non-sensitive key %q should not be rotated", r.Key)
		}
	}
	// APP_ENV is not sensitive, so at most 2 results
	if len(results) > 2 {
		t.Errorf("expected at most 2 rotated keys, got %d", len(results))
	}
}

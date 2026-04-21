package secrets

import (
	"strings"
	"testing"
)

func TestMask_FullMode(t *testing.T) {
	opts := MaskOptions{Mode: MaskFull, MaskChar: '*', VisibleLen: 2}
	got := Mask("supersecret", opts)
	if got != "***********" {
		t.Errorf("expected all asterisks, got %q", got)
	}
}

func TestMask_NoneMode(t *testing.T) {
	opts := MaskOptions{Mode: MaskNone, MaskChar: '*', VisibleLen: 2}
	got := Mask("plaintext", opts)
	if got != "plaintext" {
		t.Errorf("expected unchanged value, got %q", got)
	}
}

func TestMask_PartialMode_Normal(t *testing.T) {
	opts := DefaultMaskOptions()
	got := Mask("abcdefgh", opts)
	// expect: ab****gh
	if !strings.HasPrefix(got, "ab") || !strings.HasSuffix(got, "gh") {
		t.Errorf("unexpected partial mask: %q", got)
	}
	if strings.Contains(got[2:len(got)-2], "c") {
		t.Errorf("middle should be masked, got %q", got)
	}
}

func TestMask_PartialMode_ShortValue(t *testing.T) {
	opts := DefaultMaskOptions()
	got := Mask("ab", opts)
	// too short for partial — should be fully masked
	if got != "**" {
		t.Errorf("short value should be fully masked, got %q", got)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	opts := DefaultMaskOptions()
	got := Mask("", opts)
	if got != "" {
		t.Errorf("empty value should remain empty, got %q", got)
	}
}

func TestMaskMap_AllValuesMasked(t *testing.T) {
	secrets := map[string]string{
		"KEY_A": "valueA",
		"KEY_B": "valueB",
	}
	opts := MaskOptions{Mode: MaskFull, MaskChar: '*', VisibleLen: 2}
	result := MaskMap(secrets, opts)
	for k, v := range result {
		if strings.Contains(v, "value") {
			t.Errorf("key %s: expected masked value, got %q", k, v)
		}
	}
}

func TestMaskSensitive_OnlyMasksSensitiveKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_NAME":    "myapp",
	}
	opts := MaskOptions{Mode: MaskFull, MaskChar: '*', VisibleLen: 2}
	result := MaskSensitive(secrets, opts)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("non-sensitive key should be unchanged, got %q", result["APP_NAME"])
	}
	if result["DB_PASSWORD"] == "s3cr3t" {
		t.Errorf("sensitive key should be masked, got %q", result["DB_PASSWORD"])
	}
}

func TestDefaultMaskOptions(t *testing.T) {
	opts := DefaultMaskOptions()
	if opts.Mode != MaskPartial {
		t.Errorf("expected MaskPartial default, got %v", opts.Mode)
	}
	if opts.MaskChar != '*' {
		t.Errorf("expected '*' mask char, got %q", opts.MaskChar)
	}
	if opts.VisibleLen != 2 {
		t.Errorf("expected VisibleLen 2, got %d", opts.VisibleLen)
	}
}

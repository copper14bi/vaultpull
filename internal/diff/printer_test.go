package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintTo_Added(t *testing.T) {
	changes := []Change{{Type: Added, Key: "FOO", NewValue: "bar"}}
	var buf bytes.Buffer
	PrintTo(&buf, changes)
	if !strings.Contains(buf.String(), "+ FOO") {
		t.Errorf("expected added marker, got: %s", buf.String())
	}
}

func TestPrintTo_Removed(t *testing.T) {
	changes := []Change{{Type: Removed, Key: "OLD", OldValue: "secret"}}
	var buf bytes.Buffer
	PrintTo(&buf, changes)
	if !strings.Contains(buf.String(), "- OLD") {
		t.Errorf("expected removed marker, got: %s", buf.String())
	}
}

func TestPrintTo_Updated(t *testing.T) {
	changes := []Change{{Type: Updated, Key: "TOKEN", OldValue: "old", NewValue: "new"}}
	var buf bytes.Buffer
	PrintTo(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "~ TOKEN") {
		t.Errorf("expected updated marker, got: %s", out)
	}
	if !strings.Contains(out, "→") {
		t.Errorf("expected arrow separator, got: %s", out)
	}
}

func TestPrintTo_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	PrintTo(&buf, []Change{})
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestDisplayVal_ShortValue(t *testing.T) {
	got := displayVal("ab")
	if got != "**" {
		t.Errorf("expected '**', got %q", got)
	}
}

func TestDisplayVal_LongValue(t *testing.T) {
	got := displayVal("supersecret")
	if !strings.HasSuffix(got, "cret") {
		t.Errorf("expected suffix 'cret', got %q", got)
	}
	if !strings.HasPrefix(got, "*") {
		t.Errorf("expected leading mask, got %q", got)
	}
}

func TestDisplayVal_Empty(t *testing.T) {
	got := displayVal("")
	if got != `""` {
		t.Errorf("expected empty string repr, got %q", got)
	}
}

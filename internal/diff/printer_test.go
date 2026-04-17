package diff

import (
	"strings"
	"testing"
)

func TestPrintTo_Added(t *testing.T) {
	var buf strings.Builder
	changes := []Change{{Type: Added, Key: "FOO", NewValue: "bar"}}
	PrintTo(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "+ FOO") {
		t.Errorf("expected added marker, got: %s", out)
	}
}

func TestPrintTo_Removed(t *testing.T) {
	var buf strings.Builder
	changes := []Change{{Type: Removed, Key: "FOO", OldValue: "bar"}}
	PrintTo(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "- FOO") {
		t.Errorf("expected removed marker, got: %s", out)
	}
}

func TestPrintTo_Updated(t *testing.T) {
	var buf strings.Builder
	changes := []Change{{Type: Updated, Key: "FOO", OldValue: "old", NewValue: "new"}}
	PrintTo(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "~ FOO") {
		t.Errorf("expected updated marker, got: %s", out)
	}
	if !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Errorf("expected old and new values in output, got: %s", out)
	}
}

func TestPrintTo_NoChanges(t *testing.T) {
	var buf strings.Builder
	PrintTo(&buf, []Change{})
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message")
	}
}

func TestDisplayVal_ShortValue(t *testing.T) {
	if displayVal("hello") != "hello" {
		t.Errorf("short value should be unchanged")
	}
}

func TestDisplayVal_LongValue(t *testing.T) {
	long := strings.Repeat("x", 64)
	out := displayVal(long)
	if !strings.HasSuffix(out, "...") {
		t.Errorf("long value should be truncated with ellipsis")
	}
}

func TestDisplayVal_Empty(t *testing.T) {
	if displayVal("") != `""` {
		t.Errorf("empty value should display as quoted empty string")
	}
}

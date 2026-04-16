package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintTo_Added(t *testing.T) {
	var buf bytes.Buffer
	changes := []Change{{Key: "FOO", Type: Added, New: "bar"}}
	PrintTo(&buf, changes)
	if !strings.Contains(buf.String(), "+ FOO") {
		t.Errorf("expected added marker, got: %s", buf.String())
	}
}

func TestPrintTo_Removed(t *testing.T) {
	var buf bytes.Buffer
	changes := []Change{{Key: "FOO", Type: Removed, Old: "bar"}}
	PrintTo(&buf, changes)
	if !strings.Contains(buf.String(), "- FOO") {
		t.Errorf("expected removed marker, got: %s", buf.String())
	}
}

func TestPrintTo_Updated(t *testing.T) {
	var buf bytes.Buffer
	changes := []Change{{Key: "FOO", Type: Updated, Old: "old", New: "new"}}
	PrintTo(&buf, changes)
	if !strings.Contains(buf.String(), "~ FOO") {
		t.Errorf("expected updated marker, got: %s", buf.String())
	}
}

func TestPrintTo_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	PrintTo(&buf, []Change{})
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestDisplayVal_ShortValue(t *testing.T) {
	if got := displayVal("hello"); got != "hello" {
		t.Errorf("expected hello, got %s", got)
	}
}

func TestDisplayVal_LongValue(t *testing.T) {
	long := strings.Repeat("a", 40)
	got := displayVal(long)
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected truncation, got: %s", got)
	}
}

func TestDisplayVal_Empty(t *testing.T) {
	if got := displayVal(""); got != `""` {
		t.Errorf("expected empty quotes, got: %s", got)
	}
}

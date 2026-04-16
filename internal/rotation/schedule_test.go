package rotation

import (
	"testing"
	"time"
)

func TestShouldRotate_NeverRotated(t *testing.T) {
	p := Policy{Interval: 24 * time.Hour}
	if !p.ShouldRotate() {
		t.Error("expected ShouldRotate=true when never rotated")
	}
}

func TestShouldRotate_RecentRotation(t *testing.T) {
	p := Policy{
		Interval:    24 * time.Hour,
		LastRotated: time.Now().Add(-1 * time.Hour),
	}
	if p.ShouldRotate() {
		t.Error("expected ShouldRotate=false for recent rotation")
	}
}

func TestShouldRotate_Overdue(t *testing.T) {
	p := Policy{
		Interval:    24 * time.Hour,
		LastRotated: time.Now().Add(-25 * time.Hour),
	}
	if !p.ShouldRotate() {
		t.Error("expected ShouldRotate=true for overdue rotation")
	}
}

func TestParseInterval_Hours(t *testing.T) {
	d, err := ParseInterval("12h")
	if err != nil {
		t.Fatal(err)
	}
	if d != 12*time.Hour {
		t.Errorf("expected 12h, got %v", d)
	}
}

func TestParseInterval_Days(t *testing.T) {
	d, err := ParseInterval("7d")
	if err != nil {
		t.Fatal(err)
	}
	if d != 7*24*time.Hour {
		t.Errorf("expected 168h, got %v", d)
	}
}

func TestParseInterval_Invalid(t *testing.T) {
	_, err := ParseInterval("forever")
	if err == nil {
		t.Error("expected error for invalid interval")
	}
}

func TestNextRotation(t *testing.T) {
	last := time.Now().Add(-10 * time.Hour)
	p := Policy{Interval: 24 * time.Hour, LastRotated: last}
	next := p.NextRotation()
	expected := last.Add(24 * time.Hour)
	if next.Round(time.Second) != expected.Round(time.Second) {
		t.Errorf("NextRotation mismatch: got %v want %v", next, expected)
	}
}

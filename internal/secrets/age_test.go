package secrets

import (
	"testing"
	"time"
)

func TestCheckAge_Fresh(t *testing.T) {
	opts := DefaultAgeOptions()
	lastRotated := time.Now().Add(-10 * 24 * time.Hour) // 10 days ago

	result := CheckAge("DB_PASSWORD", lastRotated, opts)

	if result.Status != AgeFresh {
		t.Errorf("expected AgeFresh, got %s", result.Status)
	}
	if result.Key != "DB_PASSWORD" {
		t.Errorf("unexpected key: %s", result.Key)
	}
}

func TestCheckAge_Warning(t *testing.T) {
	opts := DefaultAgeOptions()
	lastRotated := time.Now().Add(-65 * 24 * time.Hour) // 65 days ago

	result := CheckAge("API_KEY", lastRotated, opts)

	if result.Status != AgeWarning {
		t.Errorf("expected AgeWarning, got %s", result.Status)
	}
	if result.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestCheckAge_Expired(t *testing.T) {
	opts := DefaultAgeOptions()
	lastRotated := time.Now().Add(-100 * 24 * time.Hour) // 100 days ago

	result := CheckAge("SECRET_TOKEN", lastRotated, opts)

	if result.Status != AgeExpired {
		t.Errorf("expected AgeExpired, got %s", result.Status)
	}
}

func TestCheckAge_ExactlyAtWarnBoundary(t *testing.T) {
	opts := DefaultAgeOptions()
	lastRotated := time.Now().Add(-60 * 24 * time.Hour)

	result := CheckAge("KEY", lastRotated, opts)

	if result.Status != AgeWarning {
		t.Errorf("expected AgeWarning at warn boundary, got %s", result.Status)
	}
}

func TestAgeStatus_String(t *testing.T) {
	cases := []struct {
		status AgeStatus
		want   string
	}{
		{AgeFresh, "fresh"},
		{AgeWarning, "warning"},
		{AgeExpired, "expired"},
		{AgeStatus(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.status.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestCheckAgeMap_ReturnsAllResults(t *testing.T) {
	opts := DefaultAgeOptions()
	timestamps := map[string]time.Time{
		"KEY_A": time.Now().Add(-5 * 24 * time.Hour),
		"KEY_B": time.Now().Add(-70 * 24 * time.Hour),
		"KEY_C": time.Now().Add(-95 * 24 * time.Hour),
	}

	results := CheckAgeMap(timestamps, opts)

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	statuses := map[string]AgeStatus{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}
	if statuses["KEY_A"] != AgeFresh {
		t.Errorf("KEY_A: expected AgeFresh, got %s", statuses["KEY_A"])
	}
	if statuses["KEY_B"] != AgeWarning {
		t.Errorf("KEY_B: expected AgeWarning, got %s", statuses["KEY_B"])
	}
	if statuses["KEY_C"] != AgeExpired {
		t.Errorf("KEY_C: expected AgeExpired, got %s", statuses["KEY_C"])
	}
}

package secrets

import (
	"testing"
	"time"
)

func TestCheckTTL_NoExpiry(t *testing.T) {
	result := CheckTTL("API_KEY", time.Time{}, DefaultTTLOptions())
	if result.Status != TTLStatusNone {
		t.Errorf("expected TTLStatusNone, got %s", result.Status)
	}
}

func TestCheckTTL_Expired(t *testing.T) {
	past := time.Now().Add(-2 * time.Hour)
	result := CheckTTL("API_KEY", past, DefaultTTLOptions())
	if result.Status != TTLStatusExpired {
		t.Errorf("expected TTLStatusExpired, got %s", result.Status)
	}
	if result.Remaining != 0 {
		t.Errorf("expected remaining=0 for expired, got %v", result.Remaining)
	}
}

func TestCheckTTL_Warning(t *testing.T) {
	soon := time.Now().Add(30 * time.Minute)
	opts := DefaultTTLOptions() // warn threshold = 24h
	result := CheckTTL("DB_PASS", soon, opts)
	if result.Status != TTLStatusWarning {
		t.Errorf("expected TTLStatusWarning, got %s", result.Status)
	}
	if result.Remaining <= 0 {
		t.Error("expected positive remaining duration")
	}
}

func TestCheckTTL_OK(t *testing.T) {
	future := time.Now().Add(72 * time.Hour)
	result := CheckTTL("SECRET_TOKEN", future, DefaultTTLOptions())
	if result.Status != TTLStatusOK {
		t.Errorf("expected TTLStatusOK, got %s", result.Status)
	}
}

func TestCheckTTL_ExactlyAtWarnBoundary(t *testing.T) {
	// Just under the threshold should be warning
	boundary := time.Now().Add(24*time.Hour - time.Second)
	result := CheckTTL("KEY", boundary, DefaultTTLOptions())
	if result.Status != TTLStatusWarning {
		t.Errorf("expected TTLStatusWarning at boundary, got %s", result.Status)
	}
}

func TestTTLStatus_String(t *testing.T) {
	cases := map[TTLStatus]string{
		TTLStatusOK:      "ok",
		TTLStatusWarning: "warning",
		TTLStatusExpired: "expired",
		TTLStatusNone:    "none",
	}
	for status, want := range cases {
		if got := status.String(); got != want {
			t.Errorf("TTLStatus(%d).String() = %q, want %q", status, got, want)
		}
	}
}

func TestCheckTTLMap_ReturnsAllResults(t *testing.T) {
	now := time.Now()
	expiries := map[string]time.Time{
		"KEY_A": now.Add(48 * time.Hour),
		"KEY_B": now.Add(-1 * time.Hour),
		"KEY_C": now.Add(6 * time.Hour),
	}
	results := CheckTTLMap(expiries, DefaultTTLOptions())
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	statuses := map[string]TTLStatus{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}
	if statuses["KEY_A"] != TTLStatusOK {
		t.Errorf("KEY_A: expected ok, got %s", statuses["KEY_A"])
	}
	if statuses["KEY_B"] != TTLStatusExpired {
		t.Errorf("KEY_B: expected expired, got %s", statuses["KEY_B"])
	}
	if statuses["KEY_C"] != TTLStatusWarning {
		t.Errorf("KEY_C: expected warning, got %s", statuses["KEY_C"])
	}
}

package secrets

import (
	"testing"
	"time"
)

func TestCheckExpiry_NoExpiry(t *testing.T) {
	result := CheckExpiry("MY_KEY", time.Time{}, DefaultExpiryOptions())
	if result.Status != ExpiryOK {
		t.Errorf("expected OK, got %s", result.Status)
	}
	if result.Message != "no expiry set" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	past := time.Now().Add(-48 * time.Hour)
	result := CheckExpiry("MY_KEY", past, DefaultExpiryOptions())
	if result.Status != ExpiryExpired {
		t.Errorf("expected Expired, got %s", result.Status)
	}
}

func TestCheckExpiry_Warning(t *testing.T) {
	soon := time.Now().Add(3 * 24 * time.Hour) // 3 days, within 7-day warn window
	result := CheckExpiry("MY_KEY", soon, DefaultExpiryOptions())
	if result.Status != ExpiryWarning {
		t.Errorf("expected Warning, got %s", result.Status)
	}
}

func TestCheckExpiry_OK(t *testing.T) {
	future := time.Now().Add(30 * 24 * time.Hour)
	result := CheckExpiry("MY_KEY", future, DefaultExpiryOptions())
	if result.Status != ExpiryOK {
		t.Errorf("expected OK, got %s", result.Status)
	}
}

func TestCheckExpiry_ExactlyAtWarnBoundary(t *testing.T) {
	opts := DefaultExpiryOptions()
	boundary := time.Now().Add(opts.WarnBefore).Add(-time.Minute)
	result := CheckExpiry("MY_KEY", boundary, opts)
	if result.Status != ExpiryWarning {
		t.Errorf("expected Warning at boundary, got %s", result.Status)
	}
}

func TestExpiryStatus_String(t *testing.T) {
	cases := []struct {
		status ExpiryStatus
		want   string
	}{
		{ExpiryOK, "ok"},
		{ExpiryWarning, "warning"},
		{ExpiryExpired, "expired"},
		{ExpiryStatus(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.status.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestCheckExpiryMap_ReturnsAllResults(t *testing.T) {
	now := time.Now()
	secrets := map[string]time.Time{
		"KEY_A": now.Add(-time.Hour),       // expired
		"KEY_B": now.Add(2 * 24 * time.Hour), // warning
		"KEY_C": now.Add(30 * 24 * time.Hour), // ok
	}
	results := CheckExpiryMap(secrets, DefaultExpiryOptions())
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	statuses := make(map[string]ExpiryStatus)
	for _, r := range results {
		statuses[r.Key] = r.Status
	}
	if statuses["KEY_A"] != ExpiryExpired {
		t.Errorf("KEY_A should be expired")
	}
	if statuses["KEY_B"] != ExpiryWarning {
		t.Errorf("KEY_B should be warning")
	}
	if statuses["KEY_C"] != ExpiryOK {
		t.Errorf("KEY_C should be ok")
	}
}

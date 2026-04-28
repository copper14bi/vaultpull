package secrets

import (
	"testing"
	"time"
)

func freshTimes() (created, rotated time.Time) {
	now := time.Now()
	return now.Add(-10 * 24 * time.Hour), now.Add(-5 * 24 * time.Hour)
}

func TestCheckDrift_NoDrift(t *testing.T) {
	created, rotated := freshTimes()
	opts := DefaultDriftOptions()
	r := CheckDrift("API_KEY", "s3cr3tV@lueXYZ123!!", created, rotated, opts)
	if r.Severity != DriftSeverityNone {
		t.Errorf("expected none, got %s: %s", r.Severity, r.Reason)
	}
}

func TestCheckDrift_EmptyValue(t *testing.T) {
	created, rotated := freshTimes()
	r := CheckDrift("DB_PASSWORD", "", created, rotated, DefaultDriftOptions())
	if r.Severity != DriftSeverityHigh {
		t.Errorf("expected high for empty value, got %s", r.Severity)
	}
}

func TestCheckDrift_NeverRotated(t *testing.T) {
	created := time.Now().Add(-5 * 24 * time.Hour)
	opts := DefaultDriftOptions()
	opts.RequireRotation = true
	r := CheckDrift("SECRET_KEY", "somevalue", created, time.Time{}, opts)
	if r.Severity != DriftSeverityMedium {
		t.Errorf("expected medium for never-rotated, got %s", r.Severity)
	}
}

func TestCheckDrift_OverMaxAge(t *testing.T) {
	created := time.Now().Add(-100 * 24 * time.Hour)
	rotated := time.Now().Add(-95 * 24 * time.Hour)
	opts := DefaultDriftOptions() // MaxAgeDays = 90
	r := CheckDrift("TOKEN", "somevalue", created, rotated, opts)
	if r.Severity == DriftSeverityNone {
		t.Error("expected non-none severity for overdue secret")
	}
}

func TestCheckDrift_CriticalAge(t *testing.T) {
	created := time.Now().Add(-200 * 24 * time.Hour)
	rotated := time.Now().Add(-185 * 24 * time.Hour)
	opts := DefaultDriftOptions()
	r := CheckDrift("TOKEN", "somevalue", created, rotated, opts)
	if r.Severity != DriftSeverityCritical {
		t.Errorf("expected critical for very old secret, got %s", r.Severity)
	}
}

func TestCheckDrift_LowEntropy(t *testing.T) {
	created, rotated := freshTimes()
	opts := DefaultDriftOptions()
	// "password" has low entropy and is a sensitive key
	r := CheckDrift("DB_PASSWORD", "aaaa", created, rotated, opts)
	if r.Severity != DriftSeverityLow {
		t.Errorf("expected low severity for low-entropy sensitive key, got %s", r.Severity)
	}
}

func TestCheckDriftMap_FiltersNone(t *testing.T) {
	created, rotated := freshTimes()
	secrets := map[string]string{
		"API_KEY": "s3cr3tV@lueXYZ123!!",
		"DB_PASSWORD": "",
	}
	results := CheckDriftMap(secrets, created, rotated, DefaultDriftOptions())
	if len(results) != 1 {
		t.Errorf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD in results, got %s", results[0].Key)
	}
}

func TestDefaultDriftOptions(t *testing.T) {
	opts := DefaultDriftOptions()
	if opts.MaxAgeDays != 90 {
		t.Errorf("expected MaxAgeDays=90, got %d", opts.MaxAgeDays)
	}
	if !opts.RequireRotation {
		t.Error("expected RequireRotation=true")
	}
}

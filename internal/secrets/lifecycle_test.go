package secrets

import (
	"testing"
	"time"
)

func TestCheckLifecycle_AllGood(t *testing.T) {
	opts := DefaultLifecycleOptions()
	createdAt := time.Now().Add(-24 * time.Hour)
	result := CheckLifecycle("API_KEY", "sUp3rStr0ngP@ssw0rd!XyZ", createdAt, opts)
	if result.Status != LifecycleOK {
		t.Errorf("expected ok, got %s: %v", result.Status, result.Messages)
	}
}

func TestCheckLifecycle_ExpiredAge(t *testing.T) {
	opts := DefaultLifecycleOptions()
	opts.Age.MaxAge = 1 * time.Hour
	createdAt := time.Now().Add(-48 * time.Hour)
	result := CheckLifecycle("DB_PASSWORD", "sUp3rStr0ngP@ssw0rd!XyZ", createdAt, opts)
	if result.Status != LifecycleExpired {
		t.Errorf("expected expired, got %s", result.Status)
	}
	if len(result.Messages) == 0 {
		t.Error("expected at least one message")
	}
}

func TestCheckLifecycle_WeakStrength(t *testing.T) {
	opts := DefaultLifecycleOptions()
	opts.MinStrength = StrengthStrong
	createdAt := time.Now().Add(-1 * time.Hour)
	result := CheckLifecycle("SECRET_KEY", "weak", createdAt, opts)
	if result.Status != LifecycleCritical {
		t.Errorf("expected critical, got %s", result.Status)
	}
}

func TestCheckLifecycle_AgeWarning(t *testing.T) {
	opts := DefaultLifecycleOptions()
	opts.Age.WarnAge = 1 * time.Hour
	opts.Age.MaxAge = 30 * 24 * time.Hour
	createdAt := time.Now().Add(-2 * time.Hour)
	result := CheckLifecycle("TOKEN", "sUp3rStr0ngP@ssw0rd!XyZ", createdAt, opts)
	if result.Status != LifecycleWarning {
		t.Errorf("expected warning, got %s", result.Status)
	}
}

func TestCheckLifecycleMap_ReturnsAllKeys(t *testing.T) {
	secrets := map[string]string{
		"KEY_A": "sUp3rStr0ngP@ssw0rd!XyZ",
		"KEY_B": "anotherStr0ngValue!99",
	}
	opts := DefaultLifecycleOptions()
	results := CheckLifecycleMap(secrets, time.Now(), opts)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestLifecycleStatus_ExpiredTakesPriorityOverCritical(t *testing.T) {
	opts := DefaultLifecycleOptions()
	opts.Age.MaxAge = 1 * time.Minute
	opts.MinStrength = StrengthExcellent
	createdAt := time.Now().Add(-1 * time.Hour)
	result := CheckLifecycle("OLD_WEAK", "weak", createdAt, opts)
	// expired is set first; critical path checks for non-expired
	if result.Status != LifecycleExpired {
		t.Errorf("expected expired to take priority, got %s", result.Status)
	}
	if len(result.Messages) < 2 {
		t.Errorf("expected messages for both age and strength, got %d", len(result.Messages))
	}
}

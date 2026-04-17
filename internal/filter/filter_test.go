package filter_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/filter"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"AWS_KEY":     "AKIA123",
		"AWS_SECRET":  "abc",
		"APP_DEBUG":   "true",
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	f := filter.New(filter.Rule{})
	out := f.Apply(baseSecrets())
	if len(out) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(out))
	}
}

func TestApply_IncludePrefix(t *testing.T) {
	f := filter.New(filter.Rule{Include: []string{"DB_*"}})
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
}

func TestApply_ExcludePrefix(t *testing.T) {
	f := filter.New(filter.Rule{Exclude: []string{"AWS_*"}})
	out := f.Apply(baseSecrets())
	if _, ok := out["AWS_KEY"]; ok {
		t.Error("AWS_KEY should be excluded")
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
}

func TestApply_ExcludeTakesPriority(t *testing.T) {
	f := filter.New(filter.Rule{
		Include: []string{"DB_*"},
		Exclude: []string{"DB_PASSWORD"},
	})
	out := f.Apply(baseSecrets())
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("DB_HOST should be included")
	}
}

func TestApply_ExactMatch(t *testing.T) {
	f := filter.New(filter.Rule{Include: []string{"APP_DEBUG"}})
	out := f.Apply(baseSecrets())
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	f := filter.New(filter.Rule{Include: []string{"DB_*"}})
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(out))
	}
}

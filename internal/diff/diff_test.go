package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	new_ := map[string]string{"FOO": "bar"}
	r := Compare(old, new_)
	if len(r.Changes) != 1 || r.Changes[0].Action != "added" {
		t.Errorf("expected 1 added change, got %+v", r.Changes)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new_ := map[string]string{}
	r := Compare(old, new_)
	if len(r.Changes) != 1 || r.Changes[0].Action != "removed" {
		t.Errorf("expected 1 removed change, got %+v", r.Changes)
	}
}

func TestCompare_Updated(t *testing.T) {
	old := map[string]string{"FOO": "old"}
	new_ := map[string]string{"FOO": "new"}
	r := Compare(old, new_)
	if len(r.Changes) != 1 || r.Changes[0].Action != "updated" {
		t.Errorf("expected 1 updated change, got %+v", r.Changes)
	}
	if r.Changes[0].Old != "old" || r.Changes[0].New != "new" {
		t.Errorf("unexpected old/new values: %+v", r.Changes[0])
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new_ := map[string]string{"FOO": "bar"}
	r := Compare(old, new_)
	if len(r.Changes) != 1 || r.Changes[0].Action != "unchanged" {
		t.Errorf("expected 1 unchanged change, got %+v", r.Changes)
	}
}

func TestHasChanges_True(t *testing.T) {
	r := &Result{Changes: []Change{{Action: "added"}}}
	if !r.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestHasChanges_False(t *testing.T) {
	r := &Result{Changes: []Change{{Action: "unchanged"}}}
	if r.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestSummary_Counts(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new_ := map[string]string{"A": "99", "C": "3"}
	r := Compare(old, new_)
	s := r.Summary()
	if s["updated"] != 1 || s["added"] != 1 || s["removed"] != 1 {
		t.Errorf("unexpected summary: %v", s)
	}
}

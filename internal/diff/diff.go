package diff

// Change represents a single secret key change.
type Change struct {
	Key    string
	Old    string
	New    string
	Action string // "added", "updated", "removed", "unchanged"
}

// Result holds the full diff between old and new env maps.
type Result struct {
	Changes []Change
}

// HasChanges returns true if any keys were added, updated, or removed.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Action != "unchanged" {
			return true
		}
	}
	return false
}

// Summary returns counts by action type.
func (r *Result) Summary() map[string]int {
	summary := map[string]int{"added": 0, "updated": 0, "removed": 0, "unchanged": 0}
	for _, c := range r.Changes {
		summary[c.Action]++
	}
	return summary
}

// Compare computes the diff between oldEnv and newEnv maps.
func Compare(oldEnv, newEnv map[string]string) *Result {
	result := &Result{}

	for key, newVal := range newEnv {
		if oldVal, exists := oldEnv[key]; !exists {
			result.Changes = append(result.Changes, Change{Key: key, Old: "", New: newVal, Action: "added"})
		} else if oldVal != newVal {
			result.Changes = append(result.Changes, Change{Key: key, Old: oldVal, New: newVal, Action: "updated"})
		} else {
			result.Changes = append(result.Changes, Change{Key: key, Old: oldVal, New: newVal, Action: "unchanged"})
		}
	}

	for key, oldVal := range oldEnv {
		if _, exists := newEnv[key]; !exists {
			result.Changes = append(result.Changes, Change{Key: key, Old: oldVal, New: "", Action: "removed"})
		}
	}

	return result
}

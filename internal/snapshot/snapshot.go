// Package snapshot captures and compares point-in-time states of secret sets,
// enabling drift detection between Vault and local .env files.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a captured state of secrets at a point in time.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Path       string            `json:"path"`
	Secrets    map[string]string `json:"secrets"`
}

// New creates a new Snapshot for the given vault path and secret map.
func New(vaultPath string, secrets map[string]string) *Snapshot {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Path:       vaultPath,
		Secrets:    copy,
	}
}

// Save writes the snapshot to a JSON file at dest.
func (s *Snapshot) Save(dest string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(dest, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", dest, err)
	}
	return nil
}

// Load reads a snapshot from a JSON file.
func Load(src string) (*Snapshot, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", src, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Drift describes the difference between a saved snapshot and a current secret map.
type Drift struct {
	Added   []string
	Removed []string
	Changed []string
}

// HasDrift returns true if any keys were added, removed, or changed.
func (d Drift) HasDrift() bool {
	return len(d.Added)+len(d.Removed)+len(d.Changed) > 0
}

// Compare returns the Drift between the snapshot's secrets and the current map.
func (s *Snapshot) Compare(current map[string]string) Drift {
	var d Drift
	for k, cv := range current {
		if sv, ok := s.Secrets[k]; !ok {
			d.Added = append(d.Added, k)
		} else if sv != cv {
			d.Changed = append(d.Changed, k)
		}
	}
	for k := range s.Secrets {
		if _, ok := current[k]; !ok {
			d.Removed = append(d.Removed, k)
		}
	}
	return d
}

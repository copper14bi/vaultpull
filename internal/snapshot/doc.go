// Package snapshot provides point-in-time capture and drift detection for
// Vault secrets.
//
// # Overview
//
// A Snapshot records the full key-value state of one or more Vault secret
// paths at a specific moment. It can be persisted to disk as a JSON file and
// later reloaded to compare against the current live state, surfacing any
// keys that were added, removed, or whose values changed.
//
// # Typical usage
//
//	// Capture and save
//	snap := snapshot.New("secret/myapp", liveSecrets)
//	snap.Save(".vaultpull.snapshot.json")
//
//	// Load and compare later
//	snap, _ := snapshot.Load(".vaultpull.snapshot.json")
//	drift := snap.Compare(currentSecrets)
//	if drift.HasDrift() {
//	    // handle added / removed / changed keys
//	}
//
// # Security
//
// Snapshot files are written with mode 0600 and should be added to .gitignore
// to prevent accidental secret exposure in version control.
package snapshot

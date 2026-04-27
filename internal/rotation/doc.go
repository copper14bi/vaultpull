// Package rotation provides backup rotation and scheduling support for vaultpull.
//
// It offers two main capabilities:
//
//  1. File rotation via Rotator — copies a .env file to a timestamped backup
//     in a configurable backup directory and prunes old backups beyond a
//     configurable maximum count.
//
//  2. Rotation scheduling via Policy — determines whether enough time has
//     elapsed since the last rotation based on a configurable interval, and
//     exposes helpers for parsing human-readable duration strings (e.g. "7d",
//     "24h", "30m").
//
// Typical usage:
//
//	r := rotation.New(".backups", 5)
//	policy := rotation.Policy{
//	    Interval:    24 * time.Hour,
//	    LastRotated: lastSyncTime,
//	}
//	if policy.ShouldRotate() {
//	    if err := r.Rotate(".env"); err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// Backup files are named using the source file's base name with a UTC
// timestamp suffix in the format "20060102T150405Z", for example:
//
//	.backups/.env.20240315T083000Z
//
// When the number of existing backups exceeds the configured maximum, the
// oldest backups are removed first, keeping only the most recent N files.
package rotation

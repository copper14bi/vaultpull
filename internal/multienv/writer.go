// Package multienv provides support for writing secrets to multiple
// .env files simultaneously, each scoped to a subset of secret paths.
package multienv

import (
	"fmt"
	"path/filepath"

	"github.com/yourusername/vaultpull/internal/env"
)

// Target maps a set of secret path prefixes to an output .env file.
type Target struct {
	// OutputFile is the path to the .env file to write.
	OutputFile string
	// Prefixes filters which secret keys to include (empty means all).
	Prefixes []string
	// BackupDir is an optional directory for backup files.
	BackupDir string
}

// Writer writes secrets to multiple .env files based on target definitions.
type Writer struct {
	targets []Target
}

// New creates a Writer for the given targets.
func New(targets []Target) *Writer {
	return &Writer{targets: targets}
}

// WriteAll writes the provided secrets to each configured target file,
// filtering keys by the target's prefix list. Returns a map of output
// file path to the number of keys written.
func (w *Writer) WriteAll(secrets map[string]string) (map[string]int, error) {
	results := make(map[string]int, len(w.targets))

	for _, t := range w.targets {
		filtered := filterByPrefixes(secrets, t.Prefixes)
		if len(filtered) == 0 {
			results[t.OutputFile] = 0
			continue
		}

		wr, err := env.NewWriter(t.OutputFile, t.BackupDir)
		if err != nil {
			return results, fmt.Errorf("multienv: create writer for %q: %w", t.OutputFile, err)
		}

		if err := wr.Write(filtered); err != nil {
			return results, fmt.Errorf("multienv: write %q: %w", t.OutputFile, err)
		}

		results[t.OutputFile] = len(filtered)
	}

	return results, nil
}

// filterByPrefixes returns a copy of secrets containing only keys that
// match at least one of the given prefixes. If prefixes is empty, all
// keys are returned.
func filterByPrefixes(secrets map[string]string, prefixes []string) map[string]string {
	if len(prefixes) == 0 {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}

	out := make(map[string]string)
	for k, v := range secrets {
		for _, p := range prefixes {
			matched, _ := filepath.Match(p+"*", k)
			if matched {
				out[k] = v
				break
			}
		}
	}
	return out
}

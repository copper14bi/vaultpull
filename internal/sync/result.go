package sync

import "fmt"

// Result holds the outcome of a single secret-to-file sync operation.
type Result struct {
	SecretPath string
	EnvFile    string
	Count      int
}

// String returns a human-readable summary of the result.
func (r Result) String() string {
	return fmt.Sprintf("synced %d key(s) from %q -> %q", r.Count, r.SecretPath, r.EnvFile)
}

// Summary prints all results to a slice of strings for display.
func Summary(results []Result) []string {
	lines := make([]string, 0, len(results))
	for _, r := range results {
		lines = append(lines, r.String())
	}
	return lines
}

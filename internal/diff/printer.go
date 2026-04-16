package diff

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const maxDisplayLen = 32

// Print writes a human-readable diff summary to stdout.
func Print(changes []Change) {
	PrintTo(os.Stdout, changes)
}

// PrintTo writes a human-readable diff summary to the given writer.
func PrintTo(w io.Writer, changes []Change) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "  + %s = %s\n", c.Key, displayVal(c.New))
		case Removed:
			fmt.Fprintf(w, "  - %s = %s\n", c.Key, displayVal(c.Old))
		case Updated:
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Key, displayVal(c.Old), displayVal(c.New))
		case Unchanged:
			fmt.Fprintf(w, "    %s (unchanged)\n", c.Key)
		}
	}
}

// displayVal truncates long values and masks secrets heuristically.
func displayVal(v string) string {
	if len(v) == 0 {
		return `""`
	}
	if len(v) > maxDisplayLen {
		return v[:maxDisplayLen] + "..."
	}
	// Mask values that look like tokens or passwords.
	lower := strings.ToLower(v)
	if strings.ContainsAny(lower, "!@#$%") || len(v) >= 20 {
		return strings.Repeat("*", min(len(v), 8))
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

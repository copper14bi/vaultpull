package diff

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	 colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	maxValLen   = 32
)

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
			fmt.Fprintf(w, "%s+ %s = %s%s\n", colorGreen, c.Key, displayVal(c.NewValue), colorReset)
		case Removed:
			fmt.Fprintf(w, "%s- %s = %s%s\n", colorRed, c.Key, displayVal(c.OldValue), colorReset)
		case Updated:
			fmt.Fprintf(w, "%s~ %s: %s → %s%s\n", colorYellow, c.Key, displayVal(c.OldValue), displayVal(c.NewValue), colorReset)
		}
	}

	added, removed, updated := 0, 0, 0
	for _, c := range changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Updated:
			updated++
		}
	}
	fmt.Fprintf(w, "\nSummary: %d added, %d removed, %d updated\n", added, removed, updated)
}

func displayVal(v string) string {
	if len(v) == 0 {
		return `""`
	}
	if len(v) > maxValLen {
		return v[:maxValLen] + "..."
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ensure strings import used
var _ = strings.Contains

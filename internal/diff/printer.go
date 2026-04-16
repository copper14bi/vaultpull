package diff

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	maskLen  = 4
	maskChar = "*"
)

// Print writes a human-readable diff summary to stdout.
func Print(changes []Change) {
	PrintTo(os.Stdout, changes)
}

// PrintTo writes a human-readable diff summary to the given writer.
func PrintTo(w io.Writer, changes []Change) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "  (no changes)")
		return
	}
	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "  + %s = %s\n", c.Key, displayVal(c.NewValue))
		case Removed:
			fmt.Fprintf(w, "  - %s = %s\n", c.Key, displayVal(c.OldValue))
		case Updated:
			fmt.Fprintf(w, "  ~ %s: %s → %s\n", c.Key, displayVal(c.OldValue), displayVal(c.NewValue))
		case Unchanged:
			fmt.Fprintf(w, "    %s (unchanged)\n", c.Key)
		}
	}
}

// displayVal masks secret values, showing only the last few characters.
func displayVal(v string) string {
	if len(v) == 0 {
		return `""`
	}
	if len(v) <= maskLen {
		return strings.Repeat(maskChar, len(v))
	}
	return strings.Repeat(maskChar, len(v)-maskLen) + v[len(v)-maskLen:]
}

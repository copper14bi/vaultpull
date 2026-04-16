package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	maskValue = "***"
)

// PrintOptions controls diff output formatting.
type PrintOptions struct {
	MaskValues bool
	Writer     io.Writer
}

// Print writes a human-readable diff to the configured writer.
func Print(r *Result, opts PrintOptions) {
	changes := make([]Change, len(r.Changes))
	copy(changes, r.Changes)
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	for _, c := range changes {
		switch c.Action {
		case "added":
			val := displayVal(c.New, opts.MaskValues)
			fmt.Fprintf(opts.Writer, "  + %s=%s\n", c.Key, val)
		case "removed":
			fmt.Fprintf(opts.Writer, "  - %s\n", c.Key)
		case "updated":
			if opts.MaskValues {
				fmt.Fprintf(opts.Writer, "  ~ %s=%s\n", c.Key, maskValue)
			} else {
				fmt.Fprintf(opts.Writer, "  ~ %s=%s -> %s\n", c.Key, c.Old, c.New)
			}
		}
	}

	s := r.Summary()
	parts := []string{}
	if s["added"] > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s["added"]))
	}
	if s["updated"] > 0 {
		parts = append(parts, fmt.Sprintf("%d updated", s["updated"]))
	}
	if s["removed"] > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s["removed"]))
	}
	if len(parts) == 0 {
		fmt.Fprintln(opts.Writer, "No changes.")
	} else {
		fmt.Fprintf(opts.Writer, "Summary: %s\n", strings.Join(parts, ", "))
	}
}

func displayVal(v string, mask bool) string {
	if mask {
		return maskValue
	}
	return v
}

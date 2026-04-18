// Package output provides formatting utilities for vaultpull CLI output.
//
// It supports two output formats:
//
//   - text: human-readable lines with optional ANSI colour coding
//   - json: one JSON object per result line, suitable for machine consumption
//
// Example usage:
//
//	f := output.New(output.FormatText, false)
//	f.Print(output.Result{
//		Path:   "secret/myapp",
//		Output: ".env",
//		Keys:   12,
//	})
//
// Pass noColor=true or set the NO_COLOR environment variable to disable
// ANSI escape sequences in text mode.
package output

// Package template provides functionality for rendering secret values
// into user-defined output templates (e.g. docker-compose, shell scripts).
package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Renderer renders a Go text/template file using a map of secret key-value pairs.
type Renderer struct {
	delimLeft  string
	delimRight string
}

// New returns a Renderer with default delimiters ({{ and }}).
func New() *Renderer {
	return &Renderer{delimLeft: "{{", delimRight: "}}"}
}

// NewWithDelims returns a Renderer with custom delimiters.
func NewWithDelims(left, right string) *Renderer {
	return &Renderer{delimLeft: left, delimRight: right}
}

// RenderFile reads the template at srcPath, executes it with secrets, and
// writes the result to dstPath with the given file permissions.
func (r *Renderer) RenderFile(srcPath, dstPath string, secrets map[string]string, perm os.FileMode) error {
	src, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("template: read source %q: %w", srcPath, err)
	}

	out, err := r.Render(string(src), secrets)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dstPath, []byte(out), perm); err != nil {
		return fmt.Errorf("template: write destination %q: %w", dstPath, err)
	}
	return nil
}

// Render executes the given template text with the provided secrets map and
// returns the rendered string.
func (r *Renderer) Render(text string, secrets map[string]string) (string, error) {
	tmpl, err := template.New("vaultpull").
		Delims(r.delimLeft, r.delimRight).
		Option("missingkey=error").
		Parse(text)
	if err != nil {
		return "", fmt.Errorf("template: parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, secrets); err != nil {
		return "", fmt.Errorf("template: execute: %w", err)
	}
	return buf.String(), nil
}

package secrets

import (
	"fmt"
	"strings"
)

// MaskMode controls how values are masked in output.
type MaskMode int

const (
	// MaskFull replaces the entire value with asterisks.
	MaskFull MaskMode = iota
	// MaskPartial shows the first and last characters with asterisks in between.
	MaskPartial
	// MaskNone does not mask the value.
	MaskNone
)

// MaskOptions configures masking behaviour.
type MaskOptions struct {
	Mode       MaskMode
	MaskChar   rune
	VisibleLen int // number of chars visible at each end in partial mode
}

// DefaultMaskOptions returns sensible masking defaults.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Mode:       MaskPartial,
		MaskChar:   '*',
		VisibleLen: 2,
	}
}

// Mask applies the given options to a single value.
func Mask(value string, opts MaskOptions) string {
	if value == "" {
		return ""
	}
	switch opts.Mode {
	case MaskNone:
		return value
	case MaskFull:
		return strings.Repeat(string(opts.MaskChar), len(value))
	case MaskPartial:
		return maskPartial(value, opts)
	default:
		return strings.Repeat(string(opts.MaskChar), len(value))
	}
}

// MaskMap applies Mask to every value in the map, returning a new map.
func MaskMap(secrets map[string]string, opts MaskOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = Mask(v, opts)
	}
	return out
}

// MaskSensitive masks only keys considered sensitive, leaving others unchanged.
func MaskSensitive(secrets map[string]string, opts MaskOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if IsSensitiveKey(k) {
			out[k] = Mask(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}

func maskPartial(value string, opts MaskOptions) string {
	n := len(value)
	visible := opts.VisibleLen
	if visible < 1 {
		visible = 1
	}
	// If the value is too short to show partial, mask fully.
	if n <= visible*2 {
		return strings.Repeat(string(opts.MaskChar), n)
	}
	midLen := n - visible*2
	return fmt.Sprintf("%s%s%s",
		value[:visible],
		strings.Repeat(string(opts.MaskChar), midLen),
		value[n-visible:],
	)
}

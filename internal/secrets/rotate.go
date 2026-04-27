package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// RotateOptions controls secret rotation generation behaviour.
type RotateOptions struct {
	Length     int
	Charset    string
	ExpiresIn  time.Duration
}

// DefaultRotateOptions returns sensible defaults for secret rotation.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		Length:    32,
		Charset:   "", // empty means base64url
		ExpiresIn: 90 * 24 * time.Hour,
	}
}

// RotateResult holds the newly generated secret and its metadata.
type RotateResult struct {
	Key       string
	OldValue  string
	NewValue  string
	ExpiresAt time.Time
	RotatedAt time.Time
}

// Rotate generates a new secret value for the given key, replacing oldValue.
// If opts.Charset is empty, a URL-safe base64 token is generated.
func Rotate(key, oldValue string, opts RotateOptions) (RotateResult, error) {
	if opts.Length <= 0 {
		return RotateResult{}, fmt.Errorf("rotate: length must be positive, got %d", opts.Length)
	}

	var newVal string
	var err error

	if opts.Charset == "" {
		newVal, err = randomBase64(opts.Length)
	} else {
		newVal, err = randomFromCharset(opts.Length, opts.Charset)
	}
	if err != nil {
		return RotateResult{}, fmt.Errorf("rotate: generate secret for %q: %w", key, err)
	}

	now := time.Now().UTC()
	return RotateResult{
		Key:       key,
		OldValue:  oldValue,
		NewValue:  newVal,
		ExpiresAt: now.Add(opts.ExpiresIn),
		RotatedAt: now,
	}, nil
}

// RotateMap rotates every key in secrets whose key passes IsSensitiveKey.
func RotateMap(secrets map[string]string, opts RotateOptions) ([]RotateResult, error) {
	var results []RotateResult
	for k, v := range secrets {
		if !IsSensitiveKey(k) {
			continue
		}
		r, err := Rotate(k, v, opts)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func randomBase64(n int) (string, error) {
	// generate enough raw bytes so that base64 gives us at least n chars
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	s := base64.URLEncoding.EncodeToString(buf)
	s = strings.TrimRight(s, "=")
	if len(s) > n {
		s = s[:n]
	}
	return s, nil
}

func randomFromCharset(n int, charset string) (string, error) {
	if len(charset) == 0 {
		return "", fmt.Errorf("charset is empty")
	}
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.Grow(n)
	for _, b := range buf {
		sb.WriteByte(charset[int(b)%len(charset)])
	}
	return sb.String(), nil
}

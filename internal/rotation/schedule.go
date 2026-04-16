package rotation

import (
	"fmt"
	"time"
)

// Policy defines when rotation should occur.
type Policy struct {
	// Interval is the minimum time between rotations.
	Interval time.Duration
	// LastRotated is the timestamp of the last rotation (zero means never).
	LastRotated time.Time
}

// ShouldRotate returns true when enough time has elapsed since the last rotation.
func (p Policy) ShouldRotate() bool {
	if p.LastRotated.IsZero() {
		return true
	}
	return time.Since(p.LastRotated) >= p.Interval
}

// NextRotation returns the time at which the next rotation is due.
func (p Policy) NextRotation() time.Time {
	if p.LastRotated.IsZero() {
		return time.Now()
	}
	return p.LastRotated.Add(p.Interval)
}

// ParseInterval parses a human-readable interval string (e.g. "24h", "7d").
func ParseInterval(s string) (time.Duration, error) {
	// Support shorthand "Nd" for days.
	var days int
	if _, err := fmt.Sscanf(s, "%dd", &days); err == nil {
		return time.Duration(days) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("rotation: invalid interval %q: %w", s, err)
	}
	return d, nil
}

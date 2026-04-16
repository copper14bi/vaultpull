package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Path      string    `json:"path,omitempty"`
	Target    string    `json:"target,omitempty"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit entries to a file in JSON-lines format.
type Logger struct {
	path string
	f    *os.File
}

// NewLogger opens (or creates) the audit log file at the given path.
func NewLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{path: path, f: f}, nil
}

// Record writes an audit entry to the log file.
func (l *Logger) Record(event, vaultPath, target string, success bool, err error) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Path:      vaultPath,
		Target:    target,
		Success:   success,
	}
	if err != nil {
		e.Error = err.Error()
	}
	line, merr := json.Marshal(e)
	if merr != nil {
		return fmt.Errorf("audit: marshal entry: %w", merr)
	}
	_, werr := fmt.Fprintf(l.f, "%s\n", line)
	return werr
}

// Close closes the underlying log file.
func (l *Logger) Close() error {
	return l.f.Close()
}

package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func newBuf(format Format) (*bytes.Buffer, *Formatter) {
	buf := &bytes.Buffer{}
	return buf, NewWithWriter(buf, format, true)
}

func TestPrint_TextSuccess(t *testing.T) {
	buf, f := newBuf(FormatText)
	f.Print(Result{Path: "secret/app", Output: ".env", Keys: 5})
	out := buf.String()
	if !strings.Contains(out, "secret/app") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, ".env") {
		t.Errorf("expected output file in output, got: %s", out)
	}
	if !strings.Contains(out, "5 keys") {
		t.Errorf("expected key count in output, got: %s", out)
	}
}

func TestPrint_TextError(t *testing.T) {
	buf, f := newBuf(FormatText)
	f.Print(Result{Path: "secret/app", Output: ".env", Err: errors.New("forbidden")})
	out := buf.String()
	if !strings.Contains(out, "forbidden") {
		t.Errorf("expected error in output, got: %s", out)
	}
	if !strings.Contains(out, "✗") {
		t.Errorf("expected failure symbol, got: %s", out)
	}
}

func TestPrint_JSONSuccess(t *testing.T) {
	buf, f := newBuf(FormatJSON)
	f.Print(Result{Path: "secret/db", Output: ".env.db", Keys: 3})
	out := buf.String()
	if !strings.Contains(out, `"path":"secret/db"`) {
		t.Errorf("expected JSON path, got: %s", out)
	}
	if !strings.Contains(out, `"keys":3`) {
		t.Errorf("expected JSON keys, got: %s", out)
	}
	if !strings.Contains(out, `"error":null`) {
		t.Errorf("expected null error, got: %s", out)
	}
}

func TestPrint_JSONError(t *testing.T) {
	buf, f := newBuf(FormatJSON)
	f.Print(Result{Path: "secret/db", Output: ".env.db", Err: errors.New("timeout")})
	out := buf.String()
	if !strings.Contains(out, `"error":"timeout"`) {
		t.Errorf("expected error field, got: %s", out)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	f := New(FormatText, true)
	if f.w == nil {
		t.Error("expected non-nil writer")
	}
	if f.format != FormatText {
		t.Errorf("expected FormatText, got %s", f.format)
	}
}

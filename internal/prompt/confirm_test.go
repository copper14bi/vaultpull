package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func newTestConfirmer(input string) (*Confirmer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	c := NewWithIO(strings.NewReader(input), out)
	return c, out
}

func TestAsk_YesInput(t *testing.T) {
	c, _ := newTestConfirmer("y\n")
	ok, err := c.Ask("Continue?", false)
	if err != nil || !ok {
		t.Fatalf("expected true, nil; got %v, %v", ok, err)
	}
}

func TestAsk_NoInput(t *testing.T) {
	c, _ := newTestConfirmer("no\n")
	ok, err := c.Ask("Continue?", true)
	if err != nil || ok {
		t.Fatalf("expected false, nil; got %v, %v", ok, err)
	}
}

func TestAsk_EmptyUsesDefault_True(t *testing.T) {
	c, _ := newTestConfirmer("\n")
	ok, err := c.Ask("Continue?", true)
	if err != nil || !ok {
		t.Fatalf("expected default true; got %v, %v", ok, err)
	}
}

func TestAsk_EmptyUsesDefault_False(t *testing.T) {
	c, _ := newTestConfirmer("\n")
	ok, err := c.Ask("Continue?", false)
	if err != nil || ok {
		t.Fatalf("expected default false; got %v, %v", ok, err)
	}
}

func TestAsk_InvalidInput(t *testing.T) {
	c, _ := newTestConfirmer("maybe\n")
	_, err := c.Ask("Continue?", false)
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
}

func TestAsk_EOFUsesDefault(t *testing.T) {
	c, _ := newTestConfirmer("") // empty reader → immediate EOF
	ok, err := c.Ask("Continue?", true)
	if err != nil || !ok {
		t.Fatalf("expected default true on EOF; got %v, %v", ok, err)
	}
}

func TestAsk_PrintsHintDefaultYes(t *testing.T) {
	c, out := newTestConfirmer("y\n")
	_, _ = c.Ask("Deploy now?", true)
	if !strings.Contains(out.String(), "Y/n") {
		t.Errorf("expected Y/n hint in output, got: %s", out.String())
	}
}

func TestAsk_PrintsHintDefaultNo(t *testing.T) {
	c, out := newTestConfirmer("n\n")
	_, _ = c.Ask("Deploy now?", false)
	if !strings.Contains(out.String(), "y/N") {
		t.Errorf("expected y/N hint in output, got: %s", out.String())
	}
}

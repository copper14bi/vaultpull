package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Confirmer handles interactive user prompts.
type Confirmer struct {
	in  io.Reader
	out io.Writer
}

// New returns a Confirmer reading from stdin and writing to stdout.
func New() *Confirmer {
	return &Confirmer{in: os.Stdin, out: os.Stdout}
}

// NewWithIO returns a Confirmer with custom reader/writer (useful for tests).
func NewWithIO(in io.Reader, out io.Writer) *Confirmer {
	return &Confirmer{in: in, out: out}
}

// Ask prints a yes/no question and returns true if the user confirms.
// Defaults to the provided defaultYes value on empty input.
func (c *Confirmer) Ask(question string, defaultYes bool) (bool, error) {
	hint := "y/N"
	if defaultYes {
		hint = "Y/n"
	}
	fmt.Fprintf(c.out, "%s [%s]: ", question, hint)

	anner(c.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("reading input: %w", err)
		}
	
		return defaultYes, nil
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	switch answer {
	case "", "\n":
		return defaultYes, nil
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("unrecognised answer %q: expected y/yes or n/no", answer)
	}
}

// AskWithRetry is like Ask but re-prompts up to maxAttempts times on
// unrecognised input, returning an error only if all attempts are exhausted.
func (c *Confirmer) AskWithRetry(question string, defaultYes bool, maxAttempts int) (bool, error) {
	for i := 0; i < maxAttempts; i++ {
		ok, err := c.Ask(question, defaultYes)
		if err == nil {
			return ok, nil
		}
		fmt.Fprintf(c.out, "Invalid input: %s. Please try again.\n", err)
	}
	return false, fmt.Errorf("no valid answer provided after %d attempt(s)", maxAttempts)
}

// MustAsk is like Ask but panics on error.
func (c *Confirmer) MustAsk(question string, defaultYes bool) bool {
	ok, err := c.Ask(question, defaultYes)
	if err != nil {
		panic(err)
	}
	return ok
}

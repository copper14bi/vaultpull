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

	scanner := bufio.NewScanner(c.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("reading input: %w", err)
		}
		// EOF — treat as default
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

// MustAsk is like Ask but panics on error.
func (c *Confirmer) MustAsk(question string, defaultYes bool) bool {
	ok, err := c.Ask(question, defaultYes)
	if err != nil {
		panic(err)
	}
	return ok
}

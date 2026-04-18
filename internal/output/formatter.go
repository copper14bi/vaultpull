package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Format controls output verbosity and style.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter writes sync result summaries to an output stream.
type Formatter struct {
	w       io.Writer
	format  Format
	noColor bool
}

// New returns a Formatter writing to stdout.
func New(format Format, noColor bool) *Formatter {
	return &Formatter{w: os.Stdout, format: format, noColor: noColor}
}

// NewWithWriter returns a Formatter writing to w.
func NewWithWriter(w io.Writer, format Format, noColor bool) *Formatter {
	return &Formatter{w: w, format: format, noColor: noColor}
}

// Result holds data for a single secret path sync result.
type Result struct {
	Path    string
	Output  string
	Keys    int
	Err     error
}

// Print writes a formatted result line.
func (f *Formatter) Print(r Result) {
	if f.format == FormatJSON {
		f.printJSON(r)
		return
	}
	f.printText(r)
}

func (f *Formatter) printText(r Result) {
	if r.Err != nil {
		label := "✗"
		if f.noColor {
			fmt.Fprintf(f.w, "%s %s → %s: %v\n", label, r.Path, r.Output, r.Err)
		} else {
			fmt.Fprintf(f.w, "%s %s\n", color.RedString(label), color.RedString(r.Path+": "+r.Err.Error()))
		}
		return
	}
	label := "✔"
	if f.noColor {
		fmt.Fprintf(f.w, "%s %s → %s (%d keys)\n", label, r.Path, r.Output, r.Keys)
	} else {
		fmt.Fprintf(f.w, "%s %s → %s %s\n", color.GreenString(label),
			r.Path, r.Output, color.CyanString(fmt.Sprintf("(%d keys)", r.Keys)))
	}
}

func (f *Formatter) printJSON(r Result) {
	errStr := ""
	if r.Err != nil {
		errStr = r.Err.Error()
	}
	errStr = strings.ReplaceAll(errStr, `"`, `\"`)
	if r.Err != nil {
		fmt.Fprintf(f.w, `{"path":%q,"output":%q,"keys":%d,"error":%q}`+"\n",
			r.Path, r.Output, r.Keys, errStr)
	} else {
		fmt.Fprintf(f.w, `{"path":%q,"output":%q,"keys":%d,"error":null}`+"\n",
			r.Path, r.Output, r.Keys)
	}
}

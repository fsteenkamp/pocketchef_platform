// package twx provides extension methods for the standard library tabwriter
package twx

import (
	"fmt"
	"io"
	"text/tabwriter"
)

func NewWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
}

func SkipLineFunc(tabs int) func(tw io.Writer) {
	var s string
	for i := 0; i < tabs; i++ {
		s += "\t"
	}

	return func(tw io.Writer) {
		fmt.Fprintln(tw, s)
	}
}

func AddLine(tw io.Writer, s ...string) {
	var out string
	for i, ss := range s {
		if i == 0 {
			out = ss
		} else {
			out = fmt.Sprintf("%s\t%s", out, ss)
		}
	}

	fmt.Fprintln(tw, out)
}

func AddHeader(tw io.Writer, s ...string) {
	var underline string

	var out string
	for i, ss := range s {
		var line string
		for i := 0; i < len(ss); i++ {
			line += "━"
		}

		if i == 0 {
			out = ss
			underline = line
		} else {
			out = fmt.Sprintf("%s\t%s", out, ss)
			underline = fmt.Sprintf("%s\t%s", underline, line)
		}

	}

	fmt.Fprintln(tw, out)
	fmt.Fprintln(tw, underline)
}

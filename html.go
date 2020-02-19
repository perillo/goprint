// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/perillo/goprint/internal/goefmt"
)

// File represents an HTML formatted Go source file.
type File struct {
	Name string
	Code template.HTML
}

// render returns an HTML fragment containing the formatted Go code for the
// specified source file.  A line number is printed at the begin of each line.
func render(name string, input []byte) template.HTML {
	buf := new(bytes.Buffer)

	n := 1
	for line := range goefmt.Format(goefmt.Scan(name, input)) {
		if line == nil {
			// Empty line
			fmt.Fprintf(buf, "<span class=\"line empty\">%3d</span>\n", n)
		} else {
			fmt.Fprintf(buf, "<span class=\"line\">%3d</span> %s\n", n,
				lineToHTML(line))
		}
		n++
	}

	return template.HTML(buf.String())
}

// spanToHTML returns an HTML representation for the code span.
func spanToHTML(s *goefmt.Span) string {
	if s.Code == "" {
		// Only horizontal white space.
		return s.Whitespace
	}
	class := strings.Join(goefmt.TokenClass(s), " ")
	code := html.EscapeString(s.Code)

	return fmt.Sprintf(`<span class="%s">%s</span>%s`, class, code, s.Whitespace)
}

// lineToHTML returns an HTML representation for the code line.  The eol is not
// included.
func lineToHTML(l goefmt.Line) string {
	if l == nil {
		// Empty line.
		return ""
	}

	spans := make([]string, len(l))
	for i, span := range l {
		spans[i] = spanToHTML(span)
	}

	return strings.Join(spans, "")
}

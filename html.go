// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// HTML formatting support.

package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/perillo/goprint/internal/goefmt"
)

// spanToHTML returns an HTML representation for the code span.
func spanToHTML(s *goefmt.Span) string {
	if s.Code == "" {
		// Only horizontal white space.
		return s.Whitespace
	}
	class := strings.Join(goefmt.TokenClass(s), " ")
	code := html.EscapeString(s.Code)

	return fmt.Sprintf(
		`<span class="%s">%s</span>%s`, class, code, s.Whitespace)
}

// lineToHTML returns an HTML representation for the code line. The eol is not
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

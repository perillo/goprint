// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// formatter.go source file is responsible for formatting tokenized Go source
// files into lines.

package main

import (
	"go/token"
	"strings"
)

// Span represents a span of Go source code around a token.
type Span struct {
	Token      token.Token
	Code       string
	Whitespace string
}

// String implements the Stringer interface.
func (s *Span) String() string {
	return s.Code + s.Whitespace
}

// Line represents a line of Go source code.
type Line []*Span

// String implements the Stringer interface.
func (l Line) String() string {
	buf := make([]string, len(l))
	for i, s := range l {
		buf[i] = s.String()
	}

	return strings.Join(buf, "")
}

type formatter struct {
	tokens chan *Token
	lines  chan Line
}

func (f *formatter) run() {
	// Avoid extra allocations.
	line := make(Line, 0, 10)
	for tok := range f.tokens {
		if isAtEOL(tok) {
			// The next token will be on a new line; add this token code and
			// emit the complete line.
			line = append(line, &Span{Token: tok.Value, Code: tok.Code})
			f.lines <- line

			// Avoid extra allocations.
			line = make(Line, 0, 10)

			// Remove eol at the end of the previous line, since it is the
			// caller responsibility to add it.
			// Emit empty lines, without eol, for the remaining newline
			// characters.
			ws, n := trimEOL(tok.Whitespace)
			for i := 0; i < n-1; i++ {
				f.lines <- Line{}
			}
			if len(ws) > 0 {
				// Add a span with only horizontal white space to the start of
				// the next line.
				line = append(line, &Span{Whitespace: ws})
			}

			continue
		}
		line = append(line, &Span{tok.Value, tok.Code, tok.Whitespace})
	}
	f.lines <- line
	close(f.lines)
}

// Format formats a tokenized Go source code returning a channel with each line
// (without the eol character) of the original source code.  Each line consists
// of a sequence of source code spans for each token.
func Format(tokens chan *Token) chan Line {
	lines := make(chan Line)
	f := formatter{
		tokens: tokens,
		lines:  lines,
	}

	go f.run()

	return lines
}

// isAtEOL returns true if the token is at the end of a line.
func isAtEOL(tok *Token) bool {
	if len(tok.Whitespace) > 0 && tok.Whitespace[0] == '\n' {
		return true
	}

	return false
}

// trimEOL trims newline characters at the start of the string.  It returns the
// updated string and the number of newline characters found.
func trimEOL(s string) (string, int) {
	// NOTE(mperillo): No need to iterate over runes.
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			break
		}
		n++
	}

	return s[n:], n
}

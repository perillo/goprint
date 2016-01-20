// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// formatter.go source file is responsible for formatting tokenized Go source
// files into lines.

package goefmt

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
	in    chan *Token
	lines chan Line
	out   chan Line
}

func (f *formatter) run1() {
	// Avoid extra allocations.
	line := make(Line, 0, 10)
	for tok := range f.in {
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

func (f *formatter) run2() {
Loop:
	for line := range f.lines {
		if len(line) == 0 {
			f.out <- line

			continue
		}
		for i, span := range line {
			pos, ok := eolPosition(span)
			if !ok {
				continue
			}

			// Split current line in three parts.
			// First emit spans on the left side, including the first line of
			// the offending comment or string.
			lhs := Span{Token: span.Token, Code: span.Code[:pos]}
			f.out <- append(line[:i], &lhs)

			// Then emit additional lines in the comment or string, excluding
			// the last one.
			extra := strings.Split(span.Code[pos+1:], "\n")
			for _, code := range extra[:len(extra)-1] {
				ent := Span{Token: span.Token, Code: code}
				f.out <- Line{&ent}
			}

			// Finally emit remaining spans on the right size, including the
			// last line of the offending comment or string, adding white
			// space.
			rhs := Span{span.Token, extra[len(extra)-1], span.Whitespace}
			f.out <- append(Line{&rhs}, line[i+1:]...)

			continue Loop
		}
		f.out <- line
	}
	close(f.out)
}

// Format formats a tokenized Go source code returning a channel with each line
// (without the eol character) of the original source code.  Each line consists
// of a sequence of source code spans for each token.
func Format(tokens chan *Token) chan Line {
	lines := make(chan Line)
	out := make(chan Line)
	f := formatter{
		in:    tokens,
		lines: lines,
		out:   out,
	}

	// In the first stage we just groups together tokens in the same line.
	go f.run1()

	// In the second stage we split raw strings and general comments in several
	// lines, in case they contains the newline character.
	go f.run2()

	return out
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

// eolPosition returns the position of the eol character for the specified
// span, if available.
func eolPosition(span *Span) (int, bool) {
	if span.Token == token.COMMENT || span.Token == token.STRING {
		// Only general comments and raw strings may span multiple
		// lines.
		idx := strings.IndexByte(span.Code, '\n')
		if idx != -1 {
			return idx, true
		}
	}

	return -1, false
}

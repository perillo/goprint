// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// lexer.go source file is responsible for scanning Go source files.

package main

import (
	"go/scanner"
	"go/token"
	"strings"
)

// Token represents a Go token and associated source code, including white
// space.
type Token struct {
	off  int
	line int
	// Code is the token source code.
	// For keywords, identifiers and basic type literals, it is the token
	// literal.
	// For the auto inserted SEMICOLON operator, it is "\n".
	// For operators, it is the operator string representation.
	// For raw strings literals and general comments, carriage return
	// characters ('\r) are discarded.
	Code string
	// Whitespace is white space after the token.
	// It contains only spaces (U+0020), horizontal tabs (U+0009) and newlines
	// (U+000A).
	Whitespace string
	Value      token.Token
}

// String implements the Stringer interface.
func (t *Token) String() string {
	return t.Code + t.Whitespace
}

type lexer struct {
	input  string
	file   *token.File
	s      scanner.Scanner
	tokens chan *Token
}

func (l *lexer) run() {
	for {
		p, tok, lit := l.s.Scan()
		if tok == token.EOF {
			close(l.tokens)

			return
		}
		pos := l.file.Position(p)
		if lit == "" {
			// Operator token, excluding SEMICOLON.
			lit = tok.String()
		} else if lit == "\n" {
			// The auto inserted SEMICOLON token.
			// Remove the "\n" character since it will be present as
			// whitespace.
			lit = ""
		}
		l.tokens <- &Token{
			off:   pos.Offset,
			line:  pos.Line,
			Code:  lit,
			Value: tok,
		}
	}
}

// Scan scans the specified Go source file and returns a channel with Token.
//
// The EOF token is not returned, and the last token does not contain the "\n"
// character.
func Scan(name string, input []byte) chan *Token {
	var s scanner.Scanner
	in := make(chan *Token)
	fset := token.NewFileSet()
	file := fset.AddFile(name, fset.Base(), len(input))

	s.Init(file, input, nil, scanner.ScanComments)
	l := lexer{
		input:  string(input),
		file:   file,
		s:      s,
		tokens: in,
	}

	// In the first stage we collect tokens, their literal code and their
	// offset in the source code.
	go l.run()

	// In the second stage we add white space after each token.
	out := make(chan *Token)

	go func() {
		prev := <-in
		for cur := range in {
			ws := l.input[prev.off+len(prev.Code) : cur.off]
			// Discard '\r' in order to provide consistent data, as it is done
			// by the Go scanner with raw string literals and general comments.
			prev.Whitespace = discardCR(ws)
			out <- prev
			prev = cur
		}
		close(out)
	}()

	return out
}

// discardCR discards carriage return characters from string.
func discardCR(s string) string {
	discard := func(r rune) rune {
		if r == '\r' {
			return -1
		}

		return r
	}

	return strings.Map(discard, s)
}

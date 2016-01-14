// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// lexer.go source file is responsible for scanning Go source files.

package main

import (
	"go/scanner"
	"go/token"
)

type Token struct {
	off  int
	line int
	// code is the token source code.
	// For keywords, identifiers and basic type literals, it is the token
	// literal.
	// For the auto inserted SEMICOLON operator, it is "\n".
	// For operators, it is the operator string representation.
	// For raw strings literals and general comments, carriage return
	// characters ('\r) are discarded.
	code string
	// whitespace is white space after the token.
	// It contains only spaces (U+0020), horizontal tabs (U+0009), and
	// newlines (U+000A).
	whitespace string
	value      token.Token
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
			// Add the newline to wsp, instead of code.
			lit = ""
		}
		l.tokens <- &Token{
			off:   pos.Offset,
			line:  pos.Line,
			code:  lit,
			value: tok,
		}
	}
}

func scan(name string, input []byte) chan *Token {
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
			prev.whitespace = l.input[prev.off+len(prev.code) : cur.off]
			out <- prev
			prev = cur
		}
		close(out)
	}()

	return out
}

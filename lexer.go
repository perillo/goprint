// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Support for scanning Go source files.

package main

import (
	"go/scanner"
	"go/token"
)

type Token struct {
	line  int
	code  string
	value token.Token
}

type lexer struct {
	input  string
	file   *token.File
	s      scanner.Scanner
	tokens chan *Token
}

func (l *lexer) run() {
	off := 0
	for {
		// Ignore the returned literal and instead find the full source code
		// around the token.
		p, tok, _ := l.s.Scan()
		if tok == token.EOF {
			close(l.tokens)

			return
		}
		pos := l.file.Position(p)
		l.tokens <- &Token{
			line:  pos.Line,
			value: tok,
			code:  l.input[off:pos.Offset],
		}
		off = pos.Offset
	}
}

func scan(name string, input []byte) chan *Token {
	var s scanner.Scanner
	tokens := make(chan *Token)
	fset := token.NewFileSet()
	file := fset.AddFile(name, fset.Base(), len(input))

	s.Init(file, input, nil, scanner.ScanComments)
	l := lexer{
		input:  string(input),
		file:   file,
		s:      s,
		tokens: tokens,
	}

	go l.run()

	return tokens
}

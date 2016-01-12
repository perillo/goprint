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
	off  int
	line int
	// chunk is the source code chunk for the token, including whitespaces.
	chunk string
	value token.Token
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
		l.tokens <- &Token{
			off:   pos.Offset,
			line:  pos.Line,
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

	// In the first stage we collect tokens and their offset in the source
	// code.
	go l.run()

	// In the second stage we add source chunk for each token.
	out := make(chan *Token)

	go func() {
		prev := <-in
		for cur := range in {
			prev.chunk = l.input[prev.off:cur.off]
			out <- prev
			prev = cur
		}
		close(out)
	}()

	return out
}

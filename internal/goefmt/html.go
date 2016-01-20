// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Despite the name, the API implemented here is not limited to HTML.

package goefmt

import (
	"go/token"
)

// TokenClass returns the HTML class for the specified code span.
func TokenClass(span *Span) []string {
	// Avoid extra allocation.
	class := make([]string, 0, 2)

	switch {
	case span.Token.IsKeyword():
		class = append(class, "keyword")
	case span.Token.IsLiteral():
		switch span.Token {
		case token.IDENT:
			class = append(class, "ident")
			if isBuiltin(span.Code) {
				// NOTE(mperillo): We ignore shadowing.
				class = append(class, "builtin")
			}
		case token.INT:
			class = append(class, "literal", "int")
		case token.FLOAT:
			class = append(class, "literal", "float")
		case token.IMAG:
			class = append(class, "literal", "imag")
		case token.CHAR:
			class = append(class, "literal", "char")
		case token.STRING:
			class = append(class, "literal", "string")
		default:
			panic("invalid literal token: " + span.Token.String())
		}
	case span.Token.IsOperator():
		class = append(class, "operator")
	case span.Token == token.COMMENT:
		class = append(class, "comment")
	case span.Token == token.ILLEGAL:
		// TODO(mperillo): Handle it in the lexer?
		class = append(class, "invalid")
	default:
		panic("unknown token type: " + span.Token.String())
	}

	return class
}

// Go predeclared identifiers.
var builtins = []string{
	// constants
	"true", "false", "iota",
	// variables
	"nil",
	// functions
	"append", "cap", "close", "complex", "copy", "delete", "imag", "len",
	"make", "new", "panic", "print", "println", "real", "recover",
	// types
	"bool", "byte", "complex128", "complex64", "error", "float32", "float64",
	"int", "int16", "int32", "int64", "int8", "rune", "string", "uint",
	"uint16", "uint32", "uint64", "uint8", "uintptr",
}

// isBuiltin returns true if the identifier is builtin.
func isBuiltin(ident string) bool {
	for _, ent := range builtins {
		if ident == ent {
			return true
		}
	}

	return false
}

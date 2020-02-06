package css

import "strings"

// mkname returns a suitable name to use for a subtest from the literal string
// used in the test case.
func mkname(s string) string {
	if s == "" {
		return "<empty>"
	}

	// Replace white spaces with one of the visible characters that can be used
	// to represent whitespace provided by Unicode.
	// Undecided between \u00B7 (Middle dot) or \u2423 (Open box).
	//
	// See https://en.wikipedia.org/wiki/Whitespace_character#Substitutes.
	return strings.Replace(s, " ", "\u2423", -1)
}

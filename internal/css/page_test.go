// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package css

import "testing"

// TestPageSize tests the Value implementation for the PageSize type, when a
// valid input is provided.
func TestPageSize(t *testing.T) {
	var tests = []struct {
		literal string
		value   PageSize
	}{
		{"A4", A4},
		{"letter", Letter},
	}

	for _, test := range tests {
		var p PageSize
		err := p.Set(test.literal)
		if err != nil {
			t.Errorf("unexpected failure for %q: %v", test.literal, err)
		} else if p != test.value {
			t.Errorf("got %q, want %q", p, test.value)
		}
	}
}

// TestInvalidPageSize tests the Value implementation for the PageSize type,
// when a invalid input is provided.
func TestInvalidPageSize(t *testing.T) {
	var tests = []string{
		"", " ", "A2",
	}

	for _, test := range tests {
		var p PageSize
		err := p.Set(test)
		if err == nil {
			t.Errorf("expected failure for %q, got %q", test, p)
		}
	}
}

// TestPageMargin tests the Value implementation for the PageMargin type, when
// a valid input is provided.
func TestPageMargin(t *testing.T) {
	var tests = []struct {
		literal string
		value   PageMargin
	}{
		{"0 0 0 0", PageMargin{
			Top:    Dimension{0, NoUnit},
			Right:  Dimension{0, NoUnit},
			Bottom: Dimension{0, NoUnit},
			Left:   Dimension{0, NoUnit},
		}},
		{"10pt 10pt 10pt 10pt", PageMargin{
			Top:    Dimension{10, Point},
			Right:  Dimension{10, Point},
			Bottom: Dimension{10, Point},
			Left:   Dimension{10, Point},
		}},
	}

	for _, test := range tests {
		var p PageMargin
		err := p.Set(test.literal)
		if err != nil {
			t.Errorf("unexpected failure for %q: %v", test.literal, err)
		} else if p != test.value {
			t.Errorf("got %q, want %q", p, test.value)
		}
	}
}

// TestInvalidPageMargin tests the Value implementation for the PageMargin
// type, when a invalid input is provided.
func TestInvalidPageMargin(t *testing.T) {
	var tests = []string{
		"", " ", "10px", "10pt", "10pt 10pt", "10pt 10", "10pt pt",
		"10pt 10pt 10pt", "10pt10pt10pt10pt",
	}

	for _, test := range tests {
		var p PageMargin
		err := p.Set(test)
		if err == nil {
			t.Errorf("expected failure for %q, got %q", test, p)
		}
	}
}

// TestPageMarginString tests the Stringer implementation for the PageMargin
// type.
func TestPageMarginString(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{"10pt 10pt 10pt 10pt", "10pt"},
		{"10pt 15pt 10pt 15pt", "10pt 15pt"},
		{"10pt 20pt 30pt 20pt", "10pt 20pt 30pt"},
		{"10pt 20pt 30pt 40pt", "10pt 20pt 30pt 40pt"},
	}

	for _, test := range tests {
		var p PageMargin
		err := p.Set(test.input)
		if err != nil {
			t.Errorf("unexpected failure for %q: %v", test.input, err)

			continue
		}
		output := p.String()
		if output != test.output {
			t.Errorf("got %q, want %q", output, test.output)
		}
	}
}

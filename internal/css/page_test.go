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
		t.Run(mkname(test.literal), func(t *testing.T) {
			var p PageSize
			err := p.Set(test.literal)
			if err != nil {
				t.Errorf("expected err == nil, got %q", err)
			} else if p != test.value {
				t.Errorf("got %q, want %q", p, test.value)
			}
		})
	}
}

// TestInvalidPageSize tests the Value implementation for the PageSize type,
// when a invalid input is provided.
func TestInvalidPageSize(t *testing.T) {
	var tests = []string{
		"", " ", "A2",
	}

	for _, test := range tests {
		t.Run(mkname(test), func(t *testing.T) {
			var p PageSize
			err := p.Set(test)
			if err == nil {
				t.Errorf("expected err != nil, got p == %q", p)
			}
		})
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
		{"10pt 20pt 30pt 40pt", PageMargin{
			Top:    Dimension{10, Point},
			Right:  Dimension{20, Point},
			Bottom: Dimension{30, Point},
			Left:   Dimension{40, Point},
		}},
		{"10pt 20pt 30pt", PageMargin{
			Top:    Dimension{10, Point},
			Right:  Dimension{20, Point},
			Bottom: Dimension{30, Point},
			Left:   Dimension{20, Point},
		}},
		{"10pt 20pt", PageMargin{
			Top:    Dimension{10, Point},
			Right:  Dimension{20, Point},
			Bottom: Dimension{10, Point},
			Left:   Dimension{20, Point},
		}},
		{"10pt", PageMargin{
			Top:    Dimension{10, Point},
			Right:  Dimension{10, Point},
			Bottom: Dimension{10, Point},
			Left:   Dimension{10, Point},
		}},
	}

	for _, test := range tests {
		t.Run(mkname(test.literal), func(t *testing.T) {
			var p PageMargin
			err := p.Set(test.literal)
			if err != nil {
				t.Fatalf("expected err == nil, got %q", err)
			}
			if p != test.value {
				t.Errorf("got %q, want %q", p, test.value)
			}
		})
	}
}

// TestInvalidPageMargin tests the Value implementation for the PageMargin
// type, when a invalid input is provided.
func TestInvalidPageMargin(t *testing.T) {
	var tests = []string{
		"", " ", "10px", "10pt 10", "10pt pt", "10pt 10pt 10pt 10pt 10pt",
		"10pt10pt10pt10pt",
	}

	for _, test := range tests {
		t.Run(mkname(test), func(t *testing.T) {
			var p PageMargin
			err := p.Set(test)
			if err == nil {
				t.Errorf("expected err != nil, got p == %q", p)
			}
		})
	}
}

// TestPageMarginFull tests the Value interface implementation for the
// PageMargin type, using both the Set and String methods.
func TestPageMarginFull(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{"10pt 10pt 10pt 10pt", "10pt"},
		{"10pt 15pt 10pt 15pt", "10pt 15pt"},
		{"10pt 20pt 30pt 20pt", "10pt 20pt 30pt"},
		{"10pt 20pt 30pt 40pt", "10pt 20pt 30pt 40pt"},
	}

	// Use input for testing Set, and the output for testing String.
	for _, test := range tests {
		t.Run(mkname(test.input), func(t *testing.T) {
			var p PageMargin
			err := p.Set(test.input)
			if err != nil {
				t.Fatalf("expected err == nil, got %q", err)
			}
			output := p.String()
			if output != test.output {
				t.Errorf("got %q, want %q", output, test.output)
			}
		})
	}

	// Now use output for testing Set, and output again for testing String.
	for _, test := range tests {
		t.Run(mkname(test.output), func(t *testing.T) {
			var p PageMargin
			err := p.Set(test.output)
			if err != nil {
				t.Fatalf("expected err == nil, got %q", err)
			}
			output := p.String()
			if output != test.output {
				t.Errorf("got %q, want %q", output, test.output)
			}
		})
	}
}

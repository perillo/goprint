// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package css

import "testing"

// TestFont tests the Value implementation for the Font type, when a valid
// input is provided.
func TestFont(t *testing.T) {
	var tests = []struct {
		literal string
		value   Font
	}{
		{`"Courier" 10pt/13pt`, Font{
			Family:     "Courier",
			Size:       Dimension{10, Point},
			LineHeight: Dimension{13, Point},
		}},
		{`"Courier" 0/13pt`, Font{
			Family:     "Courier",
			Size:       Dimension{0, NoUnit},
			LineHeight: Dimension{13, Point},
		}},
		{`"Courier" 10pt/0`, Font{
			Family:     "Courier",
			Size:       Dimension{10, Point},
			LineHeight: Dimension{0, NoUnit},
		}},
		{`"Courier New" 10pt/13pt`, Font{
			Family:     "Courier New",
			Size:       Dimension{10, Point},
			LineHeight: Dimension{13, Point},
		}},
	}

	for _, test := range tests {
		var f Font
		err := f.Set(test.literal)
		if err != nil {
			t.Errorf("unexpected failure for %q: %v", test.literal, err)
		} else if f != test.value {
			t.Errorf("got %q, want %q", f, test.value)
		}
	}
}

// TestInvalidFont tests the Value implementation for the Font type, when a
// invalid input is provided.
func TestInvalidFont(t *testing.T) {
	var tests = []string{
		"", " ", "Courier", `"Courier"`, "10pt/13pt", "   10pt/13pt",
		`"Courier" 10pt`, `"Courier" 10pt/`, `"Courier" 10em/13em`,
	}

	for _, test := range tests {
		var f Font
		err := f.Set(test)
		if err == nil {
			t.Errorf("expected failure for %q, got %q", test, f)
		}
	}
}

// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package css

import (
	"fmt"
	"testing"
)

// TestDimensionScan tests the Scanner implementation for the Dimension type,
// when a valid input is provided.
func TestDimensionScan(t *testing.T) {
	var tests = []struct {
		literal string
		value   Dimension
	}{
		{"0", Dimension{0, NoUnit}},
		{"10pt", Dimension{10, Point}},
		{"0pt", Dimension{0, Point}},
		{"10pc", Dimension{10, Pica}},
		{"10inch", Dimension{10, Inch}},
		{"10mm", Dimension{10, Millimeter}},
		{"10cm", Dimension{10, Centimeter}},
	}

	for _, test := range tests {
		t.Run(mkname(test.literal), func(t *testing.T) {
			var d Dimension
			_, err := fmt.Sscan(test.literal, &d)
			if err != nil {
				t.Fatalf("expected err == nil, got %q", err)
			}
			if d != test.value {
				t.Errorf("got %q, want %q", d, test.value)
			}
		})
	}
}

// TestInvalidDimensionScan tests the Scanner implementation for the Dimension
// type, when a invalid input is provided.
func TestInvalidDimensionScan(t *testing.T) {
	var tests = []string{
		"", " ", "10..", "10px", "10", "-1", "-10pt", "inch",
	}

	for _, test := range tests {
		t.Run(mkname(test), func(t *testing.T) {
			var d Dimension
			_, err := fmt.Sscan(test, &d)
			if err == nil {
				t.Errorf("expected err != nil, got d == %q", d)
			}
		})
	}
}

// TestDimension tests the Value implementation for the Dimension type,
// when a valid input is provided.
func TestDimension(t *testing.T) {
	var tests = []struct {
		literal string
		value   Dimension
	}{
		{"0", Dimension{0, NoUnit}},
		{"10pt", Dimension{10, Point}},
		{"0pt", Dimension{0, Point}},
		{"10pc", Dimension{10, Pica}},
		{"10inch", Dimension{10, Inch}},
		{"10mm", Dimension{10, Millimeter}},
		{"10cm", Dimension{10, Centimeter}},
	}

	for _, test := range tests {
		t.Run(mkname(test.literal), func(t *testing.T) {
			var d Dimension
			err := d.Set(test.literal)
			if err != nil {
				t.Fatalf("expected err == nil, got %q", err)
			}
			if d != test.value {
				t.Errorf("got %q, want %q", d, test.value)
			}
		})
	}
}

// TestInvalidDimension tests the Value implementation for the Dimension type,
// when a invalid input is provided.
func TestInvalidDimension(t *testing.T) {
	var tests = []string{
		"", " ", "10..", "10px", "10", "-1", "-10pt", "inch",
	}

	for _, test := range tests {
		t.Run(mkname(test), func(t *testing.T) {
			var d Dimension
			err := d.Set(test)
			if err == nil {
				t.Errorf("expected err != nil, got f == %q", d)
			}
		})
	}
}

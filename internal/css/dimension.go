// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Support for CSS dimensions, as specified in
//
//	CSS Values and Units Module Level 3
//
// Only the length type is supported, with only absolute units (excluding
// pixels and quarter-millimiters).  Only real numbers are supported, as
// float64.  For zero lengths the unit identifier is optional.

package css

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Number represents a CSS number.  Only real numbers are supported.
type Number float64

// String implements the Stringer interface.
func (n Number) String() string {
	// 'g' and 5 are used to print a number with minimal decimal digits,
	// avoiding scientific notation.
	return strconv.FormatFloat(float64(n), 'g', 5, 64)
}

// Set implements the Value interface.
func (n *Number) Set(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("invalid number: %q", s)
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// TODO(mperillo): Improve error message; it can be syntax error or
		// range error (see strconv.NumErr type).
		return fmt.Errorf("invalid number: %q: %v", s, err)
	}

	*n = Number(v)

	return nil
}

// Unit represents the unit of a CSS quantity.  Only absolute length units are
// supported, excluding pixels and quarter-millimeters.
type Unit string

// Supported units.
const (
	NoUnit     Unit = ""
	Point      Unit = "pt"
	Pica       Unit = "pc"
	Inch       Unit = "inch"
	Millimeter Unit = "mm"
	Centimeter Unit = "cm"
)

// String implements the Stringer interface.
func (u Unit) String() string {
	return string(u)
}

var units = map[string]bool{
	"":     true,
	"pt":   true,
	"pc":   true,
	"inch": true,
	"mm":   true,
	"cm":   true,
}

// Set implements the Value interface.
func (u *Unit) Set(s string) error {
	if ok := units[s]; !ok {
		return fmt.Errorf("invalid unit: %q", s)
	}

	*u = Unit(s)

	return nil
}

// Dimension represents a CSS number with unit.
type Dimension struct {
	Value Number
	Unit  Unit
}

// String implements the Stringer interface.
func (d Dimension) String() string {
	return fmt.Sprintf("%v%s", d.Value, d.Unit)
}

func numberToken(ch rune) bool {
	// No negative numbers, and no scientific notation.
	return strings.ContainsRune("0123456789.", ch)
}

func unitToken(ch rune) bool {
	return unicode.IsLetter(ch)
}

func dimensionToken(ch rune) bool {
	return numberToken(ch) || unitToken(ch)
}

// Scan implements the Scanner interface.
func (d *Dimension) Scan(state fmt.ScanState, verb rune) error {
	if verb != 'v' {
		return fmt.Errorf("Dimension.Scan: invalid verb %c", verb)
	}

	// Scan the entire string <num><unit> and delegate to the Set method.
	// Ensure leading white space are not ignored.
	//
	// Dimension implements the Scanner interface only to make it more easy for
	// Font and PageMargin to implement the Value interface.
	ws := readspace(state)
	tok, err := state.Token(false, dimensionToken)
	if err != nil {
		return fmt.Errorf("Dimension.Scan: %v", err)
	}

	return d.Set(ws + string(tok))
}

// Set implements the Value interface.
func (d *Dimension) Set(s string) error {
	var v Dimension

	// Split number and unit.
	i := strings.IndexFunc(s, unitToken)
	if i < 0 {
		if err := v.Value.Set(s); err != nil {
			return fmt.Errorf("invalid dimension: %v", err)
		}
		if v.Value == 0 {
			// 0 is a valid dimension.
			return nil
		}

		return fmt.Errorf("invalid dimension: unit is required: %q", s)
	}

	if err := v.Value.Set(s[:i]); err != nil {
		return fmt.Errorf("invalid dimension: %q: %v", s, err)
	}
	if err := v.Unit.Set(s[i:]); err != nil {
		return fmt.Errorf("invalid dimension: %q: %v", s, err)
	}

	*d = v

	return nil
}

// NOTE(mperillo): The Value interface is not implemented for Dimension.

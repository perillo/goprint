// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Support for pages, as specified in
//
//	CSS Paged Media Module Level 3
//
// Only a simple syntax is supported.

package css

import (
	"fmt"
	"strings"
)

// PageSize represents a CSS page size.  Only A4 and letter page sizes are
// supported.  The page orientation is portrait and can not be changed.
type PageSize string

// Supported page sizes.
const (
	A4     PageSize = "A4"
	Letter PageSize = "letter"
)

// String implements the Stringer interface.
func (p PageSize) String() string {
	return fmt.Sprintf("%s portrait", string(p))
}

// Set implements the Value interface.
func (p *PageSize) Set(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("invalid page size: %q", s)
	}
	v := PageSize(s)
	if v != A4 && v != Letter {
		return fmt.Errorf("invalid page size: %q", s)
	}

	*p = v

	return nil
}

// PageMargin represents CSS page margins.
type PageMargin struct {
	Top    Dimension
	Right  Dimension
	Bottom Dimension
	Left   Dimension
}

// String implements the Stringer interface.
func (p PageMargin) String() string {
	switch {
	case p.Top == p.Bottom && p.Right == p.Left && p.Top == p.Right:
		return fmt.Sprintf("%v", p.Top)
	case p.Top == p.Bottom && p.Right == p.Left:
		return fmt.Sprintf("%v %v", p.Top, p.Right)
	case p.Right == p.Left:
		return fmt.Sprintf("%v %v %v", p.Top, p.Right, p.Bottom)
	}

	return fmt.Sprintf("%v %v %v %v", p.Top, p.Right, p.Bottom, p.Left)
}

// Set implements the Value interface.
func (p *PageMargin) Set(s string) error {
	var v PageMargin

	l := strings.Fields(s)
	if len(l) == 0 {
		return fmt.Errorf("invalid page margin: %q", s)
	}
	if len(l) > 0 {
		if err := v.Top.Set(l[0]); err != nil {
			return fmt.Errorf("invalid page margin: %q: %v", s, err)
		}
		v.Right = v.Top
		v.Bottom = v.Top
		v.Left = v.Top
	}
	if len(l) > 1 {
		if err := v.Right.Set(l[1]); err != nil {
			return fmt.Errorf("invalid page margin: %q: %v", s, err)
		}
		v.Left = v.Right
	}
	if len(l) > 2 {
		if err := v.Bottom.Set(l[2]); err != nil {
			return fmt.Errorf("invalid page margin: %q: %v", s, err)
		}
	}
	if len(l) > 3 {
		if err := v.Left.Set(l[3]); err != nil {
			return fmt.Errorf("invalid page margin: %q: %v", s, err)
		}
	}
	if len(l) > 4 {
		return fmt.Errorf("invalid page margin: %q", s)
	}

	*p = v

	return nil
}

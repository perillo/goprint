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
		return fmt.Errorf("page size is required")
	}
	v := PageSize(s)
	if v != A4 && v != Letter {
		return fmt.Errorf("%q is not a valid page size", s)
	}

	*p = v

	return nil
}

// PageMargin represents CSS page margins.  Top, right, bottom and left
// dimensions must all be provided.
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

	_, err := fmt.Sscanf(s, "%v %v %v %v", &v.Top, &v.Right, &v.Bottom, &v.Left)
	if err != nil {
		// TODO(mperillo): Improve error message when Right, Bottom or Left is
		// missing.
		return fmt.Errorf("%q is not a valid page margin: %v", s, err)
	}

	*p = v

	return nil
}

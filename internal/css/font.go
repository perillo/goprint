// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Support for CSS fonts, as specified in
//
//	CSS Fonts Module Level 3
//
// Only a simple syntax is supported, and it is different from the one in the
// font shorthand property.

package css

import (
	"fmt"
	"io"
)

// Font represents a CSS font.  Only family, size and line height are
// supported, and they must all be provided.  Font family must always be
// quoted.
type Font struct {
	Family     string
	Size       Dimension
	LineHeight Dimension
}

// String implements the Stringer interface.
func (f Font) String() string {
	return fmt.Sprintf("%q %v/%v", f.Family, f.Size, f.LineHeight)
}

// Set implements the Value interface.
func (f *Font) Set(s string) error {
	var v Font

	_, err := fmt.Sscanf(s, "%q %v/%v", &v.Family, &v.Size, &v.LineHeight)
	if err == io.EOF {
		return fmt.Errorf("font is required")
	} else if err == io.ErrUnexpectedEOF {
		// The '/' character was missing.
		return fmt.Errorf("%q is not a valid font: line height is required", s)
	} else if err != nil {
		// TODO(mperillo): Improve error message when Family is not quoted and
		// when Size or LineHeight is missing.
		return fmt.Errorf("%q is not a valid font: %v", s, err)
	}

	*f = v

	return nil
}

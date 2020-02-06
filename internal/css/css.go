// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package css implements the minimal CSS syntax for allowing CSS properties to
// be defined via the standard flag package.
package css

import (
	"fmt"
	"unicode"
)

// readspace skips and returns white space from a ScanState.
func readspace(state fmt.ScanState) string {
	// Ignore the error, as it is done by ScanState.SkipSpace
	tok, _ := state.Token(false, unicode.IsSpace)

	return string(tok)
}

// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packages implements support for loading packages.
package packages

import (
	"encoding/json"
	"fmt"
	"io"
)

// A Package describes a single package found in a directory.
type Package struct {
	Dir        string // directory containing package sources
	ImportPath string // import path of package in dir
	Name       string // package name

	// Source files
	GoFiles        []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles       []string // .go sources files that import "C"
	IgnoredGoFiles []string // .go sources ignored due to build constraints

	// Test information
	TestGoFiles  []string // _test.go files in package
	XTestGoFiles []string // _test.go files outside package
}

// Load loads and return the package named by the given patterns.
//
// If more than one package matches the patterns, only the first one is
// returned.
func Load(patterns ...string) (*Package, error) {
	argv := []string{"-json"}
	argv = append(argv, patterns...)
	stdout, err := invokeGo("list", argv, nil)
	if err != nil {
		return nil, err
	}

	// Decode the first package, and ignore the rest.
	pkg := new(Package)
	if err := json.NewDecoder(stdout).Decode(pkg); err == io.EOF {
		// TODO(mperillo): Should we report a custom error message if a pattern
		// was specified?
		return nil, fmt.Errorf("cannot find package %q", patterns)
	} else if err != nil {
		return nil, err
	}

	return pkg, nil
}

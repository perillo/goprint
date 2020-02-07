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
	Dir        string `json:",omitempty"` // directory containing package sources
	ImportPath string `json:",omitempty"` // import path of package in dir
	Name       string `json:",omitempty"` // package name

	// Source files
	GoFiles        []string `json:",omitempty"` // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles       []string `json:",omitempty"` // .go sources files that import "C"
	IgnoredGoFiles []string `json:",omitempty"` // .go sources ignored due to build constraints

	// Test information
	TestGoFiles  []string `json:",omitempty"` // _test.go files in package
	XTestGoFiles []string `json:",omitempty"` // _test.go files outside package
}

// Find returns the package named by the import path.
//
// If the import path is a pattern and more than one package is matched, only
// the first one is returned.
func Find(path ...string) (*Package, error) {
	argv := []string{"-json"}
	argv = append(argv, path...)
	stdout, err := invokeGo("list", argv, nil)
	if err != nil {
		return nil, err
	}

	// Decode the first package, and ignore the rest.
	pkg := new(Package)
	if err := json.NewDecoder(stdout).Decode(pkg); err == io.EOF {
		// TODO(mperillo): Should we report a custom error message if a pattern
		// was specified?
		return nil, fmt.Errorf("cannot find package %q", path)
	} else if err != nil {
		return nil, err
	}

	return pkg, nil
}

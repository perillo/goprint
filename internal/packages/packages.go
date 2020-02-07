// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packages implements support for loading packages.
package packages

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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
//
// Load returns at least one package or an error.
func Load(patterns ...string) (*Package, error) {
	argv := []string{"-json"}
	argv = append(argv, patterns...)
	stdout, err := invokeGo("list", argv, nil)
	if err != nil {
		return nil, err
	}

	// Decode the first package, and ignore the rest.
	pkg := new(Package)
	if err := json.NewDecoder(stdout).Decode(pkg); err != nil {
		return nil, fmt.Errorf("JSON decode: %v", err)
	}

	return normalize(pkg), nil
}

// normalize ensures all the source file paths are absolute, for consistency.
func normalize(pkg *Package) *Package {
	abspaths(pkg.Dir, pkg.GoFiles)
	abspaths(pkg.Dir, pkg.CgoFiles)
	abspaths(pkg.Dir, pkg.IgnoredGoFiles)
	abspaths(pkg.Dir, pkg.TestGoFiles)
	abspaths(pkg.Dir, pkg.XTestGoFiles)

	return pkg
}

func abspaths(dir string, names []string) []string {
	for i, name := range names {
		path := filepath.Join(dir, name)
		names[i] = path
	}

	return names
}

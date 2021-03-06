// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packages implements support for loading packages.
package packages

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// A Package describes a single package found in a directory.
type Package struct {
	Dir        string  // directory containing package sources
	ImportPath string  // import path of package in dir
	Name       string  // package name
	Module     *Module // info about package's containing module, if any (can be nil)

	// Source files
	GoFiles        []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles       []string // .go sources files that import "C"
	IgnoredGoFiles []string // .go sources ignored due to build constraints

	// Test information
	TestGoFiles  []string // _test.go files in package
	XTestGoFiles []string // _test.go files outside package
}

// String implements the Stringer interface.
func (p *Package) String() string {
	return p.ImportPath
}

// SourceFiles returns all the .go files, including files ignored due to build
// constraints.
func (p *Package) SourceFiles() []string {
	return concat(p.GoFiles, p.CgoFiles, p.IgnoredGoFiles)
}

// TestFiles returns all the _test.go files.
func (p *Package) TestFiles() []string {
	return concat(p.TestGoFiles, p.XTestGoFiles)
}

// Load loads and return the package named by the given pattern.
//
// If more than one package matches the pattern, only the first one is
// returned.
//
// Load returns at least one package or an error.
func Load(pattern string) (*Package, error) {
	pkglist, err := load(pattern)
	if err != nil {
		return nil, err
	}
	if len(pkglist) > 1 {
		fmt.Fprintf(os.Stderr, "warning: %q matched multiple packages\n", pattern)
	}

	return pkglist[0], nil
}

// load loads and return the packages named by the given pattern.
func load(pattern string) ([]*Package, error) {
	argv := []string{"-json"}
	if pattern != "" {
		// Don't pass an empty argument to go list.
		// See https://github.com/golang/go/issues/37300.
		argv = append(argv, pattern)
	}
	stdout, err := invokeGo("list", argv, nil)
	if err != nil {
		return nil, err
	}

	return decode(stdout)
}

func decode(r io.Reader) ([]*Package, error) {
	pkglist := make([]*Package, 0, 10)
	for dec := json.NewDecoder(r); dec.More(); {
		pkg := new(Package)
		if err := dec.Decode(pkg); err != nil {
			return nil, fmt.Errorf("JSON decode: %w", err)
		}

		pkg = normalize(pkg)
		pkglist = append(pkglist, pkg)
	}

	return pkglist, nil
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

// concat concatenates args into a single []string.  The resulting slice is
// sorted.
func concat(args ...[]string) []string {
	var buf []string
	for _, arg := range args {
		buf = append(buf, arg...)
	}
	sort.Strings(buf)

	return buf
}

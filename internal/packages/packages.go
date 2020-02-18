// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packages implements support for loading packages.
package packages

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"time"
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

// SourceFiles returns all the .go files, including files ignored due to build
// constraints.
func (p *Package) SourceFiles() []string {
	return concat(p.GoFiles, p.CgoFiles, p.IgnoredGoFiles)
}

// TestFiles returns all the _test.go files.
func (p *Package) TestFiles() []string {
	return concat(p.TestGoFiles, p.XTestGoFiles)
}

// A Module describes a package's containing module.
type Module struct {
	Path    string     // module path
	Version string     // module version
	Time    *time.Time // time version was created
}

// Date returns the date when the module version was created, or the current
// date if it is not available.
func (m *Module) Date() string {
	const unixDate = "Mon Jan _2 2006"

	if m == nil {
		// Ensure it does not panic when modules are not supported.
		return ""
	}
	if m.Time == nil {
		// Return the current date.
		return time.Now().Format(unixDate)
	}

	return m.Time.Format(unixDate)
}

// String implements the Stringer interface.
func (m *Module) String() string {
	if m == nil {
		// Ensure it does not panic when modules are not supported.
		return ""
	}

	s := m.Path
	if m.Version != "" {
		s += "@" + m.Version
	}

	return s
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

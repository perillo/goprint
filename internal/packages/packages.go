// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packages implements support for loading packages.
package packages

import (
	"encoding/json"
	"fmt"
	"io"
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
	Path     string     // module path
	Version  string     // module version
	Time     *time.Time // time version was created
	Dir      string     // directory holding files for this module, if any
	Packages []*Package // packages belonging to the module

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

// Load loads and return the package named by the given pattern.
//
// If more than one package matches the pattern, only the first one is
// returned.
//
// Load returns at least one package or an error.
func Load(pattern string) (*Package, error) {
	argv := []string{"-json"}
	argv = append(argv, pattern)
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

// LoadModule loads and return the module named by path and all its packages.
func LoadModule(path string) (*Module, error) {
	mod, err := loadModule(path)
	if err != nil {
		return nil, err
	}

	pkglist, err := loadPackages(mod)
	if err != nil {
		return nil, err
	}
	mod.Packages = pkglist

	return mod, nil
}

// loadModule loads and return the module named by the given path.
func loadModule(path string) (*Module, error) {
	argv := []string{"-m", "-json"}
	if path != "" {
		// Don't pass an empty module path to go list -m.
		// See https://github.com/golang/go/issues/37300.
		argv = append(argv, path)
	}
	stdout, err := invokeGo("list", argv, nil)
	if err != nil {
		return nil, err
	}

	// Decode the module; there is only one.
	mod := new(Module)
	if err := json.NewDecoder(stdout).Decode(mod); err != nil {
		return nil, fmt.Errorf("JSON decode: %v", err)
	}

	return mod, nil
}

// loadPackages loads and return all the package of the given module mod.
func loadPackages(mod *Module) ([]*Package, error) {
	attr := attr{
		Dir: mod.Dir,
	}
	argv := []string{"-json", "./..."}
	stdout, err := invokeGo("list", argv, &attr)
	if err != nil {
		return nil, err
	}

	pkglist, err := decode(stdout)
	if err != nil {
		return nil, err
	}

	return pkglist, nil
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

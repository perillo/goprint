// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packages

import (
	"encoding/json"
	"fmt"
	"time"
)

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

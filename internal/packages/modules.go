// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packages

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

// LoadModule loads and return the module named by pattern and all its
// packages.
func LoadModule(pattern string) (*Module, error) {
	modlist, err := loadm(pattern)
	if err != nil {
		return nil, err
	}
	if len(modlist) > 1 {
		fmt.Fprintf(os.Stderr, "warning: %q matched multiple modules\n", pattern)
	}

	mod := modlist[0]
	pkglist, err := loadPackages(mod)
	if err != nil {
		return nil, err
	}
	mod.Packages = pkglist

	return mod, nil
}

// loadm loads and return the modules named by the given pattern.
func loadm(pattern string) ([]*Module, error) {
	argv := []string{"-m", "-json"}
	if pattern != "" {
		// Don't pass an empty argument to go list -m.
		// See https://github.com/golang/go/issues/37300.
		argv = append(argv, pattern)
	}
	stdout, err := invokeGo("list", argv, nil)
	if err != nil {
		return nil, err
	}

	return decodem(stdout)
}

func decodem(r io.Reader) ([]*Module, error) {
	modlist := make([]*Module, 0, 10)
	for dec := json.NewDecoder(r); dec.More(); {
		mod := new(Module)
		if err := dec.Decode(mod); err != nil {
			return nil, fmt.Errorf("JSON decode: %w", err)
		}

		modlist = append(modlist, mod)
	}

	return modlist, nil
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

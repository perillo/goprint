// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Support for handling Go packages.

package main

import (
	"encoding/json"
	"os/exec"
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
func Find(path string) (*Package, error) {
	pkg := new(Package)
	cmd := exec.Command("go", "list", "-json", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		println("StdoutPipe")
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		println("Start")
		return nil, err
	}
	err = json.NewDecoder(stdout).Decode(pkg)
	if err != nil {
		println("Decode")
		return nil, err
	}
	err = cmd.Wait()
	if err != nil {
		println("Wait")
		return nil, err
	}

	return pkg, nil
}

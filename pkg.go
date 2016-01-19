// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Support for handling Go packages.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

// SourceFile represents a selection of Go source files.
type SourceFile int

// Available Go source file selections.
const (
	GoFiles SourceFile = iota
	CgoFiles
	IgnoredGoFiles
	TestGoFiles
	XTestGoFiles
)

var sourceFile_string = []string{
	GoFiles:        "go",
	CgoFiles:       "cgo",
	IgnoredGoFiles: "ignored",
	TestGoFiles:    "test",
	XTestGoFiles:   "xtest",
}

// String implements the Stringer interface.
func (src SourceFile) String() string {
	return sourceFile_string[src]
}

var sourceFile_set = map[string]SourceFile{
	"go":      GoFiles,
	"cgo":     CgoFiles,
	"ignored": IgnoredGoFiles,
	"test":    TestGoFiles,
	"xtest":   XTestGoFiles,
}

// Set implements the Value interface.
func (src *SourceFile) Set(s string) error {
	v, ok := sourceFile_set[s]
	if !ok {
		return fmt.Errorf("%q is not a valid source file selection", s)
	}

	*src = v

	return nil
}

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
	pkg := new(Package)
	args := []string{"list", "-json"}
	cmd := exec.Command("go", append(args, path...)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if err := json.NewDecoder(stdout).Decode(pkg); err == io.EOF {
		// TODO(mperillo): Should we report a custom error message if a pattern
		// was specified?
		return nil, fmt.Errorf("cannot find package %q", path)
	} else if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return pkg, nil
}

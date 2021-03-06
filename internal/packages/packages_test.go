// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packages

import (
	"os"
	"testing"
)

// TestLoad tests that Load works correctly when loading a standard package.
func TestLoad(t *testing.T) {
	pkg, err := Load("flag")
	if err != nil {
		t.Error("expected err == nil")
	}
	if pkg == nil {
		t.Error("expected pkg != nil")
	}
}

// TestLoadCurrent tests that Load works correctly when loading the current
// package (packages).
func TestLoadCurrent(t *testing.T) {
	pkg, err := Load("")
	if err != nil {
		t.Error("expected err == nil")
	}
	if pkg == nil {
		t.Fatal("expected pkg != nil")
	}

	const want = "packages"
	if pkg.Name != want {
		t.Errorf("want pkg.Name = %s, got %s", want, pkg.Name)
	}
}

// TestLoadFail tests that Load returns nil and an error when loading a
// nonexistent package.
func TestLoadFail(t *testing.T) {
	pkg, err := Load("xxx")
	if err == nil {
		t.Error("expected err != nil")
	}
	if pkg != nil {
		t.Error("expected pkg == nil")
	}
}

// Allow changing the go command to use in the test.  This can be useful when
// testing older versions of go, e.g.
//  go get golang.org/dl/go1.10.8
//  go1.10.8 download
//  GOCMD=go1.10.8 go1.10.8 test ./...
func init() {
	if value := os.Getenv("GOCMD"); value != "" {
		// TODO(mperillo): Check that $GOCMD is an executable in $PATH.
		gocmd = value
	}
}

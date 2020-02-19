// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// goprint is a command used to print the source code of a Go package.
//
// The generated document is in HTML format, written on stdout and with CSS
// specialized for printing.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/perillo/goprint/internal/css"
	"github.com/perillo/goprint/internal/packages"
)

// Context is the context used by the HTML template.
type Context struct {
	// Package to print.
	Package *packages.Package
	// Package's containing module.
	Module *packages.Module
	// Source files to print.
	Files []File
	// Style configuration.
	PageSize   css.PageSize
	PageMargin css.PageMargin
	Font       css.Font
}

var tmpl *template.Template

// Command line flags.
var (
	test       = flag.Bool("test", false, "print _test.go source files")
	pageSize   = css.A4
	pageMargin = css.PageMargin{
		Top:    css.Dimension{2.5, css.Centimeter},
		Right:  css.Dimension{1, css.Centimeter},
		Bottom: css.Dimension{2.5, css.Centimeter},
		Left:   css.Dimension{1, css.Centimeter},
	}
	font = css.Font{
		Family:     "Courier",
		Size:       css.Dimension{10, css.Point},
		LineHeight: css.Dimension{12, css.Point},
	}
)

func init() {
	tmpl = template.Must(template.New("index.html").Parse(index))
	template.Must(tmpl.New("style.css").Parse(style))

	flag.Var(&pageSize, "page-size", "page size")
	flag.Var(&pageMargin, "page-margin", "page margin")
	flag.Var(&font, "font", "font")
}

func main() {
	// Setup log.
	log.SetFlags(0)

	// Parse command line.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\tgoprint [flags] importpath\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	var arg string
	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		flag.Usage()
	}
	if flag.NArg() == 1 {
		arg = flag.Arg(0)
	}

	// Get package info, and format source files.
	pkg, err := packages.Load(arg)
	if err != nil {
		log.Fatal(err)
	}

	files, err := build(pkg, *test)
	if err != nil {
		log.Fatal(err)
	}

	// Render template.
	ctx := Context{
		Package:    pkg,
		Module:     pkg.Module,
		Files:      files,
		PageSize:   pageSize,
		PageMargin: pageMargin,
		Font:       font,
	}
	err = tmpl.Execute(os.Stdout, ctx)
	if err != nil {
		log.Fatalf("executing: %v", err)
	}
}

// build returns the package pkg .go source files formatted in HTML.
//
// If test is true, build will use the package pkg _test.go files.
func build(pkg *packages.Package, test bool) ([]File, error) {
	srcfiles := pkg.SourceFiles()
	if test {
		srcfiles = pkg.TestFiles()
	}

	files := make([]File, len(srcfiles))
	for i, path := range srcfiles {
		name := filepath.Base(path)
		input, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read file %s: %v", path, err)
		}
		files[i] = File{
			Name: name,
			Code: render(name, input),
		}
	}

	return files, nil
}

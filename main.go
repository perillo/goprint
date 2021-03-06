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

// Command line flags.
var (
	test       = flag.Bool("test", false, "print _test.go source files")
	module     = flag.Bool("m", false, "print all the packages in the module")
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

	// Print the package or module.
	printer := printPackage
	if *module {
		printer = printModule
	}
	if err := printer(arg, *test); err != nil {
		log.Fatal(err)
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

// buildModules returns all the module mod packages .go source files formatted in
// HTML.
//
// If test is true, build will use each package _test.go files.
func buildModule(mod *packages.Module, test bool) ([]Package, error) {
	pkglist := make([]Package, len(mod.Packages))
	for i, pkg := range mod.Packages {
		files, err := build(pkg, test)
		if err != nil {
			return nil, err
		}

		p := Package{
			ImportPath: pkg.ImportPath,
			Name:       pkg.Name,
			Files:      files,
		}
		pkglist[i] = p
	}

	return pkglist, nil
}

// printPackage writes on stdout an HTML document with the all the .go source
// files of the package named by path.
//
// It test is true, printPackage will use the package _test.go files.
func printPackage(path string, test bool) error {
	// Get package info.
	pkg, err := packages.Load(path)
	if err != nil {
		return err
	}

	// Format source files.
	files, err := build(pkg, test)
	if err != nil {
		return err
	}

	// Load template.
	tmpl := template.Must(template.New("index.html").Parse(index))
	template.Must(tmpl.New("style.css").Parse(style))

	// Render template.
	ctx := struct {
		Package    *packages.Package
		Module     *packages.Module
		Files      []File
		PageSize   css.PageSize
		PageMargin css.PageMargin
		Font       css.Font
	}{
		pkg,
		pkg.Module,
		files,
		pageSize,
		pageMargin,
		font,
	}
	if err := tmpl.Execute(os.Stdout, ctx); err != nil {
		return fmt.Errorf("execute: %v", err)
	}

	return nil
}

// printModule writes on stdout an HTML document with all the .go source files
// of all the packages belonging to the module named by path.
//
// It test is true, printModule will use the packages _test.go files.
func printModule(path string, test bool) error {
	// Get module info.
	mod, err := packages.LoadModule(path)
	if err != nil {
		return err
	}

	// Format packages.
	pkglist, err := buildModule(mod, test)
	if err != nil {
		return err
	}

	// Load template.
	tmpl := template.Must(template.New("index.html").Parse(indexmod))
	template.Must(tmpl.New("style.css").Parse(stylemod))

	// Render template.
	ctx := struct {
		Module     *packages.Module
		Packages   []Package
		PageSize   css.PageSize
		PageMargin css.PageMargin
		Font       css.Font
	}{
		mod,
		pkglist,
		pageSize,
		pageMargin,
		font,
	}
	if err := tmpl.Execute(os.Stdout, ctx); err != nil {
		return fmt.Errorf("execute: %v", err)
	}

	return nil
}

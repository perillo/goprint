// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// goprint is a command used to print the source code of a Go package.
//
// The generated document is in HTML format, written on stdout and with CSS
// specialized for printing.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/perillo/goprint/internal/css"
	"github.com/perillo/goprint/internal/goefmt"
	"github.com/perillo/goprint/internal/packages"
)

// File represents an HTML formatted Go source file.
type File struct {
	Name string
	Code template.HTML
}

// Context is the context used by the HTML template.
type Context struct {
	// Package import path.
	ImportPath string
	// Package name.
	Name string
	// Source files to print.
	Files []File
	// Style configuration
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

	// Get package info, and format source files.
	// Only .go source files, excluding files using Cgo, are printed, to avoid
	// consuming too much paper.
	pkg, err := packages.Load(flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}

	srcfiles := pkg.SourceFiles()
	if *test {
		srcfiles = pkg.TestFiles()
	}
	files := make([]File, len(srcfiles))
	for i, path := range srcfiles {
		name := filepath.Base(path)
		input, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("reading file: %v", err)
		}
		files[i] = File{
			Name: name,
			Code: printFile(name, input),
		}
	}

	// Render template.
	ctx := Context{
		ImportPath: pkg.ImportPath,
		Name:       pkg.Name,
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

// printFile returns an HTML fragment containing the formatted Go code for the
// specified source file. A line number is printed at the begin of each line.
func printFile(name string, input []byte) template.HTML {
	buf := new(bytes.Buffer)

	n := 1
	for line := range goefmt.Format(goefmt.Scan(name, input)) {
		if line == nil {
			// Empty line
			fmt.Fprintf(buf, "<span class=\"line empty\">%3d</span>\n", n)
		} else {
			fmt.Fprintf(buf, "<span class=\"line\">%3d</span> %s\n", n,
				lineToHTML(line))
		}
		n++
	}

	return template.HTML(buf.String())
}

// spanToHTML returns an HTML representation for the code span.
func spanToHTML(s *goefmt.Span) string {
	if s.Code == "" {
		// Only horizontal white space.
		return s.Whitespace
	}
	class := strings.Join(goefmt.TokenClass(s), " ")
	code := html.EscapeString(s.Code)

	return fmt.Sprintf(`<span class="%s">%s</span>%s`, class, code, s.Whitespace)
}

// lineToHTML returns an HTML representation for the code line. The eol is not
// included.
func lineToHTML(l goefmt.Line) string {
	if l == nil {
		// Empty line.
		return ""
	}

	spans := make([]string, len(l))
	for i, span := range l {
		spans[i] = spanToHTML(span)
	}

	return strings.Join(spans, "")
}

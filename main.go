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
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// File represents an HTML formatted Go source file.
type File struct {
	Name string
	Code template.HTML
}

type Context struct {
	// Package import path.
	ImportPath string
	// Package name.
	Name string
	// Source files to print.
	Files []File
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("index.html").Parse(index))
	template.Must(tmpl.New("style.css").Parse(style))
}

func main() {
	// Setup log.
	log.SetFlags(0)

	// Parse command line.
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	// Get package info, and format source files.
	// Only GoFiles are printed, to avoid consuming too much paper.
	pkg, err := Find(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	files := make([]File, len(pkg.GoFiles))
	for i, name := range pkg.GoFiles {
		path := filepath.Join(pkg.Dir, name)
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
	for line := range Format(Scan(name, input)) {
		if line == nil {
			// Empty line
			fmt.Fprintf(buf, "<span class=\"line empty\">%3d</span>\n", n)
		} else {
			fmt.Fprintf(
				buf, "<span class=\"line\">%3d</span> %s\n", n, line.HTML())
		}
		n++
	}

	return template.HTML(buf.String())
}

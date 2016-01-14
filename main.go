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
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Setup log.
	log.SetFlags(0)

	// Parse command line.
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	name := flag.Arg(0)
	input, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	for token := range Scan(name, input) {
		fmt.Printf("%v", token)
	}
	// Add final end of line.
	fmt.Print("\n")
}

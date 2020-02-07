# goprint [![GoDoc](https://godoc.org/github.com/perillo/goprint?status.svg)](http://godoc.org/github.com/perillo/goprint)

`goprint` is a command that prints the source code of a *Go* packages.  The
generated document is in *HTML* format, suitable for printing or converting to
*PDF*.

The tool uses experimental *CSS* features for setting the page headers
and footer that, at the present time, are not fully supported by browsers.
However they are supported by *Prince* (http://www.princexml.com).

On the top left of the page is reported the package `import path`.

On the top right of the page is reported the file name.

On the bottom right of the page is reported the page number.

## Usage

    Usage: goprint [flags] importpath
    Flags:
      -files value
          files to print
      -font value
          font (default "Courier" 10pt/12pt)
      -page-margin value
          page margin (default 2.5cm 1cm)
      -page-size value
          page size (default A4 portrait)

`importpath` is interpreted as in `go list`, however `goprint` only process the
first package.

By default `goprint` will print only `.go` source files, excluding `CGo` files,
*ignored* files, *test* files and *external test* files.

### `-files`

Using the `-file` flag, it is possible to specify an alternate source file
selection (`"go"`, `"cgo`", `"ignored"`, `"test"` and `"xtest"`).  As an
example:

    goprint -files=test

In alternative, it is always possible to specify a custom list of files as
command line arguments.  As an example:

    goprint main.go dimension.go

### `-page-size`

Supported page sizes are `A4` or `letter`.  The page orientation is `portrait`
and can not be changed.

### `-page-margin`

The right, bottom and left margins can be omitted.

### `-font`

The font family, font size and line height must all be specified.  The font
family must be quoted, even if it contains no spaces.

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
      -font value
          font (default "Courier" 10pt/12pt)
      -page-margin value
          page margin (default 2.5cm 1cm)
      -page-size value
          page size (default A4 portrait)
      -test
          print _test.go source files

`importpath` is interpreted as in `go list`, however `goprint` only process the
first package.

By default `goprint` will print only `.go` source files, excluding `CGo` files,
*ignored* files, *test* files and *external test* files.

### `-page-size`

Supported page sizes are `A4` or `letter`.  The page orientation is `portrait`
and can not be changed.

### `-page-margin`

The right, bottom and left margins can be omitted.

### `-font`

The font family, font size and line height must all be specified.  The font
family must be quoted, even if it contains no white space.

### `-test`

By default, `goprint` prints all the `.go` source files, including files
ignored to build constraint.

When the `-test` flag is set, `goprint` will print all the `_test.go` source
files, instead.


## Examples

```
goprint main.go > build/pkg.html
```

```
goprint -font='"Inconsolata" 10pt/12pt' ./internal/css > build/pkg.html
prince -o build/pkg.pdf build/pkg.html
```

# Requirements

`goprint` requires at least *Go 1.8*.

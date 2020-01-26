# goprint [![GoDoc](https://godoc.org/github.com/perillo/goprint?status.svg)](http://godoc.org/github.com/perillo/goprint)

`goprint` is a command that prints the source code of a Go packages.  The
generated document is in `HTML` format, suitable for printing or converting to
`PDF`.

The tool uses experimental `CSS` features for setting the page headers
and footer that, at the present time, are not supported by browsers.  However
they are supported by `Prince` (http://www.princexml.com).

## Usage
`goprint` requires an optional package `import path`, that is interpreted as
with the `Go` *standard tools*.

* If no `import path` is specified, `goprint` will print the package in the
  current directory.

* If the `import path` is a pattern (contains the `...`), only the first
  matched package will be printed.  The same with the special `std`, `cmd` and
  `all` special `import path`.

* If one or more file are specified (all inside the same directory), `goprint`
  will print only the specified source files.

* Otherwise the specified package will be printed. As an example:

    goprint github.com/perillo/goprint

By default `goprint` will print only `.go` source files, exluding `CGo` files,
*ignored* files, *test* files and *external test* files.  Using the `-file`
flag it is possible to specify an alternate source file selection
(`"go"`, `"cgo`", `"ignored"`, `"test"` and `"xtest"`), or it is possible to
specify a custom list of files as command line arguments.  As an example:

    goprint -files=test
    goprint main.go dimension.go

By default the page size is `A4`.  Using the `--page-size` flag it is possible
to specify either `A4` or `letter`.  The page orientation is `portrait` and can
not be changed.  As an example:

    goprint -page=letter

By default the page margins are `2.5cm` for top and bottom, and `1cm` for left
and right.  This seems a reasonable minimal page margin.  It is possible to
change the page margin using the `-page-margin` flag.  As an example:

    goprint -page-margin='3cm 2cm 3cm 2cm'

Note how all the top, right, bottom and left margins must be specified.

By default the font family is `Courier`, the font size is `10pt` and the line
height is `12pt`.  Using the `-font` option it is possible to change the font.
As an example:

    goprint -font='"Inconsolata" 10pt/12pt'

Note how the font name must be quoted, even if it contains no spaces.

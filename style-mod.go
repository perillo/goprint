// vim: set filetype=css :
// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Definition of the CSS template for modules.

package main

var stylemod = `
* {
	margin: 0;
	padding: 0;
	tab-size: 4;
}

html {
	font-family: "{{ .Font.Family }}", Courier, monospace;
	font-size: {{ .Font.Size }};
	line-height: {{ .Font.LineHeight }};
}

code {
	display: block;
}

.line {
	color: #999;
}

.operator, .ident {
	font-style: normal;
	font-weight: normal;
}

.keyword {
	font-weight: bold;
}

.builtin {
	font-weight: bold;
	font-style: italic;
}

.literal {
	font-style: italic;
}

.comment {
	font-style: oblique;
}

.invalid {
	background-color: red;
}

@media print {
	@page {
		size: {{ .PageSize }};
		margin: {{ .PageMargin }};
		font-size: {{ .Font.Size }};
		counter-increment: page 1;

		@top-left {
			vertical-align: bottom;
			margin-bottom: 1.5em;
			content: string(package);
		}

		@top-right {
			vertical-align: bottom;
			margin-bottom: 1.5em;
			content: string(file);
		}

		@bottom-left {
			vertical-align: top;
			margin-top: 1.5em;
			content: "{{ .Module }}" "\2003" "{{ .Module.Date }}";
		}

		@bottom-right {
			vertical-align: top;
			margin-top: 1.5em;
			content: "page " counter(page);
		}
	}

	.package {
		page-break-after: always;
		string-set: package attr(data-package);
	}

	.package > h1 {
		display: none;
	}

	.package:last-of-type {
		page-break-after: auto;
	}

	.file {
		page-break-after: always;
		string-set: file attr(data-file);
	}

	.file:last-of-type {
		page-break-after: auto;
	}

	h1, h2, h3 {
		display: none;
	}
}

@media screen {
	h1 {
		font-size: 2.5em;
	}

	h2 {
		margin: 16px;
		font-size: 2em;
		text-align: center;
	}

	h3 {
		margin: 16px;
		font-size: 1.5em;
		text-align: center;
	}
}
`

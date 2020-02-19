// vim: set filetype=html :
// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Definition of the HTML template for modules.

package main

var indexmod = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<style type="text/css">
			{{ template "style.css" . }}
		</style>

		<title>{{ .Module }}</title>
	</head>
	<body>
	  {{ range .Packages }}
		<section class="package" data-package="{{ .ImportPath }}">
			<h1>{{ .Name }}</h1>
			{{ range .Files }}
			<section class="file" data-file="{{ .Name }}">
				<h1>{{ .Name }}</h1>
				<pre><code>{{ .Code }}</code></pre>
			</section>
			{{ end }}
		</section>
		{{ end }}
	</body>
</html>
`

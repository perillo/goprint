// vim: set filetype=html :
// Copyright 2016 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Definition of the HTML template.

package main

var index = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<style type="text/css">
			{{ template "style.css" . }}
		</style>

		<title>{{ .Package.Name }}</title>
	</head>
	<body>
		<section class="package">
			<h1>{{.Package.Name}}</h1>
			{{ range .Files }}
			<section class="file" data-file="{{.Name}}">
				<h1>{{.Name}}</h1>
				<pre><code>{{.Code}}</code></pre>
			</section>
			{{ end }}
		</section>
	</body>
</html>
`

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

var htmlHeader = `<!DOCTYPE html>
<html lang="en">
<head>
<title>Azure GO SDK MGMT REPORT</title>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
 
<style type="text/css">
html {
    font-family: sans-serif;
    -ms-text-size-adjust: 100%;
    -webkit-text-size-adjust: 100%;
}
 
body {
    margin: 10px;
}
table {
    border-collapse: collapse;
    border-spacing: 0;
}
 
td,th {
    padding: 0;
}
 
.pure-table {
    border-collapse: collapse;
    border-spacing: 0;
    empty-cells: show;
    border: 1px solid #cbcbcb;
}
 
.pure-table caption {
    color: #000;
    font: italic 85%/1 arial,sans-serif;
    padding: 1em 0;
    text-align: center;
}
 
.pure-table td,.pure-table th {
    border-left: 1px solid #cbcbcb;
    border-width: 0 0 0 1px;
    font-size: inherit;
    margin: 0;
    overflow: visible;
    padding: .5em 1em;
}
 
.pure-table thead {
    background-color: #e0e0e0;
    color: #000;
    text-align: left;
    vertical-align: bottom;
}
 
.pure-table td {
    background-color: transparent;
}
 
.pure-table-odd td {
    background-color: #f2f2f2;
}
</style>
</head>
<body>
    <table class="pure-table">
        <thead>
            <tr>
				<th align="left">module</th>
				<th align="center">latest version</th>
				<th align="center">tag</th>
				<th align="center">live test coverage</th>
				<th align="center">mock test result</th>
				<th align="center">mock test coverage</th>
            </tr>
        </thead>
    
        <tbody>`

var htmlTR = `
			<tr%s>
				<td align="left">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
			</tr>`

var htmlTail = `
        </tbody>
	</table>
</body>
</html>`

var tdBackgroundStyle = ` class="pure-table-odd"`

var mgmtReportMDHeader = `|module | latest version | tag | live test coverage | mock test result | mock test coverage |
|---|---|---|---|---|---|
`

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report_test

import (
	"os"
	"testing"
	"testing/fstest"

	"apigen/internal/loader"
	"apigen/internal/report"

	"github.com/stretchr/testify/require"
)

func TestExtractModule(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	files := fstest.MapFS{
		"go.mod": {Data: []byte(`module example.com/sample

go 1.26.0

require (
	example.com/alpha v1.2.3
	example.com/beta v1.4.0
	example.com/testonly v1.0.0
	example.com/transitive v1.0.0 // indirect
)

replace (
	example.com/alpha => ./deps/alpha
	example.com/beta => ./deps/beta
	example.com/testonly => ./deps/testonly
	example.com/transitive => ./deps/transitive
)
`)},
		"sample.go": {Data: []byte(`package sample

import (
	"example.com/beta/one"
	"example.com/beta/two"
	"example.com/sample/internal/core"
	_ "unsafe"
)

var Value = one.Value + two.Value + core.Value
`)},
		"sample_test.go": {Data: []byte(`package sample_test

import "example.com/testonly"

var Value = testonly.Value
`)},
		"internal/core/core.go": {Data: []byte(`package core

import "example.com/alpha"

var Value = alpha.Value
`)},
		"deps/alpha/go.mod": {Data: []byte(`module example.com/alpha

go 1.26.0

require example.com/transitive v1.0.0
`)},
		"deps/alpha/alpha.go": {Data: []byte(`package alpha

import "example.com/transitive"

var Value = transitive.Value
`)},
		"deps/beta/go.mod": {
			Data: []byte("module example.com/beta\n\ngo 1.26.0\n"),
		},
		"deps/beta/one/one.go": {
			Data: []byte("package one\n\nconst Value = 1\n"),
		},
		"deps/beta/two/two.go": {
			Data: []byte("package two\n\nconst Value = 2\n"),
		},
		"deps/testonly/go.mod": {
			Data: []byte("module example.com/testonly\n\ngo 1.26.0\n"),
		},
		"deps/testonly/testonly.go": {
			Data: []byte("package testonly\n\nconst Value = 3\n"),
		},
		"deps/transitive/go.mod": {
			Data: []byte("module example.com/transitive\n\ngo 1.26.0\n"),
		},
		"deps/transitive/transitive.go": {
			Data: []byte("package transitive\n\nconst Value = 4\n"),
		},
	}
	require.NoError(t, os.CopyFS(dir, files))

	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)

	module := report.ExtractModule(loaded)
	require.Equal(t, report.Module{
		Path:      "example.com/sample",
		GoVersion: "1.26.0",
		Dependencies: []report.Dependency{
			{Path: "example.com/alpha", Version: "v1.2.3"},
			{Path: "example.com/beta", Version: "v1.4.0"},
		},
	}, module)
}

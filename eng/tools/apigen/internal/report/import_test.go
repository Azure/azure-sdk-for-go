// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report_test

import (
	"bytes"
	"go/format"
	"go/printer"
	"os"
	"testing"
	"testing/fstest"

	"apigen/internal/loader"
	"apigen/internal/report"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectImports(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	files := fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(`package sample

import (
	"math"
	renamed "example.com/sample/api"
	"example.com/sample/bodyonly"
	. "example.com/sample/dotted"
	"example.com/sample/hiddenonly"
	"example.com/sample/initializers"
	"example.com/sample/internal/model"
	_ "example.com/sample/sideeffect"
)

type Surface struct {
	Public renamed.Thing
	hidden hiddenonly.Thing
}

func Exported(value renamed.Thing) renamed.Thing {
	bodyonly.Touch()
	return value
}

type DotAlias = Dotted
type PublicRecord = model.Record
type Code = model.Code

var Item = renamed.Thing{Field: 1}
var Constructed = initializers.New()
var DotValue = NewDotted()

const Pi = math.Pi
const ExternalKind = renamed.Default
const FromInternal = model.Default
const UntypedZero = 0
`)},
		"api/api.go": {Data: []byte(`package different

type Thing struct{ Field int }
type Kind int
const Default Kind = 1
`)},
		"bodyonly/body.go": {
			Data: []byte("package bodyonly\n\nfunc Touch() {}\n"),
		},
		"dotted/dotted.go": {
			Data: []byte("package dotted\n\ntype Dotted struct{}\nfunc NewDotted() Dotted { return Dotted{} }\n"),
		},
		"hiddenonly/hidden.go": {
			Data: []byte("package hiddenonly\n\ntype Thing struct{}\n"),
		},
		"initializers/initializers.go": {
			Data: []byte("package initializers\n\nfunc New() int { return 1 }\n"),
		},
		"internal/model/model.go": {Data: []byte(`package model

import (
	tm "time"
	"sync"
)

type Record struct {
	When tm.Time
	hidden sync.Mutex
}
type Code string
const Default Code = "default"
`)},
		"sideeffect/sideeffect.go": {
			Data: []byte("package sideeffect\n\nfunc init() {}\n"),
		},
	}
	require.NoError(t, os.CopyFS(dir, files))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	root := loaded.Packages[0]
	var before bytes.Buffer
	require.NoError(t, format.Node(&before, root.Fset, root.Syntax[0]))

	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	require.NotEmpty(t, packages)
	assert.Equal(t, "example.com/sample", packages[0].Path)
	assert.Equal(t, []report.Import{
		{Path: "time", Standard: true},
		{Path: "example.com/sample/api", Name: "renamed"},
		{Path: "example.com/sample/dotted", Name: "."},
		{Path: "example.com/sample/initializers"},
	}, packages[0].Imports)

	constantTypes := make(map[string]string)
	constants := append([]report.Constant{}, packages[0].Constants...)
	for _, typ := range packages[0].Types {
		constants = append(constants, typ.Constants...)
	}
	for _, value := range constants {
		require.NotNil(t, value.Comments)
		if value.TypeExpr == nil {
			constantTypes[value.Name] = ""
			continue
		}
		var output bytes.Buffer
		require.NoError(t, format.Node(&output, packages[0].FileSet, &printer.CommentedNode{
			Node: value.TypeExpr, Comments: value.Comments,
		}))
		constantTypes[value.Name] = output.String()
	}
	assert.Equal(t, map[string]string{
		"ExternalKind": "renamed.Kind",
		"FromInternal": "Code",
		"Pi":           "",
		"UntypedZero":  "",
	}, constantTypes)

	var after bytes.Buffer
	require.NoError(t, format.Node(&after, root.Fset, root.Syntax[0]))
	assert.Equal(t, before.String(), after.String(), "import collection must not mutate loaded syntax")
}

func TestCollectImportQualifierCollision(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"first.go": {Data: []byte(`package sample

import "example.com/sample/first"

func First(value models.Value) {}
`)},
		"second.go": {Data: []byte(`package sample

import "example.com/sample/second"

func Second(value models.Value) {}
`)},
		"first/value.go": {
			Data: []byte("package models\n\ntype Value string\n"),
		},
		"second/value.go": {
			Data: []byte("package models\n\ntype Value int\n"),
		},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	packages, err := report.ExtractPackages(loaded)
	require.ErrorContains(t, err, `import qualifier "models" refers to both "example.com/sample/first" and "example.com/sample/second"`)
	assert.Nil(t, packages)
}

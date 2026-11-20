// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report_test

import (
	"bytes"
	"go/constant"
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

func TestExtractValues(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	source := `package sample

type Mode uint8
type Widget struct{ Name string }

const (
	privateStart Mode = iota + 2
	Zulu
	_
	Alpha
)

const Echo = privateStart
const Mixed, hiddenConst = 7, 9
const Ratio = 1.0 / 3
const Parts = complex(1.0 / 3, -2.0 / 7)

var (
	ZuluValue = Widget{Name: "z"}
	// AlphaValue documents the initial value.
	AlphaValue Widget = Widget{Name: "a"}
	Zed, Apple = 10, 20
	BothNumber, BothText = pair()
	First, hiddenResult = pair()
	privateResult, Third = pair()
	_, Second = pair()
	MixedPublic, privateInit = 21, "secret"
	TypedA, TypedB Widget
	Opaque struct{ hidden int }
	privateOnly = 42
)

func pair() (int, string) { return 1, "one" }
`
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod":    {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(source)},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	root := loaded.Packages[0]
	var before bytes.Buffer
	require.NoError(t, format.Node(&before, root.Fset, root.Syntax[0]))

	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	require.Len(t, packages, 1)
	pkg := packages[0]
	constants := append([]report.Constant{}, pkg.Constants...)
	for _, typ := range pkg.Types {
		constants = append(constants, typ.Constants...)
	}
	variables := pkg.Variables

	t.Run("constants", func(t *testing.T) {
		var names []string
		values := make(map[string]report.Constant)
		for _, value := range constants {
			names = append(names, value.Name)
			values[value.Name] = value
		}
		require.ElementsMatch(t, []string{"Alpha", "Echo", "Mixed", "Parts", "Ratio", "Zulu"}, names)
		assert.Equal(t, "5", values["Alpha"].Value.ExactString())
		assert.Equal(t, "3", values["Zulu"].Value.ExactString())
		assert.Equal(t, "2", values["Echo"].Value.ExactString())
		assert.Equal(t, "example.com/sample.Mode", values["Alpha"].Type.String())
		assert.Equal(t, "example.com/sample.Mode", values["Echo"].Type.String())
		assert.Equal(t, "7", values["Mixed"].Value.ExactString())
		assert.Equal(t, "untyped int", values["Mixed"].Type.String())
		assert.Equal(t, "1/3", values["Ratio"].Value.ExactString())
		assert.Equal(t, "untyped float", values["Ratio"].Type.String())
		assert.Equal(t, constant.Complex, values["Parts"].Value.Kind())
		assert.Equal(t, "1/3", constant.Real(values["Parts"].Value).ExactString())
		assert.Equal(t, "-2/7", constant.Imag(values["Parts"].Value).ExactString())
	})

	t.Run("variables", func(t *testing.T) {
		var names []string
		rendered := make(map[string]string)
		bindings := make(map[string]string)
		for _, variable := range variables {
			require.NotEmpty(t, variable.Objects)
			name := variable.Objects[0].Name()
			names = append(names, name)
			for _, object := range variable.Objects {
				bindings[object.Name()] = object.Type().String()
			}
			require.NotNil(t, variable.Comments)
			var output bytes.Buffer
			require.NoError(t, format.Node(&output, pkg.FileSet, &printer.CommentedNode{
				Node:     variable.Decl,
				Comments: variable.Comments,
			}))
			_, err := format.Source(append([]byte("package sample\n"), output.Bytes()...))
			require.NoError(t, err, "invalid variable declaration for %s", name)
			rendered[name] = output.String()
		}
		assert.ElementsMatch(t, []string{
			"AlphaValue", "Apple", "BothNumber", "First", "MixedPublic", "Opaque",
			"Second", "Third", "TypedA", "TypedB", "Zed", "ZuluValue",
		}, names)
		assert.Equal(t, `var AlphaValue Widget = Widget{Name: "a"}`, rendered["AlphaValue"])
		assert.Equal(t, `var ZuluValue = Widget{Name: "z"}`, rendered["ZuluValue"])
		assert.Equal(t, "var Apple = 20", rendered["Apple"])
		assert.Equal(t, "var Zed = 10", rendered["Zed"])
		assert.Equal(t, "var BothNumber, BothText = pair()", rendered["BothNumber"])
		assert.Equal(t, "var First, _ = pair()", rendered["First"])
		assert.Equal(t, "var _, Second = pair()", rendered["Second"])
		assert.Equal(t, "var _, Third = pair()", rendered["Third"])
		assert.Equal(t, "var MixedPublic = 21", rendered["MixedPublic"])
		assert.Equal(t, "var TypedA Widget", rendered["TypedA"])
		assert.Equal(t, "var TypedB Widget", rendered["TypedB"])
		assert.Contains(t, rendered["Opaque"], "// contains filtered or unexported fields")
		assert.NotContains(t, rendered["Opaque"], "hidden")
		assert.Equal(t, "example.com/sample.Widget", bindings["ZuluValue"])
		assert.Equal(t, "int", bindings["BothNumber"])
		assert.Equal(t, "string", bindings["BothText"])
		assert.Equal(t, "int", bindings["First"])
		assert.Equal(t, "string", bindings["Second"])
		assert.Equal(t, "string", bindings["Third"])
	})

	var after bytes.Buffer
	require.NoError(t, format.Node(&after, root.Fset, root.Syntax[0]))
	assert.Equal(t, before.String(), after.String(), "value extraction must not mutate loaded syntax")
}

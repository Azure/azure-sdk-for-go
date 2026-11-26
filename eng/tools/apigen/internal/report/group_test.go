// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report_test

import (
	"os"
	"testing"
	"testing/fstest"

	"apigen/internal/loader"
	"apigen/internal/report"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupDeclarations(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(`package sample

import (
	"example.com/sample/internal/model"
	"example.com/sample/other"
)

type Widget struct{}
type Alias = Widget
type Box[T any] struct{ Value T }
type Mode string
type ExternalAlias = other.Widget
type ZContract interface{ Read() }
type AContract = ZContract
type ResponseContent = model.ResponseContent

const ZuluMode Mode = "z"
const AlphaMode Mode = "a"
const Bare = 1

var ZConfig = Widget{}
var AConfig *Widget
var First, Second = pairWidgets()
var MixedWidget, MixedNumber = mixed()
var AliasValue Alias

func NewWidget() (*Widget, error) { return &Widget{}, nil }
func NewBox[T any](value T) *Box[T] { return &Box[T]{Value: value} }
func New[T any]() *Box[T] { return &Box[T]{} }
func NewAlias() Alias { return Alias{} }
func Newspaper() Widget { return Widget{} }
func Widgets() []*Widget { return nil }
func External() other.Widget { return other.Widget{} }
func NewExternal() other.Widget { return other.Widget{} }
func NewExternalAlias() ExternalAlias { return ExternalAlias{} }
func GetResponseContent[T any](value T) ResponseContent { return ResponseContent{} }
func NewResponseContent() ResponseContent { return ResponseContent{} }
func Generic[T any](value T) T { return value }
func Pair() (Widget, Widget) { return Widget{}, Widget{} }
func Both() (Widget, Mode) { return Widget{}, "" }
func PossibleModeValues() [2]Mode { return [2]Mode{} }
func (Widget) Zulu() {}
func (*Widget) Alpha() {}
func pairWidgets() (Widget, Widget) { return Widget{}, Widget{} }
func mixed() (Widget, int) { return Widget{}, 1 }
`)},
		"other/other.go":          {Data: []byte("package other\n\ntype Widget struct{}\n")},
		"internal/model/model.go": {Data: []byte("package model\n\ntype ResponseContent struct{ HTTPStatus int }\n")},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	require.NotEmpty(t, packages)
	pkg := packages[0]
	require.Equal(t, "example.com/sample", pkg.Path)

	groups := make(map[string]report.Type)
	for _, typ := range pkg.Types {
		groups[typ.Name] = typ
	}
	functionNames := func(functions []report.Function) []string {
		var names []string
		for _, function := range functions {
			names = append(names, function.Name)
		}
		return names
	}
	assert.Equal(t, []string{"Both", "External", "Generic", "GetResponseContent", "NewExternal", "Newspaper", "Pair", "PossibleModeValues", "Widgets"}, functionNames(pkg.Functions))
	assert.Equal(t, []string{"NewWidget"}, functionNames(groups["Widget"].Functions))
	assert.Equal(t, []string{"Alpha", "Zulu"}, functionNames(groups["Widget"].Methods))
	assert.Equal(t, []string{"New", "NewBox"}, functionNames(groups["Box"].Functions))
	assert.Equal(t, []string{"NewAlias"}, functionNames(groups["Alias"].Functions))
	assert.Equal(t, []string{"NewExternalAlias"}, functionNames(groups["ExternalAlias"].Functions))
	assert.Equal(t, []string{"NewResponseContent"}, functionNames(groups["ResponseContent"].Functions))

	var typeNames []string
	for _, typ := range pkg.Types {
		typeNames = append(typeNames, typ.Name)
	}
	assert.Equal(t, []string{"AContract", "ZContract", "Alias", "Box", "ExternalAlias", "Mode", "ResponseContent", "Widget"}, typeNames)
	assert.True(t, groups["AContract"].IsInterface)
	assert.True(t, groups["ZContract"].IsInterface)

	var constantNames []string
	for _, constant := range groups["Mode"].Constants {
		constantNames = append(constantNames, constant.Name)
	}
	assert.Equal(t, []string{"AlphaMode", "ZuluMode"}, constantNames)
	require.Len(t, pkg.Constants, 1)
	assert.Equal(t, "Bare", pkg.Constants[0].Name)
	require.Len(t, groups["Mode"].Constants, 2)
	require.NotNil(t, groups["Mode"].Constants[0].TypeExpr, "grouping must retain import-ready type expressions")

	var variableNames []string
	for _, variable := range pkg.Variables {
		for _, object := range variable.Objects {
			variableNames = append(variableNames, object.Name())
		}
	}
	assert.Equal(t, []string{"AConfig", "AliasValue", "First", "Second", "MixedWidget", "MixedNumber", "ZConfig"}, variableNames)
	require.Len(t, pkg.Variables, 5)
	require.Len(t, pkg.Variables[3].Objects, 2)
	assert.Equal(t, "MixedWidget", pkg.Variables[3].Objects[0].Name())
	assert.Equal(t, "MixedNumber", pkg.Variables[3].Objects[1].Name())

	functionCount := len(pkg.Functions)
	for _, typ := range pkg.Types {
		functionCount += len(typ.Functions) + len(typ.Methods)
	}
	assert.Equal(t, 17, functionCount, "grouping must neither lose nor duplicate functions")
}

func TestGroupEnumConstants(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(`package sample

import (
	"example.com/sample/internal/model"
	"example.com/sample/other"
)

type ZEvent = model.Event
type AEvent = model.Event
type Code int
type NoValues string
type Fraction float64
type Enabled bool
type ExternalAlias = other.Kind

const ZValue ZEvent = "z"
const Inferred = model.Default
const Positive Code = 1
const Negative Code = -1
const Aliased ExternalAlias = other.Default
const External = other.Default
const Floating Fraction = 1.5
const Boolean Enabled = true
const Untyped = "value"

func PossibleCodeValues() []Code { return nil }
func PossibleZEventValues() []ZEvent { return nil }
func PossibleExternalAliasValues() []ExternalAlias { return nil }
func PossibleValues() []Code { return nil }
func PossibleNoValuesValues() []NoValues { return nil }
func PossibleFractionValues() []Fraction { return nil }
func PossibleEnabledValues() []Enabled { return nil }
func PossibleAEventValues(filter bool) []AEvent { return nil }
`)},
		"internal/model/model.go": {Data: []byte(`package model

type Event string
const Default Event = "default"
`)},
		"other/other.go": {Data: []byte(`package other

type Kind int
const Default Kind = 1
`)},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	require.NotEmpty(t, packages)
	pkg := packages[0]
	assert.Equal(t, "example.com/sample", pkg.Path)

	enumGroups := make(map[string][]string)
	valueHelpers := make(map[string][]string)
	for _, typ := range pkg.Types {
		for _, constant := range typ.Constants {
			enumGroups[typ.Name] = append(enumGroups[typ.Name], constant.Name)
		}
		for _, function := range typ.Functions {
			valueHelpers[typ.Name] = append(valueHelpers[typ.Name], function.Name)
		}
	}
	assert.Equal(t, map[string][]string{
		"AEvent":        {"Inferred"},
		"Code":          {"Negative", "Positive"},
		"ExternalAlias": {"Aliased"},
		"ZEvent":        {"ZValue"},
	}, enumGroups)
	assert.Equal(t, map[string][]string{
		"Code":          {"PossibleCodeValues"},
		"ExternalAlias": {"PossibleExternalAliasValues"},
		"ZEvent":        {"PossibleZEventValues"},
	}, valueHelpers)
	var freeFunctions []string
	for _, function := range pkg.Functions {
		freeFunctions = append(freeFunctions, function.Name)
	}
	assert.Equal(t, []string{
		"PossibleAEventValues", "PossibleEnabledValues", "PossibleFractionValues",
		"PossibleNoValuesValues", "PossibleValues",
	}, freeFunctions)
	var unassociated []string
	for _, constant := range pkg.Constants {
		unassociated = append(unassociated, constant.Name)
	}
	assert.Equal(t, []string{"Boolean", "External", "Floating", "Untyped"}, unassociated)
}

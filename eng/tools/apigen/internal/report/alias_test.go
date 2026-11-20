// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report_test

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"os"
	"testing"
	"testing/fstest"

	"apigen/internal/loader"
	"apigen/internal/report"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandInternalAliases(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	files := fstest.MapFS{
		"go.mod": {Data: []byte(`module example.com/sdk/client

go 1.26.0

require example.com/sdk/internal v1.2.3

replace example.com/sdk/internal => ./deps/internal
`)},
		"client.go": {Data: []byte(`package client

import (
	opt "example.com/sdk/client/options"
	impl "example.com/sdk/internal/model"
)

type PublicRecord = impl.Record
type Opaque = impl.Secret
type Credential = impl.Credential
type Event = impl.Event
type Forward = impl.Forward
type ClientOptions = opt.Config

func NewRecord() *PublicRecord { return &PublicRecord{} }
`)},
		"options/options.go": {Data: []byte(`package options

import "example.com/sdk/internal/model"

type RequestOptions = model.Options
type Config struct{ Enabled bool }
`)},
		"deps/internal/go.mod": {
			Data: []byte("module example.com/sdk/internal\n\ngo 1.26.0\n"),
		},
		"deps/internal/model/model.go": {Data: []byte(`package model

import tm "time"

type Record struct {
	Name string
	Next *Record
	Options Options
	When tm.Time
	secret *unexposed
}

func (r *Record) Update(other *Record, options Options) Record { return *other }
func (r Record) IsZero() bool { return r.Name == "" }
func (*Record) hidden() {}

type Options struct{ Scopes []string }
type Secret struct{ key *unexposed }
func (*Secret) Reset() {}
type unexposed struct{}

type Event string
func (e Event) String() string { return string(e) }

type Credential interface {
	Get(options Options) (Record, error)
}

type Wrapped struct{ Code int }
func (w Wrapped) Value() int { return w.Code }
type Forward = Wrapped
`)},
	}
	require.NoError(t, os.CopyFS(dir, files))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	origin := loaded.Packages[0].Imports["example.com/sdk/internal/model"]
	require.NotNil(t, origin)
	var before bytes.Buffer
	require.NoError(t, format.Node(&before, origin.Fset, origin.Syntax[0]))

	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	require.Len(t, packages, 2)
	root := packages[0]
	assert.Equal(t, "example.com/sdk/client", root.Path)

	declarations := make(map[string]report.Type)
	rendered := make(map[string]string)
	for _, typ := range root.Types {
		declarations[typ.Name] = typ
		rendered[typ.Name] = renderAliasDeclaration(t, root.FileSet, typ.Decl, typ.Comments)
		require.NotNil(t, typ.Object)
		assert.Equal(t, typ.Name, typ.Object.Name())
		assert.Equal(t, root.Path, typ.Object.Pkg().Path())
	}
	assert.Equal(t, "example.com/sdk/internal/model.Record", declarations["PublicRecord"].AliasOf)
	assert.Equal(t, "example.com/sdk/internal/model.Wrapped", declarations["Forward"].AliasOf)
	assert.Empty(t, declarations["ClientOptions"].AliasOf)
	assert.Equal(t, "type ClientOptions = opt.Config", rendered["ClientOptions"])
	assert.Equal(t, "type Event string", rendered["Event"])
	assert.Contains(t, rendered["Opaque"], "// contains filtered or unexported fields")
	assert.NotContains(t, rendered["Opaque"], "unexposed")
	assert.Contains(t, rendered["Credential"], "Get(options options.RequestOptions) (PublicRecord, error)")

	spec, ok := declarations["PublicRecord"].Decl.Specs[0].(*ast.TypeSpec)
	require.True(t, ok)
	assert.False(t, spec.Assign.IsValid(), "expanded aliases use review declarations")
	structure, ok := spec.Type.(*ast.StructType)
	require.True(t, ok)
	require.Len(t, structure.Fields.List, 4)
	fieldTypes := make(map[string]string)
	for _, field := range structure.Fields.List {
		require.Len(t, field.Names, 1)
		var output bytes.Buffer
		require.NoError(t, format.Node(&output, root.FileSet, field.Type))
		fieldTypes[field.Names[0].Name] = output.String()
	}
	assert.Equal(t, map[string]string{
		"Name": "string", "Next": "*PublicRecord",
		"Options": "options.RequestOptions", "When": "time.Time",
	}, fieldTypes)

	signatures := make(map[string]string)
	functions := append([]report.Function{}, root.Functions...)
	for _, typ := range root.Types {
		functions = append(functions, typ.Functions...)
		functions = append(functions, typ.Methods...)
	}
	for _, function := range functions {
		key := function.Name
		if function.Receiver != "" {
			key = function.Receiver + "." + function.Name
		}
		signatures[key] = renderAliasDeclaration(t, root.FileSet, function.Decl, function.Comments)
	}
	assert.Equal(t, map[string]string{
		"NewRecord":           "func NewRecord() *PublicRecord",
		"Event.String":        "func (Event) String() string",
		"Forward.Value":       "func (Forward) Value() int",
		"Opaque.Reset":        "func (*Opaque) Reset()",
		"PublicRecord.IsZero": "func (PublicRecord) IsZero() bool",
		"PublicRecord.Update": "func (*PublicRecord) Update(other *PublicRecord, options options.RequestOptions) PublicRecord",
	}, signatures)

	var after bytes.Buffer
	require.NoError(t, format.Node(&after, origin.Fset, origin.Syntax[0]))
	assert.Equal(t, before.String(), after.String(), "alias expansion must not mutate dependency syntax")
}

func TestExpandInternalAliasUnaliasedReferences(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(`package sample

import (
	"example.com/sample/internal/generated"
	"example.com/sample/internal/model"
)

type ExpiryType = model.ExpiryType
type ExpiryTypeNever = model.ExpiryTypeNever
type SetExpiryOptions = model.SetExpiryOptions
type Wrapper = model.Wrapper
const DefaultExpiry = generated.DefaultExpiry
`)},
		"internal/model/model.go": {Data: []byte(`package model

import "example.com/sample/internal/generated"

type ExpiryType interface {
	Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions)
	notPubliclyImplementable()
}
type ExpiryTypeNever struct{}
type SetExpiryOptions struct{}
func (ExpiryTypeNever) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions) {
	return generated.DefaultExpiry, nil
}
func (ExpiryTypeNever) notPubliclyImplementable() {}
type Wrapper struct{ Value generated.Payload[int] }
`)},
		"internal/generated/types.go": {Data: []byte(`package generated

type ExpiryOptions string
const DefaultExpiry ExpiryOptions = "Never"
type BlobClientSetExpiryOptions struct{ Header string }
type Payload[T any] struct{ Value T }
`)},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	require.Len(t, packages, 1)
	pkg := packages[0]
	assert.Equal(t, []report.Import{{Path: "example.com/sample/internal/generated"}}, pkg.Imports)
	types := make(map[string]report.Type)
	for _, typ := range pkg.Types {
		types[typ.Name] = typ
	}
	signature := "Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.BlobClientSetExpiryOptions)"
	contract := types["ExpiryType"]
	assert.Contains(t, renderAliasDeclaration(t, pkg.FileSet, contract.Decl, contract.Comments), signature)
	require.Len(t, types["ExpiryTypeNever"].Methods, 1)
	method := types["ExpiryTypeNever"].Methods[0]
	assert.Equal(t, "func (ExpiryTypeNever) "+signature, renderAliasDeclaration(t, pkg.FileSet, method.Decl, method.Comments))
	wrapper := types["Wrapper"]
	assert.Contains(t, renderAliasDeclaration(t, pkg.FileSet, wrapper.Decl, wrapper.Comments), "Value generated.Payload[int]")
	require.Len(t, pkg.Constants, 1)
	var constantType bytes.Buffer
	require.NoError(t, format.Node(&constantType, pkg.FileSet, pkg.Constants[0].TypeExpr))
	assert.Equal(t, "generated.ExpiryOptions", constantType.String())
}

func TestExpandInternalAliasEmbeddedNames(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(`package sample

import "example.com/sample/internal/model"

type Response = model.Response
type Renamed = model.Embedded
type RenamedPointer = model.Pointer
type RecursiveAlias = model.Recursive
`)},
		"other/other.go": {Data: []byte(`package other

import "example.com/sample/internal/model"

type Embedded = model.Embedded
`)},
		"internal/model/model.go": {Data: []byte(`package model

type Embedded struct{ Value int }
type Pointer struct{ Value string }
type Generic[T any] struct{ Value T }
type Response struct {
	Embedded
	*Pointer
	Generic[int]
	Direct Embedded
}
type Recursive struct {
	*Recursive
	Next *Recursive
}
`)},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	require.NotEmpty(t, packages)
	pkg := packages[0]
	declarations := make(map[string]report.Type)
	for _, typ := range pkg.Types {
		declarations[typ.Name] = typ
	}
	spec, ok := declarations["Response"].Decl.Specs[0].(*ast.TypeSpec)
	require.True(t, ok)
	structure, ok := spec.Type.(*ast.StructType)
	require.True(t, ok)
	require.Len(t, structure.Fields.List, 4)
	var embedded []string
	for _, field := range structure.Fields.List[:3] {
		assert.Empty(t, field.Names, "embedded fields must remain embedded")
		var output bytes.Buffer
		require.NoError(t, format.Node(&output, pkg.FileSet, field.Type))
		embedded = append(embedded, output.String())
	}
	assert.Equal(t, []string{"other.Embedded", "model.Generic[int]", "*model.Pointer"}, embedded)
	var direct bytes.Buffer
	require.NoError(t, format.Node(&direct, pkg.FileSet, structure.Fields.List[3].Type))
	assert.Equal(t, "Renamed", direct.String(), "named fields can still prefer differently named public aliases")
	recursive := declarations["RecursiveAlias"]
	assert.Contains(t, renderAliasDeclaration(t, pkg.FileSet, recursive.Decl, recursive.Comments), "*model.Recursive")
	assert.Contains(t, renderAliasDeclaration(t, pkg.FileSet, recursive.Decl, recursive.Comments), "Next *RecursiveAlias")
	assert.Equal(t, []report.Import{
		{Path: "example.com/sample/internal/model"},
		{Path: "example.com/sample/other"},
	}, pkg.Imports)
}

func TestExpandInternalAliasErrors(t *testing.T) {
	t.Setenv("GOWORK", "off")
	tests := []struct {
		name        string
		declaration string
		origin      string
		message     string
	}{
		{
			name:        "generic alias",
			declaration: "type Alias[T any] = impl.Generic[T]",
			origin:      "type Generic[T any] struct{ Value T }",
			message:     "generic alias or instantiated target",
		},
		{
			name:        "instantiated target",
			declaration: "type Alias = impl.Generic[string]",
			origin:      "type Generic[T any] struct{ Value T }",
			message:     "generic alias or instantiated target",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.CopyFS(dir, fstest.MapFS{
				"go.mod":                  {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
				"sample.go":               {Data: []byte("package sample\n\nimport impl \"example.com/sample/internal/model\"\n\n" + tt.declaration + "\n")},
				"internal/model/model.go": {Data: []byte("package model\n\n" + tt.origin + "\n")},
			}))
			loaded, err := loader.Load(t.Context(), dir)
			require.NoError(t, err)
			packages, err := report.ExtractPackages(loaded)
			require.ErrorContains(t, err, tt.message)
			assert.Nil(t, packages)
		})
	}
}

func renderAliasDeclaration(t *testing.T, fset *token.FileSet, decl ast.Decl, comments []*ast.CommentGroup) string {
	t.Helper()
	var output bytes.Buffer
	require.NoError(t, format.Node(&output, fset, &printer.CommentedNode{
		Node: decl, Comments: comments,
	}))
	_, err := format.Source(append([]byte("package sample\n"), output.Bytes()...))
	require.NoError(t, err, "invalid expanded declaration")
	return output.String()
}

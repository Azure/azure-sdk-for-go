// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report_test

import (
	"bytes"
	"go/format"
	"os"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"apigen/internal/loader"
	"apigen/internal/report"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractPackages(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	files := fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(`package sample

// Zulu returns value.
func Zulu(value string) (result string) { return value }

func Alpha[T ~int](values ...T) (sum T) {
	for _, value := range values {
		sum += value
	}
	return
}

type Client[T any] struct{}

func NewClient[T any]() *Client[T] { return &Client[T]{} }

// Send returns value.
func (c *Client[T]) Send(value T) (result T) { return value }
func (c Client[T]) Read() T { var zero T; return zero }
func (*Client[T]) Unnamed(value T) {}
func (c *Client[T]) hiddenMethod() {}
func hidden() {}

type secret struct{}

func (*secret) Hidden() {}
`)},
		"sub/sub.go": {
			Data: []byte("package sample\n\nfunc Child() {}\n"),
		},
		"internalized/public.go": {
			Data: []byte("package visible\n\nfunc Public() {}\n"),
		},
		"internal/helper/helper.go": {
			Data: []byte("package helper\n\nfunc Hidden() {}\n"),
		},
		"deep/internal/hidden/hidden.go": {
			Data: []byte("package hidden\n\nfunc Hidden() {}\n"),
		},
		"cmd/tool/main.go": {
			Data: []byte("package main\n\nfunc main() {}\nfunc Visible() {}\n"),
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
	var paths []string
	for _, pkg := range packages {
		paths = append(paths, pkg.Path)
	}
	require.Equal(t, []string{
		"example.com/sample",
		"example.com/sample/internalized",
		"example.com/sample/sub",
	}, paths)
	assert.Equal(t, "sample", packages[0].Name)
	assert.Equal(t, "visible", packages[1].Name)
	assert.Equal(t, "sample", packages[2].Name)

	var signatures []string
	var receivers []string
	functions := append([]report.Function{}, packages[0].Functions...)
	for _, typ := range packages[0].Types {
		functions = append(functions, typ.Functions...)
		functions = append(functions, typ.Methods...)
	}
	slices.SortFunc(functions, func(a, b report.Function) int {
		if order := strings.Compare(a.Receiver, b.Receiver); order != 0 {
			return order
		}
		return strings.Compare(a.Name, b.Name)
	})
	for _, function := range functions {
		assert.Nil(t, function.Decl.Body)
		assert.Nil(t, function.Decl.Doc)
		if function.Decl.Recv != nil {
			for _, receiver := range function.Decl.Recv.List {
				assert.Empty(t, receiver.Names)
			}
		}
		var signature bytes.Buffer
		require.NoError(t, format.Node(&signature, packages[0].FileSet, function.Decl))
		signatures = append(signatures, signature.String())
		receivers = append(receivers, function.Receiver)
		assert.Equal(t, function.Name, function.Decl.Name.Name)
	}
	assert.Equal(t, []string{
		"func Alpha[T ~int](values ...T) (sum T)",
		"func NewClient[T any]() *Client[T]",
		"func Zulu(value string) (result string)",
		"func (Client[T]) Read() T",
		"func (*Client[T]) Send(value T) (result T)",
		"func (*Client[T]) Unnamed(value T)",
	}, signatures)
	assert.Equal(t, []string{"", "", "", "Client", "Client", "Client"}, receivers)
	require.Len(t, packages[1].Functions, 1)
	assert.Equal(t, "Public", packages[1].Functions[0].Name)
	require.Len(t, packages[2].Functions, 1)
	assert.Equal(t, "Child", packages[2].Functions[0].Name)

	var after bytes.Buffer
	require.NoError(t, format.Node(&after, root.Fset, root.Syntax[0]))
	assert.Equal(t, before.String(), after.String(), "extraction must not mutate loaded syntax")
}

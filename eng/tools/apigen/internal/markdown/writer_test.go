// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package markdown_test

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	"apigen/internal/loader"
	"apigen/internal/markdown"
	"apigen/internal/report"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod": {Data: []byte("module example.com/sample/v2\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(`package sample

import (
	"context"
	"example.com/sample/v2/internal/model"
)

func Identity[T any](value T) T { return value }

// Client sends requests.
type Client struct {
	name string

	state bool
}

func NewClient() (*Client, error) { return &Client{}, nil }
func (c *Client) Send(ctx context.Context, value string) error { return nil }

type ETag string
func NewETag(value string) ETag { return ETag(value) }
func PossibleETagValues() []ETag { return nil }
func (e ETag) String() string { return string(e) }

const (
	ETagZulu ETag = "z"
	ETagAny ETag = "*"
)

type Code int
const CodeZulu Code = 2
const CodeAlpha Code = 1
const Version = 1

type Config struct {
	// Value is the configured value.
	Value *string ` + "`json:\"value,omitempty\"`" + `
	hidden bool
	// Enabled controls behavior.
	// It has a zero-value default.
	Enabled bool
}

var Default = Config{}
var Alternate = Config{Value: nil}

type Record = model.Record
func GetRecord[T any](value T) Record { return Record{} }
func NewRecord() Record { return Record{} }
type ZContract interface { Zulu(); Alpha() }
type AContract = ZContract
type IContract = model.Contract
`)},
		"internal/model/model.go": {Data: []byte(`package model

type Record struct {
	Name string
	hidden int
}

func (r *Record) Reset() {}
type Contract interface { Zulu(); Alpha() }
`)},
		"arm/policy/policy.go": {Data: []byte("package policy\n\ntype ARM string\n")},
		"policy/policy.go":     {Data: []byte("package policy\n\ntype Policy string\n")},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	metadata := report.ExtractModule(loaded)
	var output bytes.Buffer
	require.NoError(t, markdown.Write(&output, metadata, packages))
	expected, err := os.ReadFile(filepath.Join("testdata", "sample.golden"))
	require.NoError(t, err)
	assert.Equal(t, strings.ReplaceAll(string(expected), "\r\n", "\n"), output.String())

	var again bytes.Buffer
	require.NoError(t, markdown.Write(&again, metadata, packages))
	assert.Equal(t, output.String(), again.String(), "rendering must not mutate report declarations")
}

func TestWriteStableOrdering(t *testing.T) {
	t.Setenv("GOWORK", "off")
	public := []string{
		"func Zulu(ctx context.Context) error { return nil }",
		"func Alpha(value time.Time) {}",
		"func DefaultTimeout() time.Duration { return 0 }",
		"type Widget struct {\n\tZuluField string\n\thidden int\n\tAlphaField int\n\tZed, Bravo int\n\ttime.Timer\n\tcontext.Context\n}",
		"type ZContract interface {\n\tZulu(value int) string\n\tAlpha(first string, second int) (last string, count int)\n}",
		"type AContract = ZContract",
		"type OpaqueWrapper struct {\n\tZulu struct{ hidden int }\n\tAlpha struct{ hidden int }\n}",
		"type ZMode = model.Mode",
		"type Mode = model.Mode",
		"type Record = model.Record",
		"const ZuluMode = model.Zulu",
		"const AlphaMode = model.Alpha",
		"func PossibleModeValues() []Mode { return nil }",
		"var ZuluWidget = Widget{}",
		"var AlphaWidget *Widget",
		"func NewZulu() Widget { return Widget{} }",
		"func NewAlpha() *Widget { return nil }",
		"func (w *Widget) Zulu() {}",
		"func (w Widget) Alpha() {}",
	}
	extra := []string{
		"const ZuluConst = 2",
		"const AlphaConst = 1",
		"var ZuluValue = 2",
		"var AlphaValue = 1",
		"var ZuluResult, AlphaResult = pair()",
		"func pair() (int, int) { return 1, 2 }",
	}
	model := []string{
		"type Mode string",
		"const Zulu Mode = \"z\"",
		"const Alpha Mode = \"a\"",
		"type Record struct {\n\tMode Mode\n\tNext *Record\n\thidden int\n}",
		"func (r *Record) Zulu() {}",
		"func (r Record) Alpha() {}",
	}
	fixture := func(reversed bool) fstest.MapFS {
		declarations := make([]string, 3)
		for i, group := range [][]string{public, extra, model} {
			ordered := slices.Clone(group)
			if reversed {
				slices.Reverse(ordered)
				for j, declaration := range ordered {
					declaration = strings.ReplaceAll(declaration,
						"ZuluField string\n\thidden int\n\tAlphaField int\n\tZed, Bravo int",
						"Bravo int\n\tAlphaField int\n\thidden int\n\tZed int\n\tZuluField string")
					declaration = strings.ReplaceAll(declaration,
						"Zulu(value int) string\n\tAlpha(first string, second int) (last string, count int)",
						"Alpha(first string, second int) (last string, count int)\n\tZulu(value int) string")
					declaration = strings.ReplaceAll(declaration, "Mode Mode\n\tNext *Record", "Next *Record\n\tMode Mode")
					declaration = strings.ReplaceAll(declaration, "time.Timer\n\tcontext.Context", "context.Context\n\ttime.Timer")
					declaration = strings.ReplaceAll(declaration, "Zulu struct{ hidden int }\n\tAlpha struct{ hidden int }",
						"Alpha struct{ hidden int }\n\tZulu struct{ hidden int }")
					ordered[j] = declaration
				}
			}
			declarations[i] = strings.Join(ordered, "\n\n") + "\n"
		}
		publicFile, extraFile := "a_public.go", "z_extra.go"
		if reversed {
			publicFile, extraFile = "z_public.go", "a_extra.go"
		}
		return fstest.MapFS{
			"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			publicFile: {Data: []byte(`package sample

import (
	"time"
	"context"
	"example.com/sample/internal/model"
)

` + declarations[0])},
			extraFile: {Data: []byte("package sample\n\n" + declarations[1])},
			"internal/model/model.go": {
				Data: []byte("package model\n\n" + declarations[2]),
			},
			"zulu/zulu.go":   {Data: []byte("package zulu\n\ntype Zulu string\n")},
			"alpha/alpha.go": {Data: []byte("package alpha\n\ntype Alpha string\n")},
		}
	}
	render := func(t *testing.T, files fstest.MapFS, reversePackages bool) string {
		t.Helper()
		dir := t.TempDir()
		require.NoError(t, os.CopyFS(dir, files))
		loaded, err := loader.Load(t.Context(), dir)
		require.NoError(t, err)
		if reversePackages {
			slices.Reverse(loaded.Packages)
		}
		packages, err := report.ExtractPackages(loaded)
		require.NoError(t, err)
		var output bytes.Buffer
		require.NoError(t, markdown.Write(&output, report.ExtractModule(loaded), packages))
		return output.String()
	}
	baseline := render(t, fixture(false), false)
	require.Contains(t, baseline, "type Widget struct {\n\tcontext.Context\n\ttime.Timer\n\tAlphaField int\n")

	t.Run("source and package order", func(t *testing.T) {
		assert.Equal(t, baseline, render(t, fixture(true), true),
			"file names, declaration order, and package traversal must not affect the report")
	})
	t.Run("implementation-only changes", func(t *testing.T) {
		files := fixture(false)
		for name, file := range files {
			if !strings.HasSuffix(name, ".go") {
				continue
			}
			source := string(file.Data)
			source = strings.ReplaceAll(source, "hidden int", "hidden string\n\tanother bool")
			source = strings.ReplaceAll(source, "(w *Widget)", "(receiver *Widget)")
			source = strings.ReplaceAll(source, "(r *Record)", "(receiver *Record)")
			source = strings.ReplaceAll(source, "return 1, 2", "return 10, 20")
			source = strings.ReplaceAll(source, "type Widget", "// Widget has revised documentation.\ntype Widget")
			source = strings.ReplaceAll(source, "{ return nil }", "{\n\t_ = 42\n\treturn nil\n}")
			source += "\n// Implementation details are not part of the report.\nfunc init() {\n\t_ = 42\n}\n"
			file.Data = []byte(source)
		}
		assert.Equal(t, baseline, render(t, files, false),
			"private declarations, fields, comments, bodies, and receiver names must not add report deltas")
	})
	t.Run("public additions", func(t *testing.T) {
		files := fixture(false)
		files["additions.go"] = &fstest.MapFile{Data: []byte(`package sample

func (*Widget) Middle() {}
type NewType string
func Middle() {}
const BetaConstWithALongerName = 2
const BetaMode = ZuluMode
func PossibleCodeValues() []Mode { return nil }
var ZuluTyped = 2
var BetaWidgetWithALongerName = Widget{}
func NewBeta() Widget { return Widget{} }
`)}
		file := files["a_public.go"]
		file.Data = []byte(strings.Replace(string(file.Data), "\tZuluField string",
			"\tLongerThanExistingFieldName string\n\tZuluField string", 1))
		expected := baseline
		for _, addition := range []struct{ before, after string }{
			{"func Zulu(ctx", "func Middle()\nfunc Zulu(ctx"},
			{"func Zulu(ctx", "func PossibleCodeValues() []Mode\nfunc Zulu(ctx"},
			{"type OpaqueWrapper struct", "type NewType string\n\ntype OpaqueWrapper struct"},
			{"func (*Widget) Zulu()", "func (*Widget) Middle()\nfunc (*Widget) Zulu()"},
			{"func NewZulu()", "func NewBeta() Widget\nfunc NewZulu()"},
			{"\tZuluConst = 2", "\tBetaConstWithALongerName = 2\n\tZuluConst = 2"},
			{"\tZuluMode Mode = \"z\"", "\tBetaMode Mode = \"z\"\n\tZuluMode Mode = \"z\""},
			{"\tZuluResult,", "\tBetaWidgetWithALongerName = Widget{}\n\tZuluResult,"},
			{"\tZuluValue = 2", "\tZuluTyped = 2\n\tZuluValue = 2"},
			{"\tZed int", "\tLongerThanExistingFieldName string\n\tZed int"},
		} {
			require.Contains(t, expected, addition.before)
			expected = strings.Replace(expected, addition.before, addition.after, 1)
		}
		assert.Equal(t, expected, render(t, files, false),
			"public additions must insert at their sorted positions without changing existing declarations")
	})
	t.Run("alias addition does not regroup functions", func(t *testing.T) {
		files := fixture(false)
		files["alias.go"] = &fstest.MapFile{Data: []byte("package sample\n\nimport \"time\"\n\ntype Duration = time.Duration\n")}
		before := "type OpaqueWrapper struct"
		require.Contains(t, baseline, before)
		expected := strings.Replace(baseline, before, "type Duration = time.Duration\n\n"+before, 1)
		assert.Equal(t, expected, render(t, files, false))
	})
	t.Run("enum group addition", func(t *testing.T) {
		files := fixture(false)
		files["enum.go"] = &fstest.MapFile{Data: []byte(`package sample

type Kind int
const KindZulu Kind = 2
const KindAlpha Kind = 1
func PossibleKindValues() []Kind { return nil }
`)}
		before := "// Alias of example.com/sample/internal/model.Mode."
		require.Contains(t, baseline, before)
		expected := strings.Replace(baseline, before, "type Kind int\n\nconst (\n"+
			"\tKindAlpha Kind = 1\n\tKindZulu Kind = 2\n)\n\nfunc PossibleKindValues() []Kind\n\n"+before, 1)
		assert.Equal(t, expected, render(t, files, false),
			"a new enum must insert its type and constants together without moving existing declarations")
	})
	t.Run("singleton blocks remain blocks", func(t *testing.T) {
		files := fstest.MapFS{
			"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(`package sample

const First = 1
var Initial = "first"
type Record struct{ Field int }
type Contract interface{ Read() }
type Kind string
const KindFirst Kind = "first"
`)},
		}
		before := render(t, files, false)
		assert.Contains(t, before, "const (\n\tFirst = 1\n)")
		assert.Contains(t, before, "var (\n\tInitial = \"first\"\n)")
		assert.Contains(t, before, "type Kind string\n\nconst (\n\tKindFirst Kind = \"first\"\n)")
		file := files["sample.go"]
		source := string(file.Data)
		source += "\nconst SecondWithALongerName = 2\nvar LaterWithALongerName = \"later\"\n"
		source += "\nconst KindSecondWithALongerName Kind = \"second\"\n"
		source = strings.Replace(source, "Field int", "Field int; LongerField string", 1)
		source = strings.Replace(source, "Read()", "Read(); Write(value string) error", 1)
		file.Data = []byte(source)
		expected := strings.Replace(before, "\tFirst = 1\n", "\tFirst = 1\n\tSecondWithALongerName = 2\n", 1)
		expected = strings.Replace(expected, "\tInitial = \"first\"\n", "\tInitial = \"first\"\n\tLaterWithALongerName = \"later\"\n", 1)
		expected = strings.Replace(expected, "\tKindFirst Kind = \"first\"\n", "\tKindFirst Kind = \"first\"\n\tKindSecondWithALongerName Kind = \"second\"\n", 1)
		expected = strings.Replace(expected, "\tField int\n", "\tField int\n\tLongerField string\n", 1)
		expected = strings.Replace(expected, "\tRead()\n", "\tRead()\n\tWrite(value string) error\n", 1)
		assert.Equal(t, expected, render(t, files, false))
	})
}

func TestWritePreservesMultilineInitializers(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	source := "package sample\n\nvar Alpha = `first\t  \n\n  second \t`\nvar Beta = `last\n\ntail`\nvar Quoted = \"keep  two\\tspaces\"\n"
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod":    {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go": {Data: []byte(source)},
	}))
	loaded, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	packages, err := report.ExtractPackages(loaded)
	require.NoError(t, err)
	var output bytes.Buffer
	require.NoError(t, markdown.Write(&output, report.ExtractModule(loaded), packages))
	_, code, ok := strings.Cut(output.String(), "```go\n")
	require.True(t, ok)
	code, _, ok = strings.Cut(code, "```\n")
	require.True(t, ok)
	parsed, err := parser.ParseFile(token.NewFileSet(), "api.go", code, 0)
	require.NoError(t, err)
	values := make(map[string]string)
	for _, decl := range parsed.Decls {
		general, ok := decl.(*ast.GenDecl)
		if !ok || general.Tok != token.VAR {
			continue
		}
		for _, spec := range general.Specs {
			value, ok := spec.(*ast.ValueSpec)
			require.True(t, ok)
			require.Len(t, value.Values, 1)
			literal, ok := value.Values[0].(*ast.BasicLit)
			require.True(t, ok)
			text, err := strconv.Unquote(literal.Value)
			require.NoError(t, err)
			values[value.Names[0].Name] = text
		}
	}
	assert.Equal(t, map[string]string{
		"Alpha":  "first\t  \n\n  second \t",
		"Beta":   "last\n\ntail",
		"Quoted": "keep  two\tspaces",
	}, values)
}

func TestWriteModuleAndErrors(t *testing.T) {
	t.Run("module without root package or go directive", func(t *testing.T) {
		var output bytes.Buffer
		require.NoError(t, markdown.Write(&output, report.Module{
			Path: "example.com/sample/v2",
			Dependencies: []report.Dependency{
				{Path: "example.com/dependency", Version: "v1.2.3"},
			},
		}, nil))
		assert.Equal(t, "# sample\n\n## Module\n\n```text\nmodule example.com/sample/v2\n\n"+
			"// Direct, non-test, non-standard-library dependencies only.\nrequire (\n"+
			"\texample.com/dependency v1.2.3\n)\n```\n", output.String())
	})
	t.Run("versionless dependency", func(t *testing.T) {
		var output bytes.Buffer
		err := markdown.Write(&output, report.Module{
			Path:         "example.com/sample",
			Dependencies: []report.Dependency{{Path: "example.com/dependency"}},
		}, nil)
		require.ErrorContains(t, err, "no selected module version")
		assert.Empty(t, output.String())
	})
	t.Run("format failure leaves output untouched", func(t *testing.T) {
		var output bytes.Buffer
		err := markdown.Write(&output, report.Module{Path: "example.com/sample"}, []report.Package{
			{Path: "example.com/sample", Name: "not a package name"},
		})
		require.ErrorContains(t, err, "format Go declarations")
		assert.Empty(t, output.String())
	})
	t.Run("output failure", func(t *testing.T) {
		failure := errors.New("output failed")
		err := markdown.Write(failingWriter{err: failure}, report.Module{Path: "example.com/sample"}, nil)
		require.ErrorIs(t, err, failure)
	})
	t.Run("short output", func(t *testing.T) {
		err := markdown.Write(failingWriter{}, report.Module{Path: "example.com/sample"}, nil)
		require.ErrorIs(t, err, io.ErrShortWrite)
	})
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

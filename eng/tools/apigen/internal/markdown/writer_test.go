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
const Separator = '\x61'
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
			{"\t// Inferred type of ZuluResult: int", "\t// Inferred type of BetaWidgetWithALongerName: Widget\n\tBetaWidgetWithALongerName = Widget{}\n\t// Inferred type of ZuluResult: int"},
			{"\t// Inferred type of ZuluValue: int", "\t// Inferred type of ZuluTyped: int\n\tZuluTyped = 2\n\t// Inferred type of ZuluValue: int"},
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
	t.Run("private receiver alias", func(t *testing.T) {
		source := `package sample

type Client struct{}
type hidden = Client
type pointer = *hidden
type private struct{}

func (*hidden) Read() int { return 1 }
func (hidden) Ready() bool { return true }
func (private) NotPublic() {}
`
		files := fstest.MapFS{
			"go.mod":    {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(source)},
		}
		before := render(t, files, false)
		assert.Contains(t, before, "type Client struct{}\n\nfunc (*Client) Read() int\nfunc (Client) Ready() bool")
		assert.NotContains(t, before, "NotPublic")

		files["sample.go"].Data = []byte(strings.NewReplacer("(*hidden) Read", "(*Client) Read", "(hidden) Ready", "(Client) Ready").Replace(source))
		assert.Equal(t, before, render(t, files, false))

		files["sample.go"].Data = []byte(strings.Replace(source, "(*hidden) Read", "(pointer) Read", 1))
		assert.Equal(t, before, render(t, files, false))
	})
	t.Run("public alias receiver methods appear once per alias", func(t *testing.T) {
		const source = `package sample

type hidden struct{}
type Public = hidden
type Other = hidden

func (*Public) Read() *hidden { return nil }
`
		files := fstest.MapFS{
			"go.mod":    {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(source)},
		}
		before := render(t, files, false)
		assert.Equal(t, 1, strings.Count(before, "func (*Public) Read("))
		assert.Contains(t, before, "func (*Public) Read() *Public\n")
		assert.Equal(t, 1, strings.Count(before, "func (*Other) Read("))
		assert.Contains(t, before, "func (*Other) Read() *Other\n")

		files["sample.go"].Data = []byte(strings.Replace(source, "(*Public) Read", "(*hidden) Read", 1))
		assert.Equal(t, before, render(t, files, false))
	})
	t.Run("private interface embeddings", func(t *testing.T) {
		source := `package sample

type Visible interface{ Shared() }
type reader[T any] interface {
	Read() T
	Shared()
	Reset()
	sealed()
}
type hidden = reader[string]
type transitive interface{ hidden }
type unrelated interface{ Unrelated() }

type Contract interface {
	Visible
	interface{ transitive }
	Reset()
}
type Plain interface{ hidden }
`
		files := fstest.MapFS{
			"go.mod":    {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(source)},
		}
		before := render(t, files, false)
		assert.Contains(t, before, "// Embedded method on Contract: Read() string")
		assert.Contains(t, before, "type Plain interface {\n\thidden\n}\n\n// Embedded method on Plain: Read() string")
		assert.NotContains(t, before, "// Embedded method on Contract: Shared()")
		assert.NotContains(t, before, "// Embedded method on Contract: Reset()")
		assert.NotContains(t, before, "sealed()")
		assert.NotContains(t, before, "Unrelated")
		files["sample.go"].Data = []byte(strings.Replace(source, "Read() T", "Read() []T", 1))
		assert.Equal(t, strings.ReplaceAll(before, "Read() string", "Read() []string"), render(t, files, false))
	})
	t.Run("exported internal-package embeddings", func(t *testing.T) {
		source := `package sample

import "example.com/sample/internal/generated"

type RenameResponse struct {
	generated.Response
	generated.Details
	Explicit string
}
type PointerResponse struct{ *generated.Response }
type PublicResponse = generated.Response
type PublicWrapper struct{ PublicResponse }
`
		embedded := `package generated

type Response struct {
	ETag *string
	Date *string
}
type Details struct{ ID *string }
func (*Response) Reset() {}
`
		files := fstest.MapFS{
			"go.mod":                         {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go":                      {Data: []byte(source)},
			"internal/generated/response.go": {Data: []byte(embedded)},
		}
		before := render(t, files, false)
		assert.Contains(t, before, `type RenameResponse struct {
	generated.Details
	// Promoted fields:
	//   ID *string
	generated.Response
	// Promoted fields:
	//   Date *string
	//   ETag *string
	Explicit string
}

// Promoted method on *RenameResponse: Reset()`)
		assert.Contains(t, before, `type PointerResponse struct {
	*generated.Response
	// Promoted fields:
	//   Date *string
	//   ETag *string
}

// Promoted method on PointerResponse: Reset()`)
		assert.Equal(t, 2, strings.Count(before, "//   ETag *string"))
		assert.NotContains(t, before, "// Promoted field: ETag")
		assert.NotContains(t, before, "// Promoted method on PublicWrapper:")
		assert.NotContains(t, before, "// Promoted method on *PublicWrapper:")

		files["internal/generated/response.go"].Data = []byte(strings.Replace(embedded, "ETag *string", "ETag string", 1))
		assert.Equal(t, strings.ReplaceAll(before, "ETag *string", "ETag string"), render(t, files, false))
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
		assert.Contains(t, before, "var (\n\t// Inferred type of Initial: string\n\tInitial = \"first\"\n)")
		assert.Contains(t, before, "type Kind string\n\nconst (\n\tKindFirst Kind = \"first\"\n)")
		file := files["sample.go"]
		source := string(file.Data)
		source += "\nconst SecondWithALongerName = 2\nvar LaterWithALongerName = \"later\"\n"
		source += "\nconst KindSecondWithALongerName Kind = \"second\"\n"
		source = strings.Replace(source, "Field int", "Field int; LongerField string", 1)
		source = strings.Replace(source, "Read()", "Read(); Write(value string) error", 1)
		file.Data = []byte(source)
		expected := strings.Replace(before, "\tFirst = 1\n", "\tFirst = 1\n\tSecondWithALongerName = 2\n", 1)
		expected = strings.Replace(expected, "\tInitial = \"first\"\n", "\tInitial = \"first\"\n\t// Inferred type of LaterWithALongerName: string\n\tLaterWithALongerName = \"later\"\n", 1)
		expected = strings.Replace(expected, "\tKindFirst Kind = \"first\"\n", "\tKindFirst Kind = \"first\"\n\tKindSecondWithALongerName Kind = \"second\"\n", 1)
		expected = strings.Replace(expected, "\tField int\n", "\tField int\n\tLongerField string\n", 1)
		expected = strings.Replace(expected, "\tRead()\n", "\tRead()\n\tWrite(value string) error\n", 1)
		assert.Equal(t, expected, render(t, files, false))
	})
	t.Run("expanded alias target identity", func(t *testing.T) {
		const models = `package model

type A struct{ Name string }
type B struct{ Name string }
type Forward = A
type Alternate = A
type Number = int
type List = []int
`
		files := fstest.MapFS{
			"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(`package sample

import "example.com/sample/internal/model"

type Public = model.Forward
type Number = model.Number
type List = model.List
`)},
			"internal/model/model.go": {Data: []byte(models)},
		}
		before := render(t, files, false)
		require.Contains(t, before, "// Alias of example.com/sample/internal/model.A.\ntype Public struct")
		assert.Contains(t, before, "// Alias of int.\ntype Number int")
		assert.Contains(t, before, "// Alias of an unnamed type.\ntype List []int")

		files["internal/model/model.go"].Data = []byte(strings.Replace(models, "Forward = A", "Forward = Alternate", 1))
		assert.Equal(t, before, render(t, files, false))

		files["internal/model/model.go"].Data = []byte(strings.Replace(models, "Forward = A", "Forward = B", 1))
		expected := strings.Replace(before, "/model.A.", "/model.B.", 1)
		assert.Equal(t, expected, render(t, files, false))
	})
	t.Run("inferred variable types", func(t *testing.T) {
		files := fstest.MapFS{
			"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(`package sample

import "example.com/sample/internal/value"

var Value = makeValue()
var Label, hidden, Count = makeResults()
var Explicit bool = true

func makeValue() int { return 1 }
func makeResults() (string, bool, *value.Count) { return "", false, nil }
`)},
			"internal/value/value.go": {Data: []byte("package value\n\ntype Count int\n")},
		}
		before := render(t, files, false)
		assert.Contains(t, before, `var (
	Explicit bool = true
	// Inferred type of Label: string
	// Inferred type of Count: *example.com/sample/internal/value.Count
	Label, _, Count = makeResults()
	// Inferred type of Value: int
	Value = makeValue()
)`)
		assert.NotContains(t, before, `import (`)
		file := files["sample.go"]
		source := string(file.Data)
		file.Data = []byte(strings.Replace(source, "makeValue() int { return 1 }", `makeValue() string { return "one" }`, 1))
		assert.Equal(t, strings.Replace(before, "Inferred type of Value: int", "Inferred type of Value: string", 1),
			render(t, files, false))
		file.Data = []byte(strings.Replace(source, "(string, bool, *value.Count)", "(string, bool, []value.Count)", 1))
		assert.Equal(t, strings.Replace(before, "*example.com/sample/internal/value.Count", "[]example.com/sample/internal/value.Count", 1),
			render(t, files, false))
	})
	t.Run("function literal initializers", func(t *testing.T) {
		source := `package sample

import (
	"context"
	"strings"
)

type Callback func(string) string

var Explicit Callback = func(s string) string { return strings.TrimSpace(s) }
var Handler, Other = (func(ctx context.Context) error { return ctx.Err() }), 42
`
		files := fstest.MapFS{
			"go.mod":    {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(source)},
		}
		before := render(t, files, false)
		assert.Contains(t, before, `var (
	Explicit Callback
	Handler func(ctx context.Context) error
	// Inferred type of Other: int
	Other = 42
)`)
		assert.Contains(t, before, `"context"`)
		assert.NotContains(t, before, `"strings"`)

		bodyChange := strings.Replace(source, `return strings.TrimSpace(s)`, `return strings.ToUpper(s)`, 1)
		bodyChange = strings.Replace(bodyChange, `return ctx.Err()`, `return nil`, 1)
		files["sample.go"].Data = []byte(bodyChange)
		assert.Equal(t, before, render(t, files, false))

		files["sample.go"].Data = []byte(strings.Replace(bodyChange, "context.Context) error", "context.Context) *error", 1))
		assert.Equal(t, strings.Replace(before, "context.Context) error", "context.Context) *error", 1),
			render(t, files, false))
	})
	t.Run("nested function literal initializers", func(t *testing.T) {
		source := `package sample

import (
	"strconv"
	"strings"
	clock "time"

	probe "example.com/sample/probe"
)

var Called = (func() int { return strconv.IntSize })()
var Callback = use(func() string { return strings.TrimSpace(" text ") })
var Composite = []func() int{func() int { return 7 }}
var Explicit any = func() string { return strings.ToUpper("text") }()
var Moment, hidden, Flag = func() (clock.Time, string, bool) {
	return clock.Time{}, strings.ToLower("TEXT"), false
}()
var Ordinary = use(named)
var Resource = func() probe.Token { return 0 }()

func Now() clock.Time { return clock.Now() }
func use(cb func() string) string { return cb() }
func named() string { return "name" }
`
		files := fstest.MapFS{
			"go.mod":         {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go":      {Data: []byte(source)},
			"probe/probe.go": {Data: []byte("package clock\n\ntype Token int\n")},
		}
		before := render(t, files, false)
		assert.Contains(t, before, `var (
	Callback string
	Called int
	Composite []func() int
	Explicit any
	Flag bool
	Moment clock.Time
	// Inferred type of Ordinary: string
	Ordinary = use(named)
	Resource clock2.Token
)`)
		assert.Contains(t, before, `clock "time"`)
		assert.Contains(t, before, `clock2 "example.com/sample/probe"`)
		assert.NotContains(t, before, `"strconv"`)
		assert.NotContains(t, before, `"strings"`)

		changed := strings.Replace(source, `return strconv.IntSize`, `return strconv.IntSize + 1`, 1)
		changed = strings.Replace(changed, `" text "`, `" revised text "`, 1)
		files["sample.go"].Data = []byte(changed)
		assert.Equal(t, before, render(t, files, false))

		changed = strings.Replace(changed, "clock.Time, string, bool", "clock.Time, string, int", 1)
		changed = strings.Replace(changed, `"TEXT"), false`, `"TEXT"), 0`, 1)
		files["sample.go"].Data = []byte(changed)
		assert.Equal(t, strings.Replace(before, "Flag bool", "Flag int", 1), render(t, files, false))
	})
}

func TestWriteExpandsDefinedTypesWithHiddenBases(t *testing.T) {
	t.Setenv("GOWORK", "off")
	render := func(t *testing.T, audienceType string) string {
		t.Helper()
		dir := t.TempDir()
		require.NoError(t, os.CopyFS(dir, fstest.MapFS{
			"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(`package sample

import "example.com/sample/internal/base"

type Options base.Options
func (Options) Own() {}
type Public[T any] base.Box[T]
type hidden struct{ Size [8]byte }
type Local hidden
`)},
			"internal/base/options.go": {Data: []byte(`package base

import "time"

type Options struct {
	Audience ` + audienceType + `
	Next *Options
	When time.Time
	hidden bool
}
func (Options) BaseOnly() {}

type Box[T any] struct{ Value T }
`)},
		}))
		loaded, err := loader.Load(t.Context(), dir)
		require.NoError(t, err)
		packages, err := report.ExtractPackages(loaded)
		require.NoError(t, err)
		var output bytes.Buffer
		require.NoError(t, markdown.Write(&output, report.ExtractModule(loaded), packages))
		return output.String()
	}

	before := render(t, "string")
	assert.Contains(t, before, `type Options struct {
	Audience string
	Next *base.Options
	When time.Time
	// contains filtered or unexported fields
}

func (Options) Own()`)
	assert.Contains(t, before, "type Public[T any] struct {\n\tValue T\n}")
	assert.Contains(t, before, "type Local struct {\n\tSize [8]byte\n}")
	assert.NotContains(t, before, "BaseOnly")
	assert.NotContains(t, before, "// Alias of")
	assert.Equal(t, strings.Replace(before, "Audience string", "Audience []byte", 1), render(t, "[]byte"))
}

func TestWriteExpandsAliasToLocalUnexportedType(t *testing.T) {
	t.Setenv("GOWORK", "off")
	render := func(t *testing.T, target string) string {
		t.Helper()
		dir := t.TempDir()
		require.NoError(t, os.CopyFS(dir, fstest.MapFS{
			"go.mod":   {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"alias.go": {Data: []byte("package sample\n\ntype Public = private\ntype Forward = intermediate\n")},
			"types.go": {Data: []byte("package sample\n\ntype intermediate = private\n\n" + target)},
		}))
		loaded, err := loader.Load(t.Context(), dir)
		require.NoError(t, err)
		packages, err := report.ExtractPackages(loaded)
		require.NoError(t, err)
		var output bytes.Buffer
		require.NoError(t, markdown.Write(&output, report.ExtractModule(loaded), packages))
		return output.String()
	}
	source := `type private struct {
	Size int
	Next *private
	hidden bool
}

func (*private) Reset() *private { return nil }
func (*private) hide() {}
`
	before := render(t, source)
	require.Contains(t, before, `// Alias of example.com/sample.private.
type Public struct {
	Next *Public
	Size int
	// contains filtered or unexported fields
}

func (*Public) Reset() *Public`)
	require.Contains(t, before, `// Alias of example.com/sample.private.
type Forward struct {
	Next *Forward
	Size int
	// contains filtered or unexported fields
}

func (*Forward) Reset() *Forward`)
	assert.NotContains(t, before, "func (*private)")
	assert.NotContains(t, before, "hide()")

	t.Run("exported field change", func(t *testing.T) {
		assert.Equal(t, strings.ReplaceAll(before, "Size int", "Size int64"),
			render(t, strings.Replace(source, "Size int", "Size int64", 1)))
	})
	t.Run("exported method change", func(t *testing.T) {
		assert.Equal(t, strings.ReplaceAll(before, "Reset()", "Reset(force bool)"),
			render(t, strings.Replace(source, "Reset()", "Reset(force bool)", 1)))
	})
	t.Run("unexported field change", func(t *testing.T) {
		assert.Equal(t, before, render(t, strings.Replace(source, "hidden bool", "hidden string", 1)))
	})
}

func TestWriteEvaluatesArrayLengths(t *testing.T) {
	t.Setenv("GOWORK", "off")
	render := func(t *testing.T, width, lengthExpr, aliasSize string) string {
		t.Helper()
		dir := t.TempDir()
		require.NoError(t, os.CopyFS(dir, fstest.MapFS{
			"go.mod": {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
			"sample.go": {Data: []byte(`package sample

import (
	"example.com/sample/internal/bounds"
	"example.com/sample/internal/model"
)

const size = bounds.Width
type Buffer [` + lengthExpr + `]byte
type Block = model.Block
var Cache [bounds.Width]byte
func Read(p *[size]byte) [0x0][]byte { return [0][]byte{} }
`)},
			"internal/bounds/bounds.go": {Data: []byte("package bounds\n\nconst Width = " + width + "\n")},
			"internal/model/model.go": {Data: []byte(`package model

const size = ` + aliasSize + `
type Block struct{ Data [size]byte }
func (*Block) DataCopy() [size]byte { return [size]byte{} }
`)},
		}))
		loaded, err := loader.Load(t.Context(), dir)
		require.NoError(t, err)
		packages, err := report.ExtractPackages(loaded)
		require.NoError(t, err)
		require.Len(t, packages, 1)
		assert.Empty(t, packages[0].Imports, "imports used only to compute array lengths should not appear")
		var output bytes.Buffer
		require.NoError(t, markdown.Write(&output, report.ExtractModule(loaded), packages))
		return output.String()
	}

	before := render(t, "16", "size", "4")
	assert.Contains(t, before, "type Buffer [16]byte")
	assert.Contains(t, before, "Cache [16]byte")
	assert.Contains(t, before, "func Read(p *[16]byte) [0][]byte")
	assert.Contains(t, before, "Data [4]byte")
	assert.Contains(t, before, "func (*Block) DataCopy() [4]byte")
	assert.Equal(t, before, render(t, "0x10", "size + 0", "2 * 2"))
	assert.Equal(t, strings.ReplaceAll(before, "[16]", "[32]"), render(t, "32", "size", "4"))
	assert.Equal(t, strings.ReplaceAll(before, "[4]", "[8]"), render(t, "16", "size", "8"))
}

func TestWriteIncludesMembersPromotedThroughPrivateEmbeddings(t *testing.T) {
	t.Setenv("GOWORK", "off")
	render := func(t *testing.T, source string) string {
		t.Helper()
		dir := t.TempDir()
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
		return output.String()
	}
	source := `package sample

type base[T any] struct {
	Count T
	Local bool
	Shared int
}
func (base[T]) Read() T { var zero T; return zero }
func (*base[T]) Reset() {}
func (base[T]) Clash() {}
func (base[T]) Covered() {}

type hidden struct{ base[int] }
type second struct{ secondBase }
type secondBase struct{ Shared string }
func (secondBase) Clash() {}
type Exported struct{ Documented string }

type Public struct {
	Exported
	hidden
	second
	Local int
}
func (Public) Covered() {}
`
	before := render(t, source)
	require.Contains(t, before, `type Public struct {
	Exported
	Local int
	// contains filtered or unexported fields
}

func (Public) Covered()

// Promoted field: Count int
// Promoted method on Public: Read() int
// Promoted method on *Public: Reset()`)
	for _, ignored := range []string{
		"Promoted field: Local", "Promoted field: Shared", "Promoted field: Documented",
		": Clash()", ": Covered()",
	} {
		assert.NotContains(t, before, ignored)
	}
	expected := strings.Replace(before, "Promoted field: Count int", "Promoted field: Count string", 1)
	expected = strings.Replace(expected, "Promoted method on Public: Read() int", "Promoted method on Public: Read() string", 1)
	assert.Equal(t, expected, render(t, strings.Replace(source, "base[int]", "base[string]", 1)))
	assert.Equal(t, strings.Replace(before, "on *Public: Reset()", "on Public: Reset()", 1),
		render(t, strings.Replace(source, "\n\thidden\n", "\n\t*hidden\n", 1)))
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

func TestWriteOmitsStructTags(t *testing.T) {
	t.Setenv("GOWORK", "off")
	const source = `package sample

import "time"

type hidden struct {
	Nested struct{ When time.Time "json:\"tag-before\"" } "json:\"tag-before\""
}

func (hidden) Read() struct{ Value string "json:\"tag-before\"" } { panic("not implemented") }

type Visible struct {
	hidden
	Field string "json:\"tag-before\""
}

type private interface {
	Read() struct{ Value string "json:\"tag-before\"" }
}
type Contract interface{ private }

var Inline = struct{ Value string "json:\"tag-before\"" }{Value: "json:\"keep me\""}
`
	render := func(t *testing.T, source string) string {
		t.Helper()
		dir := t.TempDir()
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
		return output.String()
	}

	got := render(t, source)
	assert.NotContains(t, got, "tag-before")
	assert.Contains(t, got, `Value: "json:\"keep me\""`)
	assert.Contains(t, got, "// Inferred type of Inline: struct{Value string}")
	assert.Contains(t, got, "// Promoted field: Nested struct{When time.Time}")
	assert.Contains(t, got, "// Promoted method on Visible: Read() struct{Value string}")
	assert.Contains(t, got, "// Embedded method on Contract: Read() struct{Value string}")
	assert.Equal(t, got, render(t, strings.ReplaceAll(source, "tag-before", "tag-after")))
	assert.Equal(t, got, render(t, strings.ReplaceAll(source, ` "json:\"tag-before\""`, "")))
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

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package markdown renders API review reports as Markdown.
package markdown

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/scanner"
	"go/token"
	"go/types"
	"io"
	"path"
	"strconv"
	"strings"

	"apigen/internal/report"

	"golang.org/x/mod/module"
)

// Write renders module metadata and public packages using the API review format.
// Its inputs must come from report extraction. Rendering completes before any
// output is written; formatting and output errors are returned to the caller.
func Write(output io.Writer, metadata report.Module, packages []report.Package) error {
	prefix, _, ok := module.SplitPathVersion(metadata.Path)
	if !ok {
		return fmt.Errorf("invalid module path %q", metadata.Path)
	}
	title := path.Base(prefix)
	for _, pkg := range packages {
		if pkg.Path == metadata.Path {
			title = pkg.Name
			break
		}
	}

	var document bytes.Buffer
	fmt.Fprintf(&document, "# %s\n\n## Module\n\n```text\nmodule %s\n", title, metadata.Path)
	if metadata.GoVersion != "" {
		fmt.Fprintf(&document, "\ngo %s\n", metadata.GoVersion)
	}
	if len(metadata.Dependencies) > 0 {
		document.WriteString("\n// Direct, non-test, non-standard-library dependencies only.\nrequire (\n")
		for _, dependency := range metadata.Dependencies {
			if dependency.Version == "" {
				return fmt.Errorf("dependency %q has no selected module version", dependency.Path)
			}
			fmt.Fprintf(&document, "\t%s %s\n", dependency.Path, dependency.Version)
		}
		document.WriteString(")\n")
	}
	document.WriteString("```\n")

	for _, pkg := range packages {
		name := pkg.Name
		if pkg.Path != metadata.Path {
			if !strings.HasPrefix(pkg.Path, metadata.Path+"/") {
				return fmt.Errorf("package %q is outside module %q", pkg.Path, metadata.Path)
			}
			name = strings.TrimPrefix(pkg.Path, metadata.Path+"/")
		}
		source, err := packageSource(pkg)
		if err != nil {
			return fmt.Errorf("render package %q: %w", pkg.Path, err)
		}
		fmt.Fprintf(&document, "\n## Package %s\n\n```go\n", name)
		document.Write(source)
		document.WriteString("```\n")
	}
	if _, err := document.WriteTo(output); err != nil {
		return fmt.Errorf("write API report: %w", err)
	}
	return nil
}

func packageSource(pkg report.Package) ([]byte, error) {
	writer := packageWriter{fileSet: pkg.FileSet}
	fmt.Fprintf(&writer.source, "package %s\n\n", pkg.Name)
	if len(pkg.Imports) > 0 {
		writer.source.WriteString("import (\n")
		for i, entry := range pkg.Imports {
			if i > 0 && pkg.Imports[i-1].Standard != entry.Standard {
				writer.source.WriteByte('\n')
			}
			writer.source.WriteByte('\t')
			if entry.Name != "" {
				fmt.Fprintf(&writer.source, "%s ", entry.Name)
			}
			fmt.Fprintln(&writer.source, strconv.Quote(entry.Path))
		}
		writer.source.WriteString(")\n\n")
	}
	if err := writer.constants(pkg.Constants); err != nil {
		return nil, err
	}
	for _, typ := range pkg.Types {
		if len(typ.Constants) > 0 {
			if err := writer.typeDeclaration(typ); err != nil {
				return nil, err
			}
		}
	}
	if err := writer.variables(pkg.Variables); err != nil {
		return nil, err
	}
	if err := writer.functions(pkg.Functions); err != nil {
		return nil, err
	}
	for _, typ := range pkg.Types {
		if len(typ.Constants) == 0 {
			if err := writer.typeDeclaration(typ); err != nil {
				return nil, err
			}
		}
	}
	fset := token.NewFileSet()
	source, err := parser.ParseFile(fset, "api.go", writer.source.Bytes(), parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("format Go declarations: %w", err)
	}
	for node := range ast.Preorder(source) {
		switch node := node.(type) {
		case *ast.Field:
			node.Tag = nil
		case *ast.StructType:
			if len(node.Fields.List) == 1 {
				// gofmt can expand a single-field struct just because it has a tag.
				// Keep it multiline whether or not a tag was present.
				node.Fields.Opening = token.NoPos
			}
		}
	}
	var formatted bytes.Buffer
	if err := format.Node(&formatted, fset, source); err != nil {
		return nil, fmt.Errorf("format Go declarations: %w", err)
	}
	return fixedSpacing(formatted.Bytes())
}

func fixedSpacing(source []byte) ([]byte, error) {
	fileSet := token.NewFileSet()
	file := fileSet.AddFile("api.go", -1, len(source))
	var scan scanner.Scanner
	var scanError error
	scan.Init(file, source, func(pos token.Position, message string) {
		if scanError == nil {
			scanError = fmt.Errorf("scan Go declarations at %s: %s", pos, message)
		}
	}, scanner.ScanComments)
	var output bytes.Buffer
	start := 0
	for {
		pos, kind, literal := scan.Scan()
		if kind == token.SEMICOLON && literal == "\n" {
			continue
		}
		end := file.Offset(pos)
		if kind == token.EOF {
			end = len(source)
		}
		part := source[start:end]
		text := bytes.TrimRight(part, " \t\r\n")
		gap := part[len(text):]
		output.Write(text)
		if len(gap) > 0 && !bytes.ContainsAny(gap, "\r\n") {
			output.WriteByte(' ')
		} else {
			output.Write(gap)
		}
		start = end
		if kind == token.EOF {
			break
		}
	}
	if scanError != nil {
		return nil, scanError
	}
	return output.Bytes(), nil
}

type packageWriter struct {
	fileSet *token.FileSet
	source  bytes.Buffer
}

func (w *packageWriter) typeDeclaration(typ report.Type) error {
	if typ.AliasOf != "" {
		fmt.Fprintf(&w.source, "// Alias of %s.\n", typ.AliasOf)
	}
	text, err := w.node(typ.Decl, typ.Comments)
	if err != nil {
		return fmt.Errorf("type %q: %w", typ.Name, err)
	}
	text, unattached, err := annotatePromotedFields(text, typ.PromotedFields)
	if err != nil {
		return fmt.Errorf("type %q: %w", typ.Name, err)
	}
	w.source.WriteString(text)
	w.source.WriteString("\n\n")
	if err := w.constants(typ.Constants); err != nil {
		return err
	}
	if err := w.functions(typ.Functions); err != nil {
		return err
	}
	if err := w.functions(typ.Methods); err != nil {
		return err
	}
	for _, field := range unattached {
		fmt.Fprintf(&w.source, "// Promoted field: %s\n", field.Signature)
	}
	for _, method := range typ.PromotedMethods {
		receiver := typ.Name
		if method.PointerOnly {
			receiver = "*" + receiver
		}
		fmt.Fprintf(&w.source, "// Promoted method on %s: %s\n", receiver, method.Signature)
	}
	for _, method := range typ.EmbeddedMethods {
		fmt.Fprintf(&w.source, "// Embedded method on %s: %s\n", typ.Name, method)
	}
	if len(unattached)+len(typ.PromotedMethods)+len(typ.EmbeddedMethods) > 0 {
		w.source.WriteByte('\n')
	}
	return nil
}

func annotatePromotedFields(text string, fields []report.PromotedField) (string, []report.PromotedField, error) {
	byEmbedding := make(map[string][]string)
	for _, field := range fields {
		if token.IsExported(field.Embedding) {
			byEmbedding[field.Embedding] = append(byEmbedding[field.Embedding], field.Signature)
		}
	}
	if len(byEmbedding) == 0 {
		return text, fields, nil
	}

	// Use the formatted declaration's positions: normalization can reorder fields
	// and can copy them from another package.
	const prefix = "package report\n"
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, "type.go", prefix+text, parser.SkipObjectResolution)
	if err != nil {
		return "", nil, fmt.Errorf("parse declaration for promoted field placement: %w", err)
	}
	var structure *ast.StructType
	ast.Inspect(parsed, func(node ast.Node) bool {
		if spec, ok := node.(*ast.TypeSpec); ok {
			structure, _ = ast.Unparen(spec.Type).(*ast.StructType)
			return false
		}
		return true
	})
	if structure == nil {
		return text, fields, nil
	}

	var annotated strings.Builder
	start := 0
	placed := make(map[string]bool)
	for _, field := range structure.Fields.List {
		name, err := report.EmbeddedFieldName(field)
		if err != nil {
			return "", nil, err
		}
		signatures := byEmbedding[name]
		if len(signatures) == 0 {
			continue
		}
		end := fset.Position(field.End()).Offset - len(prefix)
		annotated.WriteString(text[start:end])
		annotated.WriteString("\n\t// Promoted fields:")
		for _, signature := range signatures {
			fmt.Fprint(&annotated, "\n\t//   ", signature)
		}
		placed[name] = true
		start = end
	}
	annotated.WriteString(text[start:])
	var unattached []report.PromotedField
	for _, field := range fields {
		if !placed[field.Embedding] {
			unattached = append(unattached, field)
		}
	}
	return annotated.String(), unattached, nil
}

func (w *packageWriter) node(node ast.Node, comments []*ast.CommentGroup) (string, error) {
	var output bytes.Buffer
	if err := format.Node(&output, w.fileSet, &printer.CommentedNode{Node: node, Comments: comments}); err != nil {
		return "", err
	}
	return compactDeclaration(strings.TrimSpace(output.String()))
}

func compactDeclaration(text string) (string, error) {
	fileSet := token.NewFileSet()
	file := fileSet.AddFile("declaration.go", -1, len(text))
	var scan scanner.Scanner
	var scanError error
	scan.Init(file, []byte(text), func(pos token.Position, message string) {
		if scanError == nil {
			scanError = fmt.Errorf("scan declaration at %s: %s", pos, message)
		}
	}, scanner.ScanComments)
	protected := make(map[int]bool)
	for {
		pos, kind, literal := scan.Scan()
		if kind == token.EOF {
			break
		}
		if kind == token.STRING || kind == token.CHAR || kind == token.COMMENT {
			if count := strings.Count(literal, "\n"); count > 0 {
				first := file.Line(pos)
				for line := first; line <= first+count; line++ {
					protected[line] = true
				}
			}
		}
	}
	if scanError != nil {
		return "", scanError
	}
	var lines []string
	for i, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) != "" || protected[i+1] {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n"), nil
}

func (w *packageWriter) functions(functions []report.Function) error {
	for _, function := range functions {
		text, err := w.node(function.Decl, function.Comments)
		if err != nil {
			return fmt.Errorf("function %q: %w", function.Name, err)
		}
		w.source.WriteString(text)
		w.source.WriteString("\n")
	}
	if len(functions) > 0 {
		w.source.WriteByte('\n')
	}
	return nil
}

func (w *packageWriter) constants(constants []report.Constant) error {
	var entries []string
	for _, item := range constants {
		entry := item.Name
		if item.TypeExpr != nil {
			text, err := w.node(item.TypeExpr, item.Comments)
			if err != nil {
				return fmt.Errorf("constant %q type: %w", item.Name, err)
			}
			entry += " " + text
		}
		value, err := constantExpression(item.Value, item.Type)
		if err != nil {
			return fmt.Errorf("constant %q: %w", item.Name, err)
		}
		entries = append(entries, entry+" = "+value)
	}
	w.group("const", entries)
	return nil
}

func (w *packageWriter) variables(variables []report.Variable) error {
	var entries []string
	for _, variable := range variables {
		text, err := w.node(variable.Decl, variable.Comments)
		if err != nil {
			return fmt.Errorf("variable declaration: %w", err)
		}
		entry, ok := strings.CutPrefix(text, "var ")
		if !ok {
			return fmt.Errorf("expected a var declaration, got %q", text)
		}
		if variable.Inferred {
			for _, object := range variable.Objects {
				typeName := report.TypeStringWithoutTags(object.Type(), types.RelativeTo(object.Pkg()))
				entries = append(entries, fmt.Sprintf("// Inferred type of %s: %s", object.Name(), typeName))
			}
		}
		entries = append(entries, entry)
	}
	w.group("var", entries)
	return nil
}

func (w *packageWriter) group(keyword string, entries []string) {
	if len(entries) == 0 {
		return
	}
	fmt.Fprintf(&w.source, "%s (\n", keyword)
	for _, entry := range entries {
		fmt.Fprintf(&w.source, "\t%s\n", entry)
	}
	w.source.WriteString(")\n\n")
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"apigen/internal/loader"

	"golang.org/x/tools/go/packages"
)

// Package groups report declarations for an importable Go package.
type Package struct {
	// Path is the package's import path.
	Path string

	// Name is the package name declared in its source files.
	Name string

	// FileSet describes the source positions in the package's declarations.
	FileSet *token.FileSet

	// Imports contains packages referenced by displayed declarations.
	// Standard library imports appear first; each group is sorted by path and alias.
	Imports []Import

	// Functions contains standalone exported functions, sorted by name.
	Functions []Function

	// Types contains exported type declarations: interfaces first, then other
	// types, sorted by name within each group. Interface aliases are interfaces.
	Types []Type

	// Constants contains exported constants not grouped with an enum type,
	// sorted by name.
	Constants []Constant

	// Variables contains all exported variable declarations, sorted by their
	// first exported binding name.
	Variables []Variable
}

// Function describes an exported function or method signature.
type Function struct {
	// Name is the function or method name.
	Name string

	// Object is the compiler identity of the source function or method.
	Object *types.Func

	// Receiver is the receiver's base type name, or empty for a standalone function.
	Receiver string

	// Decl contains the signature without its body, declaration comment, or
	// receiver variable names. It shares type-expression nodes with the loaded
	// source and must be treated as read-only.
	Decl *ast.FuncDecl

	// Comments is a non-nil list of report-generated comments for Decl.
	// Use it with printer.CommentedNode when formatting to suppress source comments.
	Comments []*ast.CommentGroup
}

// ExtractPackages extracts public package metadata and exported declarations.
// It excludes main packages, internal packages, and methods declared on
// unexported receiver types.
// loaded must be a successful result from loader.Load.
func ExtractPackages(loaded *loader.Result) ([]Package, error) {
	aliases, err := newAliasExpander(loaded)
	if err != nil {
		return nil, err
	}
	packageInfo := make(map[*types.Package]*packages.Package)
	packages.Visit(loaded.Packages, nil, func(pkg *packages.Package) {
		packageInfo[pkg.Types] = pkg
	})
	var result []Package
	for _, pkg := range loaded.Packages {
		if !isPublicPackage(pkg) {
			continue
		}
		refs := make(map[*ast.Ident]*types.Package)
		for ident, object := range pkg.TypesInfo.Uses {
			if name, ok := object.(*types.PkgName); ok {
				refs[ident] = name.Imported()
			}
		}
		functions, err := extractFunctions(pkg)
		if err != nil {
			return nil, fmt.Errorf("extract functions from package %q: %w", pkg.PkgPath, err)
		}
		types, methods, err := extractTypes(pkg, aliases, refs)
		if err != nil {
			return nil, fmt.Errorf("extract types from package %q: %w", pkg.PkgPath, err)
		}
		functions = append(functions, methods...)
		slices.SortFunc(functions, func(a, b Function) int {
			if order := strings.Compare(a.Receiver, b.Receiver); order != 0 {
				return order
			}
			return strings.Compare(a.Name, b.Name)
		})
		constants, variables, err := extractValues(pkg)
		if err != nil {
			return nil, fmt.Errorf("extract values from package %q: %w", pkg.PkgPath, err)
		}
		extracted := Package{
			Path:      pkg.PkgPath,
			Name:      pkg.Name,
			FileSet:   pkg.Fset,
			Functions: functions,
			Types:     types,
			Constants: constants,
			Variables: variables,
		}
		collector := importCollector{
			source: pkg, refs: refs, aliases: aliases, packages: packageInfo,
			imports: make(map[Import]struct{}), qualifiers: make(map[string]string),
		}
		extracted.Imports, err = collector.collect(&extracted)
		if err != nil {
			return nil, fmt.Errorf("collect imports for package %q: %w", pkg.PkgPath, err)
		}
		if err := groupDeclarations(&extracted); err != nil {
			return nil, fmt.Errorf("group declarations for package %q: %w", pkg.PkgPath, err)
		}
		result = append(result, extracted)
	}
	slices.SortFunc(result, func(a, b Package) int {
		return strings.Compare(a.Path, b.Path)
	})
	return result, nil
}

func extractFunctions(pkg *packages.Package) ([]Function, error) {
	var functions []Function
	for _, file := range pkg.Syntax {
		for _, node := range file.Decls {
			decl, ok := node.(*ast.FuncDecl)
			if !ok || !decl.Name.IsExported() {
				continue
			}
			receiver, err := receiverName(decl.Recv)
			if err != nil {
				return nil, fmt.Errorf("method %q: %w", decl.Name.Name, err)
			}
			if receiver != "" && !token.IsExported(receiver) {
				continue
			}
			object, ok := pkg.TypesInfo.Defs[decl.Name].(*types.Func)
			if !ok {
				return nil, fmt.Errorf("function %q has no function type information", decl.Name.Name)
			}

			normalizer := declarationNormalizer{comments: []*ast.CommentGroup{}}
			function, err := normalizeFunction(decl, object, receiver, &normalizer)
			if err != nil {
				return nil, fmt.Errorf("function %q: %w", decl.Name.Name, err)
			}
			functions = append(functions, function)
		}
	}
	return functions, nil
}

func normalizeFunction(decl *ast.FuncDecl, object *types.Func, receiver string, normalizer *declarationNormalizer) (Function, error) {
	signature := *decl
	signature.Doc = nil
	signature.Body = nil
	var err error
	signature.Type, err = normalizer.signature(decl.Type)
	if err != nil {
		return Function{}, err
	}
	signature.Recv, err = normalizer.fields(decl.Recv, false)
	if err != nil {
		return Function{}, err
	}
	if signature.Recv != nil {
		for _, field := range signature.Recv.List {
			field.Names = nil
		}
	}
	return Function{
		Name:     decl.Name.Name,
		Object:   object,
		Receiver: receiver,
		Decl:     &signature,
		Comments: normalizer.comments,
	}, nil
}

func receiverName(recv *ast.FieldList) (string, error) {
	if recv == nil {
		return "", nil
	}
	if len(recv.List) != 1 {
		return "", fmt.Errorf("expected one receiver, got %d", len(recv.List))
	}
	return baseTypeName(recv.List[0].Type)
}

func baseTypeName(expr ast.Expr) (string, error) {
	for {
		switch receiver := expr.(type) {
		case *ast.Ident:
			return receiver.Name, nil
		case *ast.SelectorExpr:
			return receiver.Sel.Name, nil
		case *ast.StarExpr:
			expr = receiver.X
		case *ast.IndexExpr:
			expr = receiver.X
		case *ast.IndexListExpr:
			expr = receiver.X
		case *ast.ParenExpr:
			expr = receiver.X
		default:
			return "", fmt.Errorf("unsupported named type expression %T", expr)
		}
	}
}

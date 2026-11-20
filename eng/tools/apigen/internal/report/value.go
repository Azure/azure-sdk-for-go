// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Constant describes an exported constant with its evaluated type and value.
type Constant struct {
	// Name is the constant name.
	Name string

	// Type is the declared or inferred Go type, including untyped constant types.
	Type types.Type

	// Value is the evaluated constant value, independent of source expressions
	// such as iota or references to unexported constants.
	Value constant.Value

	// TypeExpr is the report spelling of Type, using public aliases and package
	// qualifiers from Imports. It is nil for untyped constants and is read-only.
	TypeExpr ast.Expr

	// Comments is a non-nil list of report-generated comments for TypeExpr.
	// Use it with printer.CommentedNode when formatting to suppress source comments.
	Comments []*ast.CommentGroup
}

// Variable describes a variable declaration and its exported bindings.
// Multiple bindings are kept together when they share a multi-result initializer.
type Variable struct {
	// Objects contains the exported bindings in declaration order, retaining
	// their compiler-provided types.
	Objects []*types.Var

	// Decl preserves explicit type annotations and initializers. Unexported names
	// in shared initializers are replaced by blank identifiers.
	// It shares expression nodes with loaded source and must be treated as read-only.
	Decl *ast.GenDecl

	// Comments is a non-nil list of report-generated comments for Decl.
	// Use it with printer.CommentedNode when formatting to suppress source comments.
	Comments []*ast.CommentGroup
}

func extractValues(pkg *packages.Package) ([]Constant, []Variable, error) {
	var constants []Constant
	var variables []Variable
	for _, file := range pkg.Syntax {
		for _, node := range file.Decls {
			decl, ok := node.(*ast.GenDecl)
			if !ok || (decl.Tok != token.CONST && decl.Tok != token.VAR) {
				continue
			}
			for _, node := range decl.Specs {
				spec, ok := node.(*ast.ValueSpec)
				if !ok {
					return nil, nil, fmt.Errorf("unexpected value specification %T", node)
				}
				if decl.Tok == token.VAR {
					values, err := extractVariables(pkg, spec)
					if err != nil {
						return nil, nil, err
					}
					variables = append(variables, values...)
					continue
				}
				for _, name := range spec.Names {
					if !name.IsExported() {
						continue
					}
					object, ok := pkg.TypesInfo.Defs[name].(*types.Const)
					if !ok {
						return nil, nil, fmt.Errorf("constant %q has no constant type information", name.Name)
					}
					constants = append(constants, Constant{
						Name:  name.Name,
						Type:  object.Type(),
						Value: object.Val(),
					})
				}
			}
		}
	}
	slices.SortFunc(constants, func(a, b Constant) int {
		return strings.Compare(a.Name, b.Name)
	})
	slices.SortFunc(variables, func(a, b Variable) int {
		return strings.Compare(a.Objects[0].Name(), b.Objects[0].Name())
	})
	return constants, variables, nil
}

func extractVariables(pkg *packages.Package, spec *ast.ValueSpec) ([]Variable, error) {
	shared := len(spec.Names) > 1 && len(spec.Values) == 1
	if !shared && len(spec.Values) != 0 && len(spec.Values) != len(spec.Names) {
		return nil, fmt.Errorf("variable declaration has %d names and %d initializers", len(spec.Names), len(spec.Values))
	}

	group := *spec
	if shared {
		group.Names = slices.Clone(spec.Names)
	}
	var result []Variable
	var objects []*types.Var
	for i, name := range spec.Names {
		if !name.IsExported() {
			if shared {
				blank := *name
				blank.Name = "_"
				group.Names[i] = &blank
			}
			continue
		}
		object, ok := pkg.TypesInfo.Defs[name].(*types.Var)
		if !ok {
			return nil, fmt.Errorf("variable %q has no variable type information", name.Name)
		}
		if shared {
			objects = append(objects, object)
			continue
		}

		copy := *spec
		copy.Names = []*ast.Ident{name}
		copy.Values = nil
		if len(spec.Values) > 0 {
			copy.Values = []ast.Expr{spec.Values[i]}
		}
		variable, err := variableDeclaration(copy, object)
		if err != nil {
			return nil, fmt.Errorf("variable %q: %w", name.Name, err)
		}
		result = append(result, variable)
	}
	if shared && len(objects) > 0 {
		variable, err := variableDeclaration(group, objects...)
		if err != nil {
			return nil, fmt.Errorf("variables starting with %q: %w", objects[0].Name(), err)
		}
		result = append(result, variable)
	}
	return result, nil
}

func variableDeclaration(spec ast.ValueSpec, objects ...*types.Var) (Variable, error) {
	spec.Doc = nil
	spec.Comment = nil
	normalizer := declarationNormalizer{comments: []*ast.CommentGroup{}}
	if spec.Type != nil {
		var err error
		spec.Type, err = normalizer.typeExpr(spec.Type)
		if err != nil {
			return Variable{}, err
		}
	}
	return Variable{
		Objects: objects,
		Decl: &ast.GenDecl{
			TokPos: spec.Pos(),
			Tok:    token.VAR,
			Specs:  []ast.Spec{&spec},
		},
		Comments: normalizer.comments,
	}, nil
}

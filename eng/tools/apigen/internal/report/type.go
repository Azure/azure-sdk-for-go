// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Type describes an exported named type declaration.
type Type struct {
	// Name is the type name.
	Name string

	// Object is the compiler identity of the public type declaration.
	Object *types.TypeName

	// IsInterface indicates an interface or an alias to an interface.
	IsInterface bool

	// AliasOf identifies an expanded internal alias target by package path and
	// type name, or is empty for declarations that were not expanded.
	AliasOf string

	// Decl contains a standalone type declaration with unexported struct fields
	// and descriptive comments removed. Public aliases retain their source targets.
	// It shares type-expression nodes with loaded source and must be treated as read-only.
	Decl *ast.GenDecl

	// Comments is a non-nil list of report-generated comments for Decl.
	// Use it with printer.CommentedNode when formatting to suppress source comments.
	Comments []*ast.CommentGroup

	// Constants contains this string or integer enum's exported constants,
	// sorted by name. A nonempty list identifies an enum group.
	Constants []Constant

	// Functions contains New or NewX constructors and Possible<Type>Values enum
	// helpers associated with this type, sorted by name.
	Functions []Function

	// Methods contains exported methods on this public receiver, sorted by name.
	Methods []Function
}

func extractTypes(pkg *packages.Package, aliases *aliasExpander, refs map[*ast.Ident]*types.Package) ([]Type, []Function, error) {
	var result []Type
	var methods []Function
	for _, file := range pkg.Syntax {
		for _, node := range file.Decls {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				continue
			}
			for _, node := range decl.Specs {
				spec, ok := node.(*ast.TypeSpec)
				if !ok {
					return nil, nil, fmt.Errorf("unexpected type specification %T", node)
				}
				if !spec.Name.IsExported() {
					continue
				}
				object, ok := pkg.TypesInfo.Defs[spec.Name].(*types.TypeName)
				if !ok {
					return nil, nil, fmt.Errorf("type %q has no type information", spec.Name.Name)
				}
				expansion, err := aliases.expansion(pkg, spec, object)
				if err != nil {
					return nil, nil, err
				}
				source := spec
				normalizer := declarationNormalizer{comments: []*ast.CommentGroup{}}
				aliasOf := ""
				if expansion != nil {
					source = expansion.source.spec
					normalizer.resolve = aliases.resolver(expansion.source.pkg, pkg, expansion, refs)
					aliasOf = expansion.origin
					functions, err := aliases.aliasMethods(expansion, pkg, refs)
					if err != nil {
						return nil, nil, err
					}
					methods = append(methods, functions...)
				}
				copy := *source
				copy.Name = &ast.Ident{NamePos: source.Name.Pos(), Name: spec.Name.Name}
				copy.Doc = nil
				copy.Comment = nil
				if expansion != nil {
					copy.Assign = token.NoPos
				}
				copy.TypeParams, err = normalizer.fields(source.TypeParams, false)
				if err != nil {
					return nil, nil, fmt.Errorf("type %q parameters: %w", spec.Name.Name, err)
				}
				copy.Type, err = normalizer.typeExpr(source.Type)
				if err != nil {
					return nil, nil, fmt.Errorf("type %q: %w", spec.Name.Name, err)
				}
				_, isInterface := object.Type().Underlying().(*types.Interface)
				result = append(result, Type{
					Name:        spec.Name.Name,
					Object:      object,
					IsInterface: isInterface,
					AliasOf:     aliasOf,
					Decl: &ast.GenDecl{
						TokPos: source.Pos(),
						Tok:    token.TYPE,
						Specs:  []ast.Spec{&copy},
					},
					Comments: normalizer.comments,
				})
			}
		}
	}
	slices.SortFunc(result, func(a, b Type) int {
		if a.IsInterface != b.IsInterface {
			if a.IsInterface {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Name, b.Name)
	})
	return result, methods, nil
}

type declarationNormalizer struct {
	comments []*ast.CommentGroup
	resolve  func(ast.Expr, string) (ast.Expr, error)
}

func (n *declarationNormalizer) typeExpr(expr ast.Expr) (ast.Expr, error) {
	return n.typeExprWithName(expr, "")
}

func (n *declarationNormalizer) typeExprWithName(expr ast.Expr, requiredName string) (ast.Expr, error) {
	switch source := expr.(type) {
	case *ast.Ident, *ast.SelectorExpr:
		if n.resolve != nil {
			return n.resolve(expr, requiredName)
		}
		return expr, nil
	case *ast.StructType:
		copy := *source
		var err error
		copy.Fields, err = n.fields(source.Fields, true)
		if err != nil {
			return nil, err
		}
		if err := sortFields(copy.Fields, true); err != nil {
			return nil, err
		}
		if source.Fields.NumFields() > copy.Fields.NumFields() {
			copy.Incomplete = true
		}
		if len(source.Fields.List) > 0 {
			// Keep the layout independent of source line breaks and member count.
			copy.Fields.Opening = token.NoPos
		}
		return &copy, nil
	case *ast.InterfaceType:
		copy := *source
		var err error
		copy.Methods, err = n.fields(source.Methods, false)
		if err != nil {
			return nil, err
		}
		if err := sortFields(copy.Methods, false); err != nil {
			return nil, err
		}
		if len(source.Methods.List) > 0 {
			copy.Methods.Opening = token.NoPos
		}
		return &copy, nil
	case *ast.FuncType:
		return n.signature(source)
	case *ast.ArrayType:
		copy := *source
		var err error
		copy.Elt, err = n.typeExpr(source.Elt)
		return &copy, err
	case *ast.MapType:
		copy := *source
		var err error
		copy.Key, err = n.typeExpr(source.Key)
		if err != nil {
			return nil, err
		}
		copy.Value, err = n.typeExpr(source.Value)
		return &copy, err
	case *ast.ChanType:
		copy := *source
		var err error
		copy.Value, err = n.typeExpr(source.Value)
		return &copy, err
	case *ast.StarExpr:
		copy := *source
		var err error
		copy.X, err = n.typeExprWithName(source.X, requiredName)
		return &copy, err
	case *ast.ParenExpr:
		copy := *source
		var err error
		copy.X, err = n.typeExprWithName(source.X, requiredName)
		return &copy, err
	case *ast.Ellipsis:
		copy := *source
		var err error
		copy.Elt, err = n.typeExpr(source.Elt)
		return &copy, err
	case *ast.IndexExpr:
		copy := *source
		var err error
		copy.X, err = n.typeExprWithName(source.X, requiredName)
		if err != nil {
			return nil, err
		}
		copy.Index, err = n.typeExpr(source.Index)
		return &copy, err
	case *ast.IndexListExpr:
		copy := *source
		var err error
		copy.X, err = n.typeExprWithName(source.X, requiredName)
		if err != nil {
			return nil, err
		}
		copy.Indices = make([]ast.Expr, len(source.Indices))
		for i, index := range source.Indices {
			copy.Indices[i], err = n.typeExpr(index)
			if err != nil {
				return nil, err
			}
		}
		return &copy, nil
	case *ast.UnaryExpr:
		copy := *source
		var err error
		copy.X, err = n.typeExpr(source.X)
		return &copy, err
	case *ast.BinaryExpr:
		copy := *source
		var err error
		copy.X, err = n.typeExpr(source.X)
		if err != nil {
			return nil, err
		}
		copy.Y, err = n.typeExpr(source.Y)
		return &copy, err
	default:
		return nil, fmt.Errorf("unsupported type expression %T", expr)
	}
}

func (n *declarationNormalizer) fields(source *ast.FieldList, publicOnly bool) (*ast.FieldList, error) {
	if source == nil {
		return nil, nil
	}
	copy := *source
	copy.List = nil
	for _, field := range source.List {
		filtered := *field
		filtered.Doc = nil
		filtered.Comment = nil
		requiredName := ""
		if publicOnly {
			if len(field.Names) == 0 {
				name, err := baseTypeName(field.Type)
				if err != nil {
					return nil, fmt.Errorf("embedded field: %w", err)
				}
				if !token.IsExported(name) {
					continue
				}
				requiredName = name
			} else {
				filtered.Names = nil
				for _, name := range field.Names {
					if name.IsExported() {
						filtered.Names = append(filtered.Names, name)
					}
				}
				if len(filtered.Names) == 0 {
					continue
				}
			}
		}
		var err error
		filtered.Type, err = n.typeExprWithName(field.Type, requiredName)
		if err != nil {
			return nil, err
		}
		if publicOnly && len(filtered.Names) > 1 {
			for _, name := range filtered.Names {
				single := filtered
				single.Names = []*ast.Ident{name}
				copy.List = append(copy.List, &single)
			}
		} else {
			copy.List = append(copy.List, &filtered)
		}
	}
	return &copy, nil
}

func sortFields(fields *ast.FieldList, structure bool) error {
	names := make(map[*ast.Field]string, len(fields.List))
	for _, field := range fields.List {
		if len(field.Names) > 0 {
			names[field] = field.Names[0].Name
		} else if structure {
			name, err := baseTypeName(field.Type)
			if err != nil {
				return fmt.Errorf("sort embedded field: %w", err)
			}
			names[field] = name
		} else {
			var text bytes.Buffer
			if err := format.Node(&text, token.NewFileSet(), field.Type); err != nil {
				return fmt.Errorf("sort interface member: %w", err)
			}
			names[field] = text.String()
		}
	}
	slices.SortStableFunc(fields.List, func(a, b *ast.Field) int {
		if structure && (len(a.Names) == 0) != (len(b.Names) == 0) {
			if len(a.Names) == 0 {
				return -1
			}
			return 1
		}
		return strings.Compare(names[a], names[b])
	})
	return nil
}

func (n *declarationNormalizer) signature(source *ast.FuncType) (*ast.FuncType, error) {
	copy := *source
	var err error
	copy.TypeParams, err = n.fields(source.TypeParams, false)
	if err != nil {
		return nil, err
	}
	copy.Params, err = n.fields(source.Params, false)
	if err != nil {
		return nil, err
	}
	copy.Results, err = n.fields(source.Results, false)
	return &copy, err
}

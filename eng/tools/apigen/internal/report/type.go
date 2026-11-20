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
	"maps"
	"slices"
	"strconv"
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

	// AliasOf identifies an expanded alias target by package path and
	// type name, or is empty for declarations that were not expanded.
	AliasOf string

	// Decl contains a standalone type declaration with unexported struct fields
	// and descriptive comments removed. Unexpanded aliases retain their source targets.
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

	// PromotedFields and PromotedMethods describe exported members reached
	// through omitted private embeddings, sorted by name.
	PromotedFields  []string
	PromotedMethods []PromotedMethod
}

// PromotedMethod describes a method promoted to the pointer method set and,
// unless PointerOnly is true, the value method set.
type PromotedMethod struct {
	Signature   string
	PointerOnly bool
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
				normalizer := declarationNormalizer{comments: []*ast.CommentGroup{}, typeInfo: pkg.TypesInfo}
				aliasOf := ""
				if expansion != nil {
					source = expansion.source.spec
					normalizer.typeInfo = expansion.source.pkg.TypesInfo
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
				var promotedFields []string
				var promotedMethods []PromotedMethod
				// Unexpanded aliases to named types already show the target name.
				if _, isStruct := ast.Unparen(source.Type).(*ast.StructType); !copy.Assign.IsValid() || isStruct {
					promotedFields, promotedMethods = promotedMembers(object.Type(), pkg.Types)
				}
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
					Comments:        normalizer.comments,
					PromotedFields:  promotedFields,
					PromotedMethods: promotedMethods,
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

func promotedMembers(receiver types.Type, pkg *types.Package) (fields []string, methods []PromotedMethod) {
	structure, ok := receiver.Underlying().(*types.Struct)
	if !ok {
		return nil, nil
	}
	privateEmbeddings := make(map[int]bool)
	fieldNames := make(map[string]bool)
	seen := make(map[*types.Struct]bool)
	var collectFieldNames func(types.Type)
	collectFieldNames = func(t types.Type) {
		if pointer, ok := types.Unalias(t).(*types.Pointer); ok {
			t = pointer.Elem()
		}
		embedded, ok := t.Underlying().(*types.Struct)
		if !ok || seen[embedded] {
			return
		}
		seen[embedded] = true
		for field := range embedded.Fields() {
			if field.Exported() {
				fieldNames[field.Name()] = true
			}
			if field.Embedded() {
				collectFieldNames(field.Type())
			}
		}
	}
	for i := range structure.NumFields() {
		field := structure.Field(i)
		if field.Embedded() && !field.Exported() {
			privateEmbeddings[i] = true
			collectFieldNames(field.Type())
		}
	}
	if len(privateEmbeddings) == 0 {
		return nil, nil
	}
	throughPrivateEmbedding := func(index []int) bool {
		return len(index) > 1 && privateEmbeddings[index[0]]
	}
	// Paths in report comments avoid ambiguous qualifiers and do not need imports.
	qualify := func(other *types.Package) string {
		if other == pkg {
			return ""
		}
		return other.Path()
	}
	for _, name := range slices.Sorted(maps.Keys(fieldNames)) {
		object, index, _ := types.LookupFieldOrMethod(receiver, true, nil, name)
		if field, ok := object.(*types.Var); ok && throughPrivateEmbedding(index) {
			fields = append(fields, name+" "+types.TypeString(field.Type(), qualify))
		}
	}
	valueMethods := types.NewMethodSet(receiver)
	for method := range types.NewMethodSet(types.NewPointer(receiver)).Methods() {
		if !method.Obj().Exported() || !throughPrivateEmbedding(method.Index()) {
			continue
		}
		name := method.Obj().Name()
		signature := name + strings.TrimPrefix(types.TypeString(method.Type(), qualify), "func")
		methods = append(methods, PromotedMethod{
			Signature:   signature,
			PointerOnly: valueMethods.Lookup(nil, name) == nil,
		})
	}
	return fields, methods
}

type declarationNormalizer struct {
	comments []*ast.CommentGroup
	typeInfo *types.Info
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
		if source.Len != nil {
			if n.typeInfo == nil {
				return nil, fmt.Errorf("array length has no type information")
			}
			array, ok := n.typeInfo.TypeOf(source).(*types.Array)
			if !ok {
				return nil, fmt.Errorf("array length has no array type information")
			}
			copy.Len = &ast.BasicLit{
				ValuePos: source.Len.Pos(),
				Kind:     token.INT,
				Value:    strconv.FormatInt(array.Len(), 10),
			}
		}
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

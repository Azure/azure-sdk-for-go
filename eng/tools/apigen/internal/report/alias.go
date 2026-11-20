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

type typeSource struct {
	pkg  *packages.Package
	spec *ast.TypeSpec
}

type methodSource struct {
	pkg  *packages.Package
	decl *ast.FuncDecl
}

type aliasExpander struct {
	types   map[*types.TypeName]typeSource
	methods map[*types.Func]methodSource
	public  []*types.TypeName
}

type aliasExpansion struct {
	source typeSource
	object *types.TypeName
	target types.Type
	origin string
}

func newAliasExpander(loaded *loader.Result) (*aliasExpander, error) {
	expander := &aliasExpander{
		types:   make(map[*types.TypeName]typeSource),
		methods: make(map[*types.Func]methodSource),
	}
	var sourceError error
	packages.Visit(loaded.Packages, nil, func(pkg *packages.Package) {
		if sourceError != nil {
			return
		}
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range decl.Specs {
						spec, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						object, ok := pkg.TypesInfo.Defs[spec.Name].(*types.TypeName)
						if !ok {
							sourceError = fmt.Errorf("type %q in package %q has no type information", spec.Name.Name, pkg.PkgPath)
							return
						}
						expander.types[object] = typeSource{pkg: pkg, spec: spec}
						if isPublicPackage(pkg) && object.Exported() && object.IsAlias() && !genericType(object.Type()) {
							expander.public = append(expander.public, object)
						}
					}
				case *ast.FuncDecl:
					if decl.Recv == nil || !decl.Name.IsExported() {
						continue
					}
					object, ok := pkg.TypesInfo.Defs[decl.Name].(*types.Func)
					if !ok {
						sourceError = fmt.Errorf("method %q in package %q has no type information", decl.Name.Name, pkg.PkgPath)
						return
					}
					expander.methods[object] = methodSource{pkg: pkg, decl: decl}
				}
			}
		}
	})
	if sourceError != nil {
		return nil, sourceError
	}
	slices.SortFunc(expander.public, func(a, b *types.TypeName) int {
		if order := strings.Compare(a.Pkg().Path(), b.Pkg().Path()); order != 0 {
			return order
		}
		return strings.Compare(a.Name(), b.Name())
	})
	return expander, nil
}

func (e *aliasExpander) expansion(pkg *packages.Package, spec *ast.TypeSpec, object *types.TypeName) (*aliasExpansion, error) {
	if !spec.Assign.IsValid() {
		return nil, nil
	}
	target := referencedTypeName(pkg, spec.Type)
	if target == nil || target.Pkg() == nil || !internalPackagePath(target.Pkg().Path()) {
		return nil, nil
	}
	if genericType(object.Type()) {
		return nil, fmt.Errorf("internal alias %q uses a generic alias or instantiated target, which is not supported", object.Name())
	}

	source, ok := e.types[target]
	if !ok {
		return nil, fmt.Errorf("source for internal alias target %s.%s is unavailable", target.Pkg().Path(), target.Name())
	}
	for source.spec.Assign.IsValid() {
		next := referencedTypeName(source.pkg, source.spec.Type)
		if next == nil || next.Pkg() == nil {
			break
		}
		nextSource, ok := e.types[next]
		if !ok {
			return nil, fmt.Errorf("source for alias target %s.%s is unavailable", next.Pkg().Path(), next.Name())
		}
		source = nextSource
	}
	return &aliasExpansion{
		source: source,
		object: object,
		target: types.Unalias(object.Type()),
		origin: target.Pkg().Path() + "." + target.Name(),
	}, nil
}

func (e *aliasExpander) aliasMethods(expansion *aliasExpansion, pkg *packages.Package, refs map[*ast.Ident]*types.Package) ([]Function, error) {
	named, ok := expansion.target.(*types.Named)
	if !ok {
		return nil, nil
	}
	if _, ok := named.Underlying().(*types.Interface); ok {
		return nil, nil
	}
	var result []Function
	for i := 0; i < named.NumMethods(); i++ {
		method := named.Method(i)
		if !method.Exported() {
			continue
		}
		source, ok := e.methods[method]
		if !ok {
			return nil, fmt.Errorf("source for method %q on internal alias %q is unavailable", method.Name(), expansion.object.Name())
		}
		normalizer := declarationNormalizer{
			comments: []*ast.CommentGroup{},
			resolve:  e.resolver(source.pkg, pkg, expansion, refs),
		}
		function, err := normalizeFunction(source.decl, method, expansion.object.Name(), &normalizer)
		if err != nil {
			return nil, fmt.Errorf("method %q on internal alias %q: %w", method.Name(), expansion.object.Name(), err)
		}
		result = append(result, function)
	}
	return result, nil
}

func (e *aliasExpander) resolver(source, current *packages.Package, expansion *aliasExpansion, refs map[*ast.Ident]*types.Package) func(ast.Expr, string) (ast.Expr, error) {
	return func(expr ast.Expr, requiredName string) (ast.Expr, error) {
		object := referencedTypeName(source, expr)
		if object == nil {
			return nil, fmt.Errorf("type information for %T in package %q is unavailable", expr, source.PkgPath)
		}
		if object.Pkg() == nil {
			return expr, nil
		}
		if _, ok := object.Type().(*types.TypeParam); ok {
			return expr, nil
		}
		return e.typeReference(object, current.Types, expansion, requiredName, expr.Pos(), refs), nil
	}
}

func (e *aliasExpander) typeReference(object *types.TypeName, current *types.Package, expansion *aliasExpansion, requiredName string, pos token.Pos, refs map[*ast.Ident]*types.Package) ast.Expr {
	reference := object
	switch {
	case expansion != nil && (requiredName == "" || expansion.object.Name() == requiredName) &&
		types.Identical(types.Unalias(object.Type()), expansion.target):
		reference = expansion.object
	case object.Pkg() != nil && internalPackagePath(object.Pkg().Path()):
		if alias := e.publicAlias(object.Type(), current, requiredName); alias != nil {
			reference = alias
		}
	}
	name := &ast.Ident{NamePos: pos, Name: reference.Name()}
	if reference.Pkg() == nil || reference.Pkg() == current {
		return name
	}
	qualifier := &ast.Ident{NamePos: pos, Name: reference.Pkg().Name()}
	refs[qualifier] = reference.Pkg()
	return &ast.SelectorExpr{X: qualifier, Sel: name}
}

func (e *aliasExpander) publicAlias(target types.Type, current *types.Package, requiredName string) *types.TypeName {
	var first *types.TypeName
	for _, alias := range e.public {
		if requiredName != "" && alias.Name() != requiredName {
			continue
		}
		if !types.Identical(types.Unalias(alias.Type()), types.Unalias(target)) {
			continue
		}
		if alias.Pkg() == current {
			return alias
		}
		if first == nil {
			first = alias
		}
	}
	return first
}

func referencedTypeName(pkg *packages.Package, expr ast.Expr) *types.TypeName {
	switch expr := expr.(type) {
	case *ast.Ident:
		object, _ := pkg.TypesInfo.Uses[expr].(*types.TypeName)
		return object
	case *ast.SelectorExpr:
		object, _ := pkg.TypesInfo.Uses[expr.Sel].(*types.TypeName)
		return object
	case *ast.ParenExpr:
		return referencedTypeName(pkg, expr.X)
	case *ast.IndexExpr:
		return referencedTypeName(pkg, expr.X)
	case *ast.IndexListExpr:
		return referencedTypeName(pkg, expr.X)
	default:
		return nil
	}
}

func genericType(t types.Type) bool {
	for {
		switch typed := t.(type) {
		case *types.Alias:
			if typed.TypeParams().Len() > 0 || typed.TypeArgs().Len() > 0 {
				return true
			}
			t = typed.Rhs()
		case *types.Named:
			return typed.TypeParams().Len() > 0 || typed.TypeArgs().Len() > 0
		default:
			return false
		}
	}
}

func internalPackagePath(path string) bool {
	return slices.Contains(strings.Split(path, "/"), "internal")
}

func isPublicPackage(pkg *packages.Package) bool {
	return pkg.Name != "main" && !internalPackagePath(pkg.PkgPath)
}

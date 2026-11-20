// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Import describes a package referenced by report declarations.
type Import struct {
	// Path is the package's import path.
	Path string

	// Name is an explicit import alias, or empty for the package's declared name.
	// A dot import uses ".".
	Name string

	// Standard indicates whether the package belongs to the standard library.
	Standard bool
}

type importCollector struct {
	source     *packages.Package
	refs       map[*ast.Ident]*types.Package // synthesized qualifiers
	aliases    *aliasExpander
	packages   map[*types.Package]*packages.Package
	imports    map[Import]struct{}
	qualifiers map[string]string
	pending    []*ast.Ident
}

func (c *importCollector) collect(pkg *Package) ([]Import, error) {
	for _, function := range pkg.Functions {
		if err := c.node(function.Decl); err != nil {
			return nil, err
		}
	}
	for _, typ := range pkg.Types {
		if err := c.node(typ.Decl); err != nil {
			return nil, err
		}
	}
	for _, variable := range pkg.Variables {
		if err := c.node(variable.Decl); err != nil {
			return nil, err
		}
	}
	// Preserve source qualifiers and allocate synthesized ones in import-path
	// order, independent of the order in which declarations were collected.
	slices.SortFunc(c.pending, func(a, b *ast.Ident) int {
		return strings.Compare(c.refs[a].Path(), c.refs[b].Path())
	})
	for _, ident := range c.pending {
		name, err := c.allocateQualifier(c.refs[ident])
		if err != nil {
			return nil, err
		}
		ident.Name = name
	}
	for i := range pkg.Constants {
		value := &pkg.Constants[i]
		var err error
		value.TypeExpr, value.Comments, err = c.constantType(value)
		if err != nil {
			return nil, fmt.Errorf("constant %q: %w", value.Name, err)
		}
	}
	for i := range pkg.Variables {
		variable := &pkg.Variables[i]
		if err := c.completeVariableType(variable); err != nil {
			return nil, fmt.Errorf("variable %q type: %w", variable.Objects[0].Name(), err)
		}
	}
	var result []Import
	for entry := range c.imports {
		result = append(result, entry)
	}
	slices.SortFunc(result, func(a, b Import) int {
		if a.Standard != b.Standard {
			if a.Standard {
				return -1
			}
			return 1
		}
		if order := strings.Compare(a.Path, b.Path); order != 0 {
			return order
		}
		return strings.Compare(a.Name, b.Name)
	})
	return result, nil
}

func (c *importCollector) node(node ast.Node) error {
	selected := make(map[*ast.Ident]bool)
	var collectError error
	ast.Inspect(node, func(node ast.Node) bool {
		if collectError != nil {
			return false
		}
		switch node := node.(type) {
		case *ast.SelectorExpr:
			selected[node.Sel] = true
		case *ast.Ident:
			if c.refs[node] != nil {
				c.pending = append(c.pending, node)
			} else if !selected[node] {
				object := c.source.TypesInfo.Uses[node]
				if imported, ok := object.(*types.PkgName); ok {
					collectError = c.add(imported.Imported(), node.Name)
				} else if object != nil && object.Pkg() != nil && object.Pkg() != c.source.Types && object.Parent() == object.Pkg().Scope() {
					collectError = c.add(object.Pkg(), ".")
				}
			}
		}
		return collectError == nil
	})
	return collectError
}

func (c *importCollector) add(pkg *types.Package, qualifier string) error {
	if qualifier != "." {
		if previous, ok := c.qualifiers[qualifier]; ok && previous != pkg.Path() {
			return fmt.Errorf("import qualifier %q refers to both %q and %q", qualifier, previous, pkg.Path())
		}
		c.qualifiers[qualifier] = pkg.Path()
	}
	info := c.packages[pkg]
	if info == nil {
		return fmt.Errorf("package metadata for import %q is unavailable", pkg.Path())
	}
	name := qualifier
	if name == pkg.Name() {
		name = ""
	}
	c.imports[Import{Path: pkg.Path(), Name: name, Standard: info.Module == nil}] = struct{}{}
	return nil
}

func (c *importCollector) preferredQualifier(pkg *types.Package) string {
	var first string
	for entry := range c.imports {
		if entry.Path != pkg.Path() || entry.Name == "." {
			continue
		}
		if entry.Name == "" {
			return pkg.Name()
		}
		if first == "" || entry.Name < first {
			first = entry.Name
		}
	}
	if first != "" {
		return first
	}
	return pkg.Name()
}

func (c *importCollector) allocateQualifier(pkg *types.Package) (string, error) {
	base := c.preferredQualifier(pkg)
	name := base
	for suffix := 2; c.qualifiers[name] != "" && c.qualifiers[name] != pkg.Path(); suffix++ {
		name = base + strconv.Itoa(suffix)
	}
	return name, c.add(pkg, name)
}

func (c *importCollector) completeVariableType(variable *Variable) error {
	spec := variable.Decl.Specs[0].(*ast.ValueSpec)
	if spec.Type != nil || len(spec.Values) != 0 {
		return nil
	}
	object := variable.Objects[0]
	expr, qualified, err := parseCompilerType(c.source, "variable-"+object.Name()+".api.go", object.Type())
	if err != nil {
		return err
	}
	normalizer := declarationNormalizer{
		comments: []*ast.CommentGroup{},
		resolve: func(expr ast.Expr, _ string) (ast.Expr, error) {
			selected, ok := expr.(*ast.SelectorExpr)
			if !ok {
				return expr, nil
			}
			qualifier, ok := selected.X.(*ast.Ident)
			if !ok || qualified[qualifier.Name] == nil {
				return nil, fmt.Errorf("unknown variable type qualifier in %T", expr)
			}
			name, err := c.allocateQualifier(qualified[qualifier.Name])
			if err != nil {
				return nil, err
			}
			qualifier.Name = name
			return selected, nil
		},
	}
	spec.Type, err = normalizer.typeExpr(expr)
	if err != nil {
		return err
	}
	// The type came from a synthetic file, so don't associate its positions
	// with the source position of the binding when printing the declaration.
	variable.Decl.TokPos = token.NoPos
	spec.Names[0] = ast.NewIdent(object.Name())
	variable.Comments = normalizer.comments
	return nil
}

func (c *importCollector) constantType(value *Constant) (ast.Expr, []*ast.CommentGroup, error) {
	normalizer := declarationNormalizer{comments: []*ast.CommentGroup{}}
	if basic, ok := value.Type.(*types.Basic); ok && basic.Info()&types.IsUntyped != 0 {
		return nil, normalizer.comments, nil
	}
	expr, qualified, err := parseCompilerType(c.source, "constant-"+value.Name+".api.go", value.Type)
	if err != nil {
		return nil, nil, fmt.Errorf("parse constant type: %w", err)
	}
	normalizer.resolve = func(expr ast.Expr, requiredName string) (ast.Expr, error) {
		var object types.Object
		switch expr := expr.(type) {
		case *ast.Ident:
			object = c.source.Types.Scope().Lookup(expr.Name)
			if object == nil {
				object = types.Universe.Lookup(expr.Name)
			}
		case *ast.SelectorExpr:
			qualifier, ok := expr.X.(*ast.Ident)
			if !ok || qualified[qualifier.Name] == nil {
				return nil, fmt.Errorf("unknown constant type qualifier in %T", expr)
			}
			object = qualified[qualifier.Name].Scope().Lookup(expr.Sel.Name)
		}
		name, ok := object.(*types.TypeName)
		if !ok {
			return nil, fmt.Errorf("constant type information for %T is unavailable", expr)
		}
		reference := c.aliases.typeReference(name, c.source.Types, nil, requiredName, expr.Pos(), c.refs)
		if selector, ok := reference.(*ast.SelectorExpr); ok {
			qualifier, ok := selector.X.(*ast.Ident)
			if !ok || c.refs[qualifier] == nil {
				return nil, fmt.Errorf("unexpected constant type qualifier %T", selector.X)
			}
			name, err := c.allocateQualifier(c.refs[qualifier])
			if err != nil {
				return nil, err
			}
			qualifier.Name = name
		}
		return reference, nil
	}
	expr, err = normalizer.typeExpr(expr)
	return expr, normalizer.comments, err
}

// parseCompilerType uses distinct placeholders so packages with the same name
// remain distinguishable until type references are normalized.
func parseCompilerType(source *packages.Package, filename string, typ types.Type) (ast.Expr, map[string]*types.Package, error) {
	qualifiers := make(map[*types.Package]string)
	qualified := make(map[string]*types.Package)
	text := types.TypeString(typ, func(pkg *types.Package) string {
		if pkg == source.Types {
			return ""
		}
		name, ok := qualifiers[pkg]
		if !ok {
			name = "__apigen_pkg_" + strconv.Itoa(len(qualifiers))
			qualifiers[pkg] = name
			qualified[name] = pkg
		}
		return name
	})
	expr, err := parser.ParseExprFrom(source.Fset, filename, text, 0)
	return expr, qualified, err
}

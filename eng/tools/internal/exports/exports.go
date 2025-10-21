// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exports

import (
	"fmt"
	"go/ast"
	"strings"
)

// Content defines the set of exported constants, funcs, and structs.
type Content struct {
	// the list of exported constants.
	// key is the exported name, value is its type and value.
	Consts map[string]Const `json:"consts,omitempty"`

	// the list of exported type aliases.
	// key is the exported name, value contains underlaying type.
	TypeAliases map[string]TypeAlias `json:"typeAliases,omitempty"`

	// the list of exported functions and methods.
	// key is the exported name, for methods it's prefixed with the receiver type (e.g. "Type.Method").
	// value contains the list of params and return types.
	Funcs map[string]Func `json:"funcs,omitempty"`

	// the list of exported interfaces.
	// key is the exported name, value contains the interface definition.
	Interfaces map[string]Interface `json:"interfaces,omitempty"`

	// the list of exported struct types.
	// key is the exported name, value contains field information.
	Structs map[string]Struct `json:"structs,omitempty"`
}

// Count returns the count of items
func (c Content) Count() int {
	return len(c.Consts) + len(c.TypeAliases) + len(c.Funcs) + len(c.Interfaces) + len(c.Structs)
}

// Const is a const definition.
type Const struct {
	// the type of the constant
	Type string `json:"type"`

	// the value of the constant
	Value string `json:"value"`
}

// Param represents a function parameter with its name and type.
type Param struct {
	// Name is the parameter name, may be empty for unnamed parameters
	Name string `json:"name,omitempty"`
	
	// Type is the parameter type
	Type string `json:"type"`
}

// Func contains parameter and return types of a function/method.
type Func struct {
	// Params is the list of function parameters
	Params []Param `json:"params,omitempty"`

	// a comma-delimited list of the return types
	Returns *string `json:"returns,omitempty"`

	// func name that replace this func with breaking change
	ReplacedBy *string `json:"replacedby,omitempty"`
}

// Interface contains the list of methods for an interface.
type Interface struct {
	// a list of embedded interfaces
	AnonymousFields []string `json:"anon,omitempty"`

	// key/value pairs of the methd names and their definitions
	Methods map[string]Func
}

// Struct contains field info about a struct.
type Struct struct {
	// a list of anonymous fields
	AnonymousFields []string `json:"anon,omitempty"`

	// key/value pairs of the field names and types respectively.
	Fields map[string]string `json:"fields,omitempty"`
}

// TypeAlias contains field info about a type.
type TypeAlias struct {
	// underlaying type
	UnderlayingType string `json:"underlayingType,omitempty"`
}

// NewContent returns an initialized Content object.
func NewContent() Content {
	return Content{
		Consts:      make(map[string]Const),
		Funcs:       make(map[string]Func),
		Interfaces:  make(map[string]Interface),
		Structs:     make(map[string]Struct),
		TypeAliases: make(map[string]TypeAlias),
	}
}

// IsEmpty returns true if there is no content in any of the fields.
func (c Content) IsEmpty() bool {
	return len(c.Consts) == 0 && len(c.TypeAliases) == 0 && len(c.Funcs) == 0 && len(c.Interfaces) == 0 && len(c.Structs) == 0
}

// adds the specified const declaration to the exports list
func (c *Content) addConst(pkg Package, g *ast.GenDecl) {
	for _, s := range g.Specs {
		co := Const{}
		vs := s.(*ast.ValueSpec)
		v := ""
		// Type is nil for untyped consts
		if vs.Type != nil {
			switch x := vs.Type.(type) {
			case *ast.Ident:
				co.Type = x.Name
				switch vs.Values[0].(type) {
				case *ast.Ident:
					v = vs.Values[0].(*ast.Ident).Name
				case *ast.BasicLit:
					v = vs.Values[0].(*ast.BasicLit).Value
				default:
					panic(fmt.Sprintf("wrong type %T", vs.Values[0]))
				}
			case *ast.SelectorExpr:
				co.Type = x.Sel.Name
				v = vs.Values[0].(*ast.BasicLit).Value
			default:
				panic(fmt.Sprintf("wrong type %T", vs.Type))
			}
		} else {
			// get the type from the token type
			if bl, ok := vs.Values[0].(*ast.BasicLit); ok {
				co.Type = strings.ToLower(bl.Kind.String())
				v = bl.Value
			} else if ce, ok := vs.Values[0].(*ast.CallExpr); ok {
				// const FooConst = FooType("value")
				co.Type = pkg.getText(ce.Fun.Pos(), ce.Fun.End())
				v = pkg.getText(ce.Args[0].Pos(), ce.Args[0].End())
			} else if ce, ok := vs.Values[0].(*ast.BinaryExpr); ok {
				// const FooConst = "value" + Bar
				co.Type = "*ast.BinaryExpr"
				v = pkg.getText(ce.X.Pos(), ce.Y.End())
			} else {
				panic(fmt.Sprintf("unhandled case for adding constant: %s", pkg.getText(vs.Pos(), vs.End())))
			}
		}
		// TODO should this also be removed?
		// remove any surrounding quotes
		if v[0] == '"' {
			v = v[1 : len(v)-1]
		}
		co.Value = v
		c.Consts[vs.Names[0].Name] = co
	}
}

// adds the specified function declaration to the exports list
func (c *Content) addFunc(pkg Package, f *ast.FuncDecl) {
	// create a method sig, for methods it's a combination of the receiver type
	// with the function name e.g. "FooReceiver.Method", else just the function name.
	sig := ""
	if f.Recv != nil {
		sig = pkg.getText(f.Recv.List[0].Type.Pos(), f.Recv.List[0].Type.End())
		sig += "."
	}
	sig += f.Name.Name
	c.Funcs[sig] = pkg.buildFunc(f.Type)
}

// adds the specified interface type to the exports list.
func (c *Content) addInterface(pkg Package, name string, i *ast.InterfaceType) {
	in := Interface{Methods: map[string]Func{}}
	if i.Methods != nil {
		pkg.translateFieldList(i.Methods.List, func(n *string, t string, f *ast.Field) {
			if n == nil {
				in.AnonymousFields = append(in.AnonymousFields, t)
			} else {
				if in.Methods == nil {
					in.Methods = map[string]Func{}
				}
				in.Methods[*n] = pkg.buildFunc(f.Type.(*ast.FuncType))
			}
		})
	}
	c.Interfaces[name] = in
}

// adds the specified struct type to the exports list.
func (c *Content) addStruct(pkg Package, name string, s *ast.StructType) {
	sd := Struct{}
	// assumes all struct types have fields
	pkg.translateFieldList(s.Fields.List, func(n *string, t string, f *ast.Field) {
		if n == nil {
			sd.AnonymousFields = append(sd.AnonymousFields, t)
		} else {
			if sd.Fields == nil {
				sd.Fields = map[string]string{}
			}
			sd.Fields[*n] = t
		}
	})
	c.Structs[name] = sd
}

// adds the specified simple type to the exports list.
func (c *Content) addTypeAlias(name string, underlayingType string) {
	c.TypeAliases[name] = TypeAlias{UnderlayingType: underlayingType}
}

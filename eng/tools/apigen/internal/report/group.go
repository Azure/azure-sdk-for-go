// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

func groupDeclarations(pkg *Package) error {
	associations := typeAssociations{
		byObject: make(map[*types.TypeName]int),
	}
	byName := make(map[string]int)
	for i, typ := range pkg.Types {
		associations.byObject[typ.Object] = i
		byName[typ.Name] = i
	}

	var constants []Constant
	for _, constant := range pkg.Constants {
		if name, ok := constant.TypeExpr.(*ast.Ident); ok {
			if index, ok := byName[name.Name]; ok {
				if basic, ok := pkg.Types[index].Object.Type().Underlying().(*types.Basic); ok &&
					basic.Info()&(types.IsString|types.IsInteger) != 0 {
					pkg.Types[index].Constants = append(pkg.Types[index].Constants, constant)
					continue
				}
			}
		}
		constants = append(constants, constant)
	}
	pkg.Constants = constants

	var functions []Function
	for _, function := range pkg.Functions {
		if function.Receiver != "" {
			index, ok := byName[function.Receiver]
			if !ok {
				return fmt.Errorf("method %q has no public receiver declaration %q", function.Name, function.Receiver)
			}
			pkg.Types[index].Methods = append(pkg.Types[index].Methods, function)
			continue
		}
		signature, ok := function.Object.Type().(*types.Signature)
		if !ok {
			return fmt.Errorf("function %q has no signature", function.Name)
		}
		suffix, constructor := strings.CutPrefix(function.Name, "New")
		constructor = constructor && (suffix == "" || token.IsExported(suffix))
		index, matches := -1, 0
		for i := 0; i < signature.Results().Len(); i++ {
			if match := associations.match(signature.Results().At(i).Type()); match >= 0 {
				index = match
				matches++
			}
		}
		if matches == 1 && (constructor || isEnumValuesHelper(function, signature, pkg.Types[index])) {
			pkg.Types[index].Functions = append(pkg.Types[index].Functions, function)
		} else {
			functions = append(functions, function)
		}
	}
	pkg.Functions = functions

	return nil
}

func isEnumValuesHelper(function Function, signature *types.Signature, typ Type) bool {
	if len(typ.Constants) == 0 || function.Name != "Possible"+typ.Name+"Values" ||
		signature.Params().Len() != 0 || signature.Results().Len() != 1 {
		return false
	}
	values, ok := signature.Results().At(0).Type().(*types.Slice)
	return ok && types.Identical(values.Elem(), typ.Object.Type())
}

type typeAssociations struct {
	byObject map[*types.TypeName]int
}

func (a typeAssociations) match(t types.Type) int {
	for {
		switch typed := t.(type) {
		case *types.Alias:
			if index, ok := a.byObject[typed.Origin().Obj()]; ok {
				return index
			}
			return -1
		case *types.Named:
			if index, ok := a.byObject[typed.Origin().Obj()]; ok {
				return index
			}
			return -1
		case *types.Pointer:
			t = typed.Elem()
		case *types.Slice:
			t = typed.Elem()
		case *types.Array:
			t = typed.Elem()
		default:
			return -1
		}
	}
}

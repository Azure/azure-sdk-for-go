// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package report

import (
	"go/types"
	"maps"
	"slices"
	"strings"
)

func typeParameterNotes(owner string, params *types.TypeParamList, pkg *types.Package) []string {
	var notes []string
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		if terms := hiddenConstraintTerms(param.Constraint(), pkg); terms != "" {
			notes = append(notes, owner+"."+param.Obj().Name()+": "+terms)
		}
	}
	return notes
}

func namedTypeParameters(t types.Type) *types.TypeParamList {
	switch t := t.(type) {
	case *types.Named:
		return t.TypeParams()
	case *types.Alias:
		return t.TypeParams()
	default:
		return nil
	}
}

func hiddenConstraintTerms(constraint types.Type, pkg *types.Package) string {
	qualify := func(other *types.Package) string {
		if other == pkg {
			return ""
		}
		return other.Path()
	}
	var visit func(types.Type) (expression string, hidden bool)
	visit = func(t types.Type) (string, bool) {
		expression := func() string { return TypeStringWithoutTags(t, qualify) }
		named, isNamed := t.(interface{ Obj() *types.TypeName })
		isHidden := isNamed && hiddenTypeTarget(pkg, named.Obj())
		if types.Unalias(t) == types.Universe.Lookup("comparable").Type() {
			return "comparable", isHidden
		}
		switch underlying := t.Underlying().(type) {
		case *types.Interface:
			if underlying.IsMethodSet() {
				return "", false
			}
			if isNamed && !isHidden {
				return expression(), false
			}
			terms := make(map[string]bool)
			for embedded := range underlying.EmbeddedTypes() {
				text, hidden := visit(embedded)
				if text != "" {
					terms[text] = true
					isHidden = isHidden || hidden
				}
			}
			sorted := slices.Sorted(maps.Keys(terms))
			switch len(sorted) {
			case 0:
				return "", false
			case 1:
				return sorted[0], isHidden
			default:
				// Separate interface elements intersect; union alternatives do not.
				return "interface{" + strings.Join(sorted, "; ") + "}", isHidden
			}
		case *types.Union:
			var alternatives []string
			for i := 0; i < underlying.Len(); i++ {
				term := underlying.Term(i)
				text, hidden := visit(term.Type())
				isHidden = isHidden || hidden
				if text == "" {
					text = "any"
				}
				if term.Tilde() {
					text = "~" + text
				}
				alternatives = append(alternatives, text)
			}
			slices.Sort(alternatives)
			return strings.Join(alternatives, " | "), isHidden
		default:
			if isHidden {
				return TypeStringWithoutTags(types.Unalias(t), qualify), named.Obj().IsAlias()
			}
			return expression(), false
		}
	}
	if terms, hidden := visit(constraint); hidden {
		return terms
	}
	return ""
}

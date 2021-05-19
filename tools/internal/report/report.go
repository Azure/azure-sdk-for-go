// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/markdown"
)

// Package represents a per-package report that contains additive and breaking changes.
type Package struct {
	AdditiveChanges *delta.Content   `json:"additiveChanges,omitempty"`
	BreakingChanges *BreakingChanges `json:"breakingChanges,omitempty"`
}

// HasBreakingChanges returns true if the package report contains breaking changes.
func (p Package) HasBreakingChanges() bool {
	return p.BreakingChanges != nil && !p.BreakingChanges.IsEmpty()
}

// HasAdditiveChanges returns true if the package report contains additive changes.
func (p Package) HasAdditiveChanges() bool {
	return p.AdditiveChanges != nil && !p.AdditiveChanges.IsEmpty()
}

// IsEmpty returns true if the report contains no data (e.g. no changes in exported types).
func (p Package) IsEmpty() bool {
	return (p.AdditiveChanges == nil || p.AdditiveChanges.IsEmpty()) &&
		(p.BreakingChanges == nil || p.BreakingChanges.IsEmpty())
}

// BreakingChanges represents a set of breaking changes.
type BreakingChanges struct {
	Consts     map[string]delta.Signature    `json:"consts,omitempty"`
	Funcs      map[string]delta.FuncSig      `json:"funcs,omitempty"`
	Interfaces map[string]delta.InterfaceDef `json:"interfaces,omitempty"`
	Structs    map[string]delta.StructDef    `json:"structs,omitempty"`
	Removed    *delta.Content                `json:"removed,omitempty"`
}

// Count returns the count of breaking changes
func (bc BreakingChanges) Count() int {
	removed := 0
	if bc.Removed != nil {
		removed = bc.Removed.Count()
	}
	return len(bc.Consts) + len(bc.Funcs) + len(bc.Interfaces) + len(bc.Structs) + removed
}

// IsEmpty returns true if there are no breaking changes.
func (bc BreakingChanges) IsEmpty() bool {
	return len(bc.Consts) == 0 && len(bc.Funcs) == 0 && len(bc.Interfaces) == 0 && len(bc.Structs) == 0 &&
		(bc.Removed == nil || bc.Removed.IsEmpty())
}

// GenerationOption ...
type GenerationOption struct {
	// OnlyBreakingChanges ...
	OnlyBreakingChanges bool
	// OnlyAdditiveChanges ...
	OnlyAdditiveChanges bool
}

// Generate generates a package report based on the delta between lhs and rhs.
// onlyBreakingChanges - pass true to include only breaking changes in the report.
// onlyAdditions - pass true to include only addition changes in the report.
func Generate(lhs, rhs exports.Content, option *GenerationOption) Package {
	onlyBreakingChanges := option != nil && option.OnlyAdditiveChanges
	onlyAdditiveChanges := option != nil && option.OnlyAdditiveChanges
	r := Package{}
	if !onlyBreakingChanges {
		if adds := delta.GetExports(lhs, rhs); !adds.IsEmpty() {
			r.AdditiveChanges = &adds
		}
	}

	if !onlyAdditiveChanges {
		breaks := BreakingChanges{}
		breaks.Consts = delta.GetConstTypeChanges(lhs, rhs)
		breaks.Funcs = delta.GetFuncSigChanges(lhs, rhs)
		breaks.Interfaces = delta.GetInterfaceMethodSigChanges(lhs, rhs)
		breaks.Structs = delta.GetStructFieldChanges(lhs, rhs)
		if removed := delta.GetExports(rhs, lhs); !removed.IsEmpty() {
			breaks.Removed = &removed
		}
		if !breaks.IsEmpty() {
			r.BreakingChanges = &breaks
		}
	}
	return r
}

// ToMarkdown creates a report of the package changes in markdown format.
func (p Package) ToMarkdown() string {
	if p.IsEmpty() {
		return ""
	}
	md := markdown.Writer{}
	p.writeBreakingChanges(&md)
	p.writeNewContent(&md)
	return md.String()
}

func (p Package) writeBreakingChanges(md *markdown.Writer) {
	if !p.HasBreakingChanges() {
		return
	}
	md.WriteTopLevelHeader("Breaking Changes")
	writeRemovedContent(p.BreakingChanges.Removed, md)
	writeSigChanges(p.BreakingChanges, md)
}

func (p Package) writeNewContent(md *markdown.Writer) {
	if !p.HasAdditiveChanges() {
		return
	}
	md.WriteTopLevelHeader("Additive Changes")
	writeConsts(p.AdditiveChanges.Consts, "New Constants", md)
	writeFuncs(p.AdditiveChanges.Funcs, "New Funcs", md)
	writeStructs(p.AdditiveChanges, "New Structs", "New Struct Fields", md)
}

// writes the subset of breaking changes pertaining to removed content
func writeRemovedContent(removed *delta.Content, md *markdown.Writer) {
	if removed == nil {
		return
	}
	writeConsts(removed.Consts, "Removed Constants", md)
	writeFuncs(removed.Funcs, "Removed Funcs", md)
	writeStructs(removed, "Removed Structs", "Removed Struct Fields", md)
}

// writes the subset of breaking changes pertaining to signature changes
func writeSigChanges(bc *BreakingChanges, md *markdown.Writer) {
	if len(bc.Consts) == 0 && len(bc.Funcs) == 0 && len(bc.Structs) == 0 {
		return
	}
	md.WriteHeader("Signature Changes")
	if len(bc.Consts) > 0 {
		items := make([]string, len(bc.Consts))
		i := 0
		for k, v := range bc.Consts {
			items[i] = fmt.Sprintf("1. %s changed type from %s to %s", k, v.From, v.To)
			i++
		}
		sort.Strings(items)
		md.WriteSubheader("Const Types")
		for _, item := range items {
			md.WriteLine(item)
		}
	}
	if len(bc.Funcs) > 0 {
		// first get all the funcs so we can sort them
		items := make([]string, len(bc.Funcs))
		i := 0
		for k := range bc.Funcs {
			items[i] = k
			i++
		}
		sort.Strings(items)
		md.WriteSubheader("Funcs")
		for _, item := range items {
			// now add params/returns info
			changes := bc.Funcs[item]
			if changes.Params != nil {
				item = fmt.Sprintf("%s\n\t- Params\n\t\t- From: %s\n\t\t- To: %s", item, changes.Params.From, changes.Params.To)
			}
			if changes.Returns != nil {
				item = fmt.Sprintf("%s\n\t- Returns\n\t\t- From: %s\n\t\t- To: %s", item, changes.Returns.From, changes.Returns.To)
			}
			md.WriteLine(fmt.Sprintf("1. %s", item))
		}
	}
	if len(bc.Structs) > 0 {
		items := make([]string, 0, len(bc.Structs))
		for k, v := range bc.Structs {
			for f, d := range v.Fields {
				items = append(items, fmt.Sprintf("1. %s.%s changed type from %s to %s", k, f, d.From, d.To))
			}
		}
		sort.Strings(items)
		md.WriteSubheader("Struct Fields")
		for _, item := range items {
			md.WriteLine(item)
		}
	}
}

// writes out const information formatted as TypeName.ConstName
func writeConsts(co map[string]exports.Const, subheader string, md *markdown.Writer) {
	if len(co) == 0 {
		return
	}
	items := make([]string, len(co))
	i := 0
	for c, t := range co {
		items[i] = fmt.Sprintf("1. %s.%s", t.Type, c)
		i++
	}
	sort.Strings(items)
	md.WriteHeader(subheader)
	for _, item := range items {
		md.WriteLine(item)
	}
}

// writes out func information formatted as [receiver].FuncName([params]) [returns]
func writeFuncs(funcs map[string]exports.Func, subheader string, md *markdown.Writer) {
	if len(funcs) == 0 {
		return
	}
	items := make([]string, len(funcs))
	i := 0
	for k, v := range funcs {
		params := ""
		if v.Params != nil {
			params = *v.Params
		}
		returns := ""
		if v.Returns != nil {
			returns = *v.Returns
			if strings.Index(returns, ",") > -1 {
				returns = fmt.Sprintf("(%s)", returns)
			}
		}
		items[i] = fmt.Sprintf("1. %s(%s) %s", k, params, returns)
		i++
	}
	sort.Strings(items)
	md.WriteHeader(subheader)
	for _, item := range items {
		md.WriteLine(item)
	}
}

// writes out struct information
// sheader1 is for added/removed struct types formatted as TypeName
// sheader2 is for added/removed struct fields formatted as TypeName.FieldName
func writeStructs(content *delta.Content, sheader1, sheader2 string, md *markdown.Writer) {
	if len(content.Structs) == 0 {
		return
	}
	md.WriteHeader("Struct Changes")
	if len(content.CompleteStructs) > 0 {
		md.WriteSubheader(sheader1)
		for _, s := range content.CompleteStructs {
			md.WriteLine(fmt.Sprintf("1. %s", s))
		}
	}
	modified := content.GetModifiedStructs()
	if len(modified) > 0 {
		md.WriteSubheader(sheader2)
		items := make([]string, 0, len(content.Structs)-len(content.CompleteStructs))
		for s, f := range modified {
			for _, af := range f.AnonymousFields {
				items = append(items, fmt.Sprintf("1. %s.%s", s, af))
			}
			for f := range f.Fields {
				items = append(items, fmt.Sprintf("1. %s.%s", s, f))
			}
		}
		sort.Strings(items)
		for _, item := range items {
			md.WriteLine(item)
		}
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/markdown"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
)

// Changelog describes the changelog generated for a package.
type Changelog struct {
	// NewPackage is true if this package does not exist in the old version
	NewPackage bool
	// RemovedPackage is true if this package does not exist in the new version
	RemovedPackage bool
	// Modified contains the details of a modified package. This is nil when either NewPackage or RemovedPackage is true
	Modified *report.Package
}

// HasBreakingChanges returns if this report of changelog contains breaking changes
func (c Changelog) HasBreakingChanges() bool {
	return c.RemovedPackage || (c.Modified != nil && c.Modified.HasBreakingChanges())
}

// String ...
func (c Changelog) String() string {
	return c.ToMarkdown()
}

// ToMarkdown returns the markdown string of this changelog
func (c Changelog) ToMarkdown() string {
	if c.NewPackage {
		return ""
	}
	if c.RemovedPackage {
		return "This package was removed" // this should never be executed
	}
	return c.Modified.ToMarkdown()
}

// ToCompactMarkdown returns the markdown string of this changelog but more compact
func (c Changelog) ToCompactMarkdown() string {
	if c.NewPackage {
		return "This is a new package"
	}
	if c.RemovedPackage {
		return "This package was removed"
	}
	return writeChangelogForPackage(c.Modified)
}

// Seperate change summary text get
func (c Changelog) GetChangeSummary() string {
	if c.NewPackage || c.RemovedPackage {
		return ""
	}
	return getSummaries(c.Modified.BreakingChanges, c.Modified.AdditiveChanges)
}

// GetBreakingChangeItems returns an array of the breaking change items
func (c Changelog) GetBreakingChangeItems() []string {
	if c.RemovedPackage {
		return []string{
			fmt.Sprintf("Package was removed"),
		}
	}
	if c.Modified == nil {
		return []string{}
	}
	return getBreakingChanges(c.Modified.BreakingChanges)
}

func writeChangelogForPackage(r *report.Package) string {
	if r == nil || r.IsEmpty() {
		return "### Other Changes\n"
	}

	md := &markdown.Writer{}

	// write breaking changes
	breakings := getBreakingChanges(r.BreakingChanges)
	if len(breakings) > 0 {
		md.WriteHeader("Breaking Changes")
		for _, item := range breakings {
			md.WriteListItem(item)
		}
	}

	// write additional changes
	additives := getNewContents(r.AdditiveChanges)
	if len(additives) > 0 {
		md.WriteHeader("Features Added")
		for _, item := range additives {
			md.WriteListItem(item)
		}
	}

	return md.String()
}

func getSummaries(breaking *report.BreakingChanges, additive *delta.Content) string {
	bc := 0
	if breaking != nil {
		bc = breaking.Count()
	}
	ac := 0
	if additive != nil {
		ac = additive.Count()
	}

	return fmt.Sprintf("Total %d breaking change(s), %d additive change(s).", bc, ac)
}

func getNewContents(c *delta.Content) []string {
	if c == nil || c.IsEmpty() {
		return nil
	}

	var items []string

	if len(c.Consts) > 0 {
		for _, k := range sortChangeItem(c.Consts) {
			line := fmt.Sprintf("New const `%s`", k)
			items = append(items, line)
		}
	}
	if len(c.TypeAliases) > 0 {
		for _, k := range sortChangeItem(c.TypeAliases) {
			line := fmt.Sprintf("New type alias `%s`", k)
			items = append(items, line)
		}
	}
	if len(c.Funcs) > 0 {
		for _, k := range sortFuncItem(c.Funcs) {
			v := c.Funcs[k]
			params := ""
			if v.Params != nil {
				params = *v.Params
			}
			returns := ""
			if v.Returns != nil {
				returns = *v.Returns
				if strings.Contains(returns, ",") {
					returns = fmt.Sprintf("(%s)", returns)
				}
			}
			line := fmt.Sprintf("New function `%s(%s) %s`", k, params, returns)
			items = append(items, line)
		}
	}
	if len(c.CompleteStructs) > 0 {
		for _, v := range c.CompleteStructs {
			line := fmt.Sprintf("New struct `%s`", v)
			items = append(items, line)
		}
	}
	if len(c.Structs) > 0 {
		modified := c.GetModifiedStructs()
		for _, s := range sortChangeItem(modified) {
			f := modified[s]
			for _, af := range f.AnonymousFields {
				line := fmt.Sprintf("New anonymous field `%s` in struct `%s`", af, s)
				items = append(items, line)
			}
			for _, field := range sortChangeItem(f.Fields) {
				line := fmt.Sprintf("New field `%s` in struct `%s`", field, s)
				items = append(items, line)
			}
		}
	}

	return items
}

func getBreakingChanges(b *report.BreakingChanges) []string {
	items := make([]string, 0)
	if b == nil || b.IsEmpty() {
		return items
	}

	// get signature changes
	items = append(items, getSignatureChangeItems(b)...)

	// get removed content
	items = append(items, getRemovedContent(b.Removed)...)

	return items
}

func getSignatureChangeItems(b *report.BreakingChanges) []string {
	if b.IsEmpty() {
		return nil
	}

	var items []string

	// write const changes
	if len(b.Consts) > 0 {
		for _, k := range sortChangeItem(b.Consts) {
			v := b.Consts[k]
			line := fmt.Sprintf("Const `%s` type has been changed from `%s` to `%s`", k, v.From, v.To)
			items = append(items, line)
		}
	}
	// write type alias changes
	if len(b.TypeAliases) > 0 {
		for _, k := range sortChangeItem(b.TypeAliases) {
			v := b.TypeAliases[k]
			line := fmt.Sprintf("Type alias `%s` type has been changed from `%s` to `%s`", k, v.From, v.To)
			items = append(items, line)
		}
	}
	// write function changes
	if len(b.Funcs) > 0 {
		for _, k := range sortFuncItem(b.Funcs) {
			v := b.Funcs[k]
			if v.Params != nil {
				line := fmt.Sprintf("Function `%s` parameter(s) have been changed from `(%s)` to `(%s)`", k, v.Params.From, v.Params.To)
				items = append(items, line)
			}
			if v.Returns != nil {
				line := fmt.Sprintf("Function `%s` return value(s) have been changed from `(%s)` to `(%s)`", k, v.Returns.From, v.Returns.To)
				items = append(items, line)
			}
		}
	}
	// write struct changes
	if len(b.Structs) > 0 {
		for _, k := range sortChangeItem(b.Structs) {
			v := b.Structs[k]
			for f, d := range v.Fields {
				line := fmt.Sprintf("Type of `%s.%s` has been changed from `%s` to `%s`", k, f, d.From, d.To)
				items = append(items, line)
			}
		}
	}
	// interfaces are skipped, which are identical to some of the functions

	return items
}

func getRemovedContent(removed *delta.Content) []string {
	if removed == nil {
		return nil
	}

	var items []string
	// write constants
	if len(removed.Consts) > 0 {
		for _, k := range sortChangeItem(removed.Consts) {
			line := fmt.Sprintf("Const `%s` has been removed", k)
			items = append(items, line)
		}
	}
	// write type alias
	if len(removed.TypeAliases) > 0 {
		for _, k := range sortChangeItem(removed.TypeAliases) {
			line := fmt.Sprintf("Type alias `%s` has been removed", k)
			items = append(items, line)
		}
	}
	// write functions
	if len(removed.Funcs) > 0 {
		for _, k := range sortFuncItem(removed.Funcs) {
			line := fmt.Sprintf("Function `%s` has been removed", k)
			items = append(items, line)
		}
	}
	// write complete struct removal
	if len(removed.CompleteStructs) > 0 {
		for _, v := range removed.CompleteStructs {
			line := fmt.Sprintf("Struct `%s` has been removed", v)
			items = append(items, line)
		}
	}
	// write struct modification (some fields are removed)
	modified := removed.GetModifiedStructs()
	if len(modified) > 0 {
		for _, s := range sortChangeItem(modified) {
			f := modified[s]
			for _, af := range f.AnonymousFields {
				line := fmt.Sprintf("Field `%s` of struct `%s` has been removed", af, s)
				items = append(items, line)
			}
			for _, field := range sortChangeItem(f.Fields) {
				line := fmt.Sprintf("Field `%s` of struct `%s` has been removed", field, s)
				items = append(items, line)
			}
		}
	}

	return items
}

type sortItem interface {
	delta.Signature | delta.StructDef | exports.Const | exports.TypeAlias | exports.Struct | string
}

func sortChangeItem[T sortItem](change map[string]T) []string {
	s := make([]string, 0, len(change))
	for k := range change {
		s = append(s, k)
	}

	sort.Strings(s)
	return s
}

func sortFuncItem[T delta.FuncSig | exports.Func](change map[string]T) []string {
	s := make([]string, 0, len(change))
	for k := range change {
		s = append(s, k)
	}

	sort.Slice(s, func(i, j int) bool {
		si := removePattern(s[i], getReturnValue(change[s[i]]))
		sj := removePattern(s[j], getReturnValue(change[s[j]]))
		return si < sj
	})

	return s
}

func getReturnValue(t interface{}) string {
	switch value := t.(type) {
	case delta.FuncSig:
		if value.Returns == nil {
			return ""
		}
		return value.Returns.To
	case exports.Func:
		if value.Returns == nil {
			return ""
		}
		return *value.Returns
	}
	return ""
}

func removePattern(funcName string, returnValue string) string {
	funcName = strings.TrimLeft(strings.TrimLeft(funcName, "*"), "New")
	before, after, b := strings.Cut(funcName, ".")
	if !b {
		return funcName
	}

	if strings.Contains(returnValue, "runtime.Poller") {
		after = strings.TrimLeft(after, "Begin")
	} else if strings.Contains(returnValue, "runtime.Pager") {
		after = strings.TrimLeft(after, "New")
	}

	return fmt.Sprintf("%s.%s", before, after)
}

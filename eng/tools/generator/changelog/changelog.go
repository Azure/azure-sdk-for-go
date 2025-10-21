// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

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
	NewPackage bool `json:"newPackage,omitempty"`
	// RemovedPackage is true if this package does not exist in the new version
	RemovedPackage bool `json:"removedPackage,omitempty"`
	// Modified contains the details of a modified package. This is nil when either NewPackage or RemovedPackage is true
	Modified *report.Package `json:"modified,omitempty"`
}

// HasBreakingChanges returns if this report of changelog contains breaking changes
func (c Changelog) HasBreakingChanges() bool {
	return c.RemovedPackage || (c.Modified != nil && c.Modified.HasBreakingChanges())
}

// String convert the changelog to markdown string
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
			"Package was removed",
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
	additives := getAdditiveChanges(r.AdditiveChanges)
	if len(additives) > 0 {
		md.WriteHeader("Features Added")
		for _, item := range additives {
			md.WriteListItem(item)
		}
	}

	return md.String()
}

func getSummaries(breaking *report.BreakingChanges, additive *report.AdditiveChanges) string {
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
	if len(c.Consts) > 0 || len(c.TypeAliases) > 0 {
		newTypeAlias := make(map[string][]string)
		existedTypeAlias := make(map[string][]string)
		for _, k := range sortChangeItem(c.Consts) {
			cs := c.Consts[k]
			if _, ok := c.TypeAliases[cs.Type]; ok {
				if alias, ok := newTypeAlias[cs.Type]; ok {
					alias = append(alias, k)
					newTypeAlias[cs.Type] = alias
				} else {
					alias = []string{k}
					newTypeAlias[cs.Type] = alias
				}
			} else {
				if alias, ok := existedTypeAlias[cs.Type]; ok {
					alias = append(alias, k)
					existedTypeAlias[cs.Type] = alias
				} else {
					existedTypeAlias[cs.Type] = []string{k}
				}
			}
		}

		for _, k := range sortChangeItem(existedTypeAlias) {
			aliasValue := ""
			for _, cs := range existedTypeAlias[k] {
				aliasValue = fmt.Sprintf("%s`%s`, ", aliasValue, cs)
			}
			line := fmt.Sprintf("New value %s added to enum type `%s`", strings.TrimRight(strings.TrimSpace(aliasValue), ","), k)
			items = append(items, line)
		}

		for _, k := range sortChangeItem(newTypeAlias) {
			aliasValue := ""
			for _, cs := range newTypeAlias[k] {
				aliasValue = fmt.Sprintf("%s`%s`, ", aliasValue, cs)
			}
			line := fmt.Sprintf("New enum type `%s` with values %s", k, strings.TrimRight(strings.TrimSpace(aliasValue), ","))
			items = append(items, line)
		}
	}

	if len(c.Funcs) > 0 {
		for _, k := range SortFuncItem(c.Funcs) {
			v := c.Funcs[k]
			params := formatParams(v.Params)
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
			afs := ""
			sort.Strings(f.AnonymousFields)
			for _, af := range f.AnonymousFields {
				afs = fmt.Sprintf("%s`%s`, ", afs, af)
			}
			if afs != "" {
				line := fmt.Sprintf("New anonymous field %s in struct `%s`", strings.TrimSuffix(strings.TrimSpace(afs), ","), s)
				items = append(items, line)
			}

			newFields := ""
			for _, field := range sortChangeItem(f.Fields) {
				newFields = fmt.Sprintf("%s`%s`, ", newFields, field)
			}
			if newFields != "" {
				line := fmt.Sprintf("New field %s in struct `%s`", strings.TrimSuffix(strings.TrimSpace(newFields), ","), s)
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
	items = append(items, getSignatureChangeItems(&b.Changes)...)

	// get removed content
	items = append(items, getRemovedContent(b.Removed)...)

	return items
}

func getAdditiveChanges(a *report.AdditiveChanges) []string {
	items := make([]string, 0)
	if a == nil || a.IsEmpty() {
		return items
	}

	// get signature changes
	items = append(items, getSignatureChangeItems(&a.Changes)...)

	// get added content
	items = append(items, getNewContents(a.Added)...)

	return items
}

func getSignatureChangeItems(b *report.Changes) []string {
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
		for _, k := range SortFuncItem(b.Funcs) {
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
	items = append(items, typeTo(b)...)

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
		removedConst := make(map[string][]string)
		for _, k := range sortChangeItem(removed.Consts) {
			cs := removed.Consts[k]
			if _, ok := removed.TypeAliases[cs.Type]; !ok {
				if alias, ok := removedConst[cs.Type]; ok {
					alias = append(alias, k)
					removedConst[cs.Type] = alias
				} else {
					alias = []string{k}
					removedConst[cs.Type] = alias
				}
			}
		}

		for _, k := range sortChangeItem(removedConst) {
			consts := ""
			for _, cs := range removedConst[k] {
				consts = fmt.Sprintf("%s`%s`, ", consts, cs)
			}
			line := fmt.Sprintf("%s from enum `%s` has been removed", strings.TrimRight(strings.TrimSpace(consts), ","), k)
			items = append(items, line)
		}
	}
	// write type alias
	if len(removed.TypeAliases) > 0 {
		for _, k := range sortChangeItem(removed.TypeAliases) {
			line := fmt.Sprintf("Enum `%s` has been removed", k)
			items = append(items, line)
		}
	}
	// write functions
	if len(removed.Funcs) > 0 {
		var lroItem []string
		var paginationItem []string
		for _, k := range SortFuncItem(removed.Funcs) {
			v := removed.Funcs[k]
			if v.ReplacedBy != nil {
				var line string
				if strings.Contains(k, "Begin") || strings.Contains(*v.ReplacedBy, "Begin") {
					if !strings.Contains(k, "Begin") {
						line = fmt.Sprintf("Operation `%s` has been changed to LRO, use `%s` instead.", k, *v.ReplacedBy)
					} else {
						line = fmt.Sprintf("Operation `%s` has been changed to non-LRO, use `%s` instead.", k, *v.ReplacedBy)
					}
					lroItem = append(lroItem, line)
				} else if strings.Contains(k, "Pager") || strings.Contains(*v.ReplacedBy, "Pager") {
					if !strings.Contains(k, "Pager") {
						line = fmt.Sprintf("Operation `%s` has supported pagination, use `%s` instead.", k, *v.ReplacedBy)
					} else {
						line = fmt.Sprintf("Operation `%s` does not support pagination anymore, use `%s` instead.", k, *v.ReplacedBy)
					}
					paginationItem = append(paginationItem, line)
				}
				continue
			}
			line := fmt.Sprintf("Function `%s` has been removed", k)
			items = append(items, line)
		}
		items = append(items, lroItem...)
		items = append(items, paginationItem...)
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
			afs := ""
			sort.Strings(f.AnonymousFields)
			for _, af := range f.AnonymousFields {
				afs = fmt.Sprintf("%s`%s`, ", afs, af)
			}
			if afs != "" {
				line := fmt.Sprintf("Field %s of struct `%s` has been removed", strings.TrimSuffix(strings.TrimSpace(afs), ","), s)
				items = append(items, line)
			}

			newFields := ""
			for _, field := range sortChangeItem(f.Fields) {
				newFields = fmt.Sprintf("%s`%s`, ", newFields, field)
			}
			if newFields != "" {
				line := fmt.Sprintf("Field %s of struct `%s` has been removed", strings.TrimSuffix(strings.TrimSpace(newFields), ","), s)
				items = append(items, line)
			}
		}
	}

	return items
}

type sortItem interface {
	delta.Signature | delta.StructDef | exports.Const | exports.TypeAlias | exports.Struct | string | []string
}

func formatParams(params []exports.Param) string {
	if len(params) == 0 {
		return ""
	}
	
	var types []string
	for _, p := range params {
		types = append(types, p.Type)
	}
	return strings.Join(types, ", ")
}

func sortChangeItem[T sortItem](change map[string]T) []string {
	if len(change) == 0 {
		return nil
	}

	s := make([]string, 0, len(change))
	for k := range change {
		s = append(s, k)
	}

	sort.Strings(s)
	return s
}

func SortFuncItem[T delta.FuncSig | exports.Func](change map[string]T) []string {
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

func typeTo(b *report.Changes) []string {
	var items []string

	if b == nil || b.IsEmpty() {
		return items
	}

	if len(b.Structs) > 0 {
		for _, k := range sortChangeItem(b.Structs) {
			v := b.Structs[k]
			for _, f := range sortChangeItem(v.Fields) {
				d := v.Fields[f]
				line := fmt.Sprintf("Type of `%s.%s` has been changed from `%s` to `%s`", k, f, d.From, d.To)
				items = append(items, line)
			}
		}
	}

	return items
}

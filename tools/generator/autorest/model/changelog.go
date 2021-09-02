// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/tools/internal/markdown"
	"github.com/Azure/azure-sdk-for-go/tools/internal/report"
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
		return "No exported changes"
	}

	md := &markdown.Writer{}

	// write breaking changes
	md.WriteHeader("Breaking Changes")
	for _, item := range getBreakingChanges(r.BreakingChanges) {
		md.WriteListItem(item)
	}

	// write additional changes
	md.WriteHeader("New Content")
	for _, item := range getNewContents(r.AdditiveChanges) {
		md.WriteListItem(item)
	}

	md.EmptyLine()
	summaries := getSummaries(r.BreakingChanges, r.AdditiveChanges)
	md.WriteLine(summaries)

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
		for k := range c.Consts {
			line := fmt.Sprintf("New const `%s`", k)
			items = append(items, line)
		}
	}
	if len(c.Funcs) > 0 {
		for k, v := range c.Funcs {
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
		for s, f := range modified {
			for _, af := range f.AnonymousFields {
				line := fmt.Sprintf("New anonymous field `%s` in struct `%s`", af, s)
				items = append(items, line)
			}
			for f := range f.Fields {
				line := fmt.Sprintf("New field `%s` in struct `%s`", f, s)
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
		for k, v := range b.Consts {
			line := fmt.Sprintf("Const `%s` type has been changed from `%s` to `%s`", k, v.From, v.To)
			items = append(items, line)
		}
		// TODO -- sort?
	}
	// write function changes
	if len(b.Funcs) > 0 {
		for k, v := range b.Funcs {
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
		for k, v := range b.Structs {
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
		for k := range removed.Consts {
			line := fmt.Sprintf("Const `%s` has been removed", k)
			items = append(items, line)
		}
	}
	// write functions
	if len(removed.Funcs) > 0 {
		for k := range removed.Funcs {
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
		for s, f := range modified {
			for _, af := range f.AnonymousFields {
				line := fmt.Sprintf("Field `%s` of struct `%s` has been removed", af, s)
				items = append(items, line)
			}
			for f := range f.Fields {
				line := fmt.Sprintf("Field `%s` of struct `%s` has been removed", f, s)
				items = append(items, line)
			}
		}
	}

	return items
}

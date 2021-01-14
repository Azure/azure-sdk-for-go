package model

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/delta"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/markdown"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/report"
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

func writeChangelogForPackage(r *report.Package) string {
	if r == nil || r.IsEmpty() {
		return "No exported changes"
	}

	md := &markdown.Writer{}
	// write breaking changes
	writeBreakingChanges(r.BreakingChanges, md)
	// write additional changes
	writeNewContent(r.AdditiveChanges, md)

	writeSummaries(r.BreakingChanges, r.AdditiveChanges, md)

	return md.String()
}

func writeSummaries(breaking *report.BreakingChanges, additive *delta.Content, md *markdown.Writer) {
	bc := 0
	if breaking != nil {
		bc = breaking.Count()
	}
	ac := 0
	if additive != nil {
		ac = additive.Count()
	}
	md.WriteLine("")
	md.WriteLine(fmt.Sprintf("Total %d breaking change(s), %d additive change(s).", bc, ac))
}

func writeNewContent(c *delta.Content, md *markdown.Writer) {
	if c == nil || c.IsEmpty() {
		return
	}
	md.WriteHeader("New Content")
	if len(c.Consts) > 0 {
		for k := range c.Consts {
			line := fmt.Sprintf("- New const `%s`", k)
			md.WriteLine(line)
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
				if strings.Index(returns, ",") > -1 {
					returns = fmt.Sprintf("(%s)", returns)
				}
			}
			line := fmt.Sprintf("- New function `%s(%s) %s`", k, params, returns)
			md.WriteLine(line)
		}
	}
	if len(c.CompleteStructs) > 0 {
		for _, v := range c.CompleteStructs {
			line := fmt.Sprintf("- New struct `%s`", v)
			md.WriteLine(line)
		}
	}
	if len(c.Structs) > 0 {
		modified := c.GetModifiedStructs()
		for s, f := range modified {
			for _, af := range f.AnonymousFields {
				line := fmt.Sprintf("- New anonymous field `%s` in struct `%s`", af, s)
				md.WriteLine(line)
			}
			for f := range f.Fields {
				line := fmt.Sprintf("- New field `%s` in struct `%s`", f, s)
				md.WriteLine(line)
			}
		}
	}
}

func writeBreakingChanges(b *report.BreakingChanges, md *markdown.Writer) {
	if b == nil || b.IsEmpty() {
		return
	}
	md.WriteHeader("Breaking Changes")
	writeSignatureChanges(b, md)
	writeRemovedContent(b.Removed, md)
}

func writeSignatureChanges(b *report.BreakingChanges, md *markdown.Writer) {
	if b.IsEmpty() {
		return
	}
	// write const changes
	if len(b.Consts) > 0 {
		for k, v := range b.Consts {
			line := fmt.Sprintf("- Const `%s` type has been changed from `%s` to `%s`", k, v.From, v.To)
			md.WriteLine(line)
		}
		// TODO -- sort?
	}
	// write function changes
	if len(b.Funcs) > 0 {
		for k, v := range b.Funcs {
			if v.Params != nil {
				line := fmt.Sprintf("- Function `%s` parameter(s) have been changed from `(%s)` to `(%s)`", k, v.Params.From, v.Params.To)
				md.WriteLine(line)
			}
			if v.Returns != nil {
				line := fmt.Sprintf("- Function `%s` return value(s) have been changed from `(%s)` to `(%s)`", k, v.Returns.From, v.Returns.To)
				md.WriteLine(line)
			}
		}
	}
	// write struct changes
	if len(b.Structs) > 0 {
		for k, v := range b.Structs {
			for f, d := range v.Fields {
				line := fmt.Sprintf("- Type of `%s.%s` has been changed from `%s` to `%s`", k, f, d.From, d.To)
				md.WriteLine(line)
			}
		}
	}
	// interfaces are skipped, which are identical to some of the functions
}

func writeRemovedContent(removed *delta.Content, md *markdown.Writer) {
	if removed == nil {
		return
	}
	// write constants
	if len(removed.Consts) > 0 {
		for k := range removed.Consts {
			line := fmt.Sprintf("- Const `%s` has been removed", k)
			md.WriteLine(line)
		}
	}
	// write functions
	if len(removed.Funcs) > 0 {
		for k := range removed.Funcs {
			line := fmt.Sprintf("- Function `%s` has been removed", k)
			md.WriteLine(line)
		}
	}
	// write complete struct removal
	if len(removed.CompleteStructs) > 0 {
		for _, v := range removed.CompleteStructs {
			line := fmt.Sprintf("- Struct `%s` has been removed", v)
			md.WriteLine(line)
		}
	}
	// write struct modification (some fields are removed)
	modified := removed.GetModifiedStructs()
	if len(modified) > 0 {
		for s, f := range modified {
			for _, af := range f.AnonymousFields {
				line := fmt.Sprintf("- Field `%s` of struct `%s` has been removed", af, s)
				md.WriteLine(line)
			}
			for f := range f.Fields {
				line := fmt.Sprintf("- Field `%s` of struct `%s` has been removed", f, s)
				md.WriteLine(line)
			}
		}
	}
}

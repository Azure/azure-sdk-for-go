package changelog

import "github.com/Azure/azure-sdk-for-go/tools/apidiff/report"

type Changelog struct {
	PackageName    string
	NewPackage     bool
	RemovedPackage bool
	Modified       *report.Package
}

func (c Changelog) HasBreakingChanges() bool {
	return c.RemovedPackage || (c.Modified != nil && c.Modified.HasBreakingChanges())
}

func (c Changelog) String() string {
	return c.ToMarkdown()
}

func (c Changelog) ToMarkdown() string {
	if c.NewPackage {
		return "This is a new package"
	}
	if c.RemovedPackage {
		return "This package was removed"
	}
	r := c.Modified.ToMarkdown()
	if r == "" {
		return "No exported changes"
	}
	return r
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package report

// Status ...
type Status interface {
	IsEmpty() bool
	HasBreakingChanges() bool
	HasAdditiveChanges() bool
}

// CommitPkgReport represents a collection of per-package reports, one for each commit hash
type CommitPkgReport struct {
	// BreakingChanges includes the commit hashes that contains breaking changes
	BreakingChanges []string `json:"breakingChanges,omitempty"`
	// CommitsReports stores the package report with the key of commit hashes
	CommitsReports map[string]Package `json:"deltas"`
}

// IsEmpty returns true if the report contains no data
func (c CommitPkgReport) IsEmpty() bool {
	for _, rpt := range c.CommitsReports {
		if !rpt.IsEmpty() {
			return false
		}
	}
	return true
}

// HasBreakingChanges returns true if the report contains breaking changes
func (c CommitPkgReport) HasBreakingChanges() bool {
	for _, r := range c.CommitsReports {
		if r.HasBreakingChanges() {
			return true
		}
	}
	return false
}

// HasAdditiveChanges returns true if the report contains additive changes
func (c CommitPkgReport) HasAdditiveChanges() bool {
	for _, r := range c.CommitsReports {
		if r.HasAdditiveChanges() {
			return true
		}
	}
	return false
}

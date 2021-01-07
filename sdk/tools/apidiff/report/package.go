// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

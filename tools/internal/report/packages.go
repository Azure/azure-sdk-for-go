// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package report

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/tools/internal/markdown"
	"github.com/Azure/azure-sdk-for-go/tools/internal/repo"
	"io/ioutil"
	"os"
	"path/filepath"
)

const apiDirSuffix = "api"

// GetPackages returns all the go sdk packages under the given root directory
func GetPackages(dir string) ([]string, error) {
	var pkgDirs []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// check if leaf dir
			fi, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			hasSubDirs := false
			for _, f := range fi {
				// check if this is the interfaces subdir, if it is don't count it as a subdir
				if f.IsDir() && f.Name() != filepath.Base(path)+apiDirSuffix {
					hasSubDirs = true
					break
				}
			}
			if !hasSubDirs {
				pkgDirs = append(pkgDirs, path)
				// skip any dirs under us (i.e. interfaces subdir)
				return filepath.SkipDir
			}
		}
		return nil
	})
	return pkgDirs, err
}

// PkgsList contains a collection of packages
type PkgsList []string

// ModifiedPackages contains a collection of package reports, it's structured as "package path":pkgReport
type ModifiedPackages map[string]Package

// IsEmpty ...
func (m ModifiedPackages) IsEmpty() bool {
	return len(m) == 0
}

// HasBreakingChanges returns true if any package contained in has a breaking change
func (m ModifiedPackages) HasBreakingChanges() bool {
	for _, p := range m {
		if p.HasBreakingChanges() {
			return true
		}
	}
	return false
}

// HasAdditiveChanges returns true if any package contained in has an additive change
func (m ModifiedPackages) HasAdditiveChanges() bool {
	for _, p := range m {
		if p.HasAdditiveChanges() {
			return true
		}
	}
	return false
}

// CommitPkgsReport represents a collection of reports, one for each commit hash.
type CommitPkgsReport struct {
	// AffectedPackages stores the package list with key of commit hashes
	AffectedPackages map[string]PkgsList `json:"affectedPackages"`
	// BreakingChanges stores the commit hashes that contain breaking changes
	BreakingChanges []string `json:"breakingChanges,omitempty"`
	// CommitsReports stores the detailed reports with the key of commit hashes
	CommitsReports map[string]PkgsReport `json:"deltas"`
}

// IsEmpty returns true if the report contains no data.
func (c CommitPkgsReport) IsEmpty() bool {
	for _, r := range c.CommitsReports {
		if !r.IsEmpty() {
			return false
		}
	}
	return true
}

// HasBreakingChanges returns true if the report contains breaking changes.
func (c CommitPkgsReport) HasBreakingChanges() bool {
	for _, r := range c.CommitsReports {
		if r.HasBreakingChanges() {
			return true
		}
	}
	return false
}

// HasAdditiveChanges returns true if the package contains additive changes.
func (c CommitPkgsReport) HasAdditiveChanges() bool {
	for _, r := range c.CommitsReports {
		if r.HasAdditiveChanges() {
			return true
		}
	}
	return false
}

// UpdateAffectedPackages updates the collection of affected packages with the packages that were touched in the specified commit
func (c *CommitPkgsReport) UpdateAffectedPackages(commit string, r PkgsReport) {
	if c.AffectedPackages == nil {
		c.AffectedPackages = map[string]PkgsList{}
	}

	for _, pkg := range r.AddedPackages {
		c.AffectedPackages[commit] = append(c.AffectedPackages[commit], pkg)
	}

	for pkgName := range r.ModifiedPackages {
		c.AffectedPackages[commit] = append(c.AffectedPackages[commit], pkgName)
	}

	for _, pkg := range r.RemovedPackages {
		c.AffectedPackages[commit] = append(c.AffectedPackages[commit], pkg)
	}
}

// represents a complete report of added, removed, and modified packages
type PkgsReport struct {
	// AddedPackages contains the relative package names that are added in this new version
	AddedPackages PkgsList `json:"added,omitempty"`
	// ModifiedPackages contains the modified packages map, using relative package names as keys
	ModifiedPackages ModifiedPackages `json:"modified,omitempty"`
	// RemovedPackages contains the relative package names that are removed in this new version
	RemovedPackages    PkgsList `json:"removed,omitempty"`
	modPkgHasAdditions bool
	modPkgHasBreaking  bool
}

// returns true if the package report contains breaking changes
func (r PkgsReport) HasBreakingChanges() bool {
	return len(r.RemovedPackages) > 0 || r.modPkgHasBreaking
}

// returns true if the package report contains additive changes
func (r PkgsReport) HasAdditiveChanges() bool {
	return len(r.AddedPackages) > 0 || r.modPkgHasAdditions
}

// returns true if the report contains no data
func (r PkgsReport) IsEmpty() bool {
	return len(r.AddedPackages) == 0 && len(r.ModifiedPackages) == 0 && len(r.RemovedPackages) == 0
}

func (r PkgsReport) ToMarkdown() string {
	if r.IsEmpty() {
		return "No exports changed"
	}

	md := &markdown.Writer{}
	r.writeAddedPackages(md)
	r.writeRemovedPackages(md)
	r.writeModifiedPackages(md)
	return md.String()
}

func (r *PkgsReport) writeAddedPackages(md *markdown.Writer) {
	if len(r.AddedPackages) == 0 {
		return
	}
	md.WriteTopLevelHeader("New Packages")
	// write list
	for _, p := range r.AddedPackages {
		md.WriteLine("- " + p)
	}
}

func (r *PkgsReport) writeModifiedPackages(md *markdown.Writer) {
	if len(r.ModifiedPackages) == 0 {
		return
	}
	md.WriteTopLevelHeader("Modified Packages")
	// write list
	for n := range r.ModifiedPackages {
		md.WriteLine(fmt.Sprintf("- `%s`", n))
	}
	// write details
	for n, p := range r.ModifiedPackages {
		md.WriteTopLevelHeader(fmt.Sprintf("Modified `%s`", n))
		md.WriteLine(p.ToMarkdown())
	}
}

func (r *PkgsReport) writeRemovedPackages(md *markdown.Writer) {
	if len(r.RemovedPackages) == 0 {
		return
	}
	md.WriteTopLevelHeader("Removed Packages")
	// write list
	for _, p := range r.RemovedPackages {
		md.WriteLine("- " + p)
	}
	md.EmptyLine()
}

// generates a PkgsReport based on the delta between lhs and rhs
func GetPkgsReport(lhs, rhs repo.RepoContent, option *GenerationOption) PkgsReport {
	rpt := PkgsReport{}

	if option == nil {
		option = &GenerationOption{}
	}

	if !option.OnlyBreakingChanges {
		rpt.AddedPackages = getDiffPkgs(lhs, rhs)
	}
	if !option.OnlyAdditiveChanges {
		rpt.RemovedPackages = getDiffPkgs(rhs, lhs)
	}

	// diff packages
	for rhsPkg, rhsCnt := range rhs {
		if _, ok := lhs[rhsPkg]; !ok {
			continue
		}
		if r := Generate(lhs[rhsPkg], rhsCnt, option); !r.IsEmpty() {
			// only add an entry if the report contains data
			if rpt.ModifiedPackages == nil {
				rpt.ModifiedPackages = ModifiedPackages{}
			}
			rpt.ModifiedPackages[rhsPkg] = r
		}
	}

	return rpt
}

// returns a list of packages in rhs that aren't in lhs
func getDiffPkgs(lhs, rhs repo.RepoContent) PkgsList {
	list := PkgsList{}
	for rhsPkg := range rhs {
		if _, ok := lhs[rhsPkg]; !ok {
			list = append(list, rhsPkg)
		}
	}
	return list
}

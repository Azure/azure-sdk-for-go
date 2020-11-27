package report

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const apiDirSuffix = "api"

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

// contains a collection of packages
type PkgsList []string

// contains a collection of package reports, it's structured as "package path":pkgReport
type ModifiedPackages map[string]Package

func (m ModifiedPackages) IsEmpty() bool {
	return len(m) == 0
}

func (m ModifiedPackages) HasBreakingChanges() bool {
	for _, p := range m {
		if p.HasBreakingChanges() {
			return true
		}
	}
	return false
}

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
	AffectedPackages map[string]PkgsList   `json:"affectedPackages"`
	BreakingChanges  []string              `json:"breakingChanges,omitempty"`
	CommitsReports   map[string]PkgsReport `json:"deltas"`
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

// updates the collection of affected packages with the packages that were touched in the specified commit
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
	AddedPackages      PkgsList         `json:"added,omitempty"`
	ModifiedPackages   ModifiedPackages `json:"modified,omitempty"`
	RemovedPackages    PkgsList         `json:"removed,omitempty"`
	modPkgHasAdditions bool
	modPkgHasBreaking  bool
}

// returns true if the package report contains breaking changes
func (r PkgsReport) HasBreakingChanges() bool {
	return len(r.RemovedPackages) > 0 || (r.ModifiedPackages != nil && r.ModifiedPackages.HasBreakingChanges())
}

// returns true if the package report contains additive changes
func (r PkgsReport) HasAdditiveChanges() bool {
	return len(r.AddedPackages) > 0 || (r.ModifiedPackages != nil && r.ModifiedPackages.HasAdditiveChanges())
}

// returns true if the report contains no data
func (r PkgsReport) IsEmpty() bool {
	return len(r.AddedPackages) == 0 && len(r.ModifiedPackages) == 0 && len(r.RemovedPackages) == 0
}

func (r *PkgsReport) ToMarkdown() string {
	if r.IsEmpty() {
		return ""
	}
	md := MarkdownWriter{}
	r.writeAddedPackages(&md)
	r.writeRemovedPackages(&md)
	r.writeModifiedPackages(&md)
	return md.String()
}

func (r *PkgsReport) writeAddedPackages(md *MarkdownWriter) {
	if len(r.AddedPackages) == 0 {
		return
	}
	md.WriteHeader("New Packages")
	// write list
	for _, p := range r.AddedPackages {
		md.WriteLine("- " + p)
	}
}

func (r *PkgsReport) writeModifiedPackages(md *MarkdownWriter) {
	if len(r.ModifiedPackages) == 0 {
		return
	}
	md.WriteHeader("Modified Packages")
	// write list
	for n := range r.ModifiedPackages {
		md.WriteLine(fmt.Sprintf("- `%s`", n))
	}
	// write details
	for n, p := range r.ModifiedPackages {
		md.WriteHeader(fmt.Sprintf("Modified `%s`", n))
		md.WriteLine(p.ToMarkdown())
	}
}

func (r *PkgsReport) writeRemovedPackages(md *MarkdownWriter) {
	if len(r.RemovedPackages) == 0 {
		return
	}
	md.WriteHeader("Removed Packages")
	// write list
	for _, p := range r.RemovedPackages {
		md.WriteLine("- " + p)
	}
	md.EmptyLine()
}

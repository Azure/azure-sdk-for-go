package changelog_ext

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/markdown"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/report"
	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
	"github.com/Azure/azure-sdk-for-go/tools/pkgchk/track1"
	"github.com/ahmetb/go-linq/v3"
)

type Writer struct {
	file    string
	version string
}

func NewWriterFromFile(file string) *Writer {
	return &Writer{
		file: file,
	}
}

func (w *Writer) WithVersion(version string) *Writer {
	w.version = version
	return w
}

// Write writes the new version changelog to the changelog file, also modify the links in the previous version
func (w *Writer) Write(r *PkgsReport) error {
	// first we read the changelog, get the lines as an array
	lines, err := w.read()
	if err != nil {
		return err
	}
	// make some modification to the old ones
	// find the version number titles
	versionRange := FindVersionTitles(lines, 2)
	if len(versionRange) < 1 {
		return fmt.Errorf("cannot find a version in changelog")
	}
	// separate changelog into three parts
	title, previousVersionContent, rest := GetLinesBetween(lines, versionRange)
	// modify the previous version content
	previousVersionContent = w.modifyPreviousVersionContent(previousVersionContent, versionRange[0].Version)
	// get the new changelog
	changelog := w.WriteReport(r)
	// write the changelog back to changelog file
	return w.rewriteChangelog(title, previousVersionContent, rest, changelog)
}

func (w *Writer) modifyPreviousVersionContent(previousVersionContent []string, previousVersion string) []string {
	// iterate over the previous version part
	for i := range previousVersionContent {
		matches := previousVersionLinkRegex.FindStringSubmatch(previousVersionContent[i])
		if len(matches) >= 2 {
			link := matches[0]
			linkContent := matches[1]
			newLinkContent := fmt.Sprintf(newLinkFmt, previousVersion, strings.TrimPrefix(linkContent, absoluteLinkPrefix))
			newLink := strings.Replace(link, linkContent, newLinkContent, 1)
			previousVersionContent[i] = strings.Replace(previousVersionContent[i], link, newLink, 1)
		}
	}
	return previousVersionContent
}

func (w *Writer) rewriteChangelog(title, previousVersionContent, rest []string, changelog string) error {
	var lines []string
	lines = append(lines, title...)
	lines = append(lines, fmt.Sprintf("## `%s`\n\n%s\n", w.version, strings.TrimSpace(changelog)))
	lines = append(lines, previousVersionContent...)
	lines = append(lines, rest...)
	return w.write(lines)
}

type VersionTitleLine struct {
	LineNumber int
	Version    string
}

func GetLinesBetween(lines []string, versionRange []VersionTitleLine) ([]string, []string, []string) {
	if len(versionRange) < 2 {
		return lines[:versionRange[0].LineNumber], lines[versionRange[0].LineNumber:], nil
	}
	return lines[:versionRange[0].LineNumber], lines[versionRange[0].LineNumber:versionRange[1].LineNumber], lines[versionRange[1].LineNumber:]
}

func FindVersionTitles(lines []string, n int) []VersionTitleLine {
	var results []VersionTitleLine
	if n < 0 {
		n = len(lines)
	}
	for i := 0; i < len(lines); i++ {
		if len(results) >= n {
			break
		}
		matches := versionNumberLineRegex.FindStringSubmatch(strings.TrimSpace(lines[i]))
		if len(matches) >= 2 {
			results = append(results, VersionTitleLine{
				LineNumber: i,
				Version:    matches[1],
			})
		}
	}
	return results
}

var (
	newLinkFmt               = `https://github.com/Azure/azure-sdk-for-go/blob/%s/%s`
	versionNumberLineRegex   = regexp.MustCompile("^## `(.+)`$")
	previousVersionLinkRegex = regexp.MustCompile(fmt.Sprintf(`\[details\]\((%sservices/.+)\)`, absoluteLinkPrefix))
)

func (w *Writer) read() ([]string, error) {
	file, err := os.Open(w.file)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var lines []string
	linq.From(strings.Split(string(b), "\n")).Select(func(line interface{}) interface{} {
		return strings.TrimSpace(line.(string))
	}).ToSlice(&lines)
	return lines, nil
}

func (w *Writer) write(lines []string) error {
	file, err := os.Create(w.file)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(strings.Join(lines, "\n")); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteReport(r *PkgsReport) string {
	if r == nil || r.isEmpty() {
		return "No exports changed"
	}
	md := &markdown.Writer{}
	writeAddedPackages(r, md)
	writeRemovedPackages(r, md)
	writeModifiedPackages(r, md)

	return md.String()
}

func writeAddedPackages(r *PkgsReport, md *markdown.Writer) {
	if len(r.AddedPackages) == 0 {
		return
	}
	md.WriteHeader("New Packages")
	// write list
	var lines []string
	for _, p := range r.AddedPackages {
		lines = append(lines, fmt.Sprintf("- `%s`", getPackageImportPath(p)))
	}

	sort.Strings(lines)

	for _, line := range lines {
		md.WriteLine(line)
	}
}

func writeRemovedPackages(r *PkgsReport, md *markdown.Writer) {
	if len(r.RemovedPackages) == 0 {
		return
	}
	md.WriteHeader("Removed Packages")
	// write list
	var lines []string
	for _, p := range r.RemovedPackages {
		lines = append(lines, fmt.Sprintf("- `%s`", getPackageImportPath(p)))
	}

	sort.Strings(lines)

	for _, line := range lines {
		md.WriteLine(line)
	}
}

func writeModifiedPackages(r *PkgsReport, md *markdown.Writer) {
	// categorize the modified packages as `breaking change` and `non-breaking change`
	var updated []string
	var breaking []string
	for n, p := range r.ModifiedPackages {
		if p.HasBreakingChanges() {
			breaking = append(breaking, getPackageImportPath(n))
		} else {
			updated = append(updated, getPackageImportPath(n))
		}
	}
	sort.Strings(updated)
	writeUpdatedPackages(updated, md)
	sort.Strings(breaking)
	writeBreakingPackages(breaking, md)
}

func writeUpdatedPackages(updated []string, md *markdown.Writer) {
	if len(updated) == 0 {
		return
	}
	md.WriteHeader("Updated Packages")
	t := markdown.NewTable("lc", "Package Path", "Changelog")
	// write table
	for _, n := range updated {
		link := fmt.Sprintf("[details](%s)", getAbsoluteChangelogLink(n))
		t.AddRow(fmt.Sprintf("`%s`", n), link)
	}
	md.WriteTable(*t)
}

func writeBreakingPackages(breaking []string, md *markdown.Writer) {
	if len(breaking) == 0 {
		return
	}
	md.WriteHeader("Breaking Changes")
	t := markdown.NewTable("lc", "Package Path", "Changelog")
	// write table
	for _, n := range breaking {
		link := fmt.Sprintf("[details](%s)", getAbsoluteChangelogLink(n))
		t.AddRow(fmt.Sprintf("`%s`", n), link)
	}
	md.WriteTable(*t)
}

// getAbsoluteChangelogLink converts the path name to the absolute link of the corresponding changelog file
func getAbsoluteChangelogLink(name string) string {
	relativeLink := getChangelogLink(name)
	return fmt.Sprintf("%s%s", absoluteLinkPrefix, relativeLink)
}

// getChangelogLink gets the relative path of the package changelog file
func getChangelogLink(name string) string {
	rel := strings.TrimPrefix(name, sdkRoot+"/")
	return utils.NormalizePath(sdk.ChangelogPath(rel))
}

func getPackageImportPath(p string) string {
	return fmt.Sprintf("github.com/Azure/azure-sdk-for-go/%s", p)
}

// TODO -- I have rewrite everything since not everything in apidiff tool in azure-sdk-for-go is exported
// A good part of rewriting everything is that we could keep these functions consistent

func GetPackagesReportFromContent(lhs RepoContent, targetRoot string) (*PkgsReport, error) {
	rhs, err := GetRepoContent(targetRoot)
	if err != nil {
		return nil, err
	}
	r := GetPkgsReport(lhs, rhs)
	return &r, nil
}

// RepoContent contains repo content, it's structured as "package path (relative to the root of sdk)":content
type RepoContent map[string]exports.Content

func GetRepoContent(sdkRoot string) (RepoContent, error) {
	// we must list over the services directory, otherwise it would walk into the .git directory and panic out
	pkgs, err := track1.List(sdk.ServicesPath(sdkRoot))
	if err != nil {
		return nil, err
	}

	r, err := getExportsForPackages(pkgs, sdkRoot)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// returns repoContent based on the provided slice of package directories
func getExportsForPackages(pkgs []track1.Package, root string) (RepoContent, error) {
	exps := RepoContent{}
	for _, pkg := range pkgs {
		relativePath, err := filepath.Rel(root, pkg.FullPath())
		if err != nil {
			return nil, err
		}
		relativePath = utils.NormalizePath(relativePath)
		if _, ok := exps[relativePath]; ok {
			return nil, fmt.Errorf("duplicate package: %s", pkg.Path())
		}
		exp, err := exports.Get(pkg.FullPath())
		if err != nil {
			return nil, err
		}
		exps[relativePath] = exp
	}
	return exps, nil
}

const (
	sdkRoot            = "github.com/Azure/azure-sdk-for-go"
	absoluteLinkPrefix = "https://github.com/Azure/azure-sdk-for-go/tree/master/"
)

// contains a collection of packages
type PkgsList []string

// contains a collection of package reports, it's structured as "package path":pkgReport
type ModifiedPackages map[string]report.Package

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
func (r PkgsReport) hasBreakingChanges() bool {
	return len(r.RemovedPackages) > 0 || r.modPkgHasBreaking
}

// returns true if the package report contains additive changes
func (r PkgsReport) hasAdditiveChanges() bool {
	return len(r.AddedPackages) > 0 || r.modPkgHasAdditions
}

// returns true if the report contains no data
func (r PkgsReport) isEmpty() bool {
	return len(r.AddedPackages) == 0 && len(r.ModifiedPackages) == 0 && len(r.RemovedPackages) == 0
}

// GetPkgsReport generates a pkgsReport based on the delta between lhs and rhs
func GetPkgsReport(lhs, rhs RepoContent) PkgsReport {
	rpt := PkgsReport{}

	rpt.AddedPackages = GetDiffPackages(lhs, rhs)
	rpt.RemovedPackages = GetDiffPackages(rhs, lhs)

	// diff packages
	for rhsPkg, rhsCnt := range rhs {
		if _, ok := lhs[rhsPkg]; !ok {
			continue
		}
		if r := report.Generate(lhs[rhsPkg], rhsCnt, nil); !r.IsEmpty() {
			if r.HasBreakingChanges() {
				rpt.modPkgHasBreaking = true
			}
			if r.HasAdditiveChanges() {
				rpt.modPkgHasAdditions = true
			}
			// only add an entry if the report contains data
			if rpt.ModifiedPackages == nil {
				rpt.ModifiedPackages = ModifiedPackages{}
			}
			rpt.ModifiedPackages[rhsPkg] = r
		}
	}

	return rpt
}

func (r PkgsReport) ToMarkdown() string {
	if r.isEmpty() {
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

// GetDiffPackages returns a list of packages in rhs that aren't in lhs
func GetDiffPackages(lhs, rhs RepoContent) PkgsList {
	list := PkgsList{}
	for rhsPkg := range rhs {
		if _, ok := lhs[rhsPkg]; !ok {
			list = append(list, rhsPkg)
		}
	}
	return list
}

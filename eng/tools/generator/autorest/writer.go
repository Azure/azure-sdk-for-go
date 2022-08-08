// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/markdown"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/ahmetb/go-linq/v3"
)

const (
	sdkRoot            = "github.com/Azure/azure-sdk-for-go"
	absoluteLinkPrefix = "https://github.com/Azure/azure-sdk-for-go/tree/main/"
)

type Writer struct {
	file    string
	version string
}

func (w *Writer) WithVersion(version string) *Writer {
	w.version = version
	return w
}

func NewWriterFromFile(file string) *Writer {
	return &Writer{
		file: file,
	}
}

// Write writes the new version changelog to the changelog file, also modify the links in the previous version
func (w *Writer) Write(r *report.PkgsReport) error {
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

func (w *Writer) WriteReport(r *report.PkgsReport) string {
	if r == nil || r.IsEmpty() {
		return "No exports changed"
	}
	md := &markdown.Writer{}
	writeAddedPackages(r, md)
	writeRemovedPackages(r, md)
	writeModifiedPackages(r, md)

	return md.String()
}

func writeAddedPackages(r *report.PkgsReport, md *markdown.Writer) {
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

func writeRemovedPackages(r *report.PkgsReport, md *markdown.Writer) {
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

func writeModifiedPackages(r *report.PkgsReport, md *markdown.Writer) {
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
	return utils.NormalizePath(common.ChangelogPath(rel))
}

func getPackageImportPath(p string) string {
	return fmt.Sprintf("github.com/Azure/azure-sdk-for-go/%s", p)
}

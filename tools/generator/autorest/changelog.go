// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/cmd"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/markdown"
	"github.com/Azure/azure-sdk-for-go/tools/internal/repo"
	"github.com/Azure/azure-sdk-for-go/tools/internal/report"
	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
	"github.com/ahmetalpbalkan/go-linq"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ChangelogContext describes all necessary data that would be needed in the processing of changelogs
type ChangelogContext interface {
	SDKRoot() string
	RepoContent() map[string]exports.Content
}

// ChangelogProcessor processes the metadata and output changelog with the desired format
type ChangelogProcessor struct {
	ctx              ChangelogContext
	metadataLocation string
	readme           string
}

// NewChangelogProcessorFromContext returns a new ChangelogProcessor
func NewChangelogProcessorFromContext(ctx ChangelogContext) *ChangelogProcessor {
	return &ChangelogProcessor{
		ctx: ctx,
	}
}

// ChangelogResult describes the result of the generated changelog for one package
type ChangelogResult struct {
	Tag             string
	PackageName     string
	PackageFullPath string
	Changelog       model.Changelog
}

// ChangelogProcessError describes the errors during the processing
type ChangelogProcessError struct {
	Errors []error
}

// Error ...
func (e *ChangelogProcessError) Error() string {
	return fmt.Sprintf("total %d error(s) during processing changelog: %+v", len(e.Errors), e.Errors)
}

type changelogErrorBuilder struct {
	errors []error
}

func (b *changelogErrorBuilder) add(err error) {
	b.errors = append(b.errors, err)
}

func (b *changelogErrorBuilder) build() error {
	if len(b.errors) == 0 {
		return nil
	}
	return &ChangelogProcessError{
		Errors: b.errors,
	}
}

// Process generates the changelogs using the input metadata map.
// Please ensure the input metadata map does not contain any package that is not under the sdk root, otherwise this might give weird results.
func (p *ChangelogProcessor) Process(metadataMap map[string]model.Metadata) ([]ChangelogResult, error) {
	builder := changelogErrorBuilder{}
	var results []ChangelogResult
	for tag, metadata := range metadataMap {
		outputFolder := filepath.Clean(metadata.PackagePath())
		// format the package before getting the changelog
		if err := FormatPackage(outputFolder); err != nil {
			builder.add(err)
			continue
		}
		result, err := p.GenerateChangelog(outputFolder, tag)
		if err != nil {
			builder.add(err)
			continue
		}
		results = append(results, *result)
	}
	return results, builder.build()
}

// GenerateChangelog generates a changelog for one package
func (p *ChangelogProcessor) GenerateChangelog(packageFullPath, tag string) (*ChangelogResult, error) {
	relativePackagePath, err := p.getRelativePackagePath(packageFullPath)
	if err != nil {
		return nil, err
	}
	lhs := p.getLHSContent(relativePackagePath)
	rhs, err := getExportsForPackage(packageFullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get exports from package '%s' in the sdk '%s': %+v", relativePackagePath, p.ctx.SDKRoot(), err)
	}
	r, err := GetChangelogForPackage(lhs, rhs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate changelog for package '%s': %+v", relativePackagePath, err)
	}
	return &ChangelogResult{
		Tag:             tag,
		PackageName:     relativePackagePath,
		PackageFullPath: packageFullPath,
		Changelog:       *r,
	}, nil
}

func (p *ChangelogProcessor) getRelativePackagePath(fullPath string) (string, error) {
	relative, err := filepath.Rel(p.ctx.SDKRoot(), fullPath)
	if err != nil {
		return "", err
	}
	return utils.NormalizePath(relative), nil
}

func (p *ChangelogProcessor) getLHSContent(relativePackagePath string) *exports.Content {
	if v, ok := p.ctx.RepoContent()[relativePackagePath]; ok {
		return &v
	}
	return nil
}

func getExportsForPackage(dir string) (*exports.Content, error) {
	// The function exports.Get does not handle the circumstance that the package does not exist
	// therefore we have to check if it exists and if not exit early to ensure we do not return an error
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}
	exp, err := exports.Get(dir)
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

// GetChangelogForPackage generates the changelog report with the given two Contents
func GetChangelogForPackage(lhs, rhs *exports.Content) (*model.Changelog, error) {
	if lhs == nil && rhs == nil {
		return nil, fmt.Errorf("this package does not exist even after the generation, this should never happen")
	}
	if lhs == nil {
		// the package does not exist before the generation: this is a new package
		return &model.Changelog{
			NewPackage: true,
		}, nil
	}
	if rhs == nil {
		// the package no longer exists after the generation: this package was removed
		return &model.Changelog{
			RemovedPackage: true,
		}, nil
	}
	// lhs and rhs are both non-nil
	p := report.Generate(*lhs, *rhs, nil)
	return &model.Changelog{
		Modified: &p,
	}, nil
}

// ToMarkdown convert this ChangelogResult to markdown string
func (r ChangelogResult) ToMarkdown() string {
	return r.Changelog.ToMarkdown()
}

// Write writes this ChangelogResult to the specified writer in markdown format
func (r ChangelogResult) Write(writer io.Writer) error {
	_, err := writer.Write([]byte(r.ToMarkdown()))
	return err
}

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
	return utils.NormalizePath(sdk.ChangelogPath(rel))
}

func getPackageImportPath(p string) string {
	return fmt.Sprintf("github.com/Azure/azure-sdk-for-go/%s", p)
}

func GetPackagesReportFromContent(lhs repo.RepoContent, targetRoot string) (*report.PkgsReport, error) {
	rhs, err := repo.GetRepoContent(targetRoot)
	if err != nil {
		return nil, err
	}
	r := cmd.GetPkgsReport(lhs, rhs)
	return &r, nil
}

const (
	sdkRoot            = "github.com/Azure/azure-sdk-for-go"
	absoluteLinkPrefix = "https://github.com/Azure/azure-sdk-for-go/tree/master/"
)

const (
	// ChangelogFilename ...
	ChangelogFilename = "CHANGELOG.md"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetAllVersionTags(moduleRelativePath string, sdkRepo repo.SDKRepository) ([]string, error) {
	arr := strings.Split(moduleRelativePath, "/")
	log.Printf("Fetching all release tags from GitHub for RP: '%s' Package: '%s' ...", arr[len(arr)-2], arr[len(arr)-1])

	remoteName := "release_remote"
	fetchOpts := &git.FetchOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{"refs/tags/*:refs/tags/*"},
		Tags:       git.AllTags,
	}

	err := fetchTagsFromRemote(sdkRepo, remoteName, fetchOpts)
	if err != nil {
		return nil, err
	}

	// Get all tags
	tags, err := sdkRepo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %v", err)
	}

	var versions []string
	var result []string
	versionTag := make(map[string]string)
	semverRegex := regexp.MustCompile(semver.SemVerRegex) // Precompile the regex
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		tagName := ref.Name().String()
		if strings.Contains(tagName, moduleRelativePath+"/v") {
			matchedVersion := semverRegex.FindString(tagName)
			if matchedVersion != "" {
				versions = append(versions, matchedVersion)
				versionTag[matchedVersion] = tagName
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process tags: %v", err)
	}

	// Sort versions in descending order
	vs := make([]*semver.Version, len(versions))
	for i, r := range versions {
		v, err := semver.NewVersion(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version %s: %v", r, err)
		}
		vs[i] = v
	}
	sort.Sort(sort.Reverse(semver.Collection(vs)))

	// Build final result
	for _, v := range vs {
		result = append(result, versionTag[v.Original()])
	}
	if err := cleanupRemote(sdkRepo, remoteName); err != nil {
		return nil, err
	}

	return result, nil
}

func fetchTagsFromRemote(sdkRepo repo.SDKRepository, remoteName string, fetchOpts *git.FetchOptions) error {
	// Create remote with center sdk repo if it doesn't exist
	_, err := sdkRepo.CreateRemote(&config.RemoteConfig{Name: remoteName, URLs: []string{utils.SDKRemoteURL}})
	if err != nil && err != git.ErrRemoteExists {
		return fmt.Errorf("failed to create remote: %v", err)
	}

	// Fetch tags from remote
	err = sdkRepo.Fetch(fetchOpts)
	// It's normal to get "already up-to-date" error if tags are already fetched
	if err != nil && err != git.NoErrAlreadyUpToDate && err.Error() != "already up-to-date" {
		return fmt.Errorf("failed to fetch: %v", err)
	}

	return nil
}

func cleanupRemote(sdkRepo repo.SDKRepository, remoteName string) error {
	// remove remote
	err := sdkRepo.DeleteRemote(remoteName)
	if err != nil {
		return fmt.Errorf("failed to delete remote: %v", err)
	}
	return nil
}

func getPreviousVersionTag(isCurrentPreview bool, allReleases []string) string {
	if isCurrentPreview {
		// for preview api, always compare with latest release
		return allReleases[0]
	} else {
		// for stable api, always compare with previous stable, if no stable, then latest release
		for _, release := range allReleases {
			if !strings.Contains(release, "beta") {
				return release
			}
		}
		return allReleases[0]
	}
}

func getExportsFromTag(moduleRelativePath, tag string) (*exports.Content, error) {
	log.Printf("Get exports for '%s' from specific tag '%s' ...", moduleRelativePath, tag)

	// Extract tag name from ref/tags/ prefix if present
	tagName := strings.TrimPrefix(tag, "refs/tags/")

	// Download and extract only the module folder
	extractedPath, err := downloadAndExtractCode(tagName, moduleRelativePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download and extract code: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(extractedPath)) // Clean up temp directory

	// Check if the module path exists in the extracted source
	if _, err := os.Stat(extractedPath); os.IsNotExist(err) {
		log.Printf("Module path '%s' does not exist in tag '%s'", moduleRelativePath, tagName)
		return &exports.Content{}, nil
	}

	// Get exports from the extracted module path
	result, err := exports.Get(extractedPath)
	// bypass the error if the module doesn't contain any exports, return empty content
	if err != nil && strings.Contains(err.Error(), "doesn't contain any exports") {
		return &exports.Content{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get exports from extracted code: %v", err)
	}

	return &result, nil
}

// downloadAndExtractCode downloads the source code for a specific tag and extracts only the moduleRelativePath folder
func downloadAndExtractCode(tagName, moduleRelativePath string) (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "azure-sdk-for-go-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Download the tag archive
	downloadURL := fmt.Sprintf("https://github.com/Azure/azure-sdk-for-go/archive/refs/tags/%s.zip", tagName)
	log.Printf("Downloading tag source from: %s", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to download tag archive: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to download tag archive: HTTP %d", resp.StatusCode)
	}

	// Save the zip file to temp directory
	zipFilePath := filepath.Join(tempDir, "source.zip")
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to create zip file: %v", err)
	}

	_, err = io.Copy(zipFile, resp.Body)
	zipFile.Close()
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to save zip file: %v", err)
	}

	// Extract only the moduleRelativePath folder from the zip file
	extractedPath, err := extractCodeFromZip(zipFilePath, tempDir, moduleRelativePath)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to extract code from zip file: %v", err)
	}

	// Remove the zip file after extraction
	os.Remove(zipFilePath)

	return extractedPath, nil
}

// extractCodeFromZip extracts only the specified moduleRelativePath from the zip file
func extractCodeFromZip(src, dest, moduleRelativePath string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var extractedPath string
	repoPrefix := ""

	// First pass: find the repository prefix (e.g., "azure-sdk-for-go-v1.2.3/")
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "azure-sdk-for-go-") && strings.Contains(f.Name, "/") {
			parts := strings.SplitN(f.Name, "/", 2)
			if len(parts) > 0 {
				repoPrefix = parts[0] + "/"
				break
			}
		}
	}

	if repoPrefix == "" {
		return "", fmt.Errorf("could not find repository prefix in zip file")
	}

	// Target path within the zip file (normalize path separators to forward slashes for zip compatibility)
	normalizedPath := strings.ReplaceAll(moduleRelativePath, "\\", "/")
	targetPathInZip := repoPrefix + normalizedPath + "/"

	// Create destination directory for the code
	codeDestDir := filepath.Join(dest, moduleRelativePath)
	os.MkdirAll(codeDestDir, 0755)
	extractedPath = codeDestDir

	// Second pass: extract only files under the module path
	for _, f := range r.File {
		// Check if this file is within our target module path
		if !strings.HasPrefix(f.Name, targetPathInZip) {
			continue
		}

		// Remove the repo prefix and module path prefix to get relative path within module
		relativePath := strings.TrimPrefix(f.Name, targetPathInZip)
		if relativePath == "" {
			continue // Skip the directory itself
		}

		// Create the file path in destination
		destPath := filepath.Join(codeDestDir, relativePath)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(destPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			// Create directory
			os.MkdirAll(destPath, f.FileInfo().Mode())
			continue
		}

		// Create the directories for file
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return "", err
		}

		// Extract file
		fileReader, err := f.Open()
		if err != nil {
			return "", err
		}

		targetFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			fileReader.Close()
			return "", err
		}

		_, err = io.Copy(targetFile, fileReader)
		targetFile.Close()
		fileReader.Close()

		if err != nil {
			return "", err
		}
	}

	return extractedPath, nil
}

// getChangelog generates the changelog report with the given two Contents
func getChangelog(lhs, rhs *exports.Content) (*Changelog, error) {
	if lhs == nil && rhs == nil {
		return nil, fmt.Errorf("nothing exists after the generation, this should never happen")
	}
	if lhs == nil {
		// the module does not exist before the generation: this is a new module
		return &Changelog{
			NewPackage: true,
		}, nil
	}
	if rhs == nil {
		// the module no longer exists after the generation: this module was removed
		return &Changelog{
			RemovedPackage: true,
		}, nil
	}
	// lhs and rhs are both non-nil
	p := report.Generate(*lhs, *rhs, nil)
	return &Changelog{
		Modified: &p,
	}, nil
}

// DetermineModuleStatus determines whether a module is new or existing
func DetermineModuleStatus(modulePath string, sdkRepo repo.SDKRepository) (utils.PackageStatus, error) {
	changelogPath := filepath.Join(modulePath, utils.ChangelogFileName)

	// Check if changelog exists
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		return utils.PackageStatusNew, nil
	}

	// Get all version tags for this module
	moduleRelativePath, err := utils.GetRelativePath(modulePath, sdkRepo)
	if err != nil {
		return utils.PackageStatusNew, err
	}
	tags, err := GetAllVersionTags(moduleRelativePath, sdkRepo)
	if err != nil {
		return utils.PackageStatusNew, err
	}

	if len(tags) == 0 {
		return utils.PackageStatusNew, nil
	}

	return utils.PackageStatusExisting, nil
}

type ChangelogResult struct {
	ChangelogData   *Changelog
	OriExports      *exports.Content
	NewExports      *exports.Content
	PreviousVersion string
}

// GenerateChangelog generates changelog by comparing exports and filtering unnecessary changes
func GenerateChangelog(modulePath string, sdkRepo repo.SDKRepository, isCurrentPreview bool) (ChangelogResult, error) {
	oriExports, previousVersion, err := getPreviousVersionAndExports(modulePath, sdkRepo, isCurrentPreview)
	if err != nil {
		return ChangelogResult{}, err
	}

	newExports, err := exports.Get(modulePath)
	if err != nil {
		return ChangelogResult{}, err
	}

	changelogData, err := getChangelog(oriExports, &newExports)
	if err != nil {
		return ChangelogResult{}, err
	}

	log.Printf("Filter changelog...")
	FilterChangelog(changelogData,
		NonExportedFilter,
		MarshalUnmarshalFilter,
		EnumFilter,
		FuncFilter,
		LROFilter,
		PageableFilter,
		InterfaceToAnyFilter,
		TypeToAnyFilter,
		ParamNameToUnderscoreFilter)

	return ChangelogResult{
		ChangelogData:   changelogData,
		OriExports:      oriExports,
		NewExports:      &newExports,
		PreviousVersion: previousVersion,
	}, nil
}

// getPreviousVersionAndExports gets the previous version and exports
func getPreviousVersionAndExports(modulePath string, sdkRepo repo.SDKRepository, isCurrentPreview bool) (*exports.Content, string, error) {
	moduleRelativePath, err := utils.GetRelativePath(modulePath, sdkRepo)
	if err != nil {
		return nil, "", err
	}

	tags, err := GetAllVersionTags(moduleRelativePath, sdkRepo)
	if err != nil {
		return nil, "", err
	}

	// New module
	if len(tags) == 0 {
		return nil, "", nil
	}

	previousVersionTag := getPreviousVersionTag(isCurrentPreview, tags)

	oriExports, err := getExportsFromTag(moduleRelativePath, previousVersionTag)
	if err != nil && !strings.Contains(err.Error(), "doesn't contain any exports") {
		return nil, "", err
	}

	tagSplit := strings.Split(previousVersionTag, "/")
	previousVersion := strings.TrimLeft(tagSplit[len(tagSplit)-1], "v")

	return oriExports, previousVersion, nil
}

// AddChangelogToFileWithReplacement adds changelog to file, replacing existing version if it exists
func AddChangelogToFileWithReplacement(changelog *Changelog, version *semver.Version, modulePath, releaseDate string) (string, error) {
	path := filepath.Join(modulePath, utils.ChangelogFileName)
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	oldChangelog := string(b)
	additionalChangelog := changelog.ToCompactMarkdown()
	if releaseDate == "" {
		releaseDate = time.Now().Format("2006-01-02")
	}

	// Look for existing version entry and replace it
	versionString := version.String()

	lines := strings.Split(oldChangelog, "\n")
	var newLines []string
	var skipLines bool
	var i int

	// Add the new header
	newLines = append(newLines, "# Release History", "")
	newLines = append(newLines, fmt.Sprintf("## %s (%s)", versionString, releaseDate))
	newLines = append(newLines, additionalChangelog)
	newLines = append(newLines, "")

	// Skip the old "# Release History" header and process the rest
	for i = 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "# Release History" {
			i++
			if i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				i++
			}
			break
		}
	}

	// Process remaining lines
	for ; i < len(lines); i++ {
		line := lines[i]

		// Check if this is a version header
		if strings.HasPrefix(strings.TrimSpace(line), "## ") && strings.Contains(line, versionString) {
			skipLines = true
			continue
		}

		// If we're skipping lines (inside the version we want to replace)
		if skipLines {
			// Stop skipping when we hit the next version header
			if strings.HasPrefix(strings.TrimSpace(line), "## ") {
				skipLines = false
				newLines = append(newLines, line)
			}
			// Otherwise, continue skipping
			continue
		}

		// Normal line, add it
		newLines = append(newLines, line)
	}

	// If we didn't find the version, all content after "# Release History" has been added
	finalChangelog := strings.Join(newLines, "\n")

	if err = os.WriteFile(path, []byte(finalChangelog), 0644); err != nil {
		return "", err
	}

	return additionalChangelog, nil
}

// UpdateLatestChangelogVersion updates the version of the latest (first) changelog entry
// This is used when setting a specific version - it keeps the changelog content but updates the version number and date
func UpdateLatestChangelogVersion(modulePath string, newVersion *semver.Version, releaseDate string) error {
	changelogPath := filepath.Join(modulePath, utils.ChangelogFileName)
	b, err := os.ReadFile(changelogPath)
	if err != nil {
		return fmt.Errorf("failed to read changelog: %v", err)
	}

	if releaseDate == "" {
		releaseDate = time.Now().Format("2006-01-02")
	}

	lines := strings.Split(string(b), "\n")
	var updatedLines []string
	foundFirstVersion := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Look for the first version header (e.g., "## 1.0.0 (2025-01-01)" or "## 1.0.0 (Unreleased)")
		if !foundFirstVersion && strings.HasPrefix(trimmed, "## ") {
			// Replace the version line with the new version
			updatedLines = append(updatedLines, fmt.Sprintf("## %s (%s)", newVersion.String(), releaseDate))
			foundFirstVersion = true
			continue
		}

		updatedLines = append(updatedLines, line)
	}

	if !foundFirstVersion {
		return fmt.Errorf("no version entry found in changelog")
	}

	finalChangelog := strings.Join(updatedLines, "\n")
	if err = os.WriteFile(changelogPath, []byte(finalChangelog), 0644); err != nil {
		return fmt.Errorf("failed to write changelog: %v", err)
	}

	return nil
}

// CreateNewChangelog creates a new changelog file for a new module
func CreateNewChangelog(modulePath string, sdkRepo repo.SDKRepository, packageVersion, releaseDate string) error {
	moduleRelativePath, err := utils.GetRelativePath(modulePath, sdkRepo)
	if err != nil {
		return err
	}

	if releaseDate == "" {
		releaseDate = time.Now().Format("2006-01-02")
	}

	content := fmt.Sprintf("# Release History\n\n## %s (%s)\n### Other Changes\n\nThe package of `github.com/Azure/azure-sdk-for-go/%s` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).\n\nTo learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).", packageVersion, releaseDate, moduleRelativePath)

	// Write the changelog file
	changelogPath := filepath.Join(modulePath, utils.ChangelogFileName)
	if err = os.WriteFile(changelogPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write changelog file: %v", err)
	}

	return nil
}

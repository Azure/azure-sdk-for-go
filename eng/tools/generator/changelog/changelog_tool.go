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

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	sdk_remote_url = "https://github.com/Azure/azure-sdk-for-go.git"
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
	_, err := sdkRepo.CreateRemote(&config.RemoteConfig{Name: remoteName, URLs: []string{sdk_remote_url}})
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

func ContainsPreviewAPIVersion(packagePath string) (bool, error) {
	log.Printf("Judge whether contains preview API version from '%s' ...", packagePath)

	files, err := os.ReadDir(packagePath)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" {
			b, err := os.ReadFile(filepath.Join(packagePath, file.Name()))
			if err != nil {
				return false, err
			}

			lines := strings.Split(string(b), "\n")
			for _, line := range lines {
				if strings.Contains(line, "\"api-version\"") {
					parts := strings.Split(line, "\"")
					if len(parts) == 5 && strings.Contains(parts[3], "preview") {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

func GetPreviousVersionTag(isCurrentPreview bool, allReleases []string) string {
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

func GetExportsFromTag(relativePackagePath, tag string) (*exports.Content, error) {
	log.Printf("Get exports for '%s' from specific tag '%s' ...", relativePackagePath, tag)

	// Extract tag name from ref/tags/ prefix if present
	tagName := strings.TrimPrefix(tag, "refs/tags/")

	// Download and extract only the packagePath folder
	extractedPackagePath, err := downloadAndExtractPackagePath(tagName, relativePackagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download and extract package path: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(extractedPackagePath)) // Clean up temp directory

	// Check if the package path exists in the extracted source
	if _, err := os.Stat(extractedPackagePath); os.IsNotExist(err) {
		log.Printf("Package path '%s' does not exist in tag '%s'", relativePackagePath, tagName)
		return &exports.Content{}, nil
	}

	// Get exports from the extracted package path
	result, err := exports.Get(extractedPackagePath)
	// bypass the error if the package doesn't contain any exports, return empty content
	if err != nil && strings.Contains(err.Error(), "doesn't contain any exports") {
		return &exports.Content{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get exports from extracted package: %v", err)
	}

	return &result, nil
}

// downloadAndExtractPackagePath downloads the source code for a specific tag and extracts only the packagePath folder
func downloadAndExtractPackagePath(tagName, relativePackagePath string) (string, error) {
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

	// Extract only the packagePath folder from the zip file
	extractedPackagePath, err := extractPackagePathFromZip(zipFilePath, tempDir, relativePackagePath)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to extract package path from zip file: %v", err)
	}

	// Remove the zip file after extraction
	os.Remove(zipFilePath)

	return extractedPackagePath, nil
}

// extractPackagePathFromZip extracts only the specified packagePath from the zip file
func extractPackagePathFromZip(src, dest, relativePackagePath string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var extractedPackagePath string
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
	normalizedPackagePath := strings.ReplaceAll(relativePackagePath, "\\", "/")
	targetPathInZip := repoPrefix + normalizedPackagePath + "/"

	// Create destination directory for the package
	packageDestDir := filepath.Join(dest, relativePackagePath)
	os.MkdirAll(packageDestDir, 0755)
	extractedPackagePath = packageDestDir

	// Second pass: extract only files under the packagePath
	for _, f := range r.File {
		// Check if this file is within our target package path
		if !strings.HasPrefix(f.Name, targetPathInZip) {
			continue
		}

		// Remove the repo prefix and package path prefix to get relative path within package
		relativePath := strings.TrimPrefix(f.Name, targetPathInZip)
		if relativePath == "" {
			continue // Skip the directory itself
		}

		// Create the file path in destination
		destPath := filepath.Join(packageDestDir, relativePath)

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

	return extractedPackagePath, nil
}

// GetChangelogForPackage generates the changelog report with the given two Contents
func GetChangelogForPackage(lhs, rhs *exports.Content) (*Changelog, error) {
	if lhs == nil && rhs == nil {
		return nil, fmt.Errorf("this package does not exist even after the generation, this should never happen")
	}
	if lhs == nil {
		// the package does not exist before the generation: this is a new package
		return &Changelog{
			NewPackage: true,
		}, nil
	}
	if rhs == nil {
		// the package no longer exists after the generation: this package was removed
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

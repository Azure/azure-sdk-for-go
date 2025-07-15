// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"go/ast"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	sdk_tag_fetch_url = "https://api.github.com/repos/Azure/azure-sdk-for-go/git/refs/tags"
	sdk_remote_url    = "https://github.com/Azure/azure-sdk-for-go.git"
)

func GetAllVersionTags(moduleRelativePath string) ([]string, error) {
	arr := strings.Split(moduleRelativePath, "/")
	log.Printf("Fetching all release tags from GitHub for RP: '%s' Package: '%s' ...", arr[len(arr)-2], arr[len(arr)-1])
	client := http.Client{}
	res, err := client.Get(sdk_tag_fetch_url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := []map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %s, error: %v", string(body), err)
		return nil, err
	}
	var tags []string
	var versions []string
	versionTag := make(map[string]string)
	for _, tag := range result {
		tagName := tag["ref"].(string)
		if strings.Contains(tagName, moduleRelativePath+"/v") {
			m := regexp.MustCompile(semver.SemVerRegex).FindString(tagName)
			versions = append(versions, m)
			versionTag[m] = tagName
		}
	}

	vs := make([]*semver.Version, len(versions))
	for i, r := range versions {
		v, err := semver.NewVersion(r)
		if err != nil {
			return nil, err
		}

		vs[i] = v
	}
	sort.Sort(sort.Reverse(semver.Collection(vs)))

	for _, v := range vs {
		tags = append(tags, versionTag[v.Original()])
	}

	return tags, nil
}

func GetAllVersionTagsV2(moduleRelativePath string, sdkRepo repo.SDKRepository) ([]string, error) {
	arr := strings.Split(moduleRelativePath, "/")
	log.Printf("Fetching all release tags from GitHub for RP: '%s' Package: '%s' ...", arr[len(arr)-2], arr[len(arr)-1])

	remoteName := "release_remote"
	fetchOpts := &git.FetchOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{"refs/tags/*:refs/tags/*"},
		Tags:       git.AllTags,
	}

	err := FetchTagsFromRemote(sdkRepo, remoteName, fetchOpts)
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
	log.Printf("Get exports from specific tag '%s' ...", tag)

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

	// Target path within the zip file
	targetPathInZip := repoPrefix + relativePackagePath + "/"

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

func MarshalUnmarshalFilter(changelog *Changelog) {
	if changelog.Modified != nil {
		if changelog.Modified.AdditiveChanges != nil {
			removeMarshalUnmarshalFunc(changelog.Modified.AdditiveChanges.Added.Funcs)
		}
		if changelog.Modified.BreakingChanges != nil && changelog.Modified.BreakingChanges.Removed != nil {
			removeMarshalUnmarshalFunc(changelog.Modified.BreakingChanges.Removed.Funcs)
		}
	}
}

func removeMarshalUnmarshalFunc(funcs map[string]exports.Func) {
	for k := range funcs {
		if strings.HasSuffix(k, ".MarshalJSON") || strings.HasSuffix(k, ".UnmarshalJSON") {
			delete(funcs, k)
		}
	}
}

func FilterChangelog(changelog *Changelog, opts ...func(changelog *Changelog)) {
	if changelog.Modified != nil {
		for _, opt := range opts {
			opt(changelog)
		}
	}
}

func EnumFilter(changelog *Changelog) {
	if changelog.Modified.HasAdditiveChanges() {
		if changelog.Modified.AdditiveChanges != nil && changelog.Modified.AdditiveChanges.Added.TypeAliases != nil {
			for typeAliases := range changelog.Modified.AdditiveChanges.Added.TypeAliases {
				funcKeys, funcExist := searchKey(changelog.Modified.AdditiveChanges.Added.Funcs, typeAliases, "Possible")
				if funcExist && len(funcKeys) == 1 {
					for _, f := range funcKeys {
						delete(changelog.Modified.AdditiveChanges.Added.Funcs, f)
					}
				}
			}
		}
	}

	if changelog.Modified.HasBreakingChanges() {
		enumOperation(changelog.Modified.BreakingChanges.Removed)
	}
}

func enumOperation(content *delta.Content) {
	if content != nil && content.TypeAliases != nil {
		for typeAliases := range content.TypeAliases {
			constKeys, constExist := searchKey(content.Consts, typeAliases, "")
			funcKeys, funcExist := searchKey(content.Funcs, typeAliases, "Possible")

			if constExist && funcExist && len(funcKeys) == 1 {
				for _, c := range constKeys {
					delete(content.Consts, c)
				}
				for _, f := range funcKeys {
					delete(content.Funcs, f)
				}
			}
		}
	}
}

func searchKey[T exports.Const | exports.Func | exports.Struct](m map[string]T, key1, prefix string) ([]string, bool) {
	keys := make([]string, 0)
	for k := range m {
		if regexp.MustCompile(fmt.Sprintf("^%s%s\\w*", prefix, key1)).MatchString(k) {
			keys = append(keys, k)
		}
	}
	if len(keys) != 0 {
		return keys, true
	}
	return nil, false
}

func FuncFilter(changelog *Changelog) {
	if changelog.Modified.HasAdditiveChanges() {
		funcOperation(changelog.Modified.AdditiveChanges.Added)
	}

	if changelog.Modified.HasBreakingChanges() {
		funcOperation(changelog.Modified.BreakingChanges.Removed)

		// function operation parameters from interface{} to any is not a breaking change
		for f, v := range changelog.Modified.BreakingChanges.Funcs {
			from := strings.Split(v.Params.From, ",")
			to := strings.Split(v.Params.To, ",")
			if len(from) != len(to) {
				continue
			}

			flag := false
			for i := range from {
				if strings.TrimSpace(from[i]) != strings.TrimSpace(to[i]) {
					if strings.TrimSpace(from[i]) == "interface{}" && strings.TrimSpace(to[i]) == "any" {
						flag = true
					} else {
						flag = false
						break
					}
				}
			}

			if flag {
				delete(changelog.Modified.BreakingChanges.Funcs, f)
			}
		}
	}
}

func funcOperation(content *delta.Content) {
	if content != nil && content.Funcs != nil {
		for funcName, funcValue := range content.Funcs {
			clientFunc := strings.Split(funcName, ".")
			if len(clientFunc) == 2 {
				// the last parameter
				if funcValue.Params != nil {
					ps := strings.Split(*funcValue.Params, ",")
					clientFuncOptions := ps[len(ps)-1]
					clientFuncOptions = strings.TrimLeft(strings.TrimSpace(clientFuncOptions), "*")
					if clientFuncOptions != "" && content.CompleteStructs != nil {
						delete(content.Structs, clientFuncOptions)
						for i, v := range content.CompleteStructs {
							if v == clientFuncOptions {
								content.CompleteStructs = append(content.CompleteStructs[:i],
									content.CompleteStructs[i+1:]...)
								break
							}
						}
					}
				}

				// the first return value
				if funcValue.Returns != nil {
					rs := strings.Split(*funcValue.Returns, ",")
					clientFuncResponse := rs[0]
					if strings.Contains(clientFuncResponse, "runtime") {
						re := regexp.MustCompile(`\[(?P<response>.*)\]`)
						clientFuncResponse = re.FindString(clientFuncResponse)
						clientFuncResponse = re.ReplaceAllString(clientFuncResponse, "${response}")
					} else {
						clientFuncResponse = strings.TrimLeft(clientFuncResponse, "*")
					}
					if clientFuncResponse != "" && content.CompleteStructs != nil {
						delete(content.Structs, clientFuncResponse)
						for i, v := range content.CompleteStructs {
							if v == clientFuncResponse {
								content.CompleteStructs = append(content.CompleteStructs[:i],
									content.CompleteStructs[i+1:]...)
								break
							}
						}
					}
				}
			}
		}
	}
}

// LROFilter LROFilter after OperationFilter
func LROFilter(changelog *Changelog) {
	if changelog.Modified.HasBreakingChanges() && changelog.Modified.HasAdditiveChanges() && changelog.Modified.BreakingChanges.Removed != nil && changelog.Modified.BreakingChanges.Removed.Funcs != nil {
		removedContent := changelog.Modified.BreakingChanges.Removed
		for bFunc, v := range removedContent.Funcs {
			var beginFunc string
			clientFunc := strings.Split(bFunc, ".")
			if len(clientFunc) == 2 {
				if strings.Contains(clientFunc[1], "Begin") {
					clientFunc[1] = strings.TrimPrefix(clientFunc[1], "Begin")
					beginFunc = fmt.Sprintf("%s.%s", clientFunc[0], clientFunc[1])
				} else {
					beginFunc = fmt.Sprintf("%s.Begin%s", clientFunc[0], clientFunc[1])
				}
				if _, ok := changelog.Modified.AdditiveChanges.Added.Funcs[beginFunc]; ok {
					delete(changelog.Modified.AdditiveChanges.Added.Funcs, beginFunc)
					v.ReplacedBy = &beginFunc
					removedContent.Funcs[bFunc] = v
				}
			}
		}
	}
}

// PageableFilter PageableFilter after OperationFilter
func PageableFilter(changelog *Changelog) {
	if changelog.Modified.HasBreakingChanges() && changelog.Modified.HasAdditiveChanges() && changelog.Modified.BreakingChanges.Removed != nil && changelog.Modified.BreakingChanges.Removed.Funcs != nil {
		removedContent := changelog.Modified.BreakingChanges.Removed
		for bFunc, v := range removedContent.Funcs {
			var pagination string
			clientFunc := strings.Split(bFunc, ".")
			if len(clientFunc) == 2 {
				if strings.Contains(clientFunc[1], "New") && strings.Contains(clientFunc[1], "Pager") {
					clientFunc[1] = strings.TrimPrefix(strings.TrimSuffix(clientFunc[1], "Pager"), "New")
					pagination = fmt.Sprintf("%s.%s", clientFunc[0], clientFunc[1])
				} else {
					pagination = fmt.Sprintf("%s.New%sPager", clientFunc[0], clientFunc[1])
				}
				if _, ok := changelog.Modified.AdditiveChanges.Added.Funcs[pagination]; ok {
					delete(changelog.Modified.AdditiveChanges.Added.Funcs, pagination)
					v.ReplacedBy = &pagination
					removedContent.Funcs[bFunc] = v
				}
			}
		}
	}
}

func InterfaceToAnyFilter(changelog *Changelog) {
	if changelog.HasBreakingChanges() {
		for structName, s := range changelog.Modified.BreakingChanges.Structs {
			for k, v := range s.Fields {
				if strings.Contains(v.From, "interface{}") && strings.Contains(v.To, "any") {
					delete(s.Fields, k)
				}
			}

			if len(s.Fields) == 0 {
				delete(changelog.Modified.BreakingChanges.Structs, structName)
			}
		}
	}
}

func NonExportedFilter(changelog *Changelog) {
	if !changelog.Modified.IsEmpty() {
		if changelog.Modified.HasAdditiveChanges() {
			nonExportOperation(changelog.Modified.AdditiveChanges.Added)
		}

		if changelog.Modified.HasBreakingChanges() {
			breakingChanges := changelog.Modified.BreakingChanges
			for fName := range breakingChanges.Funcs {
				before, after, _ := strings.Cut(fName, ".")
				if !ast.IsExported(strings.TrimLeft(before, "*")) || (after != "" && !ast.IsExported(after)) {
					delete(changelog.Modified.BreakingChanges.Funcs, fName)
				}
			}

			for sName := range breakingChanges.Structs {
				if !ast.IsExported(sName) {
					delete(changelog.Modified.BreakingChanges.Structs, sName)
				}
			}

			if breakingChanges.Removed != nil && !breakingChanges.Removed.IsEmpty() {
				nonExportOperation(breakingChanges.Removed)
			}

		}
	}
}

func nonExportOperation(content *delta.Content) {
	if content.IsEmpty() {
		return
	}

	for fName := range content.Funcs {
		before, after, _ := strings.Cut(fName, ".")
		if !ast.IsExported(strings.TrimLeft(before, "*")) || (after != "" && !ast.IsExported(after)) {
			delete(content.Funcs, fName)
		}
	}

	for sName := range content.Structs {
		if !ast.IsExported(sName) {
			delete(content.Structs, sName)
		}
	}
}

func TypeToAnyFilter(changelog *Changelog) {
	if changelog.Modified.HasBreakingChanges() {
		for structName, s := range changelog.Modified.BreakingChanges.Changes.Structs {
			for k, v := range s.Fields {
				if v.To == "any" {
					delete(s.Fields, k)
					if changelog.Modified.AdditiveChanges == nil {
						changelog.Modified.AdditiveChanges = &report.AdditiveChanges{}
					}
					if changelog.Modified.AdditiveChanges.Changes.Structs == nil {
						changelog.Modified.AdditiveChanges.Changes.Structs = map[string]delta.StructDef{}
					}
					if _, ok := changelog.Modified.AdditiveChanges.Changes.Structs[structName]; !ok {
						changelog.Modified.AdditiveChanges.Changes.Structs[structName] = delta.StructDef{
							Fields: make(map[string]delta.Signature),
						}
					}
					changelog.Modified.AdditiveChanges.Changes.Structs[structName].Fields[k] = v
				}
			}
			if len(s.Fields) == 0 {
				delete(changelog.Modified.BreakingChanges.Structs, structName)
			}
		}
	}
}

func FetchTagsFromRemote(sdkRepo repo.SDKRepository, remoteName string, fetchOpts *git.FetchOptions) error {
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

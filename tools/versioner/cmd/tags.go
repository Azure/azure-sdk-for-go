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

package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/Azure/azure-sdk-for-go/tools/internal/modinfo"
	"github.com/Masterminds/semver"
)

// returns the appropriate module tag based on the package version info
// tags - list of all current tags for the module
func calculateModuleTag(baseline string, versionSetting *VersionSetting, mod modinfo.Provider, hookFunc TagsHookFunc) (string, error) {
	if mod.IsPreviewPackage() {
		return calculateTagForPreview(baseline, versionSetting.InitialVersionPreview, mod, hookFunc)
	}
	return calculateTagForStable(baseline, versionSetting.InitialVersion, mod, hookFunc)
}

func calculateTagForStable(baseline, initialStartVersion string, mod modinfo.Provider, hookFunc TagsHookFunc) (string, error) {
	if mod.IsPreviewPackage() {
		return "", errors.New("package is not stable package and should not reach this function")
	}
	if mod.BreakingChanges() && !mod.VersionSuffix() {
		return "", errors.New("package has breaking changes but directory has no version suffix")
	}
	tagPrefix, err := getTagPrefix(baseline)
	if err != nil {
		return "", fmt.Errorf("failed to get tag prefix: %+v", err)
	}
	tags, err := hookFunc(baseline, tagPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to list tags: %+v", err)
	}

	latestVersion, err := getLatestSemver(tags, tagPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %+v", err)
	}

	if latestVersion == nil {
		// this is the first module version
		if !mod.NewModule() {
			return "", fmt.Errorf("module is a not new module but no tags were found")
		}
		return fmt.Sprintf("%s/%s", tagPrefix, initialStartVersion), nil
	}

	// if this has breaking changes then it's simply the prefix as a new major version
	if mod.BreakingChanges() {
		return fmt.Sprintf("%s/v%s", tagPrefix, latestVersion.IncMajor().String()), nil
	}

	if mod.NewExports() {
		return fmt.Sprintf("%s/v%s", tagPrefix, latestVersion.IncMinor().String()), nil
	}

	return fmt.Sprintf("%s/v%s", tagPrefix, latestVersion.IncPatch().String()), nil
}

func calculateTagForPreview(baseline, initialStartVersion string, mod modinfo.Provider, hookFunc TagsHookFunc) (string, error) {
	if !mod.IsPreviewPackage() {
		return "", errors.New("package is not preview package and should not reach this function")
	}
	// preview module should not have version suffix
	if mod.VersionSuffix() {
		return "", errors.New("preview module should not have version suffix")
	}
	tagPrefix, err := getTagPrefix(baseline)
	if err != nil {
		return "", fmt.Errorf("failed to get tag prefix: %+v", err)
	}
	tags, err := hookFunc(baseline, tagPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to list tags: %+v", err)
	}

	latestVersion, err := getLatestSemver(tags, tagPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %+v", err)
	}

	if latestVersion == nil {
		// this is the first module version
		if !mod.NewModule() {
			return "", fmt.Errorf("module is not a new module but no tags were found")
		}
		return fmt.Sprintf("%s/%s", tagPrefix, initialStartVersion), nil
	}

	// preview package does not bump major version even receiving breaking changes
	if mod.BreakingChanges() || mod.NewExports() {
		return fmt.Sprintf("%s/v%s", tagPrefix, latestVersion.IncMinor().String()), nil
	}

	return fmt.Sprintf("%s/v%s", tagPrefix, latestVersion.IncPatch().String()), nil
}

func getLatestSemver(tags []string, tagPrefix string) (*semver.Version, error) {
	var versions []*semver.Version
	for _, tag := range tags {
		index := strings.Index(tag, tagPrefix)
		if index < 0 {
			return nil, fmt.Errorf("do not find '%s' in tag '%s'", tagPrefix, tag)
		}
		verString := strings.Trim(tag[index + len(tagPrefix):], "/")
		ver, err := semver.NewVersion(verString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse semver %s: %+v", verString, err)
		}
		versions = append(versions, ver)
	}

	sort.Sort(semver.Collection(versions))
	if len(versions) == 0 {
		return nil, nil
	}
 	return versions[len(versions) - 1], nil
}

func findVersionSuffixInTag(tag string) string {
	r := verSuffixRegex.FindAllString(tag, -1)
	if len(r) == 0 {
		return ""
	}
	suffix := r[len(r) - 1]
	if suffix == "v1" {
		return ""
	}
	return suffix
}

// returns a slice of tags for the specified repo and tag prefix
func getTags(repoPath, tagPrefix string) ([]string, error) {
	wt, err := repo.Get(repoPath)
	if err != nil {
		return nil, err
	}
	return wt.ListTags(tagPrefix + "*")
}

// returns the tag prefix for the specified package.
// assumes repo root of github.com/Azure/azure-sdk-for-go/
func getTagPrefix(pkgDir string) (string, error) {
	// e.g. /work/src/github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis/v2
	// would return services/redis/mgmt/2018-03-01/redis/v2
	repoRoot := filepath.Join("github.com", repoOrg, repoName)
	i := strings.Index(pkgDir, repoRoot)
	if i < 0 {
		return "", fmt.Errorf("didn't find '%s' in '%s'", repoRoot, pkgDir)
	}
	return strings.Replace(pkgDir[i+len(repoRoot)+1:], "\\", "/", -1), nil
}

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
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/Azure/azure-sdk-for-go/tools/internal/modinfo"
	"github.com/Masterminds/semver"
	"path/filepath"
	"strings"
)

// returns the appropriate module tag based on the package version info
// tags - list of all current tags for the module
func calculateModuleTag(tags []string, mod modinfo.Provider) (string, error) {
	if mod.IsPreviewPackage() {
		return calculateTagForPreview(tags, mod)
	}
	return calculateTagForStable(tags, mod)
}

func calculateTagForStable(tags []string, mod modinfo.Provider) (string, error) {
	if mod.IsPreviewPackage() {
		return "", errors.New("package is not stable package and should not reach this function")
	}
	if mod.BreakingChanges() && !mod.VersionSuffix() {
		return "", errors.New("package has breaking changes but directory has no version suffix")
	}
	tagPrefix, err := getTagPrefix(mod.DestDir())
	if err != nil {
		return "", err
	}
	// if this has breaking changes then it's simply the prefix as a new major version
	if mod.BreakingChanges() {
		return tagPrefix + ".0.0", nil
	}
	if len(tags) == 0 {
		if mod.VersionSuffix() {
			panic("module contains a version suffix but no tags were found")
		}
		// this is the first module version
		return tagPrefix + "/" + startingModVer, nil
	}
	if !mod.VersionSuffix() {
		tagPrefix = tagPrefix + "/v1"
	}
	tag := tags[len(tags)-1]
	v := semverRegex.FindString(tag)
	if v == "" {
		return "", fmt.Errorf("didn't find semver in tag '%s'", tag)
	}
	sv, err := semver.NewVersion(v)
	if err != nil {
		return "", fmt.Errorf("failed to parse semver: %v", err)
	}
	// for non-breaking changes determine if this is a minor or patch update.
	if mod.NewExports() {
		// new exports, this is a minor update so bump minor version
		n := sv.IncMinor()
		sv = &n
	} else {
		// no new exports and has changes, this is a patch update
		n := sv.IncPatch()
		sv = &n
	}
	return strings.Replace(tag, v, "v"+sv.String(), 1), nil
}

func calculateTagForPreview(tags []string, mod modinfo.Provider) (string, error) {
	if !mod.IsPreviewPackage() {
		return "", errors.New("package is not preview package and should not reach this function")
	}
	// preview module should not have version suffix
	if mod.VersionSuffix() {
		return "", errors.New("preview module should not have version suffix")
	}
	tagPrefix, err := getTagPrefix(mod.DestDir())
	if err != nil {
		return "", err
	}
	// preview package do not bump major version even receiving breaking changes
	//if mod.BreakingChanges() {
	//	return tagPrefix + ".0.0", nil
	//}
	if len(tags) == 0 {
		// this is the first module version
		return tagPrefix + "/" + startingModVerPreview, nil
	}
	//if !mod.VersionSuffix() {
	//	tagPrefix = tagPrefix + "/v1"
	//}
	tag := tags[len(tags)-1]
	v := semverRegex.FindString(tag)
	if v == "" {
		return "", fmt.Errorf("didn't find semver in tag '%s'", tag)
	}
	sv, err := semver.NewVersion(v)
	if err != nil {
		return "", fmt.Errorf("failed to parse semver: %v", err)
	}
	// for any changes determine if this is a minor or patch update.
	if mod.BreakingChanges() || mod.NewExports() {
		// breaking changes or new exports, this is a minor update so bump minor version
		n := sv.IncMinor()
		sv = &n
	} else {
		// no new exports and has changes, this is a patch update
		n := sv.IncPatch()
		sv = &n
	}
	return strings.Replace(tag, v, "v"+sv.String(), 1), nil
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
	// would return services/redis/mgmt/2018-03-01/redis/v2.0.0
	repoRoot := filepath.Join("github.com", "Azure", "azure-sdk-for-go")
	i := strings.Index(pkgDir, repoRoot)
	if i < 0 {
		return "", fmt.Errorf("didn't find '%s' in '%s'", repoRoot, pkgDir)
	}
	return strings.Replace(pkgDir[i+len(repoRoot)+1:], "\\", "/", -1), nil
}

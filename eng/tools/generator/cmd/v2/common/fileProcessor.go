// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Masterminds/semver"
)

const (
	sdk_generated_file_prefix         = "zz_generated_"
	autorest_md_file_suffix           = "readme.md"
	autorest_md_module_version_prefix = "module-version: "
	swagger_md_module_name_prefix     = "module-name: "
)

var (
	v2BeginRegex             = regexp.MustCompile("^```\\s*yaml\\s*\\$\\(go\\)\\s*&&\\s*\\$\\((track2|v2)\\)")
	v2EndRegex               = regexp.MustCompile("^\\s*```\\s*$")
	newClientMethodNameRegex = regexp.MustCompile("^New.+Client$")
)

// reads from readme.go.md, parses the `track2` section to get module and package name
func ReadV2ModuleNameToGetNamespace(path string) (map[string][]string, error) {
	result := make(map[string][]string)
	log.Printf("Reading from readme.go.md '%s'...", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	log.Printf("Parsing module and package name from readme.go.md ...")
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")

	var start []int
	var end []int
	for i, line := range lines {
		if v2BeginRegex.MatchString(line) {
			start = append(start, i)
		}
		if len(start) != len(end) && v2EndRegex.MatchString(line) {
			end = append(end, i)
		}
	}

	if len(start) == 0 {
		return nil, fmt.Errorf("cannot find any `track2` section")
	}
	if len(start) != len(end) {
		return nil, fmt.Errorf("last `track2` section does not properly end")
	}

	for i := range start {
		// get the content of the `track2` section
		section := lines[start[i]+1 : end[i]]
		// iterate over the rest lines, get module name
		for _, line := range section {
			if strings.HasPrefix(line, swagger_md_module_name_prefix) {
				modules := strings.Split(strings.TrimSpace(line[len(swagger_md_module_name_prefix):]), "/")
				if len(modules) != 4 {
					return nil, fmt.Errorf("cannot parse module name from `track2` section")
				}
				namespaceName := strings.TrimSuffix(strings.TrimSuffix(modules[3], "\n"), "\r")
				log.Printf("RP: %s Package: %s", modules[2], namespaceName)
				result[modules[2]] = append(result[modules[2]], namespaceName)
			}
		}
	}

	return result, nil
}

// remove all sdk generated files in given path
func CleanSDKGeneratedFiles(path string) error {
	log.Printf("Removing all sdk generated files in '%s'...", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), sdk_generated_file_prefix) {
			err = os.Remove(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// replace repo commit with local path in autorest.md files
func ChangeConfigWithLocalPath(path, specPath, specRPName string) error {
	log.Printf("Replacing repo commit with local path in autorest.md ...")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	specPath, err = filepath.Abs(specPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if strings.Contains(line, autorest_md_file_suffix) {
			lines[i] = fmt.Sprintf("- %s", filepath.Join(specPath, "specification", specRPName, "resource-manager", "readme.md"))
			lines[i+1] = fmt.Sprintf("- %s", filepath.Join(specPath, "specification", specRPName, "resource-manager", "readme.go.md"))
			break
		}
	}

	return ioutil.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

// replace repo URL and commit id in autorest.md files
func ChangeConfigWithCommitID(path, repoURL, commitID, specRPName string) error {
	log.Printf("Replacing repo URL and commit id in autorest.md ...")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if strings.Contains(line, autorest_md_file_suffix) {
			lines[i] = fmt.Sprintf("- %s/blob/%s/specification/%s/resource-manager/readme.md", repoURL, commitID, specRPName)
			lines[i+1] = fmt.Sprintf("- %s/blob/%s/specification/%s/resource-manager/readme.go.md", repoURL, commitID, specRPName)
			break
		}
	}

	return ioutil.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

// get latest version according to `module-version: ` prefix in autorest.md file
func GetLatestVersion(packageRootPath string) (*semver.Version, error) {
	b, err := ioutil.ReadFile(filepath.Join(packageRootPath, "autorest.md"))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, autorest_md_module_version_prefix) {
			versionString := strings.TrimSuffix(strings.TrimSuffix(line[len(autorest_md_module_version_prefix):], "\n"), "\r")
			return semver.NewVersion(versionString)
		}
	}

	return nil, fmt.Errorf("cannot parse version from autorest.md")
}

// replace version according to `module-version: ` prefix in autorest.md file
func ReplaceVersion(packageRootPath string, newVersion string) error {
	path := filepath.Join(packageRootPath, "autorest.md")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, autorest_md_module_version_prefix) {
			lines[i] = line[:len(autorest_md_module_version_prefix)] + newVersion
			break
		}
	}

	return ioutil.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

// calculate new version by changelog using semver package
func CalculateNewVersion(changelog *model.Changelog, packageRootPath string) (*semver.Version, error) {
	version, err := GetLatestVersion(packageRootPath)
	if err != nil {
		return nil, err
	}
	log.Printf("Lastest version is: %s", version.String())

	var newVersion semver.Version
	if version.Major() == 0 {
		// preview version calculation
		if changelog.HasBreakingChanges() {
			newVersion = version.IncMinor()
		} else {
			newVersion = version.IncPatch()
		}
	} else {
		// release version calculation
		if changelog.HasBreakingChanges() {
			newVersion = version.IncMajor()
		} else if changelog.Modified.HasAdditiveChanges() {
			newVersion = version.IncMinor()
		} else {
			newVersion = version.IncPatch()
		}
	}

	log.Printf("New version is: %s", newVersion.String())
	return &newVersion, nil
}

// add new changelog md to changelog file
func AddChangelogToFile(changelog *model.Changelog, version *semver.Version, packageRootPath, releaseDate string) (string, error) {
	path := filepath.Join(packageRootPath, common.ChangelogFilename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	oldChangelog := string(b)
	insertPos := strings.Index(oldChangelog, "##")
	additionalChangelog := changelog.ToCompactMarkdown()
	if releaseDate == "" {
		releaseDate = time.Now().Format("2006-01-02")
	}
	newChangelog := oldChangelog[:insertPos] + "## " + version.String() + " (" + releaseDate + ")\n" + additionalChangelog + "\n\n" + oldChangelog[insertPos:]
	err = ioutil.WriteFile(path, []byte(newChangelog), 0644)
	if err != nil {
		return "", err
	}
	return additionalChangelog, nil
}

// replace `{{NewClientName}}`` placeholder in README.md by first func name according to `^New.+Method$` pattern
func ReplaceNewClientNamePlaceholder(packageRootPath string, exports exports.Content) error {
	path := filepath.Join(packageRootPath, "README.md")
	var clientName string
	for k, v := range exports.Funcs {
		if newClientMethodNameRegex.MatchString(k) && *v.Params == "string, azcore.TokenCredential, *arm.ClientOptions" {
			clientName = k
			break
		}
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read from file '%s': %+v", path, err)
	}

	var content = strings.ReplaceAll(string(b), "{{NewClientName}}", clientName)
	return ioutil.WriteFile(path, []byte(content), 0644)
}

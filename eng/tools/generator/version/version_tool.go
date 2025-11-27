// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package version

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Masterminds/semver"
	"golang.org/x/tools/go/ast/astutil"
)

const (
	autorestMdModuleVersionPrefix = "module-version: "
)

var (
	versionLineRegex = regexp.MustCompile(`moduleVersion\s*=\s*\".*v.+"`)
)

// UpdateAllVersionFiles updates all version-related files in the package
// This includes:
// - autorest.md (if exists, for swagger-based packages)
// - version.go
// - go.mod (for v2+ modules)
// - README.md (module path for v2+)
// - import (for v2+ modules)
func UpdateAllVersionFiles(modulePath string, version *semver.Version, sdkRepo repo.SDKRepository) error {
	// Update autorest.md if it exists (swagger-based packages)
	autorestMdPath := filepath.Join(modulePath, "autorest.md")
	if _, err := os.Stat(autorestMdPath); err == nil {
		if err := UpdateAutorestMdVersion(autorestMdPath, version.String()); err != nil {
			return fmt.Errorf("failed to update autorest.md: %v", err)
		}
	}

	// Update version.go
	if err := UpdateVersionGoFile(modulePath, version); err != nil {
		return fmt.Errorf("failed to update version.go: %v", err)
	}

	if version.Major() > 1 {
		// Update go.mod for v2+ modules
		if err := UpdateModuleDefinition(modulePath, version, sdkRepo); err != nil {
			return fmt.Errorf("failed to update go.mod: %v", err)
		}

		// Update README.md module path for v2+
		if err := UpdateReadmeModule(modulePath, version, sdkRepo); err != nil {
			return fmt.Errorf("failed to update README.md: %v", err)
		}

		// Update import paths for v2+
		if err := UpdateImportPaths(modulePath, version, sdkRepo); err != nil {
			return fmt.Errorf("failed to update import paths: %v", err)
		}
	}

	return nil
}

// UpdateAutorestMdVersion updates the module version in autorest.md file
func UpdateAutorestMdVersion(autorestMdPath, newVersion string) error {
	log.Printf("Updating autorest.md version to %s...", newVersion)

	b, err := os.ReadFile(autorestMdPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, autorestMdModuleVersionPrefix) {
			lines[i] = autorestMdModuleVersionPrefix + newVersion
			break
		}
	}

	return os.WriteFile(autorestMdPath, []byte(strings.Join(lines, "\n")), 0644)
}

// UpdateVersionGoFile updates the moduleVersion const in version.go file
func UpdateVersionGoFile(modulePath string, version *semver.Version) error {
	log.Printf("Updating version.go to %s...", version.String())

	path := filepath.Join(modulePath, "version.go")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	contents := versionLineRegex.ReplaceAllString(string(data), "moduleVersion = \"v"+version.String()+"\"")
	if contents == string(data) {
		return nil
	}

	return os.WriteFile(path, []byte(contents), 0644)
}

// UpdateModuleDefinition updates module definition in go.mod file according to version
func UpdateModuleDefinition(modulePath string, version *semver.Version, sdkRepo repo.SDKRepository) error {
	log.Printf("Update module definition if v2+...")

	if version.Major() <= 1 {
		return nil
	}

	path := filepath.Join(modulePath, "go.mod")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	moduleRelativePath, err := utils.GetRelativePath(modulePath, sdkRepo)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read go.mod: %v", err)
	}

	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "module") {
			line = strings.TrimRight(line, "\r")
			parts := strings.Split(line, "/")
			if parts[len(parts)-1] != fmt.Sprintf("v%d", version.Major()) {
				lines[i] = fmt.Sprintf("module github.com/Azure/azure-sdk-for-go/%s/v%d", moduleRelativePath, version.Major())
			}
			break
		}
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

// UpdateReadmeModule updates the module path in README.md according to current version
func UpdateReadmeModule(modulePath string, version *semver.Version, sdkRepo repo.SDKRepository) error {
	log.Printf("Update README.md module path if v2+...")

	if version.Major() <= 1 {
		return nil
	}

	readmePath := filepath.Join(modulePath, "README.md")
	readmeFile, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	moduleRelativePath, err := utils.GetRelativePath(modulePath, sdkRepo)
	if err != nil {
		return err
	}

	module := fmt.Sprintf("github.com/Azure/azure-sdk-for-go/%s", moduleRelativePath)

	readmeModule := module
	match := regexp.MustCompile(fmt.Sprintf(`%s/v\d+`, module))
	if match.Match(readmeFile) {
		readmeModule = match.FindString(string(readmeFile))
	}

	newModule := fmt.Sprintf("%s/v%d", module, version.Major())

	if newModule == readmeModule {
		return nil
	}

	newReadmeFile := strings.ReplaceAll(string(readmeFile), readmeModule, newModule)
	return os.WriteFile(readmePath, []byte(newReadmeFile), 0644)
}

// UpdateImportPaths updates the import paths in Go source files according to the new version
func UpdateImportPaths(modulePath string, version *semver.Version, sdkRepo repo.SDKRepository) error {
	log.Printf("Update import paths if v2+...")

	if version.Major() <= 1 {
		return nil
	}

	relativePath, err := utils.GetRelativePath(modulePath, sdkRepo)
	if err != nil {
		return err
	}
	baseModule := fmt.Sprintf("%s/%s", "github.com/Azure/azure-sdk-for-go", relativePath)

	return filepath.WalkDir(modulePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".go") {
			if err = replaceImport(path, baseModule, version.Major()); err != nil {
				return err
			}
		}

		return nil
	})
}

func replaceImport(sourceFile string, baseModule string, majorVersion int64) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, sourceFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	rewrote := false
	for _, i := range f.Imports {
		if strings.HasPrefix(i.Path.Value, fmt.Sprintf("\"%s", baseModule)) {
			oldPath := importPath(i)
			after, _ := strings.CutPrefix(oldPath, baseModule)

			newPath := baseModule
			if after != "" {
				before, sub, _ := strings.Cut(strings.TrimLeft(after, "/"), "/")
				if regexp.MustCompile(`^v\d+$`).MatchString(before) {
					if majorVersion > 1 {
						newPath = fmt.Sprintf("%s/v%d", baseModule, majorVersion)
					}
					if sub != "" {
						newPath = fmt.Sprintf("%s/%s", newPath, sub)
					}
				} else {
					if majorVersion > 1 {
						newPath = fmt.Sprintf("%s/v%d", baseModule, majorVersion)
					}
					newPath = fmt.Sprintf("%s/%s", newPath, before)
					if sub != "" {
						newPath = fmt.Sprintf("%s/%s", newPath, sub)
					}
				}
			} else {
				if majorVersion > 1 {
					newPath = fmt.Sprintf("%s/v%d", baseModule, majorVersion)
				}
			}

			if newPath != oldPath {
				rewrote = astutil.RewriteImport(fset, f, oldPath, newPath)
			}
		}
	}

	if rewrote {
		w, err := os.Create(sourceFile)
		if err != nil {
			return err
		}
		defer w.Close()

		return printer.Fprint(w, fset, f)
	}

	return nil
}

func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err != nil {
		return ""
	}
	return t
}

// IsCurrentPreviewVersion determines if the current package version is a preview version
func IsCurrentPreviewVersion(modulePath string, sdkRepo repo.SDKRepository, override *bool) (bool, error) {
	if override != nil {
		return *override, nil
	}
	// Determine if current package contains preview API version
	isCurrentPreview, err := containsPreviewAPIVersion(modulePath)
	if err != nil {
		return false, err
	}
	return isCurrentPreview, nil
}

// containsPreviewAPIVersion checks if the package contains any preview API version calls
func containsPreviewAPIVersion(packagePath string) (bool, error) {
	log.Printf("Check if package contains preview API version from '%s' ...", packagePath)

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

// CalculateNewVersion calculates the new version based on the changelog
func CalculateNewVersion(changelog *changelog.Changelog, previousVersion string, isCurrentPreview bool) (*semver.Version, utils.PullRequestLabel, error) {
	version, err := semver.NewVersion(previousVersion)
	if err != nil {
		return nil, "", err
	}
	log.Printf("Latest version is: %s", version.String())

	var newVersion semver.Version
	var prl utils.PullRequestLabel
	if version.Major() == 0 {
		// preview version calculation
		if !isCurrentPreview {
			tempVersion, err := semver.NewVersion("1.0.0")
			if err != nil {
				return nil, "", err
			}
			newVersion = *tempVersion
			if changelog.HasBreakingChanges() {
				prl = utils.FirstGABreakingChangeLabel
			} else {
				prl = utils.FirstGALabel
			}
		} else if changelog.HasBreakingChanges() {
			newVersion = version.IncMinor()
			prl = utils.BetaBreakingChangeLabel
		} else if changelog.Modified.HasAdditiveChanges() {
			newVersion = version.IncMinor()
			prl = utils.BetaLabel
		} else {
			newVersion = version.IncPatch()
			prl = utils.BetaLabel
		}
	} else {
		if isCurrentPreview {
			if strings.Contains(previousVersion, "beta") {
				betaNumber, err := strconv.Atoi(strings.Split(version.Prerelease(), "beta.")[1])
				if err != nil {
					return nil, "", err
				}
				newVersion, err = version.SetPrerelease("beta." + strconv.Itoa(betaNumber+1))
				if err != nil {
					return nil, "", err
				}
				if changelog.HasBreakingChanges() {
					prl = utils.BetaBreakingChangeLabel
				} else {
					prl = utils.BetaLabel
				}
			} else {
				if changelog.HasBreakingChanges() {
					newVersion = version.IncMajor()
					prl = utils.BetaBreakingChangeLabel
				} else if changelog.Modified.HasAdditiveChanges() {
					newVersion = version.IncMinor()
					prl = utils.BetaLabel
				} else {
					newVersion = version.IncPatch()
					prl = utils.BetaLabel
				}
				newVersion, err = newVersion.SetPrerelease("beta.1")
				if err != nil {
					return nil, "", err
				}
			}
		} else {
			if strings.Contains(previousVersion, "beta") {
				return nil, "", fmt.Errorf("must have stable previous version")
			}
			// release version calculation
			if changelog.HasBreakingChanges() {
				newVersion = version.IncMajor()
				prl = utils.StableBreakingChangeLabel
			} else if changelog.Modified.HasAdditiveChanges() {
				newVersion = version.IncMinor()
				prl = utils.StableLabel
			} else {
				newVersion = version.IncPatch()
				prl = utils.StableLabel
			}
		}
	}

	log.Printf("New version is: %s", newVersion.String())
	return &newVersion, prl, nil
}

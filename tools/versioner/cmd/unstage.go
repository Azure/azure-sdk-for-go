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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/files"
	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/Azure/azure-sdk-for-go/tools/internal/modinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func unstageCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unstage <staging dir>",
		Short: "Creates or updates the latest major version for a package from staged content.",
		Long: `This tool will compare a staged package against its latest major version to detect
breaking changes.  If there are no breaking changes the latest major version is updated
with the staged content.  If there are breaking changes the staged content becomes the
next latest major version and the go.mod file is updated.
The default version for new modules is v1.0.0, and for preview modules is v0.0.0.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := args[0]
			startingVer := viper.GetString("starting-version")
			if !modinfo.IsValidModuleVersion(startingVer) {
				return fmt.Errorf("the string '%s' is not a valid module version", startingVer)
			}
			startingVerPreview := viper.GetString("starting-preview-version")
			if !modinfo.IsValidModuleVersion(startingVerPreview) {
				return fmt.Errorf("the string '%s' is not a valid module version", startingVerPreview)
			}
			repoRoot = viper.GetString("gomod-root")
			_, _, err := ExecuteUnstage(root, Flags{
				RepoRoot:       viper.GetString("gomod-root"),
				VersionSetting: &VersionSetting{
					InitialVersion:        startingVer,
					InitialVersionPreview: startingVerPreview,
				},
				GetTagsHook:    getTags,
			})
			return err
		},
	}
}

var (
	semverRegex    = regexp.MustCompile(`v\d+\.\d+\.\d+$`)
	versionGoRegex = regexp.MustCompile(`\d+\.\d+\.\d+`)
	verSuffixRegex = regexp.MustCompile(`v\d+`)
)

const (
	goModFilename   = "go.mod"
	changeLogName   = "CHANGELOG.md"
	versionFilename = "version.go"
	interfacesName  = "interfaces.go"
)

// TagsHookFunc is a func used for get tags from remote
type TagsHookFunc func(root string, tagPrefix string) ([]string, error)

// ExecuteUnstage executes the unstage command
func ExecuteUnstage(s string, flags Flags) (string, string, error) {
	flags.apply()
	stage, err := filepath.Abs(s)
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path from '%s': %+v", s, err)
	}
	// format the stage directory to avoid unexpected diff
	if err := formatCode(stage); err != nil {
		return "", "", fmt.Errorf("failed to format directory '%s': %+v", stage, err)
	}
	lmv, err := findLatestMajorVersion(stage)
	if err != nil {
		return "", "", fmt.Errorf("failed to find latest major version directory in '%s': %+v", stage, err)
	}
	log.Debugf("Latest major version directory: %s", lmv)
	log.Debug("Comparing exports in latest major directory and stage content...")
	mod, err := modinfo.GetModuleInfo(lmv, stage)
	if err != nil {
		return "", "", fmt.Errorf("failed to create module info: %+v", err)
	}

	// preview packages do not do side by side update, they always get updated in-placed
	var tag string
	if !mod.IsPreviewPackage() && mod.BreakingChanges() {
		tag, err = forSideBySideRelease(stage, mod)
	} else {
		// check if lmv and stage are identical
		log.Debugf("Checking if '%s' and '%s' are identical...", lmv, stage)
		if err := formatCode(lmv); err != nil {
			return "", "", fmt.Errorf("failed to format directory '%s': %+v", lmv, err)
		}
		identical, err2 := checkIdentical(lmv, stage)
		if err2 != nil {
			return "", "", fmt.Errorf("failed to check identical: %+v", err2)
		}
		if identical {
			tag, err = forIdenticalPackage(stage)
		} else {
			tag, err = forInPlaceUpdate(lmv, stage, mod)
		}
	}

	return mod.DestDir(), tag, err
}

func forSideBySideRelease(stage string, mod modinfo.Provider) (string, error) {
	log.Debug("This is a side by side update")
	// calculate module tag
	tag, err := calculateModuleTag(filepath.Dir(stage), versionSetting, repoRoot, mod, getTagsHook)
	if err != nil {
		return "", fmt.Errorf("failed to calculate module tag: %+v", err)
	}
	log.Debugf("Tag calculated for stage '%s' is %s", stage, tag)
	log.Debug("Updating go.mod...")
	if err := updateGoModFile(stage, tag); err != nil {
		return "", fmt.Errorf("failed to update go.mod file: %+v", err)
	}
	log.Debug("Updating import statements...")
	if err := updateImportStatement(stage, mod.DestDir()); err != nil {
		return "", fmt.Errorf("failed to replace import statement: %+v", err)
	}
	log.Debug("Updating version.go...")
	if err := updateVersionFile(stage, tag); err != nil {
		return "", fmt.Errorf("failed to update version.go file: %+v", err)
	}
	log.Debug("Writing CHANGELOG...")
	if err := writeChangelog(stage, mod); err != nil {
		return "", fmt.Errorf("failed to write CHANGELOG.md: %+v", err)
	}
	log.Debugf("Renaming stage directory to '%s'...", mod.DestDir())
	if err := os.Rename(stage, mod.DestDir()); err != nil {
		return "", fmt.Errorf("failed to rename '%s' to '%s': %+v", stage, mod.DestDir(), err)
	}
	return tag, nil
}

func forInPlaceUpdate(lmv, stage string, mod modinfo.Provider) (string, error) {
	log.Debug("This is a in-placed update")
	// calculate module tag
	tag, err := calculateModuleTag(filepath.Dir(stage), versionSetting, repoRoot, mod, getTagsHook)
	if err != nil {
		return "", fmt.Errorf("failed to calculate module tag: %+v", err)
	}
	log.Debugf("Tag calculated for stage '%s' is %s", stage, tag)
	log.Debug("Updating go.mod...")
	if err := updateGoModFile(stage, tag); err != nil {
		return "", fmt.Errorf("failed to update go.mod file: %+v", err)
	}
	log.Debug("Updating import statements...")
	if err := updateImportStatement(stage, mod.DestDir()); err != nil {
		return "", fmt.Errorf("failed to replace import statement: %+v", err)
	}
	log.Debug("Updating version.go...")
	if err := updateVersionFile(stage, tag); err != nil {
		return "", fmt.Errorf("failed to update version.go file: %+v", err)
	}
	log.Debug("Writing CHANGELOG...")
	if err := writeChangelog(stage, mod); err != nil {
		return "", fmt.Errorf("failed to write CHANGELOG.md: %+v", err)
	}
	log.Debugf("Overriding stage directory to '%s'...", lmv)
	if err := overrideLMVFromStageDirectory(lmv, stage); err != nil {
		return "", fmt.Errorf("failed to override stage to lmv '%s': %+v", lmv, err)
	}
	return tag, nil
}

// forIdenticalPackage will update the code in baseline directory
func forIdenticalPackage(stage string) (string, error) {
	log.Debug("Latest major version and stage content are identical")
	baseline := filepath.Dir(stage)
	tagPrefix, err := getTagPrefix(baseline, repoRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get tag prefix: %+v", err)
	}
	tags, err := getTagsHook(baseline, tagPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to list tags: %+v", err)
	}
	latestVersion, err := getLatestSemver(tags)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %+v", err)
	}

	if latestVersion == nil {
		return "", fmt.Errorf("lmv and stage are identical, but do not find any latest version of this module")
	}

	tag := fmt.Sprintf("%s/v%s", tagPrefix, latestVersion.String())
	// remove the stage directory
	if err := os.RemoveAll(stage); err != nil {
		return "", fmt.Errorf("failed to remove stage directory: %+v", err)
	}

	return tag, nil
}

func checkIdentical(baseline, stage string) (bool, error) {
	// first list all the files and directories in baseline and stage
	// and since stage should be a subdirectory of baseline, we need to escape stage when list all file and directories
	fileListInBaseline, err := listAllFiles(baseline)
	if err != nil {
		return false, err
	}
	// strip out the stage directory and all its children
	fileListInBaseline = escapeStage(fileListInBaseline, stage)
	fileListInBaseline = escapeSpecialFiles(fileListInBaseline)
	// no need to escape the major sub-directories. if we are comparing the v1 and stage, there must not be any major sub-directories.
	// on the other hand, if we are comparing one major sub-directory to stage, the other sub-directories should not be in this directory.
	fileListInStage, err := listAllFiles(stage)
	if err != nil {
		return false, err
	}
	fileListInStage = escapeSpecialFiles(fileListInStage)
	if len(fileListInBaseline) != len(fileListInStage) {
		return false, nil
	}
	// filepath.Walk follows the lexical order, therefore the two lists should all be in lexical order
	for i, file := range fileListInBaseline {
		fileInStage := fileListInStage[i]
		same, err := files.DeepCompare(file, fileInStage)
		if err != nil {
			return false, err
		}
		if !same {
			return false, nil
		}
	}
	return true, nil
}

func listAllFiles(root string) ([]string, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("the root path '%s' does not exist", root)
	}
	var results []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			results = append(results, path)
		}
		return nil
	})
	return results, err
}

// escape the stage directory and all its sub-directories from the list
func escapeStage(fileList []string, stage string) []string {
	var result []string
	for _, path := range fileList {
		if !strings.HasPrefix(path, stage) {
			result = append(result, path)
		}
	}
	return result
}

// escape the special files -- go.mod, changelog and version.go
func escapeSpecialFiles(fileList []string) []string {
	var result []string
	for _, path := range fileList {
		base := filepath.Base(path)
		if base != goModFilename && base != changeLogName && base != versionFilename && base != interfacesName {
			result = append(result, path)
		}
	}
	return result
}

func overrideLMVFromStageDirectory(baseline, stage string) error {
	// in our case, stage should always be a child of baseline
	// first move stage to a temp directory outside of baseline, then remove the whole baseline directory, and finally move temp back to baseline
	temp := filepath.Join(filepath.Dir(baseline), "temp")
	if err := os.Rename(stage, temp); err != nil {
		return fmt.Errorf("failed to rename '%s' to '%s': %+v", stage, temp, err)
	}
	if err := os.RemoveAll(baseline); err != nil {
		return fmt.Errorf("failed to delete '%s': %+v", baseline, err)
	}
	if err := os.Rename(temp, baseline); err != nil {
		return fmt.Errorf("failed to rename '%s' to '%s': %+v", temp, baseline, err)
	}
	return nil
}

func updateGoModFile(directory, tag string) error {
	path := filepath.Join(directory, goModFilename)
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open for read '%s': %+v", directory, err)
	}
	defer file.Close()
	// find if the new module will have a major version suffix (/v2, /v3 etc)
	ver := findVersionSuffixInTag(tag)
	if err := updateGoMod(file, ver); err != nil {
		return fmt.Errorf("failed to update go.mod file: %+v", err)
	}
	return nil
}

// here we only update the version number in version.go, the api-version in User-Agent method will be taken care of
// when this file is generated by autorest
func updateVersionFile(directory, tag string) error {
	version := semverRegex.FindString(tag)
	log.Debugf("Updating version.go file in %s with version %s", directory, version)
	// version.go file must exists
	file := filepath.Join(directory, versionFilename)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return errors.New("version.go file does not exist")
	}
	verFile, err := os.OpenFile(file, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to version.go file: %+v", err)
	}
	defer verFile.Close()
	scanner := bufio.NewScanner(verFile)
	scanner.Split(bufio.ScanLines)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	_, err = verFile.Seek(0, io.SeekStart) // set pointer to start of the file
	if err != nil {
		return fmt.Errorf("failed to update version.go file: %+v", err)
	}
	hasTag := false
	for _, line := range lines {
		if !strings.HasPrefix(line, "// ") && versionGoRegex.MatchString(line) {
			line = versionGoRegex.ReplaceAllString(line, version)
		}
		if strings.HasPrefix(line, "// tag: ") {
			line = fmt.Sprintf("// tag: %s", tag)
			hasTag = true
		}
		if _, err := fmt.Fprintln(verFile, line); err != nil {
			return err
		}
	}
	if !hasTag {
		if _, err := fmt.Fprintf(verFile, "\n// tag: %s\n", tag); err != nil {
			return err
		}
	}
	return nil
}

// returns the absolute path to the latest major version based on the provided staging directory.
// it's assumed that the staging directory is a subdirectory of the actual package directory.
func findLatestMajorVersion(stage string) (string, error) {
	// example input:
	// ~/work/src/github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis/stage
	// finds:
	// ~/work/src/github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis
	// ~/work/src/github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis/v2
	// returns:
	// ~/work/src/github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis/v2
	parent := filepath.Dir(stage)
	dirs, err := modinfo.GetModuleSubdirs(parent)
	if err != nil {
		return "", fmt.Errorf("failed to get module subdirs '%s': %v", parent, err)
	}
	// no dirs means this is a v1 package
	if len(dirs) == 0 {
		return parent, nil
	}
	sort.Strings(dirs)
	// last one in the slice is the largest
	v := filepath.Join(parent, dirs[len(dirs)-1])
	log.Infof("Latest major version: %v", v)
	return v, nil
}

// updates the module version inside the go.mod file
func updateGoMod(goMod io.ReadWriteSeeker, newVer string) error {
	if newVer == "" {
		return nil
	}
	lines := files.GetLines(goMod)

	if _, err := goMod.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to start: %v", err)
	}
	for _, line := range lines {
		if strings.Index(line, "module") > -1 {
			if modinfo.HasVersionSuffix(line) {
				line = strings.Replace(line, "/"+modinfo.FindVersionSuffix(line), "/"+newVer, 1)
			} else {
				line = line + "/" + newVer
			}
		}
		if _, err := fmt.Fprintln(goMod, line); err != nil {
			return fmt.Errorf("failed to write line: %s", line)
		}
	}
	return nil
}

// traversal all go source files in the stage folder, and replace the import statement with new ones
func updateImportStatement(stage, dest string) error {
	newImport, err := importPathFromAbsPath(dest, repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get import from '%s': %+v", dest, err)
	}
	baseline := filepath.Dir(stage)
	oldImport, err := importPathFromAbsPath(baseline, repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get import from '%s': %+v", baseline, err)
	}
	log.Debugf("Attempting to replace import statement from '%s' to '%s'", oldImport, newImport)
	goFiles, err := findAllGoSourceFiles(stage)
	if err != nil {
		return err
	}
	log.Debugf("Found %d go source files: \n%s", len(goFiles), strings.Join(goFiles, "\n"))
	for _, file := range goFiles {
		if err := replaceImportInFile(file, oldImport, newImport); err != nil {
			return fmt.Errorf("failed to replace import statement in file '%s': %+v", file, err)
		}
	}
	return nil
}

func importPathFromAbsPath(path, repoRoot string) (string, error) {
	path = strings.ReplaceAll(path, "\\", "/")
	index := strings.Index(path, repoRoot)
	if index < 0 {
		return "", fmt.Errorf("do not find '%s' in path '%s'", repoRoot, path)
	}
	return path[index:], nil
}

func findAllGoSourceFiles(path string) ([]string, error) {
	var fileList []string // fileList stores filenames for those content contained the given import statements
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			fileList = append(fileList, path)
			return nil
		}
		return nil
	})
	return fileList, err
}

func replaceImportInFile(filepath, oldContent, newContent string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := updateGoSourceFile(file, oldContent, newContent); err != nil {
		return fmt.Errorf("failed to update import statement in file '%s': %+v", filepath, err)
	}
	return nil
}

func updateGoSourceFile(file io.ReadWriteSeeker, oldImport, newImport string) error {
	lines := files.GetLines(file)
	content := strings.Join(lines, "\n")
	importStatements := findImportStatements(content)
	newImportStatements := strings.ReplaceAll(importStatements, oldImport, newImport)
	newFileContent := strings.ReplaceAll(content, importStatements, newImportStatements)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	for _, line := range strings.Split(newFileContent, "\n") {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return err
		}
	}
	return nil
}

func findImportStatements(content string) string {
	oneLineImport := regexp.MustCompile(`import ".*"`)
	if oneLineImport.MatchString(content) {
		return oneLineImport.FindString(content)
	}
	multiLineRegex := `import \(\n(\s*\".*\"\n)+\s*\)`
	multiLineImport := regexp.MustCompile(multiLineRegex)
	if multiLineImport.MatchString(content) {
		return multiLineImport.FindString(content)
	}
	return ""
}

func writeChangelog(stage string, mod modinfo.Provider) error {
	// don't write a changelog for a new module
	if mod.NewModule() {
		return nil
	}
	rpt := mod.GenerateReport()
	changelog, err := os.Create(filepath.Join(stage, changeLogName))
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", changeLogName, err)
	}
	defer changelog.Close()
	if rpt.IsEmpty() {
		_, err = changelog.WriteString("No changes to exported content compared to the previous release.\n")
		return err
	}
	_, err = changelog.WriteString(rpt.ToMarkdown())
	return err
}

func formatCode(folder string) error {
	c := exec.Command("gofmt", "-w", folder)
	if output, err := c.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}
	return nil
}

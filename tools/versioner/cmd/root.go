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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/delta"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/exports"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "versioner <staging dir>",
	Short: "Creates or updates the latest major version for a package from staged content.",
	Long: `This tool will compare a staged package against its latest major version to detect
breaking changes.  If there are no breaking changes the latest major version is updated
with the staged content.  If there are breaking changes the staged content becomes the
next latest major vesion and the go.mod file is updated.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return theCommand(args)
	},
}

var (
	lhsExports exports.Content
	rhsExports exports.Content
)

// Execute executes the specified command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func theCommand(args []string) error {
	stage := filepath.Clean(args[0])
	lmv, err := findLatestMajorVersion(stage)
	if err != nil {
		return fmt.Errorf("failed to find latest major version: %v", err)
	}
	if err = loadExports(lmv, stage); err != nil {
		return fmt.Errorf("failed to load exports: %v", err)
	}
	if err = writeChangelog(stage); err != nil {
		return fmt.Errorf("failed to write changelog: %v", err)
	}
	hasBreaks := hasBreakingChanges(lmv, stage)
	// check if the LMV has a v2+ directory suffix
	hasVer, err := regexp.MatchString(`v\d+$`, lmv)
	if err != nil {
		return fmt.Errorf("failed to check for major version suffix: %v", err)
	}
	dest := filepath.Dir(stage)
	if hasBreaks {
		// move staging to new LMV directory
		v := 2
		if hasVer {
			s := string(lmv[len(lmv)-1])
			v, err = strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("failed to convert '%s' to int: %v", s, err)
			}
			v++
		}
		// update the go.mod file with the new major version
		ver := fmt.Sprintf("v%d", v)
		if err = updateGoModVer(stage, ver); err != nil {
			return fmt.Errorf("failed to update go.mod file: %v", err)
		}
		dest = filepath.Join(dest, ver)
		if err = os.Rename(stage, dest); err != nil {
			return fmt.Errorf("failed to rename '%s' to '%s': %v", stage, dest, err)
		}
		var tag string
		if tag, err = calculateModuleTag(true, true, dest); err != nil {
			return fmt.Errorf("failed to calculate module tag: %v", err)
		}
		fmt.Printf("tag: %s\n", tag)
		return nil
	}
	// move staging directory over the LMV by first deleting LMV then renaming stage
	if hasVer {
		if err = os.RemoveAll(lmv); err != nil {
			return fmt.Errorf("failed to delete '%s': %v", lmv, err)
		}
		if err = os.Rename(stage, lmv); err != nil {
			return fmt.Errorf("failed to rename '%s' toi '%s': %v", stage, lmv, err)
		}
		var tag string
		if tag, err = calculateModuleTag(false, true, lmv); err != nil {
			return fmt.Errorf("failed to calculate module tag: %v", err)
		}
		fmt.Printf("tag: %s\n", tag)
		return nil
	}
	// for v1 it's a bit more complicated since stage is a subdir of LMV.
	// first move stage to a temp dir outside of LMV, then remove LMV, then move temp to LMV
	temp := dest + "1temp"
	if err = os.Rename(stage, temp); err != nil {
		return fmt.Errorf("failed to rename '%s' to '%s': %v", stage, temp, err)
	}
	if err = os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to delete '%s': %v", dest, err)
	}
	if err = os.Rename(temp, dest); err != nil {
		return fmt.Errorf("failed to rename '%s' to '%s': %v", temp, dest, err)
	}
	var tag string
	if tag, err = calculateModuleTag(false, false, dest); err != nil {
		return fmt.Errorf("failed to calculate module tag: %v", err)
	}
	fmt.Printf("tag: %s\n", tag)
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
	f, err := os.Open(parent)
	if err != nil {
		return "", fmt.Errorf("failed to open '%s': %v", parent, err)
	}
	defer f.Close()
	names, err := f.Readdirnames(0)
	if err != nil {
		return "", fmt.Errorf("failed to read dir contents: %v", err)
	}
	dirs := []string{}
	for _, name := range names {
		fi, err := os.Lstat(filepath.Join(parent, name))
		if err != nil {
			return "", fmt.Errorf("failed to get file info: %v", err)
		}
		// only include major version subdirs, v2, v3, etc...
		if fi.IsDir() && fi.Name()[0] == 'v' {
			dirs = append(dirs, filepath.Join(parent, fi.Name()))
		}
	}
	// no dirs means this is a v1 package
	if len(dirs) == 0 {
		return parent, nil
	}
	sort.Strings(dirs)
	// last one in the slice is the largest
	return dirs[len(dirs)-1], nil
}

// loads exported content in lhsExports and rhsExports vars
func loadExports(lmv, stage string) error {
	var err error
	lhsExports, err = exports.Get(lmv)
	if err != nil {
		return fmt.Errorf("failed to get exports for package '%s': %s", lmv, err)
	}
	rhsExports, err = exports.Get(stage)
	if err != nil {
		return fmt.Errorf("failed to get exports for package '%s': %s", stage, err)
	}
	return nil
}

// returns true if the package in stage contains breaking changes
func hasBreakingChanges(lmv, stage string) bool {
	// check for changed content
	if len(delta.GetConstTypeChanges(lhsExports, rhsExports)) > 0 ||
		len(delta.GetFuncSigChanges(lhsExports, rhsExports)) > 0 ||
		len(delta.GetInterfaceMethodSigChanges(lhsExports, rhsExports)) > 0 ||
		len(delta.GetStructFieldChanges(lhsExports, rhsExports)) > 0 {
		return true
	}
	// check for removed content
	if removed := delta.GetExports(rhsExports, lhsExports); !removed.IsEmpty() {
		return true
	}
	return false
}

// updates the module version inside the go.mod file
func updateGoModVer(stage, newVer string) error {
	goMod := filepath.Join(stage, "go.mod")
	file, err := os.Open(goMod)
	if err != nil {
		return fmt.Errorf("failed to open for read '%s': %v", goMod, err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()
	file, err = os.Create(goMod)
	if err != nil {
		return fmt.Errorf("failed to open for write '%s': %v", goMod, err)
	}
	defer file.Close()
	for _, line := range lines {
		if strings.Index(line, "module") > -1 {
			line = line + "/" + newVer
		}
		fmt.Fprintln(file, line)
	}
	return nil
}

func writeChangelog(stage string) error {
	// TODO
	return nil
}

// returns the appropriate module tag based on the package version info
func calculateModuleTag(hasBreaks, hasVer bool, dest string) (string, error) {
	// if this has breaking changes then it's simply the dest minus the repo path
	// e.g. ~/work/src/github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis/v2
	// would return services/redis/mgmt/2018-03-01/redis/v2.0.0
	tagPrefix := dest[strings.Index(dest, "github.com")+34:]
	if hasBreaks {
		return tagPrefix + ".0.0", nil
	}
	if !hasVer {
		tagPrefix = tagPrefix + "/v1"
	}
	wt, err := repo.Get(dest)
	if err != nil {
		return "", err
	}
	tags, err := wt.ListTags(tagPrefix + "*")
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		// this is v1.0.0
		return tagPrefix + ".0.0", nil
	}
	regex := regexp.MustCompile(`v\d+\.\d+\.\d+$`)
	sort.SliceStable(tags, func(i, j int) bool {
		l := regex.FindString(tags[i])
		r := regex.FindString(tags[j])
		if l == "" || r == "" {
			panic("semver missing in module tag!")
		}
		lv, err := semver.NewVersion(l)
		if err != nil {
			panic(err)
		}
		rv, err := semver.NewVersion(r)
		if err != nil {
			panic(err)
		}
		return lv.LessThan(rv)
	})
	tag := tags[len(tags)-1]
	v := regex.FindString(tag)
	sv, _ := semver.NewVersion(v)
	// for non-breaking changes determine if this is a minor or patch update.
	if adds := delta.GetExports(lhsExports, rhsExports); !adds.IsEmpty() {
		// new exports, this is a minor update so bump minor version
		n := sv.IncMinor()
		sv = &n
	} else {
		// no new exports, this is a patch update
		n := sv.IncPatch()
		sv = &n
	}
	return strings.Replace(tag, v, sv.String(), 1), nil
}

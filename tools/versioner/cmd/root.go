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
	hasBreaks, err := hasBreakingChanges(lmv, stage)
	if err != nil {
		return fmt.Errorf("failed to detect breaking changes: %v", err)
	}
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
		err = updateGoModVer(stage, ver)
		if err != nil {
			return fmt.Errorf("failed to update go.mod file: %v", err)
		}
		dest = filepath.Join(dest, ver)
		err = os.Rename(stage, dest)
		if err != nil {
			err = fmt.Errorf("failed to rename '%s' to '%s': %v", stage, dest, err)
		}
		return err
	}
	// move staging directory over the LMV by first deleting LMV then renaming stage
	if hasVer {
		err = os.RemoveAll(lmv)
		if err != nil {
			return fmt.Errorf("failed to delete '%s': %v", lmv, err)
		}
		err = os.Rename(stage, lmv)
		if err != nil {
			err = fmt.Errorf("failed to rename '%s' toi '%s': %v", stage, lmv, err)
		}
		return err
	}
	// for v1 it's a bit more complicated since stage is a subdir of LMV.
	// first move stage to a temp dir outside of LMV, then remove LMV, then move temp to LMV
	temp := dest + "1temp"
	err = os.Rename(stage, temp)
	if err != nil {
		return fmt.Errorf("failed to rename '%s' to '%s': %v", stage, temp, err)
	}
	err = os.RemoveAll(dest)
	if err != nil {
		return fmt.Errorf("failed to delete '%s': %v", dest, err)
	}
	err = os.Rename(temp, dest)
	if err != nil {
		err = fmt.Errorf("failed to rename '%s' to '%s': %v", temp, dest, err)
	}
	return err
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

// returns true if the package in stage contains breaking changes
func hasBreakingChanges(lmv, stage string) (bool, error) {
	lhs, err := exports.Get(lmv)
	if err != nil {
		return false, fmt.Errorf("failed to get exports for package '%s': %s", lmv, err)
	}
	rhs, err := exports.Get(stage)
	if err != nil {
		return false, fmt.Errorf("failed to get exports for package '%s': %s", stage, err)
	}
	// check for changed content
	if len(delta.GetConstTypeChanges(lhs, rhs)) > 0 ||
		len(delta.GetFuncSigChanges(lhs, rhs)) > 0 ||
		len(delta.GetInterfaceMethodSigChanges(lhs, rhs)) > 0 ||
		len(delta.GetStructFieldChanges(lhs, rhs)) > 0 {
		return true, nil
	}
	// check for removed content
	if removed := delta.GetExports(rhs, lhs); !removed.IsEmpty() {
		return true, nil
	}
	return false, nil
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

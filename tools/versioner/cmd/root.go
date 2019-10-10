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
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "versioner <searching dir> [initial module version]",
	Short: "Creates or updates the latest major version for a package from staged content.",
	Long: `This tool will compare a staged package against its latest major version to detect
breaking changes.  If there are no breaking changes the latest major version is updated
with the staged content.  If there are breaking changes the staged content becomes the
next latest major vesion and the go.mod file is updated.
The default version for new modules is v1.0.0 or the value specified for [initial module version].
`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		_, err := theCommand(args)
		return err
	},
}

var (
	quietFlag   bool
	verboseFlag bool
)

const (
	stageName = "stage"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "quiet output")
}

// Execute executes the specified command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ExecuteVersioner is used for programmatically call in other tools
func ExecuteVersioner(root string, tagsHook TagsHookFunc) ([]string, error) {
	if tagsHook != nil {
		getTagsHook = tagsHook
	}
	return theCommand([]string{root})
}

// wrapper for cobra, prints tag to stdout
func theCommand(args []string) ([]string, error) {
	root, err := filepath.Abs(args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path from '%s': %v", args[0], err)
	}
	stages, err := findAllSubDirectories(root, stageName)
	printf("Found %d stage folder(s)", len(stages))
	vprintf("Stage folders: \n%s\n", strings.Join(stages, "\n"))
	tags := make([]string, 0)
	for _, stage := range stages {
		args[0] = stage
		tag, err := theUnstageCommand(args)
		if err != nil {
			return tags, fmt.Errorf("failed to get tag in stage folder '%s': %v", stage, err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func findAllSubDirectories(root, target string) ([]string, error) {
	// check if root exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("the root path '%s' does not exist", root)
	}
	stages := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == target {
			stages = append(stages, path)
			return nil
		}
		return nil
	})
	return stages, err
}

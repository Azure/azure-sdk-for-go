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
	"github.com/Azure/azure-sdk-for-go/tools/internal/modinfo"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	// the root command
	root := &cobra.Command{
		Use:   "versioner <searching dir> [initial module version]",
		Short: "Creates or updates the latest version for a package from the staged content.",
		Long: `This tool will compare a staged package against its latest version to detect
breaking changes.  If there are no breaking changes, the latest version is updated
with the staged content.  If there are breaking changes the staged content becomes the
next latest major vesion and the go.mod file is updated.
The default version for new modules is v1.0.0 or the value specified for [initial module version].
`,
		Args: cobra.RangeArgs(1, 2),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel("warn")
			if verbose := viper.GetBool("verbose"); verbose {
				log.SetLevel("debug")
			}
			if quiet := viper.GetBool("quiet"); quiet {
				log.SetLevel("error")
			}
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			root := args[0]
			versionSetting, err := parseVersionSetting(args[1:]...)
			if err != nil {
				return err
			}
			_, err = ExecuteVersioner(root, versionSetting, getTags)
			return err
		},
	}
	// register flags
	pFlags := root.PersistentFlags()
	pFlags.Bool("verbose", false, "verbose output")
	if err := viper.BindPFlag("verbose", pFlags.Lookup("verbose")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}
	pFlags.Bool("quiet", false, "suppress all outputs")
	if err := viper.BindPFlag("quiet", pFlags.Lookup("quiet")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}

	// other sub-commands
	root.AddCommand(unstageCommand())

	root.AddCommand(initCommand())

	return root
}

const (
	stageName = "stage"
	// default version to start a module at if not specified
	startingModVer = "v1.0.0"
	// default version for a new preview module
	startingModVerPreview = "v0.0.0"

	repoOrg = "Azure"
	repoName = "azure-sdk-for-go"
)

type VersionSetting struct {
	InitialVersion        string
	InitialVersionPreview string
}

func parseVersionSetting(args ...string) (*VersionSetting, error) {
	initialVersion := startingModVer
	if len(args) > 0 {
		if !modinfo.IsValidModuleVersion(args[0]) {
			return nil, fmt.Errorf("the string '%s' is not a valid module version", args[0])
		}
		initialVersion = args[1]
	}
	initialPreviewVersion := startingModVerPreview
	if len(args) > 1 {
		if !modinfo.IsValidModuleVersion(args[1]) {
			return nil, fmt.Errorf("the string '%s' is not a valid module version", args[1])
		}
		initialPreviewVersion = args[1]
	}
	return &VersionSetting{
		InitialVersion:        initialVersion,
		InitialVersionPreview: initialPreviewVersion,
	}, nil
}

// wrapper for cobra, prints tag to stdout
func ExecuteVersioner(r string, versionSetting *VersionSetting, hookFunc TagsHookFunc) ([]string, error) {
	root, err := filepath.Abs(r)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path from '%s': %v", root, err)
	}
	// get all the stage sub-directories
	subDirectories, err := findAllSubDirectories(root, stageName)
	if err != nil {
		return nil, fmt.Errorf("failed to list all sub-directories under '%s': %+v", root, err)
	}

	tags := make([]string, 0)
	for _, dir := range subDirectories {
		_, tag, err := ExecuteUnstage(dir, versionSetting, hookFunc)
		if err != nil {
			return tags, fmt.Errorf("failed to get tag in stage folder '%s': %v", dir, err)
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

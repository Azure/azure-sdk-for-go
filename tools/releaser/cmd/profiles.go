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
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/ioext"
	"github.com/Azure/azure-sdk-for-go/tools/internal/dirs"
	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/spf13/cobra"
)

func profilesCommand() *cobra.Command {
	profile := &cobra.Command{
		Use:   "profiles <sdk directory>",
		Short: "Release new modules of profiles from the profiles' root directory",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return executeProfiles(path)
		},
	}

	return profile
}

const (
	profilesName = "profiles"
	tempName     = "profiles1temp"
)

func executeProfiles(path string) error {
	// first we need to backup the current state of profiles
	root, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of '%s': %+v", path, err)
	}
	// check if root exists
	profilesPath := filepath.Join(root, profilesName)
	temp := filepath.Join(root, tempName)
	log.Debugf("Copying from `%s` to `%s`", profilesPath, temp)
	if err := ioext.CopyDir(profilesPath, temp); err != nil {
		return fmt.Errorf("failed to copy '%s' to '%s'", root, temp)
	}

	defer func() {
		log.Debugf("Removing temp directory '%s'...", temp)
		if err := os.RemoveAll(temp); err != nil {
			log.Errorf("failed to remove temp directory '%s': %+v", temp, err)
		}
	}()
	// format old code to avoid CRLF difference
	if err := formatCode(temp); err != nil {
		return fmt.Errorf("failed to format temp directory: %+v", err)
	}

	// runs go generate
	log.Debug("Regenerating profiles...")
	profiles, err := dirs.GetSubdirs(profilesPath)
	if err != nil {
		return fmt.Errorf("failed to list all profiles: %+v", err)
	}
	for _, profile := range profiles {
		identical, err := compareProfile(root, filepath.Base(profile))
		if err != nil {
			return fmt.Errorf("failed to compare profile '%s': %+v", profile, err)
		}
		log.Debugf("Profile %s identical: %t", profile, identical)
	}
	// format both of them to avoid CRLF difference
	if err := formatCode(profilesPath); err != nil {
		return fmt.Errorf("failed to format '%s': %+v", profilesPath, err)
	}
	// compares differences and bump version if necessary
	// iterate all profiles
	return nil
}

func compareProfile(root, profile string) (bool, error) {
	base := filepath.Join(root, tempName, profile)
	target := filepath.Join(root, profilesName, profile)
	return dirs.DeepCompare(base, target)
}

func executeGoGenerate(path string) error {
	c := exec.Command("go", "generate", "./" + profilesName)
	c.Dir = path
	output, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(output))
	}
	log.Debug(string(output))
	return nil
}

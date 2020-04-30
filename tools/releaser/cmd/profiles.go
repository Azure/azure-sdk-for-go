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
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/ioext"
	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/spf13/cobra"
	"os/exec"
	"path/filepath"
)

func profilesCommand() *cobra.Command {
	profile := &cobra.Command{
		Use:   "profiles <sdk directory>",
		Short: "Release new modules of profiles from the profiles' root directory",
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return executeProfiles(path)
		},
	}

	return profile
}

const (
	profilesName = "profiles"
	tempName = "profiles1temp"
)

func executeProfiles(path string) error {
	// first we need to backup the current state of profiles
	root, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of '%s': %+v", path, err)
	}
	// check if root exists
	profiles := filepath.Join(root, profilesName)
	temp := filepath.Join(root, tempName)
	log.Debugf("Copying from `%s` to `%s`", profiles, temp)
	if err := ioext.CopyDir(profiles, temp); err != nil {
		return fmt.Errorf("failed to copy '%s' to '%s'", root, temp)
	}
	// runs go generate
	if err := executeGoGenerate(root); err != nil {
		return fmt.Errorf("failed to run go generate: %+v", err)
	}
	// format both of them to avoid CRLF difference
	// compares differences and bump version if necessary
	// remove the temp directory
	return nil
}

func executeGoGenerate(path string) error {
	c := exec.Command("go", "generate", "./" + profilesName)
	c.Dir = path
	output, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(output))
	}
	log.Info(string(output))
	return nil
}

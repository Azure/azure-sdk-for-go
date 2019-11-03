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
	"strings"
)

// check if this update has nothing changed other than a version number
// this is necessary because autorest uses multiapi for golang generation, whether this tag is updated or not,
// the stage directory will generate anyway.
func checkIdentical(dest, stage string) (bool, error) {
	// get files with changes
	files, err := getFilesWithRealChanges(dest, stage)
	if err != nil {
		return false, err
	}
	// check if there is no change
	if len(files) == 0 {
		return true, nil
	}
	// there is only version.go file changed
	if len(files) == 1 && filepath.Base(files[0]) == "version.go" {
		return true, nil
	}
	return false, nil
}

func getFilesWithRealChanges(dest, stage string) ([]string, error) {
	defer reset(dest)
	c := exec.Command("git", "add", dest)
	output, err := c.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to call git add: %s", string(output))
	}
	c = exec.Command("git", "status", "-s", "--no-renames", dest)
	output, err = c.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to call git status: %s", string(output))
	}
	changes, err := analyzeOutput(string(output), stage)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze git status output: %v", err)
	}
	return changes, nil
}

func reset(dest string) error {
	c := exec.Command("git", "reset", "HEAD", dest)
	output, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to call git reset: %s", string(output))
	}
	return nil
}

func analyzeOutput(output, stage string) ([]string, error) {
	lines := strings.Split(output, "\n")
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %v", err)
	}
	var files []string
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		path := filepath.Join(pwd, line[3:])
		// ignore stage folder and everything other than go source files
		if strings.HasPrefix(path, stage) || !strings.HasSuffix(path, ".go") {
			continue
		}
		files = append(files, path)
	}
	return files, nil
}
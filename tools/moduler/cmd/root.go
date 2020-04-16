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
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/spf13/cobra"
)

const (
	tagPrefix   = "// tag: "
	versionFile = "version.go"
)

func Command() *cobra.Command {
	// the root command
	root := &cobra.Command{
		Use: "moduler <searching dir>",
		Short: "Release a new module by pushing new tags to github",
		Long: `This tool search the whole SDK folder for tags in version.go files 
which are produced by the versioner tool, and push the new tags to github.`,
		Args: cobra.ExactArgs(1),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel("warn")
			if verbose := viper.GetBool("verbose"); verbose {
				log.SetLevel("debug")
			}
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			flags := flags{
				dryRunFlag: viper.GetBool("dry-run"),
				addOnly: viper.GetBool("add-only"),
			}
			return ExecuteModuler(args[0], flags, getTags)
		},
	}
	// register flags
	pFlags := root.PersistentFlags()
	pFlags.Bool("verbose", false, "verbose output")
	if err := viper.BindPFlag("verbose", pFlags.Lookup("verbose")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}
	flags := root.Flags()
	flags.BoolP("dry-run", "d", false, "dry run, only list the detected tags, do not add or push them")
	if err := viper.BindPFlag("dry-run", pFlags.Lookup("dry-run")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}
	flags.BoolP("add-only", "a", false, "only add tags but do not push them")
	if err := viper.BindPFlag("add-only", pFlags.Lookup("add-only")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}

	return root
}

type flags struct {
	dryRunFlag bool
	addOnly bool
}

// TagsHookFunc is a func used for get tags from remote
type TagsHookFunc func(root string) ([]string, error)

func ExecuteModuler(r string, flags flags, getTagsHook TagsHookFunc) error {
	root, err := filepath.Abs(r)
	if err != nil {
		return fmt.Errorf("failed to get absolute root: %v", err)
	}
	newTags, err := readNewTags(root, getTagsHook)
	if err != nil {
		return err
	}
	log.Infoln("Found new tags: ")
	log.Infoln(strings.Join(newTags, "\n"))
	if !flags.dryRunFlag {
		// add new tags
		for _, tag := range newTags {
			if err := addNewTag(tag); err != nil {
				return err
			}
		}
		// push new tags
		if !flags.addOnly {
			if err := pushTags(); err != nil {
				return fmt.Errorf("failed to push tags: %v", err)
			}
		}
	}
	return nil
}

func readNewTags(root string, getTagsHook TagsHookFunc) ([]string, error) {
	// get list of all existing tags
	tags, err := getTagsHook(root)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %v", err)
	}
	log.Infof("Get %d tags from remote", len(tags))
	log.Debugln(strings.Join(tags, "\n"))
	// find all version.go files
	files, err := findAllFiles(root, versionFile)
	if err != nil {
		return nil, fmt.Errorf("failed to find all version.go files: %v", err)
	}
	// read new tags in version.go
	newTags, err := readTags(tags, files)
	if err != nil {
		return nil, fmt.Errorf("failed during reading new tags: %v", err)
	}
	return newTags, nil
}

func getTags(repoPath string) ([]string, error) {
	wt, err := repo.Get(repoPath)
	if err != nil {
		return nil, err
	}
	return wt.ListTags("*")
}

func tagExists(tags []string, tag string) bool {
	for _, item := range tags {
		if item == tag {
			return true
		}
	}
	return false
}

func findAllFiles(root, filename string) ([]string, error) {
	// check if root exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("the root path '%s' does not exist", root)
	}
	files := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Name() == filename {
			files = append(files, path)
			return nil
		}
		return nil
	})
	return files, err
}

func readTags(tags []string, files []string) ([]string, error) {
	var newTags []string
	for _, file := range files {
		newTag, err := readNewTagInFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read tag in file '%s': %v", file, err)
		}
		if newTag == "" {
			log.Debugf("Found empty tag in file '%s'", file)
			continue
		} else if tagExists(tags, newTag) {
			log.Debugf("Found existed tag '%s' in file '%s'", newTag, file)
			continue
		}
		// add new tag to list
		tags = append(tags, newTag)
		newTags = append(newTags, newTag)
	}
	return newTags, nil
}

func readNewTagInFile(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	tag := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, tagPrefix) {
			tag = strings.ReplaceAll(line, tagPrefix, "")
			break
		}
	}
	return tag, nil
}

func addNewTag(newTag string) error {
	log.Debugf("Adding new tag '%s'", newTag)
	if err := executeCommand("git", "tag", "-a", newTag, "-m", newTag); err != nil {
		return err
	}
	return nil
}

func pushTags() error {
	log.Debugln("Pushing new tags")
	return executeCommand("git", "push", "origin", "--tags")
}

func executeCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}
	return nil
}

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
	"strings"
	"unicode"

	"github.com/Azure/azure-sdk-for-go/tools/internal/pkgs"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var exceptFileFlag string

var rootCmd = &cobra.Command{
	Use:   "pkgchk <dir>",
	Short: "Performs package validation tasks against all packages found under the specified directory.",
	Long: `This tool will perform various package validation checks against all of the packages
found under the specified directory.  Failures can be baselined and thus ignored by
copying the failure text verbatim, pasting it into a text file then specifying that
file via the optional exceptions flag.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return theCommand(args)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&exceptFileFlag, "exceptions", "e", "", "text file containing the list of exceptions")
}

// Execute executes the specified command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func theCommand(args []string) error {
	rootDir, err := filepath.Abs(args[0])
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path")
	}
	ps, err := pkgs.GetPkgs(rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to get packages")
	}
	exceptions, err := loadExceptions(exceptFileFlag)
	if err != nil {
		return errors.Wrap(err, "failed to load exceptions")
	}
	verifiers := getVerifiers()
	count := 0
	for _, pkg := range ps {
		for _, v := range verifiers {
			if err = v(pkg); err != nil && !contains(exceptions, err.Error()) {
				fmt.Fprintln(os.Stderr, err)
				count++
			}
		}
	}
	if count > 0 {
		return fmt.Errorf("found %d errors", count)
	}
	return nil
}

func contains(items []string, item string) bool {
	if items == nil {
		return false
	}
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func loadExceptions(exceptFile string) ([]string, error) {
	if exceptFile == "" {
		return nil, nil
	}
	f, err := os.Open(exceptFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	exceps := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		exceps = append(exceps, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return exceps, nil
}

type verifier func(p pkgs.Pkg) error

// returns a list of verifiers to execute
func getVerifiers() []verifier {
	return []verifier{
		verifyPkgMatchesDir,
		verifyLowerCase,
		verifyDirectoryStructure,
	}
}

// ensures that the leaf directory name matches the package name
// new to modules: if the last leaf is version suffix, find its parent as leaf folder
func verifyPkgMatchesDir(p pkgs.Pkg) error {
	leaf := findPackageFolderInPath(p.Dest)
	if !strings.EqualFold(leaf, p.Package.Name) {
		return fmt.Errorf("leaf directory of '%s' doesn't match package name '%s'", p.Dest, p.Package.Name)
	}
	return nil
}

func findPackageFolderInPath(path string) string {
	regex := regexp.MustCompile(`/v\d+$`)
	if regex.MatchString(path) {
		// folder path ends with version suffix
		path = path[:strings.LastIndex(path, "/")]
	}
	result := path[strings.LastIndex(path, "/")+1:]
	return result
}

// ensures that there are no upper-case letters in a package's directory
func verifyLowerCase(p pkgs.Pkg) error {
	// walk the package directory looking for upper-case characters
	for _, r := range p.Dest {
		if r == '/' {
			continue
		}
		if unicode.IsUpper(r) {
			return fmt.Errorf("found upper-case character in directory '%s'", p.Dest)
		}
	}
	return nil
}

// ensures that the package's directory hierarchy is properly formed
func verifyDirectoryStructure(p pkgs.Pkg) error {
	// for ARM the package directory structure is highly deterministic:
	// /redis/mgmt/2015-08-01/redis
	// /resources/mgmt/2017-06-01-preview/policy
	// /preview/signalr/mgmt/2018-03-01-preview/signalr
	// /preview/security/mgmt/v2.0/security (version scheme for composite packages)
	// /network/mgmt/2019-10-01/network/v2 (new with modules)
	if !p.IsARMPkg() {
		return nil
	}
	regexStr := strings.Join([]string{
		`^(?:/preview)?`,
		`[a-z0-9\-]+`,
		`mgmt`,
		`((?:\d{4}-\d{2}-\d{2}(?:-preview)?)|(?:v\d{1,2}\.\d{1,2}))`,
		`[a-z0-9]+`,
	}, "/")
	regexStr = regexStr + `(/v\d+)?$`
	regex := regexp.MustCompile(regexStr)
	if !regex.MatchString(p.Dest) {
		return fmt.Errorf("bad directory structure '%s'", p.Dest)
	}
	return nil
}

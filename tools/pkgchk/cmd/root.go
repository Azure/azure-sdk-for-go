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

	"github.com/Azure/azure-sdk-for-go/tools/pkgchk/track1"
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
	return VerifyPackages(rootDir, exceptFileFlag)
}

// VerifyPackages verifies all the track 1 packages under the root directory
func VerifyPackages(root, exceptionFile string) error {
	pkgs, err := track1.List(root)
	if err != nil {
		return fmt.Errorf("failed to get packages: %+v", err)
	}
	exceptions, err := loadExceptions(exceptionFile)
	if err != nil {
		return fmt.Errorf("failed to load exceptions: %+v", err)
	}
	verifier := track1.GetDefaultVerifier()
	count := 0
	for _, pkg := range pkgs {
		for _, err := range verifier.Verify(pkg) {
			if _, ok := exceptions[err.Error()]; !ok {
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

func loadExceptions(exceptFile string) (map[string]bool, error) {
	if exceptFile == "" {
		return nil, nil
	}
	f, err := os.Open(exceptFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	exceptions := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		exceptions[scanner.Text()] = true
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return exceptions, nil
}

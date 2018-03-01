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

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/exports"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package [package dir] [base commit] [target commit]",
	Short: "Generates report for the package in the specified directory.",
	Long: `The package command generates a report for the package in the directory specified in [package dir].
The package content in [target commit] is compared against the package content in [base commit]
to determine what changes were introduced in [target commit].`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgDir, cloneRepo, cleanupFn, err := processArgsAndClone(args)
		if err != nil {
			return err
		}
		defer cleanupFn()

		baseCommit := args[1]
		targetCommit := args[2]

		// lhs
		vprintf("checking out base commit %s and gathering exports\n", baseCommit)
		lhs, err := getContentForCommit(cloneRepo, pkgDir, baseCommit)
		if err != nil {
			return err
		}

		// rhs
		vprintf("checking out target commit %s and gathering exports\n", targetCommit)
		rhs, err := getContentForCommit(cloneRepo, pkgDir, targetCommit)
		if err != nil {
			return err
		}

		return printReport(getPkgReport(lhs, rhs))
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
}

func getContentForCommit(wt repo.WorkingTree, dir, commit string) (cnt exports.Content, err error) {
	err = wt.Checkout(commit)
	if err != nil {
		err = fmt.Errorf("failed to check out commit '%s': %s", commit, err)
		return
	}

	cnt, err = exports.Get(dir)
	if err != nil {
		err = fmt.Errorf("failed to get exports for commit '%s': %s", commit, err)
	}
	return
}

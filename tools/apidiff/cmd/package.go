// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/repo"
	"github.com/Azure/azure-sdk-for-go/tools/internal/report"
	"github.com/spf13/cobra"
)

var (
	dirMode    bool
	asMarkdown bool
)

var packageCmd = &cobra.Command{
	Use:   "package <package dir> (<base commit> <target commit(s)>) | (<commit sequence>)",
	Short: "Generates a report for the package in the specified directory containing the delta between commits.",
	Long: `The package command generates a report for the package in the directory specified in <package dir>.
Commits can be specified as either a base and one or more target commits or a sequence of commits.
For a base/target pair each target commit is compared against the base commit.
For a commit sequence each commit N in the sequence is compared against commit N+1.
Commit sequences must be comma-delimited.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rpt, err := thePackageCmd(args)
		if err != nil {
			return err
		}
		evalReportStatus(rpt)
		return nil
	},
}

// split into its own func as we can't call os.Exit from it (the defer won't get executed)
func thePackageCmd(args []string) (rs report.Status, err error) {
	if dirMode {
		return packageCmdDirMode(args)
	}
	return packageCmdCommitMode(args)
}

func init() {
	packageCmd.PersistentFlags().BoolVarP(&dirMode, "directories", "i", false, "compares packages in two different directories")
	packageCmd.PersistentFlags().BoolVarP(&asMarkdown, "markdown", "m", false, "emits the report in markdown format")
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

func packageCmdCommitMode(args []string) (rs report.Status, err error) {
	cloneRepo, err := processArgsAndClone(args)
	if err != nil {
		return
	}

	var rpt report.CommitPkgReport
	rpt.CommitsReports = map[string]report.Package{}
	worker := func(pkgDir string, cloneRepo repo.WorkingTree, baseCommit, targetCommit string) error {
		// lhs
		vprintf("checking out base commit %s and gathering exports\n", baseCommit)
		var lhs exports.Content
		lhs, err = getContentForCommit(cloneRepo, pkgDir, baseCommit)
		if err != nil {
			return err
		}

		// rhs
		vprintf("checking out target commit %s and gathering exports\n", targetCommit)
		var rhs exports.Content
		rhs, err = getContentForCommit(cloneRepo, pkgDir, targetCommit)
		if err != nil {
			return err
		}
		r := report.Generate(lhs, rhs, &report.GenerationOption{
			OnlyBreakingChanges: onlyBreakingChangesFlag,
			OnlyAdditiveChanges: onlyAdditiveChangesFlag,
		})
		if r.HasBreakingChanges() {
			rpt.BreakingChanges = append(rpt.BreakingChanges, targetCommit)
		}
		rpt.CommitsReports[fmt.Sprintf("%s:%s", baseCommit, targetCommit)] = r
		return nil
	}

	err = generateReports(args, cloneRepo, worker)
	if err != nil {
		return
	}

	err = PrintReport(rpt)
	return
}

func packageCmdDirMode(args []string) (rs report.Status, err error) {
	if len(args) != 2 {
		return nil, errors.New("directories mode requires two arguments")
	}
	lhs, err := exports.Get(args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get exports for package '%s': %s", args[0], err)
	}
	rhs, err := exports.Get(args[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get exports for package '%s': %s", args[1], err)
	}
	r := report.Generate(lhs, rhs, &report.GenerationOption{
		OnlyBreakingChanges: onlyBreakingChangesFlag,
		OnlyAdditiveChanges: onlyAdditiveChangesFlag,
	})
	if asMarkdown && !suppressReport {
		println(r.ToMarkdown())
	} else {
		err = PrintReport(r)
	}
	return r, err
}

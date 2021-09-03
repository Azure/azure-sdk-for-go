// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/repo"
	"github.com/Azure/azure-sdk-for-go/tools/internal/report"
	"github.com/spf13/cobra"
)

var packagesCmd = &cobra.Command{
	Use:   "packages <package search dir> (<base commit> <target commit(s)>) | (<commit sequence>)",
	Short: "Generates a report for all packages under the specified directory containing the delta between commits.",
	Long: `The packages command generates a report for all of the packages under the directory specified in <package dir>.
Commits can be specified as either a base and one or more target commits or a sequence of commits.
For a base/target pair each target commit is compared against the base commit.
For a commit sequence each commit N in the sequence is compared against commit N+1.
Commit sequences must be comma-delimited.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rpt, err := thePackagesCmd(args)
		if err != nil {
			return err
		}
		err = PrintReport(rpt)
		if err != nil {
			return err
		}
		evalReportStatus(rpt)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)
}

// ExecPackagesCmd is the programmatic interface for the packages command.
func ExecPackagesCmd(pkgDir string, commitSeq string, flags CommandFlags) (report.CommitPkgsReport, error) {
	flags.apply()
	return thePackagesCmd([]string{pkgDir, commitSeq})
}

// split into its own func as we can't call os.Exit from it (the defer won't get executed)
func thePackagesCmd(args []string) (rpt report.CommitPkgsReport, err error) {
	cloneRepo, err := processArgsAndClone(args)
	if err != nil {
		return
	}

	rpt.CommitsReports = map[string]report.PkgsReport{}
	worker := func(rootDir string, cloneRepo repo.WorkingTree, baseCommit, targetCommit string) error {
		vprintf("generating diff between %s and %s\n", baseCommit, targetCommit)
		// get for lhs
		dprintf("checking out base commit %s and gathering exports\n", baseCommit)
		lhs, err := getRepoContentForCommit(&cloneRepo, rootDir, baseCommit)
		if err != nil {
			return err
		}

		// get for rhs
		dprintf("checking out target commit %s and gathering exports\n", targetCommit)
		rhs, err := getRepoContentForCommit(&cloneRepo, rootDir, targetCommit)
		if err != nil {
			return err
		}
		r := getPkgsReport(lhs, rhs)
		rpt.UpdateAffectedPackages(targetCommit, r)
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

	return
}

func getRepoContentForCommit(wt *repo.WorkingTree, dir, commit string) (r RepoContent, err error) {
	err = wt.Checkout(commit)
	if err != nil {
		err = fmt.Errorf("failed to check out commit '%s': %s", commit, err)
		return
	}

	return getRepoContent(wt, dir)
}

func getRepoContent(wt *repo.WorkingTree, dir string) (RepoContent, error) {
	pkgDirs, err := report.GetPackages(dir)
	if err != nil {
		return nil, err
	}

	if debugFlag {
		fmt.Println("found the following package directories")
		for _, d := range pkgDirs {
			fmt.Printf("\t%s\n", d)
		}
	}

	r, err := getExportsForPackages(wt.Root(), pkgDirs)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// RepoContent contains repo content, it's structured as "package path":content
type RepoContent map[string]exports.Content

// returns RepoContent based on the provided slice of package directories
func getExportsForPackages(root string, pkgDirs []string) (RepoContent, error) {
	exps := RepoContent{}
	for _, pkgDir := range pkgDirs {
		dprintf("getting exports for %s\n", pkgDir)
		// pkgDir = "C:\Users\somebody\AppData\Local\Temp\apidiff-1529437978\services\addons\mgmt\2017-05-15\addons"
		// convert to package path "github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2016-05-16/analysisservices"
		pkgPath := strings.Replace(pkgDir, root, "github.com/Azure/azure-sdk-for-go", -1)
		pkgPath = strings.Replace(pkgPath, string(os.PathSeparator), "/", -1)
		if _, ok := exps[pkgPath]; ok {
			return nil, fmt.Errorf("duplicate package: %s", pkgPath)
		}
		exp, err := exports.Get(pkgDir)
		if err != nil {
			return nil, err
		}
		exps[pkgPath] = exp
	}
	return exps, nil
}

// Print prints the RepoContent to a Writer as JSON string
func (r *RepoContent) Print(o io.Writer) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}
	_, err = o.Write(b)
	return err
}

// generates a PkgsReport based on the delta between lhs and rhs
func getPkgsReport(lhs, rhs RepoContent) report.PkgsReport {
	rpt := report.PkgsReport{}

	if !onlyBreakingChangesFlag {
		rpt.AddedPackages = getDiffPkgs(lhs, rhs)
	}
	if !onlyAdditiveChangesFlag {
		rpt.RemovedPackages = getDiffPkgs(rhs, lhs)
	}

	// diff packages
	for rhsPkg, rhsCnt := range rhs {
		if _, ok := lhs[rhsPkg]; !ok {
			continue
		}
		if r := report.Generate(lhs[rhsPkg], rhsCnt, &report.GenerationOption{
			OnlyBreakingChanges: onlyBreakingChangesFlag,
			OnlyAdditiveChanges: onlyAdditiveChangesFlag,
		}); !r.IsEmpty() {
			// only add an entry if the report contains data
			if rpt.ModifiedPackages == nil {
				rpt.ModifiedPackages = report.ModifiedPackages{}
			}
			rpt.ModifiedPackages[rhsPkg] = r
		}
	}

	return rpt
}

// returns a list of packages in rhs that aren't in lhs
func getDiffPkgs(lhs, rhs RepoContent) report.PkgsList {
	list := report.PkgsList{}
	for rhsPkg := range rhs {
		if _, ok := lhs[rhsPkg]; !ok {
			list = append(list, rhsPkg)
		}
	}
	return list
}

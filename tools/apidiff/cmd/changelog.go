// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/markdown"
	"github.com/Azure/azure-sdk-for-go/tools/internal/report"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog <package search dir> <base commit> <target commit> <release tag version>",
	Short: "Generates a CHANGELOG report in markdown format for the packages under the specified directory.",
	Long: `The changelog command generates a CHANGELOG for all of the packages under the directory specified in <package dir>.
A table for added, removed, updated, and breaking changes will be created as required.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// there should be exactly four args, a directory, two commit hashes and the release tag version
		if err := cobra.ExactArgs(4)(cmd, args); err != nil {
			return err
		}
		if strings.Index(args[2], ",") > -1 {
			return errors.New("sequence of target commits is not supported")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return theChangelogCmd(args)
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
}

func theChangelogCmd(args []string) error {
	// TODO: refactor so that we don't depend on the packages command
	rpt, err := thePackagesCmd(args[:3])
	if err != nil {
		return err
	}
	if rpt.IsEmpty() {
		return nil
	}

	// there should only be one report, the delta between the base and target commits
	if len(rpt.CommitsReports) > 1 {
		panic("expected only one report")
	}
	for _, cr := range rpt.CommitsReports {
		changelog, err := writePackageChangelog(cr, args[3])
		if err != nil {
			return err
		}
		fmt.Println(changelog)
	}
	return nil
}

func writePackageChangelog(pr report.PkgsReport, version string) (string, error) {
	md := &markdown.Writer{}
	// write out the changelog's title and the release tag header before populating with other changes.
	md.WriteTitle("Release History")
	md.WriteTopLevelHeader(fmt.Sprintf("%s (Released)", version))
	if err := reportAddedPkgs(pr, md); err != nil {
		return "", fmt.Errorf("failed to write table for added packages: %+v", err)
	}
	if err := reportUpdatedPkgs(pr, md); err != nil {
		return "", fmt.Errorf("failed to write table for updated packages: %+v", err)
	}
	if err := reportBreakingPkgs(pr, md); err != nil {
		return "", fmt.Errorf("failed to write table for breaking change packages: %+v", err)
	}
	if err := reportRemovedPkgs(pr, md); err != nil {
		return "", fmt.Errorf("failed to write table for removed packages: %+v", err)
	}
	return md.String(), nil
}

func reportAddedPkgs(pr report.PkgsReport, md *markdown.Writer) error {
	if len(pr.AddedPackages) == 0 {
		return nil
	}
	t, err := createPackageTable(pr.AddedPackages)
	if err != nil {
		return err
	}
	if t.Rows() > 0 {
		md.WriteSubheader("New Packages")
		md.WriteTable(*t)
	}
	return nil
}

func reportUpdatedPkgs(pr report.PkgsReport, md *markdown.Writer) error {
	if pr.ModifiedPackages == nil || !pr.ModifiedPackages.HasAdditiveChanges() {
		return nil
	}
	var updated []string
	for pkgName, pkgRpt := range pr.ModifiedPackages {
		if pkgRpt.HasAdditiveChanges() && !pkgRpt.HasBreakingChanges() {
			updated = append(updated, pkgName)
		}
	}
	t, err := createPackageTable(updated)
	if err != nil {
		return err
	}
	if t.Rows() > 0 {
		md.WriteSubheader("Updated Packages")
		md.WriteTable(*t)
	}
	return nil
}

func reportBreakingPkgs(pr report.PkgsReport, md *markdown.Writer) error {
	if pr.ModifiedPackages == nil || !pr.ModifiedPackages.HasBreakingChanges() {
		return nil
	}
	var breaking []string
	for pkgName, pkgRpt := range pr.ModifiedPackages {
		if pkgRpt.HasBreakingChanges() {
			breaking = append(breaking, pkgName)
		}
	}
	t, err := createPackageTable(breaking)
	if err != nil {
		return err
	}
	if t.Rows() > 0 {
		md.WriteSubheader("Breaking Changes")
		md.WriteTable(*t)
	}
	return nil
}

func reportRemovedPkgs(pr report.PkgsReport, md *markdown.Writer) error {
	if len(pr.RemovedPackages) == 0 {
		return nil
	}
	t, err := createPackageTable(pr.RemovedPackages)
	if err != nil {
		return err
	}
	if t.Rows() > 0 {
		md.WriteSubheader("Removed Packages")
		md.WriteTable(*t)
	}
	return nil
}

type tableRow struct {
	pkgName     string
	apiVersions []string
}

func convertFullPackagePathToPackageNameAndAPIVersion(packageName string) (string, string, error) {
	// packageName is a string like "github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2018-08-31/consumption"
	segments := strings.Split(packageName, "/")
	if len(segments) < 2 {
		return "", "", fmt.Errorf("expecting package name '%s' to have at least two segments", packageName)
	}
	return segments[len(segments)-1], segments[len(segments)-2], nil
}

func createPackageTable(pkgs []string) (*markdown.Table, error) {
	t := markdown.NewTable("rc", "Package Name", "API Version")
	rows, err := categorizePackageAPIVersions(pkgs)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		t.AddRow(row.pkgName, strings.Join(row.apiVersions, "<br/>"))
	}
	return t, nil
}

func categorizePackageAPIVersions(pkgs []string) ([]tableRow, error) {
	entries := make(map[string][]string)
	for _, pkg := range pkgs {
		pkgName, apiVer, err := convertFullPackagePathToPackageNameAndAPIVersion(pkg)
		if err != nil {
			return nil, err
		}
		if apis, ok := entries[pkgName]; ok {
			entries[pkgName] = append(apis, apiVer)
		} else {
			entries[pkgName] = []string{apiVer}
		}
	}
	// convert the map to a slice of tableRows
	var rows []tableRow
	for pkgName, apiVers := range entries {
		sort.Strings(apiVers)
		rows = append(rows, tableRow{
			pkgName:     pkgName,
			apiVersions: apiVers,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].pkgName < rows[j].pkgName
	})
	return rows, nil
}

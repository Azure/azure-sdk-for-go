// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import "github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"

// GetPkgsReport generates a PkgsReport based on the delta between lhs and rhs
func GetPkgsReport(lhs, rhs RepoContent, option *report.GenerationOption) report.PkgsReport {
	rpt := report.PkgsReport{}

	if option == nil {
		option = &report.GenerationOption{}
	}

	if !option.OnlyBreakingChanges {
		rpt.AddedPackages = getDiffPkgs(lhs, rhs)
	}
	if !option.OnlyAdditiveChanges {
		rpt.RemovedPackages = getDiffPkgs(rhs, lhs)
	}

	// diff packages
	for rhsPkg, rhsCnt := range rhs {
		if _, ok := lhs[rhsPkg]; !ok {
			continue
		}
		if r := report.Generate(lhs[rhsPkg], rhsCnt, option); !r.IsEmpty() {
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

func GetPackagesReportFromContent(lhs RepoContent, targetRoot string) (*report.PkgsReport, error) {
	rhs, err := GetRepoContent(targetRoot)
	if err != nil {
		return nil, err
	}
	r := GetPkgsReport(lhs, rhs, nil)
	return &r, nil
}

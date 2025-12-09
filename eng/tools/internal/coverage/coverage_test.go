// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"testing"
)

func Test_ParseCoverageNoMatch(t *testing.T) {
	coverageLine := "<coverage no match>"
	_, err := parseCoveragePercent([]byte(coverageLine))
	if err == nil {
		t.Error("Expected error to be thrown on non-matching coverage line.")
	}
}

func Test_ParseCoverageMaximum(t *testing.T) {
	coverageLine := "total:                                                          (statements)    100.0%"

	coveragePercent, err := parseCoveragePercent([]byte(coverageLine))
	if err != nil {
		t.Error(err)
	}
	if coveragePercent != 1 {
		t.Errorf("Expected coverage of 1 to be parsed as 1, found %f", coveragePercent)
	}
}

func Test_ParseCoverageFloat(t *testing.T) {
	coverageLine := "total:                                                                                          (statements)                                      80.5%"

	coveragePercent, err := parseCoveragePercent([]byte(coverageLine))
	if err != nil {
		t.Error(err)
	}

	expected := .805
	if coveragePercent != expected {
		t.Errorf("Expected coverage percent of .805 to be parsed as %f, found %f", expected, coveragePercent)
	}
}

func Test_FindCoverageGoal(t *testing.T) {
	configData := &codeCoverage{
		Packages: []coveragePackage{
			{Name: "module", CoverageGoal: 1},
			{Name: "module/submodule", CoverageGoal: 2},
			{Name: "module/submodule/submodule_2", CoverageGoal: 3},
		},
	}
	for _, test := range []struct {
		covFile string
		want    float64
	}{
		{"default", 0.95},
		{"module", 1},
		{"module/foo", 1},
		{`C:\prefix\sdk\module`, 1},
		{"/prefix/sdk/module", 1},
		{"/prefix/sdk/module/foo/bar", 1},
		{"/prefix/sdk/module/submodule", 2},
		{"/prefix/sdk/module/submodule/submodule_2/submodule", 3},
	} {
		if got := findCoverageGoal([]string{test.covFile}, configData); got != test.want {
			t.Errorf("findCoverageGoal(%v) = %.2f; want %.2f", test.covFile, got, test.want)
		}
	}
}

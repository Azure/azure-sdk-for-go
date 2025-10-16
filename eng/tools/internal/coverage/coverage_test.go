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

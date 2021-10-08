package main

import (
	"strconv"
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
	coverageLine := "<coverage line-rate=\"1\" branch-rate=\"0\" lines-covered=\"3\" lines-valid=\"3\" branches-covered=\"0\" branches-valid=\"0\" complexity=\"0\" version=\"\" timestamp=\"1633558939111\">"

	coveragePercent, err := parseCoveragePercent([]byte(coverageLine))
	if err != nil {
		t.Error(err)
	}
	if coveragePercent != 1 {
		t.Errorf("Expected coverage of 1 to be parsed as 1, found %f", coveragePercent)
	}
}

func Test_ParseCoverageFloat(t *testing.T) {
	coverageLine := "<coverage line-rate=\"0.23893805\" branch-rate=\"0\" lines-covered=\"216\" lines-valid=\"904\" branches-covered=\"0\" branches-valid=\"0\" complexity=\"0\" version=\"\" timestamp=\"1633570973824\">"

	coveragePercent, err := parseCoveragePercent([]byte(coverageLine))
	if err != nil {
		t.Error(err)
	}

	expected, _ := strconv.ParseFloat("0.23893805", 32)
	if coveragePercent != expected {
		t.Errorf("Expected coverage percent of .23893805 to be parsed as %f, found %f", expected, coveragePercent)
	}
}

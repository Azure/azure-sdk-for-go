package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func Test_findAllSubDirectories(t *testing.T) {
	cleanTestData()
	defer cleanTestData()
	expected := []string{
		"../../testdata/scenarioa/foo/stage",
		"../../testdata/scenariob/foo/stage",
		"../../testdata/scenarioc/foo/stage",
		"../../testdata/scenariod/foo/stage",
		"../../testdata/scenarioe/foo/stage",
		"../../testdata/scenariof/foo/stage",
		"../../testdata/scenariog/foo/mgmt/2019-10-11/foo/stage",
		"../../testdata/scenarioh/foo/mgmt/2019-10-11/foo/stage",
		"../../testdata/scenarioi/foo/mgmt/2019-10-23/foo/stage",
		"../../testdata/scenarioj/foo/mgmt/2019-10-23/foo/stage",
		"../../testdata/scenariok/foo/mgmt/2019-11-01-preview/foo/stage",
		"../../testdata/scenariol/foo/mgmt/2019-11-01-preview/foo/stage",
		"../../testdata/scenariom/foo/mgmt/2019-11-01-preview/foo/stage",
		"../../testdata/scenarion/foo/mgmt/2019-11-01-preview/foo/stage",
		"../../testdata/scenarioo/foo/mgmt/2019-11-01-preview/foo/stage",
	}
	root, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatalf("error when get absolute path of root: %+v", err)
	}
	stages, err := findAllSubDirectories(root, "stage")
	if err != nil {
		t.Fatalf("error when listing all stage folders: %+v", err)
	}
	if len(stages) != len(expected) {
		t.Fatalf("expected %d stages folders, but got %d", len(expected), len(stages))
	}
	for i, stage := range stages {
		e, err := filepath.Abs(expected[i])
		if err != nil {
			t.Fatalf("error when parsing expected results '%s'(%d)", expected[i], i)
		}
		if !strings.EqualFold(stage, e) {
			t.Fatalf("expected folder '%s', but got '%s'", e, stage)
		}
	}
}

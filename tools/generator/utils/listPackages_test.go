package utils_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/generator/utils"
)

func TestListTrack1SDKPackages(t *testing.T) {
	pkgs, err := utils.ListTrack1SDKPackages("./testdata")
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expect 1 package, but got %d", len(pkgs))
	}
	expected := []string{
		"testdata/track1package",
	}
	for i, p := range pkgs {
		p = utils.NormalizePath(p)
		if p != expected[i] {
			t.Fatalf("expect package name '%s', but got '%s'", expected[i], p)
		}
	}
}

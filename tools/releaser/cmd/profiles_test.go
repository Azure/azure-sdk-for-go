package cmd

import (
	"fmt"
	"testing"
)

func TestGetExports(t *testing.T) {
	pkg, err := GetExports("../../../profiles/2017-03-09/compute/mgmt/compute")
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	fmt.Println(pkg)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package util_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/indexer/util"
)

func Test_GetIndexedPackages(t *testing.T) {
	td, err := os.OpenFile("./testdata/testdata.html", 0, 0)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer td.Close()
	ps, err := util.GetIndexedPackages(td)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if l := len(ps); l != 250 {
		t.Logf("wrong number of packages, got %v, wanted 250\n", l)
		t.Fail()
	}
	// spot-check a few packages
	pkgs := []string{
		"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization",
		"github.com/Azure/azure-sdk-for-go/services/classic/management",
		"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/luis/programmatic",
		"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2015-11-01-preview/operationalinsights",
		"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2016-09-01/web",
	}
	for _, pkg := range pkgs {
		if indexed, ok := ps[pkg]; !ok {
			t.Logf("didn't find package '%s'\n", pkg)
			t.Fail()
		} else if !indexed {
			t.Logf("package '%s' wasn't flagged as indexed", pkg)
			t.Fail()
		}
	}
}

func Test_GetPackagesForIndexing(t *testing.T) {
	ps, err := util.GetPackagesForIndexing("./testdata/github.com/Azure/azure-sdk-for-go/services")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if l := len(ps); l != 4 {
		t.Logf("wrong number of packages, got %v, wanted 4\n", l)
		t.Fail()
	}
	// verify packages
	pkgs := []string{
		"github.com/Azure/azure-sdk-for-go/services/zed/mgmt/2017-06-01/zed",
		"github.com/Azure/azure-sdk-for-go/services/zed/v1.0/zeddp",
		"github.com/Azure/azure-sdk-for-go/services/bar/mgmt/2017-01-01/bar",
		"github.com/Azure/azure-sdk-for-go/services/foo/mgmt/2018-01-01/foo",
	}
	for _, pkg := range pkgs {
		if indexed, ok := ps[pkg]; !ok {
			t.Logf("didn't find package '%s'\n", pkg)
			t.Fail()
		} else if indexed {
			t.Logf("package '%s' was flagged as indexed", pkg)
			t.Fail()
		}
	}
}

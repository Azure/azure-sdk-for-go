// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/stretchr/testify/assert"
)

func TestEnumFilter(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/enum")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/enum")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := autorest.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.EnumFilter)

	excepted := fmt.Sprint("### Breaking Changes\n\n- Type alias `EnumRemove` has been removed\n\n### Features Added\n\n- New value `EnumExistB` added to type alias `EnumExist`\n- New type alias `EnumAdd` with values `EnumAddA`, `EnumAddB`\n")
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestFuncFilter(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/operation")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/operation")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := autorest.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.FuncFilter)

	excepted := fmt.Sprint("### Breaking Changes\n\n- Function `*Client.Update` has been removed\n\n### Features Added\n\n- New function `*Client.BeginCreateOrUpdate(string, *ClientBeginCreateOrUpdateOptions) (ClientBeginCreateOrUpdateResponse, error)`\n")
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestLROFilter(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/lro")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/lro")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := autorest.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.FuncFilter, common.LROFilter)

	excepted := fmt.Sprint("### Breaking Changes\n\n- Operation `*Client.CreateOrUpdate` has been changed to LRO, use `*Client.BeginCreateOrUpdate` instead.\n")
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

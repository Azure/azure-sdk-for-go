// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common_test

import (
	"testing"

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

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.EnumFilter)

	excepted := "### Breaking Changes\n\n- Enum `EnumRemove` has been removed\n\n### Features Added\n\n- New value `EnumExistB` added to enum type `EnumExist`\n- New enum type `EnumAdd` with values `EnumAddA`, `EnumAddB`\n"
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

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.FuncFilter)

	excepted := "### Breaking Changes\n\n- Function `*Client.BeingDelete` has been removed\n- Function `*Client.NewListPager` has been removed\n- Function `*Client.Update` has been removed\n\n### Features Added\n\n- New function `*Client.BeginCreateOrUpdate(string, *ClientBeginCreateOrUpdateOptions) (ClientBeginCreateOrUpdateResponse, error)`\n- New function `*Client.NewListBySubscriptionPager(*ClientListBySubscriptionOptions) *runtime.Pager[ClientListBySubscriptionResponse]`\n"
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

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.FuncFilter, common.LROFilter)

	excepted := "### Breaking Changes\n\n- Operation `*Client.BeginDelete` has been changed to non-LRO, use `*Client.Delete` instead.\n- Operation `*Client.CreateOrUpdate` has been changed to LRO, use `*Client.BeginCreateOrUpdate` instead.\n"
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestPageableFilter(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/page")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/page")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.FuncFilter, common.PageableFilter)

	excepted := "### Breaking Changes\n\n- Operation `*Client.GetLog` has supported pagination, use `*Client.NewGetLogPager` instead.\n- Operation `*Client.NewListPager` does not support pagination anymore, use `*Client.List` instead.\n"
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestInterfaceToAnyFilter(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/interfacetoany")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/interfacetoany")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.InterfaceToAnyFilter)

	excepted := "### Breaking Changes\n\n- Type of `Interface2Any.NewType` has been changed from `interface{}` to `string`\n"
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestTypeToAny(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/toany")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/toany")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	excepted := "### Breaking Changes\n\n- Type of `Client.M` has been changed from `map[string]string` to `map[string]any`\n\n### Features Added\n\n- Type of `Client.A` has been changed from `*int` to `any`\n"
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestFuncParameterChange(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/parameter")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/parameter")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.FuncFilter)

	excepted := "### Breaking Changes\n\n- Function `*Client.AfterAny` parameter(s) have been changed from `(context.Context, string, string, interface{}, ClientOption)` to `(context.Context, string, string, any, Option)`\n- Function `*Client.BeforeAny` parameter(s) have been changed from `(context.Context, string, string, interface{}, ClientOption)` to `(context.Context, string, any, any, ClientOption)`\n"
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestNonExportedFilter(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/non-export")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/non-export")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := common.GetChangelogForPackage(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	common.FilterChangelog(changelog, common.NonExportedFilter)

	excepted := "### Breaking Changes\n\n- Function `*Public.PublicMethod` has been removed\n\n### Features Added\n\n- New function `*Public.NewPublicMethos() `\n"
	assert.Equal(t, excepted, changelog.ToCompactMarkdown())
}

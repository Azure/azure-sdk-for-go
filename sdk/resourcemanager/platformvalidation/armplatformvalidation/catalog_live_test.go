// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armplatformvalidation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlatformValidationOperations(t *testing.T) {
	clients, ctx := newTestClients(t)
	pager := clients.factory.NewOperationsClient().NewListPager(nil)
	operations := map[string]bool{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		requireNoError(t, err)
		for _, operation := range page.Value {
			operations[strings.ToLower(stringValue(operation.Name))] = true
		}
	}
	for _, resource := range []string{
		"cloudValidations",
		"cloudValidations/validationExecutionPlans",
		"cloudValidations/validationExecutionPlans/executionPlanRuns",
	} {
		require.True(t, operations[strings.ToLower("Microsoft.PlatformValidation/"+resource+"/read")])
	}
}

func TestPlatformValidationCatalogTests(t *testing.T) {
	clients, ctx := newTestClients(t)
	client := clients.factory.NewValidationTestsClient()
	pager := client.NewListBySubscriptionPager(nil)
	listedID := ""
	for pager.More() {
		page, err := pager.NextPage(ctx)
		requireNoError(t, err)
		for _, entry := range page.Value {
			if stringValue(entry.Name) == bootTestName {
				listedID = stringValue(entry.ID)
			}
		}
	}
	require.NotEmpty(t, listedID, "The approved boot validation must be discoverable")
	test, err := client.Get(ctx, bootTestName, nil)
	requireNoError(t, err)
	require.Equal(t, listedID, stringValue(test.ID))
	require.NotNil(t, test.Properties)
	require.NotEmpty(t, test.Properties.DisplayName)
	require.Equal(t, strings.ToLower("/subscriptions/"+clients.subscriptionID+bootTestID), strings.ToLower(stringValue(test.ID)))
}

func TestPlatformValidationCatalogCategories(t *testing.T) {
	clients, ctx := newTestClients(t)
	client := clients.factory.NewValidationTestCategoriesClient()
	pager := client.NewListBySubscriptionPager(nil)
	name, listedID := "", ""
	for pager.More() && name == "" {
		page, err := pager.NextPage(ctx)
		requireNoError(t, err)
		if len(page.Value) > 0 {
			name, listedID = stringValue(page.Value[0].Name), stringValue(page.Value[0].ID)
		}
	}
	require.NotEmpty(t, name)
	category, err := client.Get(ctx, name, nil)
	requireNoError(t, err)
	require.Equal(t, listedID, stringValue(category.ID))
	require.Equal(t, name, stringValue(category.Name))
	require.Equal(t, strings.ToLower("/subscriptions/"+clients.subscriptionID+
		"/providers/Microsoft.PlatformValidation/validationTestCategories/"+name), strings.ToLower(stringValue(category.ID)))
}

func TestPlatformValidationCatalogVersions(t *testing.T) {
	clients, ctx := newTestClients(t)
	client := clients.factory.NewValidationTestVersionsClient()
	pager := client.NewListPager(bootTestName, nil)
	name, listedID := "", ""
	for pager.More() && name == "" {
		page, err := pager.NextPage(ctx)
		requireNoError(t, err)
		if len(page.Value) > 0 {
			name, listedID = stringValue(page.Value[0].Name), stringValue(page.Value[0].ID)
		}
	}
	require.NotEmpty(t, name)
	version, err := client.Get(ctx, bootTestName, name, nil)
	requireNoError(t, err)
	require.Equal(t, listedID, stringValue(version.ID))
	require.Equal(t, name, stringValue(version.Name))
	require.Equal(t, strings.ToLower("/subscriptions/"+clients.subscriptionID+bootTestID+"/versions/"+name), strings.ToLower(stringValue(version.ID)))
}

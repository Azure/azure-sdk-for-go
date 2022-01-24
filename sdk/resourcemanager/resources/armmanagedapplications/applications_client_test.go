//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmanagedapplications_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armmanagedapplications"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplicationsClient_BeginCreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create application definition

	// create application
	applicationsClient := armmanagedapplications.NewApplicationsClient(subscriptionID, cred, opt)
	applicationName, _ := createRandomName(t, "application")
	_, err := applicationsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		applicationName,
		armmanagedapplications.Application{
			Location: to.StringPtr(location),
			Properties: &armmanagedapplications.ApplicationProperties{
				ManagedResourceGroupID: to.StringPtr(*rg.ID + "test"),
				//ApplicationDefinitionID: to.StringPtr(),
			},
			Kind: to.StringPtr("ServiceCatalog"),
		},
		nil,
	)
	require.Error(t, err)
	//require.NoError(t, err)
	//applicationResp, err := applicationPollerResp.PollUntilDone(ctx, 10*time.Second)
	//require.NoError(t, err)
	//require.Equal(t, applicationName, *applicationResp.Name)
}

func TestApplicationClient_ListOperations(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	applicationClient := armmanagedapplications.NewApplicationClient(cred, opt)
	pager := applicationClient.ListOperations(nil)
	require.NoError(t, pager.Err())
}

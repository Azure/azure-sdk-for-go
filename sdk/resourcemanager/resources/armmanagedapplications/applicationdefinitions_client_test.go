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

func TestApplicationDefinitionsClient_BeginCreateOrUpdate(t *testing.T) {
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
	applicationDefinitionClient := armmanagedapplications.NewApplicationDefinitionsClient(subscriptionID, cred, opt)
	applicationDefinitionName, _ := createRandomName(t, "definition")
	_, err := applicationDefinitionClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		applicationDefinitionName,
		armmanagedapplications.ApplicationDefinition{
			Location: to.StringPtr(location),
			Properties: &armmanagedapplications.ApplicationDefinitionProperties{
				Authorizations: []*armmanagedapplications.ApplicationProviderAuthorization{
					{
						//PrincipalID: to.StringPtr(),
						//RoleDefinitionID: to.StringPtr(),
					},
				},
				LockLevel:      armmanagedapplications.ApplicationLockLevelNone.ToPtr(),
				PackageFileURI: to.StringPtr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-managed-application-with-linked-templates/artifacts/ManagedAppZip/pkg.zip"),
			},
		},
		nil,
	)
	require.Error(t, err)
	//require.NoError(t, err)
	//adResp, err := adPollerResp.PollUntilDone(ctx, 10*time.Second)
	//require.NoError(t, err)
	//require.Equal(t, applicationDefinitionName, *adResp.Name)
}

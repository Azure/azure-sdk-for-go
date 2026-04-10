// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmsi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/stretchr/testify/suite"
)

type ManagedIdentityTestSuite struct {
	suite.Suite

	ctx                                     context.Context
	cred                                    azcore.TokenCredential
	options                                 *arm.ClientOptions
	armEndpoint                             string
	federatedIdentityCredentialResourceName string
	resourceName                            string
	location                                string
	resourceGroupName                       string
	subscriptionId                          string
}

func (testsuite *ManagedIdentityTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.federatedIdentityCredentialResourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "federate", 14, false)
	testsuite.resourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "resource", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ManagedIdentityTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestManagedIdentityTestSuite(t *testing.T) {
	suite.Run(t, new(ManagedIdentityTestSuite))
}

func (testsuite *ManagedIdentityTestSuite) Prepare() {
	var err error
	// From step UserAssignedIdentities_CreateOrUpdate
	fmt.Println("Call operation: UserAssignedIdentities_CreateOrUpdate")
	userAssignedIdentitiesClient, err := armmsi.NewUserAssignedIdentitiesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = userAssignedIdentitiesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armmsi.Identity{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}
func (testsuite *ManagedIdentityTestSuite) TestUserAssignedIdentities() {
	var err error
	// From step UserAssignedIdentities_ListBySubscription
	fmt.Println("Call operation: UserAssignedIdentities_ListBySubscription")
	userAssignedIdentitiesClient, err := armmsi.NewUserAssignedIdentitiesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	userAssignedIdentitiesClientNewListBySubscriptionPager := userAssignedIdentitiesClient.NewListBySubscriptionPager(nil)
	for userAssignedIdentitiesClientNewListBySubscriptionPager.More() {
		_, err := userAssignedIdentitiesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step UserAssignedIdentities_ListByResourceGroup
	fmt.Println("Call operation: UserAssignedIdentities_ListByResourceGroup")
	userAssignedIdentitiesClientNewListByResourceGroupPager := userAssignedIdentitiesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for userAssignedIdentitiesClientNewListByResourceGroupPager.More() {
		_, err := userAssignedIdentitiesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step UserAssignedIdentities_Get
	fmt.Println("Call operation: UserAssignedIdentities_Get")
	_, err = userAssignedIdentitiesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step UserAssignedIdentities_Update
	fmt.Println("Call operation: UserAssignedIdentities_Update")
	_, err = userAssignedIdentitiesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armmsi.IdentityUpdate{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}/federatedIdentityCredentials/{federatedIdentityCredentialResourceName}
func (testsuite *ManagedIdentityTestSuite) TestFederatedIdentityCredentials() {
	var err error
	// From step FederatedIdentityCredentials_CreateOrUpdate
	fmt.Println("Call operation: FederatedIdentityCredentials_CreateOrUpdate")
	federatedIdentityCredentialsClient, err := armmsi.NewFederatedIdentityCredentialsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = federatedIdentityCredentialsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.federatedIdentityCredentialResourceName, armmsi.FederatedIdentityCredential{
		Properties: &armmsi.FederatedIdentityCredentialProperties{
			Audiences: []*string{
				to.Ptr("api://AzureADTokenExchange")},
			Issuer:  to.Ptr("https://oidc.prod-aks.azure.com/TenantGUID/IssuerGUID"),
			Subject: to.Ptr("system:serviceaccount:ns:svcaccount"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step FederatedIdentityCredentials_List
	fmt.Println("Call operation: FederatedIdentityCredentials_List")
	federatedIdentityCredentialsClientNewListPager := federatedIdentityCredentialsClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, &armmsi.FederatedIdentityCredentialsClientListOptions{Top: nil,
		Skiptoken: nil,
	})
	for federatedIdentityCredentialsClientNewListPager.More() {
		_, err := federatedIdentityCredentialsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FederatedIdentityCredentials_Get
	fmt.Println("Call operation: FederatedIdentityCredentials_Get")
	_, err = federatedIdentityCredentialsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.federatedIdentityCredentialResourceName, nil)
	testsuite.Require().NoError(err)

	// From step FederatedIdentityCredentials_Delete
	fmt.Println("Call operation: FederatedIdentityCredentials_Delete")
	_, err = federatedIdentityCredentialsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.federatedIdentityCredentialResourceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ManagedIdentity/operations
func (testsuite *ManagedIdentityTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armmsi.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *ManagedIdentityTestSuite) Cleanup() {
	var err error
	// From step UserAssignedIdentities_Delete
	fmt.Println("Call operation: UserAssignedIdentities_Delete")
	userAssignedIdentitiesClient, err := armmsi.NewUserAssignedIdentitiesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = userAssignedIdentitiesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
}

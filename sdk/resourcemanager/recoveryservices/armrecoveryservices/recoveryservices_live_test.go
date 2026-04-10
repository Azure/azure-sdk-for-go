// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armrecoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices/v2"
	"github.com/stretchr/testify/suite"
)

type RecoveryservicesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	vaultName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *RecoveryservicesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.vaultName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vaultnam", 8+6, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *RecoveryservicesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestRecoveryservicesTestSuite(t *testing.T) {
	suite.Run(t, new(RecoveryservicesTestSuite))
}

// Microsoft.RecoveryServices/vaults/{vaultName}
func (testsuite *RecoveryservicesTestSuite) TestRecoveryServices() {
	var err error
	// From step Vaults_CreateOrUpdate
	fmt.Println("Call operation: Vaults_CreateOrUpdate")
	vaultsClient, err := armrecoveryservices.NewVaultsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vaultsClientCreateOrUpdateResponsePoller, err := vaultsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vaultName, armrecoveryservices.Vault{
		Location: to.Ptr(testsuite.location),
		Identity: &armrecoveryservices.IdentityData{
			Type: to.Ptr(armrecoveryservices.ResourceIdentityTypeSystemAssigned),
		},
		Properties: &armrecoveryservices.VaultProperties{
			PublicNetworkAccess: to.Ptr(armrecoveryservices.PublicNetworkAccessEnabled),
		},
		SKU: &armrecoveryservices.SKU{
			Name: to.Ptr(armrecoveryservices.SKUNameStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, vaultsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Vaults_ListBySubscriptionId
	fmt.Println("Call operation: Vaults_ListBySubscriptionID")
	vaultsClientNewListBySubscriptionIDPager := vaultsClient.NewListBySubscriptionIDPager(nil)
	for vaultsClientNewListBySubscriptionIDPager.More() {
		_, err := vaultsClientNewListBySubscriptionIDPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Vaults_ListByResourceGroup
	fmt.Println("Call operation: Vaults_ListByResourceGroup")
	vaultsClientNewListByResourceGroupPager := vaultsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for vaultsClientNewListByResourceGroupPager.More() {
		_, err := vaultsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Vaults_Get
	fmt.Println("Call operation: Vaults_Get")
	_, err = vaultsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vaultName, nil)
	testsuite.Require().NoError(err)

	// From step Vaults_Update
	fmt.Println("Call operation: Vaults_Update")
	vaultsClientUpdateResponsePoller, err := vaultsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vaultName, armrecoveryservices.PatchVault{
		Tags: map[string]*string{
			"PatchKey": to.Ptr("PatchKeyUpdated"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, vaultsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VaultExtendedInfo_CreateOrUpdate
	fmt.Println("Call operation: VaultExtendedInfo_CreateOrUpdate")
	vaultExtendedInfoClient, err := armrecoveryservices.NewVaultExtendedInfoClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = vaultExtendedInfoClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vaultName, armrecoveryservices.VaultExtendedInfoResource{
		Properties: &armrecoveryservices.VaultExtendedInfo{
			IntegrityKey: to.Ptr("myIntegrityKey"),
			Algorithm:    to.Ptr("None"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step VaultExtendedInfo_Get
	fmt.Println("Call operation: VaultExtendedInfo_Get")
	_, err = vaultExtendedInfoClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vaultName, nil)
	testsuite.Require().NoError(err)

	// From step ReplicationUsages_List
	fmt.Println("Call operation: ReplicationUsages_List")
	replicationUsagesClient, err := armrecoveryservices.NewReplicationUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	replicationUsagesClientNewListPager := replicationUsagesClient.NewListPager(testsuite.resourceGroupName, testsuite.vaultName, nil)
	for replicationUsagesClientNewListPager.More() {
		_, err := replicationUsagesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Usages_ListByVaults
	fmt.Println("Call operation: Usages_ListByVaults")
	usagesClient, err := armrecoveryservices.NewUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usagesClientNewListByVaultsPager := usagesClient.NewListByVaultsPager(testsuite.resourceGroupName, testsuite.vaultName, nil)
	for usagesClientNewListByVaultsPager.More() {
		_, err := usagesClientNewListByVaultsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armrecoveryservices.NewOperationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RecoveryServices_Capabilities
	fmt.Println("Call operation: RecoveryServices_Capabilities")
	client, err := armrecoveryservices.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.Capabilities(testsuite.ctx, testsuite.location, armrecoveryservices.ResourceCapabilities{
		Type: to.Ptr("Microsoft.RecoveryServices/Vaults"),
		Properties: &armrecoveryservices.CapabilitiesProperties{
			DNSZones: []*armrecoveryservices.DNSZone{
				{
					SubResource: to.Ptr(armrecoveryservices.VaultSubResourceTypeAzureBackup),
				},
				{
					SubResource: to.Ptr(armrecoveryservices.VaultSubResourceTypeAzureSiteRecovery),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	var privateLinkResourceName string
	// From step PrivateLinkResources_List
	fmt.Println("Call operation: PrivateLinkResources_List")
	privateLinkResourcesClient, err := armrecoveryservices.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListPager := privateLinkResourcesClient.NewListPager(testsuite.resourceGroupName, testsuite.vaultName, nil)
	for privateLinkResourcesClientNewListPager.More() {
		result, err := privateLinkResourcesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		privateLinkResourceName = *result.Value[0].Name
		break
	}

	// From step PrivateLinkResources_Get
	fmt.Println("Call operation: PrivateLinkResources_Get")
	_, err = privateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vaultName, privateLinkResourceName, nil)
	testsuite.Require().NoError(err)

	// From step Vaults_Delete
	fmt.Println("Call operation: Vaults_Delete")
	vaultsClientDeleteResponsePoller, err := vaultsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.vaultName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, vaultsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

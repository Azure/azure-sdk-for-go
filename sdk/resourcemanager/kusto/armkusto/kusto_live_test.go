// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type KustoTestSuite struct {
	suite.Suite

	ctx                               context.Context
	cred                              azcore.TokenCredential
	options                           *arm.ClientOptions
	armEndpoint                       string
	attachedClusterName               string
	attachedDatabaseConfigurationName string
	attachedDatabaseName              string
	clusterName                       string
	clusterPrincipalAssignmentName    string
	dataConnectionName                string
	databaseName                      string
	databasePrincipalAssignmentName   string
	eventhubName                      string
	eventhubNamespace                 string
	kustoId                           string
	managedPrivateEndpointName        string
	scriptName                        string
	storageAccountName                string
	azureClientId                     string
	azureTenantId                     string
	location                          string
	resourceGroupName                 string
	subscriptionId                    string
}

func (testsuite *KustoTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.attachedClusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "attachedclustern", 22, false)
	testsuite.attachedDatabaseConfigurationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "attached", 14, false)
	testsuite.attachedDatabaseName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "attacheddatabase", 22, false)
	testsuite.clusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clustern", 14, false)
	testsuite.clusterPrincipalAssignmentName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clusterprincipa", 21, false)
	testsuite.dataConnectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "dataconn", 14, false)
	testsuite.databaseName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "database", 14, false)
	testsuite.databasePrincipalAssignmentName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "databaseprincipa", 22, false)
	testsuite.eventhubName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventhubn", 15, false)
	testsuite.eventhubNamespace, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventhubnamespace", 23, false)
	testsuite.managedPrivateEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "managedp", 14, false)
	testsuite.scriptName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "scriptna", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storagea", 14, true)
	testsuite.azureClientId = recording.GetEnvVariable("AZURE_CLIENT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.azureTenantId = recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *KustoTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestKustoTestSuite(t *testing.T) {
	suite.Run(t, new(KustoTestSuite))
}

func (testsuite *KustoTestSuite) Prepare() {
	var err error
	// From step Clusters_CheckNameAvailability
	fmt.Println("Call operation: Clusters_CheckNameAvailability")
	clustersClient, err := armkusto.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = clustersClient.CheckNameAvailability(testsuite.ctx, testsuite.location, armkusto.ClusterCheckNameRequest{
		Name: to.Ptr(testsuite.clusterName),
		Type: to.Ptr("Microsoft.Kusto/clusters"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Clusters_CreateOrUpdate
	fmt.Println("Call operation: Clusters_CreateOrUpdate")
	clustersClientCreateOrUpdateResponsePoller, err := clustersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.Cluster{
		Location: to.Ptr(testsuite.location),
		Identity: &armkusto.Identity{
			Type: to.Ptr(armkusto.IdentityTypeSystemAssigned),
		},
		Properties: &armkusto.ClusterProperties{
			AllowedIPRangeList: []*string{
				to.Ptr("0.0.0.0/0")},
			EnableAutoStop:         to.Ptr(true),
			EnableDoubleEncryption: to.Ptr(false),
			EnablePurge:            to.Ptr(true),
			EnableStreamingIngest:  to.Ptr(true),
			LanguageExtensions: &armkusto.LanguageExtensionsList{
				Value: []*armkusto.LanguageExtension{
					{
						LanguageExtensionImageName: to.Ptr(armkusto.LanguageExtensionImageNamePython3108),
						LanguageExtensionName:      to.Ptr(armkusto.LanguageExtensionNamePYTHON),
					},
					{
						LanguageExtensionImageName: to.Ptr(armkusto.LanguageExtensionImageNameR),
						LanguageExtensionName:      to.Ptr(armkusto.LanguageExtensionNameR),
					}},
			},
			PublicIPType:        to.Ptr(armkusto.PublicIPTypeDualStack),
			PublicNetworkAccess: to.Ptr(armkusto.PublicNetworkAccessEnabled),
		},
		SKU: &armkusto.AzureSKU{
			Name:     to.Ptr(armkusto.AzureSKUNameStandardL16AsV3),
			Capacity: to.Ptr[int32](2),
			Tier:     to.Ptr(armkusto.AzureSKUTierStandard),
		},
	}, &armkusto.ClustersClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	var clustersClientCreateOrUpdateResponse *armkusto.ClustersClientCreateOrUpdateResponse
	clustersClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.kustoId = *clustersClientCreateOrUpdateResponse.ID

	// From step Databases_CheckNameAvailability
	fmt.Println("Call operation: Databases_CheckNameAvailability")
	databasesClient, err := armkusto.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = databasesClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.CheckNameRequest{
		Name: to.Ptr(testsuite.databaseName),
		Type: to.Ptr(armkusto.TypeMicrosoftKustoClustersDatabases),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Databases_CreateOrUpdate
	fmt.Println("Call operation: Databases_CreateOrUpdate")
	databasesClientCreateOrUpdateResponsePoller, err := databasesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, &armkusto.ReadWriteDatabase{
		Kind:     to.Ptr(armkusto.KindReadWrite),
		Location: to.Ptr(testsuite.location),
		Properties: &armkusto.ReadWriteDatabaseProperties{
			SoftDeletePeriod: to.Ptr("P1D"),
		},
	}, &armkusto.DatabasesClientBeginCreateOrUpdateOptions{CallerRole: to.Ptr(armkusto.CallerRoleAdmin)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/clusters/{clusterName}
func (testsuite *KustoTestSuite) TestClusters() {
	var err error
	// From step Clusters_List
	fmt.Println("Call operation: Clusters_List")
	clustersClient, err := armkusto.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientNewListPager := clustersClient.NewListPager(nil)
	for clustersClientNewListPager.More() {
		_, err := clustersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_ListSkus
	fmt.Println("Call operation: Clusters_ListSkus")
	clustersClientNewListSKUsPager := clustersClient.NewListSKUsPager(nil)
	for clustersClientNewListSKUsPager.More() {
		_, err := clustersClientNewListSKUsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_ListByResourceGroup
	fmt.Println("Call operation: Clusters_ListByResourceGroup")
	clustersClientNewListByResourceGroupPager := clustersClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for clustersClientNewListByResourceGroupPager.More() {
		_, err := clustersClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_Get
	fmt.Println("Call operation: Clusters_Get")
	_, err = clustersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)

	// From step Clusters_ListSkusByResource
	fmt.Println("Call operation: Clusters_ListSkusByResource")
	clustersClientNewListSKUsByResourcePager := clustersClient.NewListSKUsByResourcePager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for clustersClientNewListSKUsByResourcePager.More() {
		_, err := clustersClientNewListSKUsByResourcePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_ListOutboundNetworkDependenciesEndpoints
	fmt.Println("Call operation: Clusters_ListOutboundNetworkDependenciesEndpoints")
	clustersClientNewListOutboundNetworkDependenciesEndpointsPager := clustersClient.NewListOutboundNetworkDependenciesEndpointsPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for clustersClientNewListOutboundNetworkDependenciesEndpointsPager.More() {
		_, err := clustersClientNewListOutboundNetworkDependenciesEndpointsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_Update
	fmt.Println("Call operation: Clusters_Update")
	clustersClientUpdateResponsePoller, err := clustersClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.ClusterUpdate{
		Location: to.Ptr(testsuite.location),
	}, &armkusto.ClustersClientBeginUpdateOptions{IfMatch: to.Ptr("*")})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_ListFollowerDatabases
	fmt.Println("Call operation: Clusters_ListFollowerDatabases")
	clustersClientNewListFollowerDatabasesPager := clustersClient.NewListFollowerDatabasesPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for clustersClientNewListFollowerDatabasesPager.More() {
		_, err := clustersClientNewListFollowerDatabasesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_ListLanguageExtensions
	fmt.Println("Call operation: Clusters_ListLanguageExtensions")
	clustersClientNewListLanguageExtensionsPager := clustersClient.NewListLanguageExtensionsPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for clustersClientNewListLanguageExtensionsPager.More() {
		_, err := clustersClientNewListLanguageExtensionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_DiagnoseVirtualNetwork
	fmt.Println("Call operation: Clusters_DiagnoseVirtualNetwork")
	clustersClientDiagnoseVirtualNetworkResponsePoller, err := clustersClient.BeginDiagnoseVirtualNetwork(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientDiagnoseVirtualNetworkResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_AddLanguageExtensions
	fmt.Println("Call operation: Clusters_AddLanguageExtensions")
	clustersClientAddLanguageExtensionsResponsePoller, err := clustersClient.BeginAddLanguageExtensions(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.LanguageExtensionsList{
		Value: []*armkusto.LanguageExtension{
			{
				LanguageExtensionName: to.Ptr(armkusto.LanguageExtensionNamePYTHON),
			},
			{
				LanguageExtensionName: to.Ptr(armkusto.LanguageExtensionNameR),
			}},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientAddLanguageExtensionsResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_RemoveLanguageExtensions
	fmt.Println("Call operation: Clusters_RemoveLanguageExtensions")
	clustersClientRemoveLanguageExtensionsResponsePoller, err := clustersClient.BeginRemoveLanguageExtensions(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.LanguageExtensionsList{
		Value: []*armkusto.LanguageExtension{
			{
				LanguageExtensionName: to.Ptr(armkusto.LanguageExtensionNamePYTHON),
			},
			{
				LanguageExtensionName: to.Ptr(armkusto.LanguageExtensionNameR),
			}},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientRemoveLanguageExtensionsResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_Stop
	fmt.Println("Call operation: Clusters_Stop")
	clustersClientStopResponsePoller, err := clustersClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_Start
	fmt.Println("Call operation: Clusters_Start")
	clustersClientStartResponsePoller, err := clustersClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientStartResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/clusters/{clusterName}/attachedDatabaseConfigurations/{attachedDatabaseConfigurationName}
func (testsuite *KustoTestSuite) TestAttachedDatabaseConfigurations() {
	var attachedKustoId string
	var err error
	// From step Clusters_CreateOrUpdate
	fmt.Println("Call operation: Clusters_CreateOrUpdate")
	clustersClient, err := armkusto.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientCreateOrUpdateResponsePoller, err := clustersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.attachedClusterName, armkusto.Cluster{
		Location: to.Ptr(testsuite.location),
		Identity: &armkusto.Identity{
			Type: to.Ptr(armkusto.IdentityTypeSystemAssigned),
		},
		Properties: &armkusto.ClusterProperties{
			EnableStreamingIngest: to.Ptr(true),
		},
		SKU: &armkusto.AzureSKU{
			Name:     to.Ptr(armkusto.AzureSKUNameStandardL16AsV3),
			Capacity: to.Ptr[int32](2),
			Tier:     to.Ptr(armkusto.AzureSKUTierStandard),
		},
	}, &armkusto.ClustersClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	var clustersClientCreateOrUpdateResponse *armkusto.ClustersClientCreateOrUpdateResponse
	clustersClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	attachedKustoId = *clustersClientCreateOrUpdateResponse.ID

	// From step Databases_CreateOrUpdate
	fmt.Println("Call operation: Databases_CreateOrUpdate")
	databasesClient, err := armkusto.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databasesClientCreateOrUpdateResponsePoller, err := databasesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.attachedClusterName, testsuite.attachedDatabaseName, &armkusto.ReadWriteDatabase{
		Kind:     to.Ptr(armkusto.KindReadWrite),
		Location: to.Ptr(testsuite.location),
		Properties: &armkusto.ReadWriteDatabaseProperties{
			SoftDeletePeriod: to.Ptr("P1D"),
		},
	}, &armkusto.DatabasesClientBeginCreateOrUpdateOptions{CallerRole: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AttachedDatabaseConfigurations_CheckNameAvailability
	fmt.Println("Call operation: AttachedDatabaseConfigurations_CheckNameAvailability")
	attachedDatabaseConfigurationsClient, err := armkusto.NewAttachedDatabaseConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = attachedDatabaseConfigurationsClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.AttachedDatabaseConfigurationsCheckNameRequest{
		Name: to.Ptr("adc1"),
		Type: to.Ptr("Microsoft.Kusto/clusters/attachedDatabaseConfigurations"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step AttachedDatabaseConfigurations_CreateOrUpdate
	fmt.Println("Call operation: AttachedDatabaseConfigurations_CreateOrUpdate")
	attachedDatabaseConfigurationsClientCreateOrUpdateResponsePoller, err := attachedDatabaseConfigurationsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.attachedDatabaseConfigurationName, armkusto.AttachedDatabaseConfiguration{
		Location: to.Ptr(testsuite.location),
		Properties: &armkusto.AttachedDatabaseConfigurationProperties{
			ClusterResourceID:                 to.Ptr(attachedKustoId),
			DatabaseName:                      to.Ptr(testsuite.attachedDatabaseName),
			DatabaseNameOverride:              to.Ptr("overridekustodatabase"),
			DefaultPrincipalsModificationKind: to.Ptr(armkusto.DefaultPrincipalsModificationKindUnion),
			TableLevelSharingProperties: &armkusto.TableLevelSharingProperties{
				ExternalTablesToExclude: []*string{
					to.Ptr("ExternalTable2")},
				ExternalTablesToInclude: []*string{
					to.Ptr("ExternalTable1")},
				MaterializedViewsToExclude: []*string{
					to.Ptr("MaterializedViewTable2")},
				MaterializedViewsToInclude: []*string{
					to.Ptr("MaterializedViewTable1")},
				TablesToExclude: []*string{
					to.Ptr("Table2")},
				TablesToInclude: []*string{
					to.Ptr("Table1")},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, attachedDatabaseConfigurationsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AttachedDatabaseConfigurations_ListByCluster
	fmt.Println("Call operation: AttachedDatabaseConfigurations_ListByCluster")
	attachedDatabaseConfigurationsClientNewListByClusterPager := attachedDatabaseConfigurationsClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for attachedDatabaseConfigurationsClientNewListByClusterPager.More() {
		_, err := attachedDatabaseConfigurationsClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AttachedDatabaseConfigurations_Get
	fmt.Println("Call operation: AttachedDatabaseConfigurations_Get")
	_, err = attachedDatabaseConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.attachedDatabaseConfigurationName, nil)
	testsuite.Require().NoError(err)

	// From step AttachedDatabaseConfigurations_Delete
	fmt.Println("Call operation: AttachedDatabaseConfigurations_Delete")
	attachedDatabaseConfigurationsClientDeleteResponsePoller, err := attachedDatabaseConfigurationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.attachedDatabaseConfigurationName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, attachedDatabaseConfigurationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/clusters/{clusterName}/principalAssignments/{principalAssignmentName}
func (testsuite *KustoTestSuite) TestClusterPrincipalAssignments() {
	principalAssignmentName := testsuite.clusterPrincipalAssignmentName
	var err error
	// From step ClusterPrincipalAssignments_CheckNameAvailability
	fmt.Println("Call operation: ClusterPrincipalAssignments_CheckNameAvailability")
	clusterPrincipalAssignmentsClient, err := armkusto.NewClusterPrincipalAssignmentsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = clusterPrincipalAssignmentsClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.ClusterPrincipalAssignmentCheckNameRequest{
		Name: to.Ptr("kustoprincipal1"),
		Type: to.Ptr("Microsoft.Kusto/clusters/principalAssignments"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step ClusterPrincipalAssignments_CreateOrUpdate
	fmt.Println("Call operation: ClusterPrincipalAssignments_CreateOrUpdate")
	clusterPrincipalAssignmentsClientCreateOrUpdateResponsePoller, err := clusterPrincipalAssignmentsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, principalAssignmentName, armkusto.ClusterPrincipalAssignment{
		Properties: &armkusto.ClusterPrincipalProperties{
			PrincipalID:   to.Ptr(testsuite.azureClientId),
			PrincipalType: to.Ptr(armkusto.PrincipalTypeApp),
			Role:          to.Ptr(armkusto.ClusterPrincipalRoleAllDatabasesAdmin),
			TenantID:      to.Ptr(testsuite.azureTenantId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clusterPrincipalAssignmentsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ClusterPrincipalAssignments_List
	fmt.Println("Call operation: ClusterPrincipalAssignments_List")
	clusterPrincipalAssignmentsClientNewListPager := clusterPrincipalAssignmentsClient.NewListPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for clusterPrincipalAssignmentsClientNewListPager.More() {
		_, err := clusterPrincipalAssignmentsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ClusterPrincipalAssignments_Get
	fmt.Println("Call operation: ClusterPrincipalAssignments_Get")
	_, err = clusterPrincipalAssignmentsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, principalAssignmentName, nil)
	testsuite.Require().NoError(err)

	// From step ClusterPrincipalAssignments_Delete
	fmt.Println("Call operation: ClusterPrincipalAssignments_Delete")
	clusterPrincipalAssignmentsClientDeleteResponsePoller, err := clusterPrincipalAssignmentsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, principalAssignmentName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clusterPrincipalAssignmentsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}
func (testsuite *KustoTestSuite) TestDatabases() {
	var err error
	// From step Databases_ListByCluster
	fmt.Println("Call operation: Databases_ListByCluster")
	databasesClient, err := armkusto.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databasesClientNewListByClusterPager := databasesClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, &armkusto.DatabasesClientListByClusterOptions{Top: nil,
		Skiptoken: nil,
	})
	for databasesClientNewListByClusterPager.More() {
		_, err := databasesClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Databases_Get
	fmt.Println("Call operation: Databases_Get")
	_, err = databasesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)

	// From step Databases_Update
	fmt.Println("Call operation: Databases_Update")
	databasesClientUpdateResponsePoller, err := databasesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, &armkusto.ReadWriteDatabase{
		Kind: to.Ptr(armkusto.KindReadWrite),
		Properties: &armkusto.ReadWriteDatabaseProperties{
			HotCachePeriod: to.Ptr("P1D"),
		},
	}, &armkusto.DatabasesClientBeginUpdateOptions{CallerRole: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Databases_ListPrincipals
	fmt.Println("Call operation: Databases_ListPrincipals")
	databasesClientNewListPrincipalsPager := databasesClient.NewListPrincipalsPager(testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	for databasesClientNewListPrincipalsPager.More() {
		_, err := databasesClientNewListPrincipalsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Database_InviteFollower
	fmt.Println("Call operation: Database_InviteFollower")
	databaseClient, err := armkusto.NewDatabaseClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = databaseClient.InviteFollower(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, armkusto.DatabaseInviteFollowerRequest{
		InviteeEmail: to.Ptr("invitee@contoso.com"),
		TableLevelSharingProperties: &armkusto.TableLevelSharingProperties{
			ExternalTablesToExclude: []*string{},
			ExternalTablesToInclude: []*string{
				to.Ptr("ExternalTable*")},
			FunctionsToExclude: []*string{
				to.Ptr("functionsToExclude2")},
			FunctionsToInclude: []*string{
				to.Ptr("functionsToInclude1")},
			MaterializedViewsToExclude: []*string{
				to.Ptr("MaterializedViewTable2")},
			MaterializedViewsToInclude: []*string{
				to.Ptr("MaterializedViewTable1")},
			TablesToExclude: []*string{
				to.Ptr("Table2")},
			TablesToInclude: []*string{
				to.Ptr("Table1")},
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/dataConnections/{dataConnectionName}
func (testsuite *KustoTestSuite) TestDataConnections() {
	var eventhubId string
	var err error
	// From step Create_Eventhub
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"eventhubId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.EventHub/namespaces/eventhubs', parameters('eventhubNamespace'), parameters('eventhubName'))]",
			},
		},
		"parameters": map[string]any{
			"eventhubName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.eventhubName,
			},
			"eventhubNamespace": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.eventhubNamespace,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('eventhubNamespace')]",
				"type":       "Microsoft.EventHub/namespaces",
				"apiVersion": "2022-10-01-preview",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"disableLocalAuth":       false,
					"isAutoInflateEnabled":   false,
					"kafkaEnabled":           true,
					"maximumThroughputUnits": float64(0),
					"minimumTlsVersion":      "1.2",
					"publicNetworkAccess":    "Enabled",
					"zoneRedundant":          true,
				},
				"sku": map[string]any{
					"name":     "Standard",
					"capacity": float64(1),
					"tier":     "Standard",
				},
			},
			map[string]any{
				"name":       "[concat(parameters('eventhubNamespace'), '/', parameters('eventhubName'))]",
				"type":       "Microsoft.EventHub/namespaces/eventhubs",
				"apiVersion": "2022-10-01-preview",
				"dependsOn": []any{
					"[resourceId('Microsoft.EventHub/namespaces', parameters('eventhubNamespace'))]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"messageRetentionInDays": float64(4),
					"partitionCount":         float64(4),
					"retentionDescription": map[string]any{
						"cleanupPolicy":        "Delete",
						"retentionTimeInHours": float64(4),
					},
					"status": "Active",
				},
			},
			map[string]any{
				"name":       "[concat(parameters('eventhubNamespace'), '/', parameters('eventhubName'), '/$Default')]",
				"type":       "Microsoft.EventHub/namespaces/eventhubs/consumergroups",
				"apiVersion": "2022-10-01-preview",
				"dependsOn": []any{
					"[resourceId('Microsoft.EventHub/namespaces/eventhubs', parameters('eventhubNamespace'), parameters('eventhubName'))]",
					"[resourceId('Microsoft.EventHub/namespaces', parameters('eventhubNamespace'))]",
				},
				"location":   "[parameters('location')]",
				"properties": map[string]any{},
			},
		},
		"variables": map[string]any{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_Eventhub", &deployment)
	testsuite.Require().NoError(err)
	eventhubId = deploymentExtend.Properties.Outputs.(map[string]interface{})["eventhubId"].(map[string]interface{})["value"].(string)

	// From step DataConnections_CheckNameAvailability
	fmt.Println("Call operation: DataConnections_CheckNameAvailability")
	dataConnectionsClient, err := armkusto.NewDataConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = dataConnectionsClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, armkusto.DataConnectionCheckNameRequest{
		Name: to.Ptr("DataConnections8"),
		Type: to.Ptr("Microsoft.Kusto/clusters/databases/dataConnections"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step DataConnections_CreateOrUpdate
	fmt.Println("Call operation: DataConnections_CreateOrUpdate")
	dataConnectionsClientCreateOrUpdateResponsePoller, err := dataConnectionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, testsuite.dataConnectionName, &armkusto.EventHubDataConnection{
		Kind:     to.Ptr(armkusto.DataConnectionKindEventHub),
		Location: to.Ptr(testsuite.location),
		Properties: &armkusto.EventHubConnectionProperties{
			ConsumerGroup:      to.Ptr("$Default"),
			EventHubResourceID: to.Ptr(eventhubId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dataConnectionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DataConnections_ListByDatabase
	fmt.Println("Call operation: DataConnections_ListByDatabase")
	dataConnectionsClientNewListByDatabasePager := dataConnectionsClient.NewListByDatabasePager(testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	for dataConnectionsClientNewListByDatabasePager.More() {
		_, err := dataConnectionsClientNewListByDatabasePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataConnections_Get
	fmt.Println("Call operation: DataConnections_Get")
	_, err = dataConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, testsuite.dataConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step DataConnections_Update
	fmt.Println("Call operation: DataConnections_Update")
	dataConnectionsClientUpdateResponsePoller, err := dataConnectionsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, testsuite.dataConnectionName, &armkusto.EventHubDataConnection{
		Kind:     to.Ptr(armkusto.DataConnectionKindEventHub),
		Location: to.Ptr(testsuite.location),
		Properties: &armkusto.EventHubConnectionProperties{
			ConsumerGroup:      to.Ptr("$Default"),
			EventHubResourceID: to.Ptr(eventhubId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dataConnectionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DataConnections_Delete
	fmt.Println("Call operation: DataConnections_Delete")
	dataConnectionsClientDeleteResponsePoller, err := dataConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, testsuite.dataConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dataConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/principalAssignments/{principalAssignmentName}
func (testsuite *KustoTestSuite) TestDatabasePrincipalAssignments() {
	principalAssignmentName := testsuite.databasePrincipalAssignmentName
	var err error
	// From step DatabasePrincipalAssignments_CheckNameAvailability
	fmt.Println("Call operation: DatabasePrincipalAssignments_CheckNameAvailability")
	databasePrincipalAssignmentsClient, err := armkusto.NewDatabasePrincipalAssignmentsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = databasePrincipalAssignmentsClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, armkusto.DatabasePrincipalAssignmentCheckNameRequest{
		Name: to.Ptr("kustoprincipal1"),
		Type: to.Ptr("Microsoft.Kusto/clusters/databases/principalAssignments"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step DatabasePrincipalAssignments_CreateOrUpdate
	fmt.Println("Call operation: DatabasePrincipalAssignments_CreateOrUpdate")
	databasePrincipalAssignmentsClientCreateOrUpdateResponsePoller, err := databasePrincipalAssignmentsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, principalAssignmentName, armkusto.DatabasePrincipalAssignment{
		Properties: &armkusto.DatabasePrincipalProperties{
			PrincipalID:   to.Ptr(testsuite.azureClientId),
			PrincipalType: to.Ptr(armkusto.PrincipalTypeApp),
			Role:          to.Ptr(armkusto.DatabasePrincipalRoleUser),
			TenantID:      to.Ptr(testsuite.azureTenantId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasePrincipalAssignmentsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DatabasePrincipalAssignments_List
	fmt.Println("Call operation: DatabasePrincipalAssignments_List")
	databasePrincipalAssignmentsClientNewListPager := databasePrincipalAssignmentsClient.NewListPager(testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	for databasePrincipalAssignmentsClientNewListPager.More() {
		_, err := databasePrincipalAssignmentsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DatabasePrincipalAssignments_Get
	fmt.Println("Call operation: DatabasePrincipalAssignments_Get")
	_, err = databasePrincipalAssignmentsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, principalAssignmentName, nil)
	testsuite.Require().NoError(err)

	// From step DatabasePrincipalAssignments_Delete
	fmt.Println("Call operation: DatabasePrincipalAssignments_Delete")
	databasePrincipalAssignmentsClientDeleteResponsePoller, err := databasePrincipalAssignmentsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, principalAssignmentName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasePrincipalAssignmentsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/clusters/{clusterName}/managedPrivateEndpoints/{managedPrivateEndpointName}
func (testsuite *KustoTestSuite) TestManagedPrivateEndpoints() {
	var storageAccountId string
	var err error
	// From step Create_StorageAccount
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"storageAccountId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"storageAccountName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.storageAccountName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('storageAccountName')]",
				"type":       "Microsoft.Storage/storageAccounts",
				"apiVersion": "2022-05-01",
				"kind":       "StorageV2",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"accessTier":                   "Hot",
					"allowBlobPublicAccess":        true,
					"allowCrossTenantReplication":  true,
					"allowSharedKeyAccess":         true,
					"defaultToOAuthAuthentication": false,
					"dnsEndpointType":              "Standard",
					"encryption": map[string]any{
						"keySource":                       "Microsoft.Storage",
						"requireInfrastructureEncryption": false,
						"services": map[string]any{
							"blob": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
							"file": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
						},
					},
					"minimumTlsVersion": "TLS1_2",
					"networkAcls": map[string]any{
						"bypass":              "AzureServices",
						"defaultAction":       "Allow",
						"ipRules":             []any{},
						"virtualNetworkRules": []any{},
					},
					"publicNetworkAccess":      "Enabled",
					"supportsHttpsTrafficOnly": true,
				},
				"sku": map[string]any{
					"name": "Standard_RAGRS",
					"tier": "Standard",
				},
			},
		},
		"variables": map[string]any{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_StorageAccount", &deployment)
	testsuite.Require().NoError(err)
	storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)

	// From step ManagedPrivateEndpoints_CheckNameAvailability
	fmt.Println("Call operation: ManagedPrivateEndpoints_CheckNameAvailability")
	managedPrivateEndpointsClient, err := armkusto.NewManagedPrivateEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managedPrivateEndpointsClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armkusto.ManagedPrivateEndpointsCheckNameRequest{
		Name: to.Ptr("pme1"),
		Type: to.Ptr("Microsoft.Kusto/clusters/managedPrivateEndpoints"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step ManagedPrivateEndpoints_CreateOrUpdate
	fmt.Println("Call operation: ManagedPrivateEndpoints_CreateOrUpdate")
	managedPrivateEndpointsClientCreateOrUpdateResponsePoller, err := managedPrivateEndpointsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.managedPrivateEndpointName, armkusto.ManagedPrivateEndpoint{
		Properties: &armkusto.ManagedPrivateEndpointProperties{
			GroupID:               to.Ptr("blob"),
			PrivateLinkResourceID: to.Ptr(storageAccountId),
			RequestMessage:        to.Ptr("Please Approve."),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedPrivateEndpointsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ManagedPrivateEndpoints_List
	fmt.Println("Call operation: ManagedPrivateEndpoints_List")
	managedPrivateEndpointsClientNewListPager := managedPrivateEndpointsClient.NewListPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for managedPrivateEndpointsClientNewListPager.More() {
		_, err := managedPrivateEndpointsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagedPrivateEndpoints_Get
	fmt.Println("Call operation: ManagedPrivateEndpoints_Get")
	_, err = managedPrivateEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.managedPrivateEndpointName, nil)
	testsuite.Require().NoError(err)

	// From step ManagedPrivateEndpoints_Update
	fmt.Println("Call operation: ManagedPrivateEndpoints_Update")
	managedPrivateEndpointsClientUpdateResponsePoller, err := managedPrivateEndpointsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.managedPrivateEndpointName, armkusto.ManagedPrivateEndpoint{
		Properties: &armkusto.ManagedPrivateEndpointProperties{
			GroupID:               to.Ptr("blob"),
			PrivateLinkResourceID: to.Ptr(storageAccountId),
			RequestMessage:        to.Ptr("Please Approve."),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedPrivateEndpointsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ManagedPrivateEndpoints_Delete
	fmt.Println("Call operation: ManagedPrivateEndpoints_Delete")
	managedPrivateEndpointsClientDeleteResponsePoller, err := managedPrivateEndpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.managedPrivateEndpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedPrivateEndpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Kusto/operations
func (testsuite *KustoTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armkusto.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Kusto/locations/{location}/skus
func (testsuite *KustoTestSuite) TestSkus() {
	var err error
	// From step Skus_List
	fmt.Println("Call operation: Skus_List")
	sKUsClient, err := armkusto.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListPager := sKUsClient.NewListPager(testsuite.location, nil)
	for sKUsClientNewListPager.More() {
		_, err := sKUsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *KustoTestSuite) Cleanup() {
	var err error
	// From step Databases_Delete
	fmt.Println("Call operation: Databases_Delete")
	databasesClient, err := armkusto.NewDatabasesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	databasesClientDeleteResponsePoller, err := databasesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.databaseName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, databasesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_Delete
	fmt.Println("Call operation: Clusters_Delete")
	clustersClient, err := armkusto.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientDeleteResponsePoller, err := clustersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

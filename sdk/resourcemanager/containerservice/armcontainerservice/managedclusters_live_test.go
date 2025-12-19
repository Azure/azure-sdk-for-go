// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcontainerservice_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ManagedclustersTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	agentPoolName     string
	agentpoolId       string
	configName        string
	resourceName      string
	azureClientId     string
	azureClientSecret string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ManagedclustersTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.agentPoolName = "agentpoolna"
	testsuite.configName = "configna"
	testsuite.resourceName = "containerservicena"
	testsuite.azureClientId = recording.GetEnvVariable("AZURE_CLIENT_ID", "000000000000")
	testsuite.azureClientSecret = recording.GetEnvVariable("AZURE_CLIENT_SECRET", "000000000000")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	testutil.StartRecording(testsuite.T(), pathToPackage)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ManagedclustersTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestManagedclustersTestSuite(t *testing.T) {
	suite.Run(t, new(ManagedclustersTestSuite))
}

func (testsuite *ManagedclustersTestSuite) Prepare() {
	var err error
	// From step ManagedClusters_CreateOrUpdate
	fmt.Println("Call operation: ManagedClusters_CreateOrUpdate")
	managedClustersClient, err := armcontainerservice.NewManagedClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	managedClustersClientCreateOrUpdateResponsePoller, err := managedClustersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armcontainerservice.ManagedCluster{
		Location: to.Ptr(testsuite.location),
		Properties: &armcontainerservice.ManagedClusterProperties{
			AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
				{
					Type:               to.Ptr(armcontainerservice.AgentPoolTypeVirtualMachineScaleSets),
					Count:              to.Ptr[int32](3),
					EnableAutoScaling:  to.Ptr(true),
					EnableNodePublicIP: to.Ptr(true),
					MaxCount:           to.Ptr[int32](10),
					MaxPods:            to.Ptr[int32](110),
					MinCount:           to.Ptr[int32](1),
					Mode:               to.Ptr(armcontainerservice.AgentPoolModeSystem),
					OSType:             to.Ptr(armcontainerservice.OSTypeLinux),
					VMSize:             to.Ptr("Standard_DS2_v2"),
					Name:               to.Ptr("nodepool1"),
				}},
			DNSPrefix: to.Ptr("dnsprefix1"),
			ServicePrincipalProfile: &armcontainerservice.ManagedClusterServicePrincipalProfile{
				ClientID: to.Ptr(testsuite.azureClientId),
				Secret:   to.Ptr(testsuite.azureClientSecret),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedClustersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AgentPools_CreateOrUpdate
	fmt.Println("Call operation: AgentPools_CreateOrUpdate")
	agentPoolsClient, err := armcontainerservice.NewAgentPoolsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	agentPoolsClientCreateOrUpdateResponsePoller, err := agentPoolsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.agentPoolName, armcontainerservice.AgentPool{
		Properties: &armcontainerservice.ManagedClusterAgentPoolProfileProperties{
			Count: to.Ptr[int32](1),
			Mode:  to.Ptr(armcontainerservice.AgentPoolModeUser),
			NodeLabels: map[string]*string{
				"key1": to.Ptr("val1"),
			},
			NodeTaints: []*string{
				to.Ptr("Key1=Value1:NoSchedule")},
			OrchestratorVersion:    to.Ptr(""),
			OSType:                 to.Ptr(armcontainerservice.OSTypeLinux),
			ScaleSetEvictionPolicy: to.Ptr(armcontainerservice.ScaleSetEvictionPolicyDelete),
			ScaleSetPriority:       to.Ptr(armcontainerservice.ScaleSetPrioritySpot),
			Tags: map[string]*string{
				"name1": to.Ptr("val1"),
			},
			VMSize: to.Ptr("Standard_DS2_v2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var agentPoolsClientCreateOrUpdateResponse *armcontainerservice.AgentPoolsClientCreateOrUpdateResponse
	agentPoolsClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.agentpoolId = *agentPoolsClientCreateOrUpdateResponse.ID

	// From step Snapshots_CreateOrUpdate
	fmt.Println("Call operation: Snapshots_CreateOrUpdate")
	snapshotsClient, err := armcontainerservice.NewSnapshotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = snapshotsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armcontainerservice.Snapshot{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("val1"),
			"key2": to.Ptr("val2"),
		},
		Properties: &armcontainerservice.SnapshotProperties{
			CreationData: &armcontainerservice.CreationData{
				SourceResourceID: to.Ptr(testsuite.agentpoolId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerService/managedClusters
func (testsuite *ManagedclustersTestSuite) TestContainserservice() {
	var commandId string
	var err error

	// From step ManagedClusters_List
	fmt.Println("Call operation: ManagedClusters_List")
	managedClustersClient, err := armcontainerservice.NewManagedClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	managedClustersClientNewListPager := managedClustersClient.NewListPager(nil)
	for managedClustersClientNewListPager.More() {
		_, err := managedClustersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagedClusters_ListByResourceGroup
	fmt.Println("Call operation: ManagedClusters_ListByResourceGroup")
	managedClustersClientNewListByResourceGroupPager := managedClustersClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for managedClustersClientNewListByResourceGroupPager.More() {
		_, err := managedClustersClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagedClusters_ListOutboundNetworkDependenciesEndpoints
	fmt.Println("Call operation: ManagedClusters_ListOutboundNetworkDependenciesEndpoints")
	managedClustersClientNewListOutboundNetworkDependenciesEndpointsPager := managedClustersClient.NewListOutboundNetworkDependenciesEndpointsPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for managedClustersClientNewListOutboundNetworkDependenciesEndpointsPager.More() {
		_, err := managedClustersClientNewListOutboundNetworkDependenciesEndpointsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagedClusters_Get
	fmt.Println("Call operation: ManagedClusters_Get")
	_, err = managedClustersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step ManagedClusters_GetUpgradeProfile
	fmt.Println("Call operation: ManagedClusters_GetUpgradeProfile")
	_, err = managedClustersClient.GetUpgradeProfile(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step ManagedClusters_UpdateTags
	fmt.Println("Call operation: ManagedClusters_UpdateTags")
	managedClustersClientUpdateTagsResponsePoller, err := managedClustersClient.BeginUpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armcontainerservice.TagsObject{
		Tags: map[string]*string{
			"archv3": to.Ptr(""),
			"tier":   to.Ptr("testing"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedClustersClientUpdateTagsResponsePoller)
	testsuite.Require().NoError(err)

	// From step ManagedClusters_GetAccessProfile
	fmt.Println("Call operation: ManagedClusters_GetAccessProfile")
	_, err = managedClustersClient.GetAccessProfile(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, "clusterUser", nil)
	testsuite.Require().NoError(err)

	// From step ManagedClusters_ListClusterUserCredentials
	fmt.Println("Call operation: ManagedClusters_ListClusterUserCredentials")
	_, err = managedClustersClient.ListClusterUserCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, &armcontainerservice.ManagedClustersClientListClusterUserCredentialsOptions{ServerFqdn: nil,
		Format: nil,
	})
	testsuite.Require().NoError(err)

	// From step ManagedClusters_ListClusterAdminCredentials
	fmt.Println("Call operation: ManagedClusters_ListClusterAdminCredentials")
	_, err = managedClustersClient.ListClusterAdminCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, &armcontainerservice.ManagedClustersClientListClusterAdminCredentialsOptions{ServerFqdn: nil})
	testsuite.Require().NoError(err)

	// From step ManagedClusters_ListClusterMonitoringUserCredentials
	fmt.Println("Call operation: ManagedClusters_ListClusterMonitoringUserCredentials")
	_, err = managedClustersClient.ListClusterMonitoringUserCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, &armcontainerservice.ManagedClustersClientListClusterMonitoringUserCredentialsOptions{ServerFqdn: nil})
	testsuite.Require().NoError(err)

	// From step ManagedClusters_RunCommand
	fmt.Println("Call operation: ManagedClusters_RunCommand")
	managedClustersClientRunCommandResponsePoller, err := managedClustersClient.BeginRunCommand(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armcontainerservice.RunCommandRequest{
		ClusterToken: to.Ptr(""),
		Command:      to.Ptr("kubectl version"),
		Context:      to.Ptr(""),
	}, nil)
	testsuite.Require().NoError(err)
	var managedClustersClientRunCommandResponse *armcontainerservice.ManagedClustersClientRunCommandResponse
	managedClustersClientRunCommandResponse, err = testutil.PollForTest(testsuite.ctx, managedClustersClientRunCommandResponsePoller)
	testsuite.Require().NoError(err)
	commandId = *managedClustersClientRunCommandResponse.ID

	// From step ManagedClusters_GetCommandResult
	fmt.Println("Call operation: ManagedClusters_GetCommandResult")
	_, err = managedClustersClient.GetCommandResult(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, commandId, nil)
	testsuite.Require().NoError(err)

	// From step ManagedClusters_Start
	fmt.Println("Call operation: ManagedClusters_Start")
	managedClustersClientStartResponsePoller, err := managedClustersClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedClustersClientStartResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerService/managedClusters/agentPools
func (testsuite *ManagedclustersTestSuite) TestAgentpools() {
	var err error
	// From step AgentPools_List
	fmt.Println("Call operation: AgentPools_List")
	agentPoolsClient, err := armcontainerservice.NewAgentPoolsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	agentPoolsClientNewListPager := agentPoolsClient.NewListPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for agentPoolsClientNewListPager.More() {
		_, err := agentPoolsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AgentPools_GetUpgradeProfile
	fmt.Println("Call operation: AgentPools_GetUpgradeProfile")
	_, err = agentPoolsClient.GetUpgradeProfile(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_Get
	fmt.Println("Call operation: AgentPools_Get")
	_, err = agentPoolsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_GetAvailableAgentPoolVersions
	fmt.Println("Call operation: AgentPools_GetAvailableAgentPoolVersions")
	_, err = agentPoolsClient.GetAvailableAgentPoolVersions(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_UpgradeNodeImageVersion
	fmt.Println("Call operation: AgentPools_UpgradeNodeImageVersion")
	agentPoolsClientUpgradeNodeImageVersionResponsePoller, err := agentPoolsClient.BeginUpgradeNodeImageVersion(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientUpgradeNodeImageVersionResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerService/snapshots
func (testsuite *ManagedclustersTestSuite) TestSnapshots() {
	var err error
	// From step Snapshots_List
	fmt.Println("Call operation: Snapshots_List")
	snapshotsClient, err := armcontainerservice.NewSnapshotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	snapshotsClientNewListPager := snapshotsClient.NewListPager(nil)
	for snapshotsClientNewListPager.More() {
		_, err := snapshotsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Snapshots_ListByResourceGroup
	fmt.Println("Call operation: Snapshots_ListByResourceGroup")
	snapshotsClientNewListByResourceGroupPager := snapshotsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for snapshotsClientNewListByResourceGroupPager.More() {
		_, err := snapshotsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Snapshots_Get
	fmt.Println("Call operation: Snapshots_Get")
	_, err = snapshotsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerService/managedClusters/maintenanceConfigurations
func (testsuite *ManagedclustersTestSuite) TestMaintenanceconfigurations() {
	var err error
	// From step MaintenanceConfigurations_CreateOrUpdate
	fmt.Println("Call operation: MaintenanceConfigurations_CreateOrUpdate")
	maintenanceConfigurationsClient, err := armcontainerservice.NewMaintenanceConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = maintenanceConfigurationsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.configName, armcontainerservice.MaintenanceConfiguration{
		Properties: &armcontainerservice.MaintenanceConfigurationProperties{
			NotAllowedTime: []*armcontainerservice.TimeSpan{
				{
					End:   to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-30T12:00:00Z"); return t }()),
					Start: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-26T03:00:00Z"); return t }()),
				}},
			TimeInWeek: []*armcontainerservice.TimeInWeek{
				{
					Day: to.Ptr(armcontainerservice.WeekDayMonday),
					HourSlots: []*int32{
						to.Ptr[int32](1),
						to.Ptr[int32](2)},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step MaintenanceConfigurations_ListByManagedCluster
	fmt.Println("Call operation: MaintenanceConfigurations_ListByManagedCluster")
	maintenanceConfigurationsClientNewListByManagedClusterPager := maintenanceConfigurationsClient.NewListByManagedClusterPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for maintenanceConfigurationsClientNewListByManagedClusterPager.More() {
		_, err := maintenanceConfigurationsClientNewListByManagedClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MaintenanceConfigurations_Get
	fmt.Println("Call operation: MaintenanceConfigurations_Get")
	_, err = maintenanceConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.configName, nil)
	testsuite.Require().NoError(err)

	// From step MaintenanceConfigurations_Delete
	fmt.Println("Call operation: MaintenanceConfigurations_Delete")
	_, err = maintenanceConfigurationsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.configName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *ManagedclustersTestSuite) Cleanup() {
	var err error
	// From step Snapshots_Delete
	fmt.Println("Call operation: Snapshots_Delete")
	snapshotsClient, err := armcontainerservice.NewSnapshotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = snapshotsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_Delete
	fmt.Println("Call operation: AgentPools_Delete")
	agentPoolsClient, err := armcontainerservice.NewAgentPoolsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	agentPoolsClientDeleteResponsePoller, err := agentPoolsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step ManagedClusters_Stop
	fmt.Println("Call operation: ManagedClusters_Stop")
	managedClustersClient, err := armcontainerservice.NewManagedClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	managedClustersClientStopResponsePoller, err := managedClustersClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedClustersClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step ManagedClusters_Delete
	fmt.Println("Call operation: ManagedClusters_Delete")
	managedClustersClientDeleteResponsePoller, err := managedClustersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managedClustersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

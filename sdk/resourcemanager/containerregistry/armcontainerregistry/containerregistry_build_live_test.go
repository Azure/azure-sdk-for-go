// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcontainerregistry_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ContainerregistryBuildTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	agentPoolName     string
	registryName      string
	taskName          string
	taskRunName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ContainerregistryBuildTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.agentPoolName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "agentpooln", 16, false)
	testsuite.registryName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "registryna2", 17, false)
	testsuite.taskName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "taskname", 14, false)
	testsuite.taskRunName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "taskrunnam", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ContainerregistryBuildTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestContainerregistryBuildTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerregistryBuildTestSuite))
}

func (testsuite *ContainerregistryBuildTestSuite) Prepare() {
	var err error
	// From step Registries_Create2
	fmt.Println("Call operation: Registries_Create")
	registriesClient, err := armcontainerregistry.NewRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	registriesClientCreateResponsePoller, err := registriesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, armcontainerregistry.Registry{
		Location: to.Ptr(testsuite.location),
		Properties: &armcontainerregistry.RegistryProperties{
			AdminUserEnabled: to.Ptr(true),
		},
		SKU: &armcontainerregistry.SKU{
			Name: to.Ptr(armcontainerregistry.SKUNamePremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, registriesClientCreateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/registries/agentPools
func (testsuite *ContainerregistryBuildTestSuite) TestAgentpools() {
	var err error
	// From step AgentPools_Create
	fmt.Println("Call operation: AgentPools_Create")
	agentPoolsClient, err := armcontainerregistry.NewAgentPoolsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	agentPoolsClientCreateResponsePoller, err := agentPoolsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, armcontainerregistry.AgentPool{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
		Properties: &armcontainerregistry.AgentPoolProperties{
			Count: to.Ptr[int32](1),
			OS:    to.Ptr(armcontainerregistry.OSLinux),
			Tier:  to.Ptr("S1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AgentPools_List
	fmt.Println("Call operation: AgentPools_List")
	agentPoolsClientNewListPager := agentPoolsClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for agentPoolsClientNewListPager.More() {
		_, err := agentPoolsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AgentPools_Get
	fmt.Println("Call operation: AgentPools_Get")
	_, err = agentPoolsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_Update
	fmt.Println("Call operation: AgentPools_Update")
	agentPoolsClientUpdateResponsePoller, err := agentPoolsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, armcontainerregistry.AgentPoolUpdateParameters{
		Properties: &armcontainerregistry.AgentPoolPropertiesUpdateParameters{
			Count: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AgentPools_GetQueueStatus
	fmt.Println("Call operation: AgentPools_GetQueueStatus")
	_, err = agentPoolsClient.GetQueueStatus(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_Delete
	fmt.Println("Call operation: AgentPools_Delete")
	agentPoolsClientDeleteResponsePoller, err := agentPoolsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

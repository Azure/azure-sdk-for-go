//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

type AgentTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	agentName         string
	armEndpoint       string
	storageMoverName  string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AgentTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.agentName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "agentnam", 14, false)
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.storageMoverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storagem", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AgentTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}

func (testsuite *AgentTestSuite) Prepare() {
	var err error
	// From step StorageMovers_CreateOrUpdate
	fmt.Println("Call operation: StorageMovers_CreateOrUpdate")
	storageMoversClient, err := armstoragemover.NewStorageMoversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = storageMoversClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, armstoragemover.StorageMover{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
		Properties: &armstoragemover.Properties{
			Description: to.Ptr("Example Storage Mover Description"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageMover/storageMovers/{storageMoverName}/agents/{agentName}
func (testsuite *AgentTestSuite) TestAgents() {
	var err error
	myUUID, err := uuid.New()
	testsuite.Require().NoError(err)

	// From step Agents_CreateOrUpdate
	fmt.Println("Call operation: Agents_CreateOrUpdate")
	agentsClient, err := armstoragemover.NewAgentsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = agentsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.agentName, armstoragemover.Agent{
		Properties: &armstoragemover.AgentProperties{
			Description:   to.Ptr("Example Agent Description"),
			ArcResourceID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
			ArcVMUUID:     to.Ptr(myUUID.String()),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Agents_List
	fmt.Println("Call operation: Agents_List")
	agentsClientNewListPager := agentsClient.NewListPager(testsuite.resourceGroupName, testsuite.storageMoverName, nil)
	for agentsClientNewListPager.More() {
		_, err := agentsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Agents_Get
	fmt.Println("Call operation: Agents_Get")
	_, err = agentsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.agentName, nil)
	testsuite.Require().NoError(err)

	// From step Agents_Update
	fmt.Println("Call operation: Agents_Update")
	_, err = agentsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.agentName, armstoragemover.AgentUpdateParameters{
		Properties: &armstoragemover.AgentUpdateProperties{
			Description: to.Ptr("Updated Agent Description"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Agents_Delete
	fmt.Println("Call operation: Agents_Delete")
	agentsClientDeleteResponsePoller, err := agentsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.agentName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *AgentTestSuite) Cleanup() {
	var err error
	// From step StorageMovers_Delete
	fmt.Println("Call operation: StorageMovers_Delete")
	storageMoversClient, err := armstoragemover.NewStorageMoversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	storageMoversClientDeleteResponsePoller, err := storageMoversClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageMoversClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

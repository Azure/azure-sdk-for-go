// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ChatModelDeploymentsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	workspaceName     string
	deploymentName    string
}

func (testsuite *ChatModelDeploymentsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.workspaceName = "test-wrksp-go01"
	testsuite.deploymentName = "test-cmd-go01"
}

func (testsuite *ChatModelDeploymentsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestChatModelDeploymentsTestSuite(t *testing.T) {
	suite.Run(t, new(ChatModelDeploymentsTestSuite))
}

// Test listing chat model deployments by workspace
func (testsuite *ChatModelDeploymentsTestSuite) TestChatModelDeploymentsListByWorkspace() {
	fmt.Println("Call operation: ChatModelDeployments_ListByWorkspace")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewChatModelDeploymentsClient().NewListByWorkspacePager(testsuite.resourceGroupName, testsuite.workspaceName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a chat model deployment
func (testsuite *ChatModelDeploymentsTestSuite) TestChatModelDeploymentsGet() {
	fmt.Println("Call operation: ChatModelDeployments_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewChatModelDeploymentsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, testsuite.deploymentName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a chat model deployment (skipped - parent workspace CreateOrUpdate not yet recorded)
func (testsuite *ChatModelDeploymentsTestSuite) SkipTestChatModelDeploymentsCreateOrUpdate() {
	fmt.Println("Call operation: ChatModelDeployments_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	chatModelClient := clientFactory.NewChatModelDeploymentsClient()
	poller, err := chatModelClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.workspaceName,
		testsuite.deploymentName,
		armdiscovery.ChatModelDeployment{
			Location: to.Ptr(testsuite.location),
			Properties: &armdiscovery.ChatModelDeploymentProperties{
				ModelFormat: to.Ptr("OpenAI"),
				ModelName:   to.Ptr("gpt-4o"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	cmd, err := poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(cmd.ID)
	fmt.Println("Created chat model deployment:", *cmd.Name)
}

// Test deleting a chat model deployment
func (testsuite *ChatModelDeploymentsTestSuite) SkipTestChatModelDeploymentsDelete() {
	fmt.Println("Call operation: ChatModelDeployments_Delete")
	// Requires existing chat model deployment
}

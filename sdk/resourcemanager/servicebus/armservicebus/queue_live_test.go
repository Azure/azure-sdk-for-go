// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armservicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	namespaceName         string
	queueName             string
	authorizationRuleName string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *QueueTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.queueName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "queuenam", 14, false)
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "queueauthoriz", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *QueueTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

func (testsuite *QueueTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.SBNamespace{
		Location: to.Ptr(testsuite.location),
		SKU: &armservicebus.SBSKU{
			Name: to.Ptr(armservicebus.SKUNamePremium),
			Tier: to.Ptr(armservicebus.SKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Queues_CreateOrUpdate
	fmt.Println("Call operation: Queues_CreateOrUpdate")
	queuesClient, err := armservicebus.NewQueuesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = queuesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, armservicebus.SBQueue{
		Properties: &armservicebus.SBQueueProperties{
			EnablePartitioning: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/queues/{queueName}
func (testsuite *QueueTestSuite) TestQueues() {
	var err error
	// From step Queues_ListByNamespace
	fmt.Println("Call operation: Queues_ListByNamespace")
	queuesClient, err := armservicebus.NewQueuesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	queuesClientNewListByNamespacePager := queuesClient.NewListByNamespacePager(testsuite.resourceGroupName, testsuite.namespaceName, &armservicebus.QueuesClientListByNamespaceOptions{Skip: nil,
		Top: nil,
	})
	for queuesClientNewListByNamespacePager.More() {
		_, err := queuesClientNewListByNamespacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Queues_Get
	fmt.Println("Call operation: Queues_Get")
	_, err = queuesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/queues/{queueName}/authorizationRules/{authorizationRuleName}
func (testsuite *QueueTestSuite) TestQueuesAuthorization() {
	var err error
	// From step Queues_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Queues_CreateOrUpdateAuthorizationRule")
	queuesClient, err := armservicebus.NewQueuesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = queuesClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, testsuite.authorizationRuleName, armservicebus.SBAuthorizationRule{
		Properties: &armservicebus.SBAuthorizationRuleProperties{
			Rights: []*armservicebus.AccessRights{
				to.Ptr(armservicebus.AccessRightsListen),
				to.Ptr(armservicebus.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Queues_GetAuthorizationRule
	fmt.Println("Call operation: Queues_GetAuthorizationRule")
	_, err = queuesClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Queues_RegenerateKeys
	fmt.Println("Call operation: Queues_RegenerateKeys")
	_, err = queuesClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, testsuite.authorizationRuleName, armservicebus.RegenerateAccessKeyParameters{
		KeyType: to.Ptr(armservicebus.KeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Queues_ListKeys
	fmt.Println("Call operation: Queues_ListKeys")
	_, err = queuesClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Queues_ListAuthorizationRules
	fmt.Println("Call operation: Queues_ListAuthorizationRules")
	queuesClientNewListAuthorizationRulesPager := queuesClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, nil)
	for queuesClientNewListAuthorizationRulesPager.More() {
		_, err := queuesClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Queues_DeleteAuthorizationRule
	fmt.Println("Call operation: Queues_DeleteAuthorizationRule")
	_, err = queuesClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *QueueTestSuite) Cleanup() {
	var err error
	// From step Queues_Delete
	fmt.Println("Call operation: Queues_Delete")
	queuesClient, err := armservicebus.NewQueuesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = queuesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.queueName, nil)
	testsuite.Require().NoError(err)
}

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

type SubscriptionsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	namespaceName     string
	topicName         string
	subscriptionName  string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SubscriptionsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.topicName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "topicnam", 14, false)
	testsuite.subscriptionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "subscrip", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SubscriptionsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSubscriptionsTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionsTestSuite))
}

func (testsuite *SubscriptionsTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.SBNamespace{
		Location: to.Ptr(testsuite.location),
		SKU: &armservicebus.SBSKU{
			Name: to.Ptr(armservicebus.SKUNameStandard),
			Tier: to.Ptr(armservicebus.SKUTierStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Topics_CreateOrUpdate
	fmt.Println("Call operation: Topics_CreateOrUpdate")
	topicsClient, err := armservicebus.NewTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = topicsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, armservicebus.SBTopic{
		Properties: &armservicebus.SBTopicProperties{
			EnableExpress: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}/subscriptions/{subscriptionName}
func (testsuite *SubscriptionsTestSuite) TestSubscriptions() {
	var err error
	// From step Subscriptions_CreateOrUpdate
	fmt.Println("Call operation: Subscriptions_CreateOrUpdate")
	subscriptionsClient, err := armservicebus.NewSubscriptionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = subscriptionsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.subscriptionName, armservicebus.SBSubscription{
		Properties: &armservicebus.SBSubscriptionProperties{
			EnableBatchedOperations: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Subscriptions_ListByTopic
	fmt.Println("Call operation: Subscriptions_ListByTopic")
	subscriptionsClientNewListByTopicPager := subscriptionsClient.NewListByTopicPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, &armservicebus.SubscriptionsClientListByTopicOptions{Skip: nil,
		Top: nil,
	})
	for subscriptionsClientNewListByTopicPager.More() {
		_, err := subscriptionsClientNewListByTopicPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Subscriptions_Get
	fmt.Println("Call operation: Subscriptions_Get")
	_, err = subscriptionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.subscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step Subscriptions_Delete
	fmt.Println("Call operation: Subscriptions_Delete")
	_, err = subscriptionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.subscriptionName, nil)
	testsuite.Require().NoError(err)
}

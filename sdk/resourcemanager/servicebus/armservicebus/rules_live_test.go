//go:build go1.18
// +build go1.18

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
	"github.com/stretchr/testify/suite"
)

type RulesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	namespaceName     string
	subscriptionName  string
	topicName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *RulesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/servicebus/armservicebus/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namespaceName = testutil.GenerateAlphaNumericID(testsuite.T(), "namespac", 6)
	testsuite.subscriptionName = testutil.GenerateAlphaNumericID(testsuite.T(), "subscrip", 6)
	testsuite.topicName = testutil.GenerateAlphaNumericID(testsuite.T(), "topicnam", 6)
	testsuite.location = testutil.GetEnv("LOCATION", "westus")
	testsuite.resourceGroupName = testutil.GetEnv("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *RulesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestRulesTestSuite(t *testing.T) {
	suite.Run(t, new(RulesTestSuite))
}

func (testsuite *RulesTestSuite) Prepare() {
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
}

// providers/Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}/subscriptions/{subscriptionName}/rules/{ruleName}
func (testsuite *RulesTestSuite) TestRules() {
	ruleName := testutil.GenerateAlphaNumericID(testsuite.T(), "rulename", 6)
	var err error
	// From step Rules_CreateOrUpdate
	fmt.Println("Call operation: Rules_CreateOrUpdate")
	rulesClient, err := armservicebus.NewRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = rulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.subscriptionName, ruleName, armservicebus.Rule{}, nil)
	testsuite.Require().NoError(err)

	// From step Rules_ListBySubscriptions
	fmt.Println("Call operation: Rules_ListBySubscriptions")
	rulesClientNewListBySubscriptionsPager := rulesClient.NewListBySubscriptionsPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.subscriptionName, &armservicebus.RulesClientListBySubscriptionsOptions{Skip: nil,
		Top: nil,
	})
	for rulesClientNewListBySubscriptionsPager.More() {
		_, err := rulesClientNewListBySubscriptionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Rules_Get
	fmt.Println("Call operation: Rules_Get")
	_, err = rulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.subscriptionName, ruleName, nil)
	testsuite.Require().NoError(err)

	// From step Rules_Delete
	fmt.Println("Call operation: Rules_Delete")
	_, err = rulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.subscriptionName, ruleName, nil)
	testsuite.Require().NoError(err)
}

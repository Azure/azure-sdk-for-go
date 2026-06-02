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

type TopicsTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	namespaceName         string
	topicName             string
	authorizationRuleName string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *TopicsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.topicName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "topicnam", 14, false)
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "topicauthoriz", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *TopicsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestTopicsTestSuite(t *testing.T) {
	suite.Run(t, new(TopicsTestSuite))
}

func (testsuite *TopicsTestSuite) Prepare() {
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

// Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}
func (testsuite *TopicsTestSuite) TestTopics() {
	var err error
	// From step Topics_ListByNamespace
	fmt.Println("Call operation: Topics_ListByNamespace")
	topicsClient, err := armservicebus.NewTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	topicsClientNewListByNamespacePager := topicsClient.NewListByNamespacePager(testsuite.resourceGroupName, testsuite.namespaceName, &armservicebus.TopicsClientListByNamespaceOptions{Skip: nil,
		Top: nil,
	})
	for topicsClientNewListByNamespacePager.More() {
		_, err := topicsClientNewListByNamespacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Topics_Get
	fmt.Println("Call operation: Topics_Get")
	_, err = topicsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}/authorizationRules/{authorizationRuleName}
func (testsuite *TopicsTestSuite) TestTopicsAuthorization() {
	var err error
	// From step Topics_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Topics_CreateOrUpdateAuthorizationRule")
	topicsClient, err := armservicebus.NewTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = topicsClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.authorizationRuleName, armservicebus.SBAuthorizationRule{
		Properties: &armservicebus.SBAuthorizationRuleProperties{
			Rights: []*armservicebus.AccessRights{
				to.Ptr(armservicebus.AccessRightsListen),
				to.Ptr(armservicebus.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Topics_ListAuthorizationRules
	fmt.Println("Call operation: Topics_ListAuthorizationRules")
	topicsClientNewListAuthorizationRulesPager := topicsClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, nil)
	for topicsClientNewListAuthorizationRulesPager.More() {
		_, err := topicsClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Topics_GetAuthorizationRule
	fmt.Println("Call operation: Topics_GetAuthorizationRule")
	_, err = topicsClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Topics_ListKeys
	fmt.Println("Call operation: Topics_ListKeys")
	_, err = topicsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Topics_RegenerateKeys
	fmt.Println("Call operation: Topics_RegenerateKeys")
	_, err = topicsClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.authorizationRuleName, armservicebus.RegenerateAccessKeyParameters{
		KeyType: to.Ptr(armservicebus.KeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Topics_DeleteAuthorizationRule
	fmt.Println("Call operation: Topics_DeleteAuthorizationRule")
	_, err = topicsClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *TopicsTestSuite) Cleanup() {
	var err error
	// From step Topics_Delete
	fmt.Println("Call operation: Topics_Delete")
	topicsClient, err := armservicebus.NewTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = topicsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.topicName, nil)
	testsuite.Require().NoError(err)
}

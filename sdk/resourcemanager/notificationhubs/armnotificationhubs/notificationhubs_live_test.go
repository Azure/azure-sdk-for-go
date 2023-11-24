//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnotificationhubs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/notificationhubs/armnotificationhubs"
	"github.com/stretchr/testify/suite"
)

type NotificationhubsTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	authorizationRuleName string
	namespaceName         string
	notificationHubName   string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *NotificationhubsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/notificationhubs/armnotificationhubs/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "authoriz", 14, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.notificationHubName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "notifica", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NotificationhubsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestNotificationhubsTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationhubsTestSuite))
}

func (testsuite *NotificationhubsTestSuite) Prepare() {
	var err error
	// From step Namespaces_CheckAvailability
	fmt.Println("Call operation: Namespaces_CheckAvailability")
	namespacesClient, err := armnotificationhubs.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = namespacesClient.CheckAvailability(testsuite.ctx, armnotificationhubs.CheckAvailabilityParameters{
		Name: to.Ptr(testsuite.namespaceName),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	_, err = namespacesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armnotificationhubs.NamespaceCreateOrUpdateParameters{
		Location: to.Ptr(testsuite.location),
		SKU: &armnotificationhubs.SKU{
			Name: to.Ptr(armnotificationhubs.SKUNameStandard),
			Tier: to.Ptr("Standard"),
		},
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_CheckNotificationHubAvailability
	fmt.Println("Call operation: NotificationHubs_CheckNotificationHubAvailability")
	client, err := armnotificationhubs.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.CheckNotificationHubAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armnotificationhubs.CheckAvailabilityParameters{
		Name:     to.Ptr(testsuite.notificationHubName),
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_CreateOrUpdate
	fmt.Println("Call operation: NotificationHubs_CreateOrUpdate")
	_, err = client.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, armnotificationhubs.NotificationHubCreateOrUpdateParameters{
		Location:   to.Ptr(testsuite.location),
		Properties: &armnotificationhubs.NotificationHubProperties{},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/operations
func (testsuite *NotificationhubsTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armnotificationhubs.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.NotificationHubs/namespaces/{namespaceName}
func (testsuite *NotificationhubsTestSuite) TestNamespaces() {
	var err error
	// From step Namespaces_ListAll
	fmt.Println("Call operation: Namespaces_ListAll")
	namespacesClient, err := armnotificationhubs.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientNewListAllPager := namespacesClient.NewListAllPager(nil)
	for namespacesClientNewListAllPager.More() {
		_, err := namespacesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_Get
	fmt.Println("Call operation: Namespaces_Get")
	_, err = namespacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_List
	fmt.Println("Call operation: Namespaces_List")
	namespacesClientNewListPager := namespacesClient.NewListPager(testsuite.resourceGroupName, nil)
	for namespacesClientNewListPager.More() {
		_, err := namespacesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_Patch
	fmt.Println("Call operation: Namespaces_Patch")
	_, err = namespacesClient.Patch(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armnotificationhubs.NamespacePatchParameters{
		SKU: &armnotificationhubs.SKU{
			Name: to.Ptr(armnotificationhubs.SKUNameStandard),
			Tier: to.Ptr("Standard"),
		},
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/namespaces/{namespaceName}/AuthorizationRules/{authorizationRuleName}
func (testsuite *NotificationhubsTestSuite) TestNamespacesAuthorizationRule() {
	var err error

	// From step Namespaces_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Namespaces_CreateOrUpdateAuthorizationRule")
	namespacesClient, err := armnotificationhubs.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = namespacesClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, "namespace"+testsuite.authorizationRuleName, armnotificationhubs.SharedAccessAuthorizationRuleCreateOrUpdateParameters{
		Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
			Rights: []*armnotificationhubs.AccessRights{
				to.Ptr(armnotificationhubs.AccessRightsListen),
				to.Ptr(armnotificationhubs.AccessRightsSend),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListAuthorizationRules
	fmt.Println("Call operation: Namespaces_ListAuthorizationRules")
	namespacesClientNewListAuthorizationRulesPager := namespacesClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for namespacesClientNewListAuthorizationRulesPager.More() {
		_, err = namespacesClientNewListAuthorizationRulesPager.NextPage(context.Background())
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_GetAuthorizationRule
	fmt.Println("Call operation: Namespaces_GetAuthorizationRule")
	_, err = namespacesClient.GetAuthorizationRule(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, "namespace"+testsuite.authorizationRuleName, &armnotificationhubs.NamespacesClientGetAuthorizationRuleOptions{})
	testsuite.Require().NoError(err)

	// From step Namespaces_ListKeys
	fmt.Println("Call operation: Namespaces_ListKeys")
	_, err = namespacesClient.ListKeys(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, "namespace"+testsuite.authorizationRuleName, &armnotificationhubs.NamespacesClientListKeysOptions{})
	testsuite.Require().NoError(err)

	// From step Namespaces_RegenerateKeys
	fmt.Println("Call operation: Namespaces_RegenerateKeys")
	_, err = namespacesClient.RegenerateKeys(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, "namespace"+testsuite.authorizationRuleName, armnotificationhubs.PolicykeyResource{
		PolicyKey: to.Ptr("PrimaryKey"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_DeleteAuthorizationRule
	fmt.Println("Call operation: Namespaces_DeleteAuthorizationRule")
	_, err = namespacesClient.DeleteAuthorizationRule(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, "namespace"+testsuite.authorizationRuleName, &armnotificationhubs.NamespacesClientDeleteAuthorizationRuleOptions{})
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/namespaces/{namespaceName}/notificationHubs/{notificationHubName}
func (testsuite *NotificationhubsTestSuite) TestNotificationHubs() {
	var err error
	// From step NotificationHubs_List
	fmt.Println("Call operation: NotificationHubs_List")
	client, err := armnotificationhubs.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientNewListPager := client.NewListPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for clientNewListPager.More() {
		_, err := clientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NotificationHubs_Get
	fmt.Println("Call operation: NotificationHubs_Get")
	_, err = client.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_Patch
	fmt.Println("Call operation: NotificationHubs_Patch")
	_, err = client.Patch(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, &armnotificationhubs.ClientPatchOptions{
		Parameters: &armnotificationhubs.NotificationHubPatchParameters{
			Tags: map[string]*string{
				"key": to.Ptr("value"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step NotificationHubs_GetPnsCredentials
	fmt.Println("Call operation: NotificationHubs_GetPnsCredentials")
	_, err = client.GetPnsCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.NotificationHubs/namespaces/{namespaceName}/notificationHubs/{notificationHubName}/AuthorizationRules/{authorizationRuleName}
func (testsuite *NotificationhubsTestSuite) TestNotificationHubsAuthorizationRule() {
	var err error

	// From step NotificationHubs_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: NotificationHubs_CreateOrUpdateAuthorizationRule")
	client, err := armnotificationhubs.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.CreateOrUpdateAuthorizationRule(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, "notification"+testsuite.authorizationRuleName, armnotificationhubs.SharedAccessAuthorizationRuleCreateOrUpdateParameters{
		Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
			Rights: []*armnotificationhubs.AccessRights{
				to.Ptr(armnotificationhubs.AccessRightsListen),
				to.Ptr(armnotificationhubs.AccessRightsSend),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_ListAuthorizationRules
	fmt.Println("Call operation: NotificationHubs_ListAuthorizationRules")
	clientListAuthorizationRulesResponse := client.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, &armnotificationhubs.ClientListAuthorizationRulesOptions{})
	for clientListAuthorizationRulesResponse.More() {
		_, err = clientListAuthorizationRulesResponse.NextPage(context.Background())
		testsuite.Require().NoError(err)
		break
	}

	// From step NotificationHubs_GetAuthorizationRule
	fmt.Println("Call operation: NotificationHubs_GetAuthorizationRule")
	_, err = client.GetAuthorizationRule(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, "notification"+testsuite.authorizationRuleName, &armnotificationhubs.ClientGetAuthorizationRuleOptions{})
	testsuite.Require().NoError(err)

	// From step NotificationHubs_ListKeys
	fmt.Println("Call operation: NotificationHubs_ListKeys")
	_, err = client.ListKeys(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, "notification"+testsuite.authorizationRuleName, &armnotificationhubs.ClientListKeysOptions{})
	testsuite.Require().NoError(err)

	// From step NotificationHubs_RegenerateKeys
	fmt.Println("Call operation: NotificationHubs_RegenerateKeys")
	_, err = client.RegenerateKeys(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, "notification"+testsuite.authorizationRuleName, armnotificationhubs.PolicykeyResource{
		PolicyKey: to.Ptr("PrimaryKey"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step NotificationHubs_DeleteAuthorizationRule
	fmt.Println("Call operation: NotificationHubs_DeleteAuthorizationRule")
	_, err = client.DeleteAuthorizationRule(context.Background(), testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, "notification"+testsuite.authorizationRuleName, &armnotificationhubs.ClientDeleteAuthorizationRuleOptions{})
	testsuite.Require().NoError(err)
}

func (testsuite *NotificationhubsTestSuite) Cleanup() {
	var err error
	// From step NotificationHubs_Delete
	fmt.Println("Call operation: NotificationHubs_Delete")
	client, err := armnotificationhubs.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = client.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.notificationHubName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_Delete
	fmt.Println("Call operation: Namespaces_Delete")
	namespacesClient, err := armnotificationhubs.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientDeleteResponsePoller, err := namespacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

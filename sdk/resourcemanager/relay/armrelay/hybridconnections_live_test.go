// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armrelay_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/relay/armrelay"
	"github.com/stretchr/testify/suite"
)

type HybridConnectionsTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	authorizationRuleName string
	hybridConnectionName  string
	namespaceName         string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *HybridConnectionsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "authoriz", 14, false)
	testsuite.hybridConnectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "hybridco", 14, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *HybridConnectionsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestHybridConnectionsTestSuite(t *testing.T) {
	suite.Run(t, new(HybridConnectionsTestSuite))
}

func (testsuite *HybridConnectionsTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armrelay.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armrelay.Namespace{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		SKU: &armrelay.SKU{
			Name: to.Ptr(armrelay.SKUNameStandard),
			Tier: to.Ptr(armrelay.SKUTierStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step HybridConnections_CreateOrUpdate
	fmt.Println("Call operation: HybridConnections_CreateOrUpdate")
	hybridConnectionsClient, err := armrelay.NewHybridConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = hybridConnectionsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, armrelay.HybridConnection{
		Properties: &armrelay.HybridConnectionProperties{
			RequiresClientAuthorization: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Relay/namespaces/{namespaceName}/hybridConnections/{hybridConnectionName}
func (testsuite *HybridConnectionsTestSuite) TestHybridConnections() {
	var err error
	// From step HybridConnections_ListByNamespace
	fmt.Println("Call operation: HybridConnections_ListByNamespace")
	hybridConnectionsClient, err := armrelay.NewHybridConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	hybridConnectionsClientNewListByNamespacePager := hybridConnectionsClient.NewListByNamespacePager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for hybridConnectionsClientNewListByNamespacePager.More() {
		_, err := hybridConnectionsClientNewListByNamespacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step HybridConnections_Get
	fmt.Println("Call operation: HybridConnections_Get")
	_, err = hybridConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Relay/namespaces/{namespaceName}/hybridConnections/{hybridConnectionName}/authorizationRules/{authorizationRuleName}
func (testsuite *HybridConnectionsTestSuite) TestHybridConnectionsAuthorization() {
	var err error
	// From step HybridConnections_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: HybridConnections_CreateOrUpdateAuthorizationRule")
	hybridConnectionsClient, err := armrelay.NewHybridConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = hybridConnectionsClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, testsuite.authorizationRuleName, armrelay.AuthorizationRule{
		Properties: &armrelay.AuthorizationRuleProperties{
			Rights: []*armrelay.AccessRights{
				to.Ptr(armrelay.AccessRightsListen),
				to.Ptr(armrelay.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step HybridConnections_GetAuthorizationRule
	fmt.Println("Call operation: HybridConnections_GetAuthorizationRule")
	_, err = hybridConnectionsClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step HybridConnections_ListAuthorizationRules
	fmt.Println("Call operation: HybridConnections_ListAuthorizationRules")
	hybridConnectionsClientNewListAuthorizationRulesPager := hybridConnectionsClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, nil)
	for hybridConnectionsClientNewListAuthorizationRulesPager.More() {
		_, err := hybridConnectionsClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step HybridConnections_ListKeys
	fmt.Println("Call operation: HybridConnections_ListKeys")
	_, err = hybridConnectionsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step HybridConnections_RegenerateKeys
	fmt.Println("Call operation: HybridConnections_RegenerateKeys")
	_, err = hybridConnectionsClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, testsuite.authorizationRuleName, armrelay.RegenerateAccessKeyParameters{
		KeyType: to.Ptr(armrelay.KeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step HybridConnections_DeleteAuthorizationRule
	fmt.Println("Call operation: HybridConnections_DeleteAuthorizationRule")
	_, err = hybridConnectionsClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *HybridConnectionsTestSuite) Cleanup() {
	var err error
	// From step HybridConnections_Delete
	fmt.Println("Call operation: HybridConnections_Delete")
	hybridConnectionsClient, err := armrelay.NewHybridConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = hybridConnectionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.hybridConnectionName, nil)
	testsuite.Require().NoError(err)
}

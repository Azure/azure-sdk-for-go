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

type WcfRelaysTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	authorizationRuleName string
	namespaceName         string
	relayName             string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *WcfRelaysTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "authoriz", 14, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.relayName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "relaynam", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *WcfRelaysTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestWcfRelaysTestSuite(t *testing.T) {
	suite.Run(t, new(WcfRelaysTestSuite))
}

func (testsuite *WcfRelaysTestSuite) Prepare() {
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

	// From step WCFRelays_CreateOrUpdate
	fmt.Println("Call operation: WCFRelays_CreateOrUpdate")
	wCFRelaysClient, err := armrelay.NewWCFRelaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = wCFRelaysClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, armrelay.WcfRelay{
		Properties: &armrelay.WcfRelayProperties{
			RelayType:                   to.Ptr(armrelay.RelaytypeNetTCP),
			RequiresClientAuthorization: to.Ptr(true),
			RequiresTransportSecurity:   to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Relay/namespaces/{namespaceName}/wcfRelays/{relayName}
func (testsuite *WcfRelaysTestSuite) TestWcfRelays() {
	var err error
	// From step WCFRelays_ListByNamespace
	fmt.Println("Call operation: WCFRelays_ListByNamespace")
	wCFRelaysClient, err := armrelay.NewWCFRelaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	wCFRelaysClientNewListByNamespacePager := wCFRelaysClient.NewListByNamespacePager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for wCFRelaysClientNewListByNamespacePager.More() {
		_, err := wCFRelaysClientNewListByNamespacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WCFRelays_Get
	fmt.Println("Call operation: WCFRelays_Get")
	_, err = wCFRelaysClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Relay/namespaces/{namespaceName}/wcfRelays/{relayName}/authorizationRules/{authorizationRuleName}
func (testsuite *WcfRelaysTestSuite) TestWcfRelaysAuthorization() {
	var err error
	// From step WCFRelays_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: WCFRelays_CreateOrUpdateAuthorizationRule")
	wCFRelaysClient, err := armrelay.NewWCFRelaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = wCFRelaysClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, testsuite.authorizationRuleName, armrelay.AuthorizationRule{
		Properties: &armrelay.AuthorizationRuleProperties{
			Rights: []*armrelay.AccessRights{
				to.Ptr(armrelay.AccessRightsListen),
				to.Ptr(armrelay.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step WCFRelays_GetAuthorizationRule
	fmt.Println("Call operation: WCFRelays_GetAuthorizationRule")
	_, err = wCFRelaysClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step WCFRelays_ListAuthorizationRules
	fmt.Println("Call operation: WCFRelays_ListAuthorizationRules")
	wCFRelaysClientNewListAuthorizationRulesPager := wCFRelaysClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, nil)
	for wCFRelaysClientNewListAuthorizationRulesPager.More() {
		_, err := wCFRelaysClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WCFRelays_RegenerateKeys
	fmt.Println("Call operation: WCFRelays_RegenerateKeys")
	_, err = wCFRelaysClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, testsuite.authorizationRuleName, armrelay.RegenerateAccessKeyParameters{
		KeyType: to.Ptr(armrelay.KeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step WCFRelays_ListKeys
	fmt.Println("Call operation: WCFRelays_ListKeys")
	_, err = wCFRelaysClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step WCFRelays_DeleteAuthorizationRule
	fmt.Println("Call operation: WCFRelays_DeleteAuthorizationRule")
	_, err = wCFRelaysClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *WcfRelaysTestSuite) Cleanup() {
	var err error
	// From step WCFRelays_Delete
	fmt.Println("Call operation: WCFRelays_Delete")
	wCFRelaysClient, err := armrelay.NewWCFRelaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = wCFRelaysClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.relayName, nil)
	testsuite.Require().NoError(err)
}

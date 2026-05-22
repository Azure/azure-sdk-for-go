// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armeventhub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type DisasterrecoveryconfigsTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	alias                 string
	authorizationRuleName string
	namespaceName         string
	namespaceNameSecond   string
	secondNamespaceId     string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *DisasterrecoveryconfigsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.alias, _ = recording.GenerateAlphaNumericID(testsuite.T(), "alias", 11, false)
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "authorizat", 16, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespacen", 16, false)
	testsuite.namespaceNameSecond, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespacensecond", 22, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *DisasterrecoveryconfigsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestDisasterrecoveryconfigsTestSuite(t *testing.T) {
	suite.Run(t, new(DisasterrecoveryconfigsTestSuite))
}

func (testsuite *DisasterrecoveryconfigsTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armeventhub.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armeventhub.EHNamespace{
		Location: to.Ptr(testsuite.location),
		SKU: &armeventhub.SKU{
			Name: to.Ptr(armeventhub.SKUNameStandard),
			Tier: to.Ptr(armeventhub.SKUTierStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Namespaces_CreateOrUpdate_Second
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClientCreateOrUpdateResponsePoller, err = namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceNameSecond, armeventhub.EHNamespace{
		Location: to.Ptr("westus2"),
		SKU: &armeventhub.SKU{
			Name: to.Ptr(armeventhub.SKUNameStandard),
			Tier: to.Ptr(armeventhub.SKUTierStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var namespacesClientCreateOrUpdateResponse *armeventhub.NamespacesClientCreateOrUpdateResponse
	namespacesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.secondNamespaceId = *namespacesClientCreateOrUpdateResponse.ID

	// From step Namespaces_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Namespaces_CreateOrUpdateAuthorizationRule")
	_, err = namespacesClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armeventhub.AuthorizationRule{
		Properties: &armeventhub.AuthorizationRuleProperties{
			Rights: []*armeventhub.AccessRights{
				to.Ptr(armeventhub.AccessRightsListen),
				to.Ptr(armeventhub.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventHub/namespaces/disasterRecoveryConfigs
func (testsuite *DisasterrecoveryconfigsTestSuite) TTestDisasterrecoveryconfig() {
	var err error
	// From step DisasterRecoveryConfigs_CheckNameAvailability
	fmt.Println("Call operation: DisasterRecoveryConfigs_CheckNameAvailability")
	disasterRecoveryConfigsClient, err := armeventhub.NewDisasterRecoveryConfigsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = disasterRecoveryConfigsClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armeventhub.CheckNameAvailabilityParameter{
		Name: to.Ptr("sdk-DisasterRecovery-9474"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step DisasterRecoveryConfigs_List
	fmt.Println("Call operation: DisasterRecoveryConfigs_List")
	disasterRecoveryConfigsClientNewListPager := disasterRecoveryConfigsClient.NewListPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for disasterRecoveryConfigsClientNewListPager.More() {
		_, err := disasterRecoveryConfigsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DisasterRecoveryConfigs_GetAuthorizationRule
	fmt.Println("Call operation: DisasterRecoveryConfigs_GetAuthorizationRule")
	_, err = disasterRecoveryConfigsClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.alias, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step DisasterRecoveryConfigs_Get
	fmt.Println("Call operation: DisasterRecoveryConfigs_Get")
	_, err = disasterRecoveryConfigsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.alias, nil)
	testsuite.Require().NoError(err)

	// From step DisasterRecoveryConfigs_ListAuthorizationRules
	fmt.Println("Call operation: DisasterRecoveryConfigs_ListAuthorizationRules")
	disasterRecoveryConfigsClientNewListAuthorizationRulesPager := disasterRecoveryConfigsClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.alias, nil)
	for disasterRecoveryConfigsClientNewListAuthorizationRulesPager.More() {
		_, err := disasterRecoveryConfigsClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DisasterRecoveryConfigs_ListKeys
	fmt.Println("Call operation: DisasterRecoveryConfigs_ListKeys")
	_, err = disasterRecoveryConfigsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.alias, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step DisasterRecoveryConfigs_FailOver
	fmt.Println("Call operation: DisasterRecoveryConfigs_FailOver")
	_, err = disasterRecoveryConfigsClient.FailOver(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceNameSecond, testsuite.alias, nil)
	testsuite.Require().NoError(err)
}

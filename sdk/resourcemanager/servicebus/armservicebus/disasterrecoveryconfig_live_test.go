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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
	"github.com/stretchr/testify/suite"
)

type DisasterRecoveryConfigTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	authorizationRuleName string
	namespaceName         string
	primaryNamespaceId    string
	primaryNamespaceName  string
	alias                 string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *DisasterRecoveryConfigTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespaceauthoriz", 23, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.primaryNamespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "promarynamespac", 21, false)
	testsuite.alias, _ = recording.GenerateAlphaNumericID(testsuite.T(), "drcalias", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DisasterRecoveryConfigTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDisasterRecoveryConfigTestSuite(t *testing.T) {
	suite.Run(t, new(DisasterRecoveryConfigTestSuite))
}

func (testsuite *DisasterRecoveryConfigTestSuite) Prepare() {
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

	// From step Namespaces_CreateOrUpdate_Primary
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClientCreateOrUpdateResponsePoller, err = namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.primaryNamespaceName, armservicebus.SBNamespace{
		Location: to.Ptr("westus2"),
		SKU: &armservicebus.SBSKU{
			Name: to.Ptr(armservicebus.SKUNamePremium),
			Tier: to.Ptr(armservicebus.SKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var namespacesClientCreateOrUpdateResponse *armservicebus.NamespacesClientCreateOrUpdateResponse
	namespacesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.primaryNamespaceId = *namespacesClientCreateOrUpdateResponse.ID

	// From step Namespaces_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Namespaces_CreateOrUpdateAuthorizationRule")
	_, err = namespacesClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armservicebus.SBAuthorizationRule{
		Properties: &armservicebus.SBAuthorizationRuleProperties{
			Rights: []*armservicebus.AccessRights{
				to.Ptr(armservicebus.AccessRightsListen),
				to.Ptr(armservicebus.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/disasterRecoveryConfigs/{alias}
func (testsuite *DisasterRecoveryConfigTestSuite) TestDisasterRecoveryConfigs() {
	var err error
	// From step DisasterRecoveryConfigs_CheckNameAvailability
	fmt.Println("Call operation: DisasterRecoveryConfigs_CheckNameAvailability")
	disasterRecoveryConfigsClient, err := armservicebus.NewDisasterRecoveryConfigsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = disasterRecoveryConfigsClient.CheckNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.CheckNameAvailability{
		Name: to.Ptr("sdk-DisasterRecovery-9474"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step DisasterRecoveryConfigs_CreateOrUpdate
	fmt.Println("Call operation: DisasterRecoveryConfigs_CreateOrUpdate")
	_, err = disasterRecoveryConfigsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.alias, armservicebus.ArmDisasterRecovery{
		Properties: &armservicebus.ArmDisasterRecoveryProperties{
			PartnerNamespace: to.Ptr(testsuite.primaryNamespaceId),
		},
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

	// From step DisasterRecoveryConfigs_GetAuthorizationRule
	fmt.Println("Call operation: DisasterRecoveryConfigs_GetAuthorizationRule")
	_, err = disasterRecoveryConfigsClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.alias, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step DisasterRecoveryConfigs_ListKeys
	fmt.Println("Call operation: DisasterRecoveryConfigs_ListKeys")
	_, err = disasterRecoveryConfigsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.alias, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

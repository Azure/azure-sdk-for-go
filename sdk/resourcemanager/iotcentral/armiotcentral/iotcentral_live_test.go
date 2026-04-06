// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armiotcentral_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotcentral/armiotcentral"
	"github.com/stretchr/testify/suite"
)

type IotcentralTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	resourceName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *IotcentralTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.resourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "resource", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *IotcentralTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestIotcentralTestSuite(t *testing.T) {
	suite.Run(t, new(IotcentralTestSuite))
}

// Microsoft.IoTCentral/iotApps/{resourceName}
func (testsuite *IotcentralTestSuite) TestApps() {
	var err error
	// From step Apps_CheckNameAvailability
	fmt.Println("Call operation: Apps_CheckNameAvailability")
	appsClient, err := armiotcentral.NewAppsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = appsClient.CheckNameAvailability(testsuite.ctx, armiotcentral.OperationInputs{
		Name: to.Ptr("myiotcentralapp"),
		Type: to.Ptr("IoTApps"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Apps_CreateOrUpdate
	fmt.Println("Call operation: Apps_CreateOrUpdate")
	appsClientCreateOrUpdateResponsePoller, err := appsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armiotcentral.App{
		Location: to.Ptr(testsuite.location),
		Identity: &armiotcentral.SystemAssignedServiceIdentity{
			Type: to.Ptr(armiotcentral.SystemAssignedServiceIdentityTypeSystemAssigned),
		},
		Properties: &armiotcentral.AppProperties{
			DisplayName: to.Ptr("My IoT Central App"),
			Subdomain:   to.Ptr("my-iotcentral-app"),
			Template:    to.Ptr("iotc-pnp-preview@1.0.0"),
		},
		SKU: &armiotcentral.AppSKUInfo{
			Name: to.Ptr(armiotcentral.AppSKUST2),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, appsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Apps_ListBySubscription
	fmt.Println("Call operation: Apps_ListBySubscription")
	appsClientNewListBySubscriptionPager := appsClient.NewListBySubscriptionPager(nil)
	for appsClientNewListBySubscriptionPager.More() {
		_, err := appsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Apps_ListByResourceGroup
	fmt.Println("Call operation: Apps_ListByResourceGroup")
	appsClientNewListByResourceGroupPager := appsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for appsClientNewListByResourceGroupPager.More() {
		_, err := appsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Apps_Get
	fmt.Println("Call operation: Apps_Get")
	_, err = appsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step Apps_Update
	fmt.Println("Call operation: Apps_Update")
	appsClientUpdateResponsePoller, err := appsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armiotcentral.AppPatch{
		Identity: &armiotcentral.SystemAssignedServiceIdentity{
			Type: to.Ptr(armiotcentral.SystemAssignedServiceIdentityTypeSystemAssigned),
		},
		Properties: &armiotcentral.AppProperties{
			DisplayName: to.Ptr("My IoT Central App 2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, appsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Apps_CheckSubdomainAvailability
	fmt.Println("Call operation: Apps_CheckSubdomainAvailability")
	_, err = appsClient.CheckSubdomainAvailability(testsuite.ctx, armiotcentral.OperationInputs{
		Name: to.Ptr("myiotcentralapp"),
		Type: to.Ptr("IoTApps"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Apps_ListTemplates
	fmt.Println("Call operation: Apps_ListTemplates")
	appsClientNewListTemplatesPager := appsClient.NewListTemplatesPager(nil)
	for appsClientNewListTemplatesPager.More() {
		_, err := appsClientNewListTemplatesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Apps_Delete
	fmt.Println("Call operation: Apps_Delete")
	appsClientDeleteResponsePoller, err := appsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, appsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.IoTCentral/operations
func (testsuite *IotcentralTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armiotcentral.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

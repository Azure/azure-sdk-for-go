// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v8"
	"github.com/stretchr/testify/suite"
)

type ServiceEndpointPolicyTestSuite struct {
	suite.Suite

	ctx                                 context.Context
	cred                                azcore.TokenCredential
	options                             *arm.ClientOptions
	serviceEndpointPolicyName           string
	serviceEndpointPolicyDefinitionName string
	location                            string
	resourceGroupName                   string
	subscriptionId                      string
}

func (testsuite *ServiceEndpointPolicyTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceEndpointPolicyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceend", 16, false)
	testsuite.serviceEndpointPolicyDefinitionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceenddef", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ServiceEndpointPolicyTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestServiceEndpointPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceEndpointPolicyTestSuite))
}

func (testsuite *ServiceEndpointPolicyTestSuite) Prepare() {
	var err error
	// From step ServiceEndpointPolicies_CreateOrUpdate
	fmt.Println("Call operation: ServiceEndpointPolicies_CreateOrUpdate")
	serviceEndpointPoliciesClient, err := armnetwork.NewServiceEndpointPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceEndpointPoliciesClientCreateOrUpdateResponsePoller, err := serviceEndpointPoliciesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, armnetwork.ServiceEndpointPolicy{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceEndpointPoliciesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/serviceEndpointPolicies/{serviceEndpointPolicyName}
func (testsuite *ServiceEndpointPolicyTestSuite) TestServiceEndpointPolicies() {
	var err error
	// From step ServiceEndpointPolicies_List
	fmt.Println("Call operation: ServiceEndpointPolicies_List")
	serviceEndpointPoliciesClient, err := armnetwork.NewServiceEndpointPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceEndpointPoliciesClientNewListPager := serviceEndpointPoliciesClient.NewListPager(nil)
	for serviceEndpointPoliciesClientNewListPager.More() {
		_, err := serviceEndpointPoliciesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ServiceEndpointPolicies_ListByResourceGroup
	fmt.Println("Call operation: ServiceEndpointPolicies_ListByResourceGroup")
	serviceEndpointPoliciesClientNewListByResourceGroupPager := serviceEndpointPoliciesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for serviceEndpointPoliciesClientNewListByResourceGroupPager.More() {
		_, err := serviceEndpointPoliciesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ServiceEndpointPolicies_Get
	fmt.Println("Call operation: ServiceEndpointPolicies_Get")
	_, err = serviceEndpointPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, &armnetwork.ServiceEndpointPoliciesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step ServiceEndpointPolicies_UpdateTags
	fmt.Println("Call operation: ServiceEndpointPolicies_UpdateTags")
	_, err = serviceEndpointPoliciesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/serviceEndpointPolicies/{serviceEndpointPolicyName}/serviceEndpointPolicyDefinitions/{serviceEndpointPolicyDefinitionName}
func (testsuite *ServiceEndpointPolicyTestSuite) TestServiceEndpointPolicyDefinitions() {
	var err error
	// From step ServiceEndpointPolicyDefinitions_CreateOrUpdate
	fmt.Println("Call operation: ServiceEndpointPolicyDefinitions_CreateOrUpdate")
	serviceEndpointPolicyDefinitionsClient, err := armnetwork.NewServiceEndpointPolicyDefinitionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceEndpointPolicyDefinitionsClientCreateOrUpdateResponsePoller, err := serviceEndpointPolicyDefinitionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, testsuite.serviceEndpointPolicyDefinitionName, armnetwork.ServiceEndpointPolicyDefinition{
		Properties: &armnetwork.ServiceEndpointPolicyDefinitionPropertiesFormat{
			Description: to.Ptr("Storage Service EndpointPolicy Definition"),
			Service:     to.Ptr("Microsoft.Storage"),
			ServiceResources: []*string{
				to.Ptr("/subscriptions/" + testsuite.subscriptionId)},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceEndpointPolicyDefinitionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ServiceEndpointPolicyDefinitions_ListByResourceGroup
	fmt.Println("Call operation: ServiceEndpointPolicyDefinitions_ListByResourceGroup")
	serviceEndpointPolicyDefinitionsClientNewListByResourceGroupPager := serviceEndpointPolicyDefinitionsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, nil)
	for serviceEndpointPolicyDefinitionsClientNewListByResourceGroupPager.More() {
		_, err := serviceEndpointPolicyDefinitionsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ServiceEndpointPolicyDefinitions_Get
	fmt.Println("Call operation: ServiceEndpointPolicyDefinitions_Get")
	_, err = serviceEndpointPolicyDefinitionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, testsuite.serviceEndpointPolicyDefinitionName, nil)
	testsuite.Require().NoError(err)

	// From step ServiceEndpointPolicyDefinitions_Delete
	fmt.Println("Call operation: ServiceEndpointPolicyDefinitions_Delete")
	serviceEndpointPolicyDefinitionsClientDeleteResponsePoller, err := serviceEndpointPolicyDefinitionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, testsuite.serviceEndpointPolicyDefinitionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceEndpointPolicyDefinitionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *ServiceEndpointPolicyTestSuite) Cleanup() {
	var err error
	// From step ServiceEndpointPolicies_Delete
	fmt.Println("Call operation: ServiceEndpointPolicies_Delete")
	serviceEndpointPoliciesClient, err := armnetwork.NewServiceEndpointPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceEndpointPoliciesClientDeleteResponsePoller, err := serviceEndpointPoliciesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceEndpointPolicyName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceEndpointPoliciesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

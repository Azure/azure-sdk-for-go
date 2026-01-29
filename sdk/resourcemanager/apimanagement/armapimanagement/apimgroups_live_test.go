// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armapimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimgroupsTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	groupId			string
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimgroupsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.groupId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "groupid", 13, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicegroup", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimgroupsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimgroupsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimgroupsTestSuite))
}

func (testsuite *ApimgroupsTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name":	to.Ptr("Contoso"),
			"Test":	to.Ptr("User"),
		},
		Location:	to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail:	to.Ptr("foo@contoso.com"),
			PublisherName:	to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:		to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity:	to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/groups
func (testsuite *ApimgroupsTestSuite) TestGroup() {
	var err error
	// From step Group_CreateOrUpdate
	fmt.Println("Call operation: Group_CreateOrUpdate")
	groupClient, err := armapimanagement.NewGroupClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = groupClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.groupId, armapimanagement.GroupCreateParameters{
		Properties: &armapimanagement.GroupCreateParametersProperties{
			DisplayName: to.Ptr("temp group"),
		},
	}, &armapimanagement.GroupClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Group_GetEntityTag
	fmt.Println("Call operation: Group_GetEntityTag")
	_, err = groupClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.groupId, nil)
	testsuite.Require().NoError(err)

	// From step Group_ListByService
	fmt.Println("Call operation: Group_ListByService")
	groupClientNewListByServicePager := groupClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.GroupClientListByServiceOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for groupClientNewListByServicePager.More() {
		_, err := groupClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Group_Get
	fmt.Println("Call operation: Group_Get")
	_, err = groupClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.groupId, nil)
	testsuite.Require().NoError(err)

	// From step Group_Update
	fmt.Println("Call operation: Group_Update")
	_, err = groupClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.groupId, "*", armapimanagement.GroupUpdateParameters{
		Properties: &armapimanagement.GroupUpdateParametersProperties{
			DisplayName: to.Ptr("temp group"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Group_Delete
	fmt.Println("Call operation: Group_Delete")
	_, err = groupClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.groupId, "*", nil)
	testsuite.Require().NoError(err)
}

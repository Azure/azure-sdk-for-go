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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimtagsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	serviceName       string
	tagId             string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimtagsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicetags", 17, false)
	testsuite.tagId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "tagsid", 12, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimtagsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimtagsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimtagsTestSuite))
}

func (testsuite *ApimtagsTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name": to.Ptr("Contoso"),
			"Test": to.Ptr("User"),
		},
		Location: to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail: to.Ptr("foo@contoso.com"),
			PublisherName:  to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:     to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/tags
func (testsuite *ApimtagsTestSuite) TestTags() {
	var err error
	// From step Tag_CreateOrUpdate
	fmt.Println("Call operation: Tag_CreateOrUpdate")
	tagClient, err := armapimanagement.NewTagClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = tagClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, armapimanagement.TagCreateUpdateParameters{
		Properties: &armapimanagement.TagContractProperties{
			DisplayName: to.Ptr(testsuite.tagId),
		},
	}, &armapimanagement.TagClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Tag_ListByService
	fmt.Println("Call operation: Tag_ListByService")
	tagClientNewListByServicePager := tagClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.TagClientListByServiceOptions{Filter: nil,
		Top:   nil,
		Skip:  nil,
		Scope: nil,
	})
	for tagClientNewListByServicePager.More() {
		_, err := tagClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Tag_GetEntityState
	fmt.Println("Call operation: Tag_GetEntityState")
	_, err = tagClient.GetEntityState(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_Get
	fmt.Println("Call operation: Tag_Get")
	_, err = tagClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_Update
	fmt.Println("Call operation: Tag_Update")
	_, err = tagClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, "*", armapimanagement.TagCreateUpdateParameters{
		Properties: &armapimanagement.TagContractProperties{
			DisplayName: to.Ptr("update_" + testsuite.tagId),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Tag_Delete
	fmt.Println("Call operation: Tag_Delete")
	_, err = tagClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, "*", nil)
	testsuite.Require().NoError(err)
}

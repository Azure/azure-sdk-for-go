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

type ApimapiversionsetsTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	serviceName		string
	versionSetId		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimapiversionsetsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceapiversionsets", 27, false)
	testsuite.versionSetId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "versionset", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimapiversionsetsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimapiversionsetsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimapiversionsetsTestSuite))
}

func (testsuite *ApimapiversionsetsTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/apiVersionSets
func (testsuite *ApimapiversionsetsTestSuite) TestApiversionset() {
	var err error
	// From step ApiVersionSet_CreateOrUpdate
	fmt.Println("Call operation: ApiVersionSet_CreateOrUpdate")
	aPIVersionSetClient, err := armapimanagement.NewAPIVersionSetClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIVersionSetClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.versionSetId, armapimanagement.APIVersionSetContract{
		Properties: &armapimanagement.APIVersionSetContractProperties{
			DisplayName:		to.Ptr("api set 1"),
			VersioningScheme:	to.Ptr(armapimanagement.VersioningSchemeSegment),
			Description:		to.Ptr("Version configuration"),
		},
	}, &armapimanagement.APIVersionSetClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step ApiVersionSet_GetEntityTag
	fmt.Println("Call operation: ApiVersionSet_GetEntityTag")
	_, err = aPIVersionSetClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.versionSetId, nil)
	testsuite.Require().NoError(err)

	// From step ApiVersionSet_ListByService
	fmt.Println("Call operation: ApiVersionSet_ListByService")
	aPIVersionSetClientNewListByServicePager := aPIVersionSetClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.APIVersionSetClientListByServiceOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for aPIVersionSetClientNewListByServicePager.More() {
		_, err := aPIVersionSetClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiVersionSet_Get
	fmt.Println("Call operation: ApiVersionSet_Get")
	_, err = aPIVersionSetClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.versionSetId, nil)
	testsuite.Require().NoError(err)

	// From step ApiVersionSet_Update
	fmt.Println("Call operation: ApiVersionSet_Update")
	_, err = aPIVersionSetClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.versionSetId, "*", armapimanagement.APIVersionSetUpdateParameters{
		Properties: &armapimanagement.APIVersionSetUpdateParametersProperties{
			Description:		to.Ptr("Version configuration"),
			DisplayName:		to.Ptr("api set 1"),
			VersioningScheme:	to.Ptr(armapimanagement.VersioningSchemeSegment),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ApiVersionSet_Delete
	fmt.Println("Call operation: ApiVersionSet_Delete")
	_, err = aPIVersionSetClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.versionSetId, "*", nil)
	testsuite.Require().NoError(err)
}

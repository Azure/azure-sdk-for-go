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

type ApimnamedvaluesTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	namedValueId		string
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimnamedvaluesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namedValueId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namedvalue", 16, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicenamed", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimnamedvaluesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimnamedvaluesTestSuite(t *testing.T) {
	suite.Run(t, new(ApimnamedvaluesTestSuite))
}

func (testsuite *ApimnamedvaluesTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/namedValues
func (testsuite *ApimnamedvaluesTestSuite) TestNamedvalue() {
	var err error
	// From step NamedValue_CreateOrUpdate
	fmt.Println("Call operation: NamedValue_CreateOrUpdate")
	namedValueClient, err := armapimanagement.NewNamedValueClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namedValueClientCreateOrUpdateResponsePoller, err := namedValueClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.namedValueId, armapimanagement.NamedValueCreateContract{
		Properties: &armapimanagement.NamedValueCreateContractProperties{
			Secret:	to.Ptr(false),
			Tags: []*string{
				to.Ptr("foo"),
				to.Ptr("bar")},
			DisplayName:	to.Ptr("prop3name"),
			Value:		to.Ptr("propValue"),
		},
	}, &armapimanagement.NamedValueClientBeginCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namedValueClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step NamedValue_GetEntityTag
	fmt.Println("Call operation: NamedValue_GetEntityTag")
	_, err = namedValueClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.namedValueId, nil)
	testsuite.Require().NoError(err)

	// From step NamedValue_ListByService
	fmt.Println("Call operation: NamedValue_ListByService")
	namedValueClientNewListByServicePager := namedValueClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.NamedValueClientListByServiceOptions{Filter: nil,
		Top:				nil,
		Skip:				nil,
		IsKeyVaultRefreshFailed:	nil,
	})
	for namedValueClientNewListByServicePager.More() {
		_, err := namedValueClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NamedValue_Get
	fmt.Println("Call operation: NamedValue_Get")
	_, err = namedValueClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.namedValueId, nil)
	testsuite.Require().NoError(err)

	// From step NamedValue_Update
	fmt.Println("Call operation: NamedValue_Update")
	namedValueClientUpdateResponsePoller, err := namedValueClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.namedValueId, "*", armapimanagement.NamedValueUpdateParameters{
		Properties: &armapimanagement.NamedValueUpdateParameterProperties{
			Secret:	to.Ptr(false),
			Tags: []*string{
				to.Ptr("foo"),
				to.Ptr("bar2")},
			DisplayName:	to.Ptr("prop3name"),
			Value:		to.Ptr("propValue"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namedValueClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step NamedValue_ListValue
	fmt.Println("Call operation: NamedValue_ListValue")
	_, err = namedValueClient.ListValue(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.namedValueId, nil)
	testsuite.Require().NoError(err)

	// From step NamedValue_Delete
	fmt.Println("Call operation: NamedValue_Delete")
	_, err = namedValueClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.namedValueId, "*", nil)
	testsuite.Require().NoError(err)
}

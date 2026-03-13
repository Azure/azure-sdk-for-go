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

type ApimportalrevisionsTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	portalRevisionId	string
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimportalrevisionsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.portalRevisionId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "portalrevi", 16, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicerevision", 21, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimportalrevisionsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimportalrevisionsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimportalrevisionsTestSuite))
}

func (testsuite *ApimportalrevisionsTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/portalRevisions
func (testsuite *ApimportalrevisionsTestSuite) TestPortalrevision() {
	var err error
	// From step PortalRevision_CreateOrUpdate
	fmt.Println("Call operation: PortalRevision_CreateOrUpdate")
	portalRevisionClient, err := armapimanagement.NewPortalRevisionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	portalRevisionClientCreateOrUpdateResponsePoller, err := portalRevisionClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.portalRevisionId, armapimanagement.PortalRevisionContract{
		Properties: &armapimanagement.PortalRevisionContractProperties{
			Description:	to.Ptr("portal revision 1"),
			IsCurrent:	to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, portalRevisionClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PortalRevision_GetEntityTag
	fmt.Println("Call operation: PortalRevision_GetEntityTag")
	_, err = portalRevisionClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.portalRevisionId, nil)
	testsuite.Require().NoError(err)

	// From step PortalRevision_ListByService
	fmt.Println("Call operation: PortalRevision_ListByService")
	portalRevisionClientNewListByServicePager := portalRevisionClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.PortalRevisionClientListByServiceOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for portalRevisionClientNewListByServicePager.More() {
		_, err := portalRevisionClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PortalRevision_Get
	fmt.Println("Call operation: PortalRevision_Get")
	_, err = portalRevisionClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.portalRevisionId, nil)
	testsuite.Require().NoError(err)

	// From step PortalRevision_Update
	fmt.Println("Call operation: PortalRevision_Update")
	portalRevisionClientUpdateResponsePoller, err := portalRevisionClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.portalRevisionId, "*", armapimanagement.PortalRevisionContract{
		Properties: &armapimanagement.PortalRevisionContractProperties{
			Description: to.Ptr("portal revision update"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, portalRevisionClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

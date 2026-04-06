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

type ApimpoliciesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	serviceName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimpoliciesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicepolicy", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimpoliciesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimpoliciesTestSuite(t *testing.T) {
	suite.Run(t, new(ApimpoliciesTestSuite))
}

func (testsuite *ApimpoliciesTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/policies
func (testsuite *ApimpoliciesTestSuite) TestPolicy() {
	var err error
	// From step Policy_CreateOrUpdate
	fmt.Println("Call operation: Policy_CreateOrUpdate")
	policyClient, err := armapimanagement.NewPolicyClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = policyClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.PolicyIDNamePolicy, armapimanagement.PolicyContract{
		Properties: &armapimanagement.PolicyContractProperties{
			Value:  to.Ptr("<policies>\r\n  <inbound />\r\n  <backend>\r\n    <forward-request />\r\n  </backend>\r\n  <outbound />\r\n</policies>"),
			Format: to.Ptr(armapimanagement.PolicyContentFormatXML),
		},
	}, &armapimanagement.PolicyClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Policy_GetEntityTag
	fmt.Println("Call operation: Policy_GetEntityTag")
	_, err = policyClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.PolicyIDNamePolicy, nil)
	testsuite.Require().NoError(err)

	// From step Policy_Get
	fmt.Println("Call operation: Policy_Get")
	_, err = policyClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.PolicyIDNamePolicy, &armapimanagement.PolicyClientGetOptions{Format: nil})
	testsuite.Require().NoError(err)

	// From step Policy_Delete
	fmt.Println("Call operation: Policy_Delete")
	_, err = policyClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.PolicyIDNamePolicy, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/policyDescriptions
func (testsuite *ApimpoliciesTestSuite) TestPolicydescription() {
	var err error
	// From step PolicyDescription_ListByService
	fmt.Println("Call operation: PolicyDescription_ListByService")
	policyDescriptionClient, err := armapimanagement.NewPolicyDescriptionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = policyDescriptionClient.ListByService(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.PolicyDescriptionClientListByServiceOptions{Scope: to.Ptr(armapimanagement.PolicyScopeContractAPI)})
	testsuite.Require().NoError(err)
}

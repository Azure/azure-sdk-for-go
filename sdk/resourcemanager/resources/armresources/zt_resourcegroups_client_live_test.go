//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ResourceGroupsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *ResourceGroupsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armresources/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *ResourceGroupsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestResourceGroupsClient(t *testing.T) {
	suite.Run(t, new(ResourceGroupsClientTestSuite))
}

func (testsuite *ResourceGroupsClientTestSuite) TestResourceGroupsCRUD() {
	// create resource group
	rgName := "go-test-rg"
	rgClient := armresources.NewResourceGroupsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(rgName, *rg.Name)

	// check existence resource group
	check, err := rgClient.CheckExistence(context.Background(), rgName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().True(check.Success)

	// get resource group
	getResp, err := rgClient.Get(context.Background(), rgName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(rgName, *getResp.Name)

	// list resource group
	listPager := rgClient.List(nil)
	testsuite.Require().NoError(listPager.Err())
	testsuite.Require().True(listPager.NextPage(testsuite.ctx))

	// update resource group
	updateResp, err := rgClient.Update(context.Background(), rgName, armresources.ResourceGroupPatchable{
		Tags: map[string]*string{
			"key": to.StringPtr("value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("value", *updateResp.Tags["key"])

	// export template resource group
	pollerResp, err := rgClient.BeginExportTemplate(context.Background(), rgName, armresources.ExportTemplateRequest{
		Resources: []*string{
			to.StringPtr("*"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	templateResp, err := pollerResp.PollUntilDone(context.Background(), 10*time.Second)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(templateResp.Template)

	// clean resource group
	delPollerResp, err := rgClient.BeginDelete(context.Background(), rgName, nil)
	testsuite.Require().NoError(err)
	delResp, err := delPollerResp.PollUntilDone(context.Background(), 10*time.Second)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}

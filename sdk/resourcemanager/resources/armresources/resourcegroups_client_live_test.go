//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
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
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
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
	fmt.Println("Call operation: ResourceGroups_CreateOrUpdate")
	rgName := "go-test-rg"
	rgClient, err := armresources.NewResourceGroupsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	rg, err := rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.Ptr("eastus"),
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(rgName, *rg.Name)

	// check existence resource group
	fmt.Println("Call operation: ResourceGroups_CheckExistence")
	check, err := rgClient.CheckExistence(context.Background(), rgName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().True(check.Success)

	// get resource group
	fmt.Println("Call operation: RResourceGroups_Get")
	getResp, err := rgClient.Get(context.Background(), rgName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(rgName, *getResp.Name)

	// list resource group
	fmt.Println("Call operation: ResourceGroups_List")
	listPager := rgClient.NewListPager(nil)
	testsuite.Require().True(listPager.More())

	// update resource group
	fmt.Println("Call operation: ResourceGroups_Update")
	updateResp, err := rgClient.Update(context.Background(), rgName, armresources.ResourceGroupPatchable{
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("value", *updateResp.Tags["key"])

	// export template resource group
	fmt.Println("Call operation: ResourceGroups_ExportTemplate")
	pollerResp, err := rgClient.BeginExportTemplate(context.Background(), rgName, armresources.ExportTemplateRequest{
		Resources: []*string{
			to.Ptr("*"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	templateResp, err := pollerResp.PollUntilDone(context.Background(), nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(templateResp.Template)

	// clean resource group
	fmt.Println("Call operation: ResourceGroups_Delete")
	delPollerResp, err := rgClient.BeginDelete(context.Background(), rgName, nil)
	testsuite.Require().NoError(err)
	_, err = delPollerResp.PollUntilDone(context.Background(), nil)
	testsuite.Require().NoError(err)
}

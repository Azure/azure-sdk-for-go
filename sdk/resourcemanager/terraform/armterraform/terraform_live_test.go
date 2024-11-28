//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armterraform_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/terraform/armterraform"
	"github.com/stretchr/testify/suite"
)

const (
	// ResourceLocation = "eastus2"
	ResourceLocation = "westus"
)

type TerraformTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *TerraformTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	fmt.Println("testsuite.resourceGroupName:", testsuite.resourceGroupName)
	testsuite.Prepare()
}

func (testsuite *TerraformTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	time.Sleep(time.Second * 3)
	testutil.StopRecording(testsuite.T())
}

func TestTerraformTestSuite(t *testing.T) {
	suite.Run(t, new(TerraformTestSuite))
}

func (testsuite *TerraformTestSuite) TestOptionList() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armterraform.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(testsuite.ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			fmt.Println("v.ActionType", v.ActionType)
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		page = armterraform.OperationsClientListResponse{
			OperationListResult: armterraform.OperationListResult{
				Value: []*armterraform.Operation{
					{
						Name: to.Ptr("Microsoft.AzureTerraform/operations/read"),
						Display: &armterraform.OperationDisplay{
							Provider:    to.Ptr("Microsoft AzureTerraform"),
							Resource:    to.Ptr("Azure Terraform Resource Provider"),
							Operation:   to.Ptr("ListOperations"),
							Description: to.Ptr("Lists all of the available RP operations."),
						},
					},
					{
						Name: to.Ptr("Microsoft.AzureTerraform/exportTerraform/action"),
						Display: &armterraform.OperationDisplay{
							Provider:    to.Ptr("Microsoft AzureTerraform"),
							Resource:    to.Ptr("Azure Terraform Resource Provider"),
							Operation:   to.Ptr("ExportTerraform"),
							Description: to.Ptr("Exports the Terraform configuration used for the specified scope."),
						},
					},
				},
			},
		}
	}
}

func (testsuite *TerraformTestSuite) Prepare() {
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	testsuite.Require().NoError(err)
	// new client factory

	fmt.Println("subscriptionId", testsuite.subscriptionId, "groupName", testsuite.resourceGroupName, "location", testsuite.location)
	clientFactory, err := armresources.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	testsuite.Require().NoError(err)
	client := clientFactory.NewResourceGroupsClient()
	ctx := context.Background()

	testsuite.Require().NoError(err)
	// check whether create new group successfully
	res, err := client.CheckExistence(ctx, testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)
	if !res.Success {
		_, err = client.CreateOrUpdate(ctx, testsuite.resourceGroupName, armresources.ResourceGroup{
			Location: to.Ptr(testsuite.location),
		}, nil)
		testsuite.Require().NoError(err)
	}

	fmt.Println("create new resource group ", testsuite.resourceGroupName, " of ", testsuite.subscriptionId, "successfully")
}

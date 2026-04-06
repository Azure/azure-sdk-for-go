// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armtemplatespecs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armtemplatespecs"
	"github.com/stretchr/testify/suite"
)

type TemplateSpecsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *TemplateSpecsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), pathToPackage)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *TemplateSpecsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestTemplateSpecsClient(t *testing.T) {
	suite.Run(t, new(TemplateSpecsClientTestSuite))
}

func (testsuite *TemplateSpecsClientTestSuite) TestTemplateSpecsCRUD() {
	// create template spec
	fmt.Println("Call operation: TemplateSpecs_CreateOrUpdate")
	templateSpecName := "go-test-template"
	templateSpecsClient, err := armtemplatespecs.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = templateSpecsClient.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		templateSpecName,
		armtemplatespecs.TemplateSpec{
			Location: to.Ptr(testsuite.location),
			Properties: &armtemplatespecs.TemplateSpecProperties{
				Description: to.Ptr("template spec properties description."),
				Metadata: map[string]string{
					"live": "test",
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)

	// update template spec
	fmt.Println("Call operation: TemplateSpecs_Update")
	_, err = templateSpecsClient.Update(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, &armtemplatespecs.ClientUpdateOptions{
		TemplateSpec: &armtemplatespecs.TemplateSpecUpdateModel{
			Tags: map[string]*string{
				"test": to.Ptr("live"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// get template spec
	fmt.Println("Call operation: TemplateSpecs_Get")
	_, err = templateSpecsClient.Get(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, nil)
	testsuite.Require().NoError(err)

	// list template spec by resource group
	fmt.Println("Call operation: TemplateSpecs_ListByResourceGroup")
	listByResourceGroup := templateSpecsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	testsuite.Require().True(listByResourceGroup.More())

	// list template spec by subscription
	fmt.Println("Call operation: TemplateSpecs_ListBySubscription")
	listBySubscription := templateSpecsClient.NewListBySubscriptionPager(nil)
	testsuite.Require().True(listBySubscription.More())

	// delete template spec
	fmt.Println("Call operation: TemplateSpecs_Delete")
	_, err = templateSpecsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, nil)
	testsuite.Require().NoError(err)
}

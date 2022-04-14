//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armtemplatespecs_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
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
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armtemplatespecs/testdata")
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
	templateSpecName := "go-test-template"
	templateSpecsClient, err := armtemplatespecs.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resp, err := templateSpecsClient.CreateOrUpdate(
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
	testsuite.Require().Equal(templateSpecName, *resp.Name)

	// update template spec
	updateResp, err := templateSpecsClient.Update(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, &armtemplatespecs.ClientUpdateOptions{
		TemplateSpec: &armtemplatespecs.TemplateSpecUpdateModel{
			Tags: map[string]*string{
				"test": to.Ptr("live"),
			},
		},
	})
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("live", *updateResp.Tags["test"])

	// get template spec
	getResp, err := templateSpecsClient.Get(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(templateSpecName, *getResp.Name)

	// list template spec by resource group
	listByResourceGroup := templateSpecsClient.ListByResourceGroup(testsuite.resourceGroupName, nil)
	testsuite.Require().True(listByResourceGroup.More())

	// list template spec by subscription
	listBySubscription := templateSpecsClient.ListBySubscription(nil)
	testsuite.Require().True(listBySubscription.More())

	// delete template spec
	_, err = templateSpecsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, nil)
	testsuite.Require().NoError(err)
}

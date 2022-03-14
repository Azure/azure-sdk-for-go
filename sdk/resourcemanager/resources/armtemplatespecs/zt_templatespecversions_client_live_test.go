//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armtemplatespecs_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armtemplatespecs"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TemplateSpecVersionsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *TemplateSpecVersionsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armtemplatespecs/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *TemplateSpecVersionsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestTemplateSpecVersionsClient(t *testing.T) {
	suite.Run(t, new(TemplateSpecVersionsClientTestSuite))
}

func (testsuite *TemplateSpecVersionsClientTestSuite) TestTemplateSpecVersionsCRUD() {
	// create template spec
	templateSpecName := "go-test-template"
	templateSpecsClient := armtemplatespecs.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	resp, err := templateSpecsClient.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		templateSpecName,
		armtemplatespecs.TemplateSpec{
			Location: to.StringPtr(testsuite.location),
			Properties: &armtemplatespecs.TemplateSpecProperties{
				Description: to.StringPtr("template spec properties description."),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(templateSpecName, *resp.Name)

	// create template version
	templateSpecVersion := "go-test-template-version"
	templateSpecVersionsClient := armtemplatespecs.NewTemplateSpecVersionsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	vResp, err := templateSpecVersionsClient.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		templateSpecName,
		templateSpecVersion,
		armtemplatespecs.TemplateSpecVersion{
			Location: to.StringPtr(testsuite.location),
			Properties: &armtemplatespecs.TemplateSpecVersionProperties{
				Description: to.StringPtr("<description>"),
				MainTemplate: map[string]interface{}{
					"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
					"contentVersion": "1.0.0.0",
					"parameters":     map[string]interface{}{},
					"resources":      []interface{}{},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(vResp)

	// update
	updateResp, err := templateSpecVersionsClient.Update(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, templateSpecVersion, &armtemplatespecs.TemplateSpecVersionsClientUpdateOptions{
		TemplateSpecVersionUpdateModel: &armtemplatespecs.TemplateSpecVersionUpdateModel{
			Tags: map[string]*string{
				"test": to.StringPtr("live"),
			},
		},
	})
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("live", *updateResp.Tags["test"])

	// get
	getResp, err := templateSpecVersionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, templateSpecVersion, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(templateSpecVersion, *getResp.Name)

	// list
	pager := templateSpecVersionsClient.List(testsuite.resourceGroupName, templateSpecName, nil)
	testsuite.Require().NoError(pager.Err())

	// delete
	delResp, err := templateSpecVersionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, templateSpecVersion, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}

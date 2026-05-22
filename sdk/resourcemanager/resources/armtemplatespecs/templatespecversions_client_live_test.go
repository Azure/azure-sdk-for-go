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
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), pathToPackage)
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

	// create template version
	fmt.Println("Call operation: TemplateSpecVersions_CreateOrUpdate")
	templateSpecVersion := "go-test-template-version"
	templateSpecVersionsClient, err := armtemplatespecs.NewTemplateSpecVersionsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vResp, err := templateSpecVersionsClient.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		templateSpecName,
		templateSpecVersion,
		armtemplatespecs.TemplateSpecVersion{
			Location: to.Ptr(testsuite.location),
			Properties: &armtemplatespecs.TemplateSpecVersionProperties{
				Description: to.Ptr("<description>"),
				MainTemplate: map[string]interface{}{
					"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
					"contentVersion": "1.0.0.0",
					"parameters":     map[string]interface{}{},
					"resources":      []interface{}{},
				},
				Metadata: map[string]string{
					"live": "test",
				},
				UIFormDefinition: map[string]string{
					"live": "test",
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(vResp)

	// update
	fmt.Println("Call operation: TemplateSpecVersions_Update")
	_, err = templateSpecVersionsClient.Update(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, templateSpecVersion, &armtemplatespecs.TemplateSpecVersionsClientUpdateOptions{
		TemplateSpecVersionUpdateModel: &armtemplatespecs.TemplateSpecVersionUpdateModel{
			Tags: map[string]*string{
				"test": to.Ptr("live"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// get
	fmt.Println("Call operation: TemplateSpecVersions_Get")
	_, err = templateSpecVersionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, templateSpecVersion, nil)
	testsuite.Require().NoError(err)

	// list
	fmt.Println("Call operation: TemplateSpecVersions_List")
	pager := templateSpecVersionsClient.NewListPager(testsuite.resourceGroupName, templateSpecName, nil)
	testsuite.Require().True(pager.More())

	// delete
	fmt.Println("Call operation: TemplateSpecVersions_Delete")
	_, err = templateSpecVersionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, templateSpecName, templateSpecVersion, nil)
	testsuite.Require().NoError(err)
}

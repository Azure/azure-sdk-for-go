//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type DeploymentsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *DeploymentsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armresources/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *DeploymentsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDeploymentsClient(t *testing.T) {
	suite.Run(t, new(DeploymentsClientTestSuite))
}

var template = `
{
 "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
 "contentVersion": "1.0.0.0",
 "parameters": {
   "location": {
     "type": "string",
     "allowedValues": [
       "East US",
       "West US",
       "West Europe",
       "East Asia",
       "South East Asia"
     ],
     "metadata": {
       "description": "Location to deploy to"
     }
   }
 },
 "resources": [
   {
     "type": "Microsoft.Compute/availabilitySets",
     "name": "availabilitySet1",
     "apiVersion": "2019-07-01",
     "location": "[parameters('location')]",
     "properties": {}
   }
 ],
 "outputs": {
   "myparameter": {
     "type": "object",
     "value": "[reference('Microsoft.Compute/availabilitySets/availabilitySet1')]"
   }
 }
}
`

func unmarshalTemplate(data string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal template error:%v", err)
	}
	return result, nil
}

func (testsuite *DeploymentsClientTestSuite) TestDeploymentsCRUD() {
	// check existence deployment
	deploymentsClient, err := armresources.NewDeploymentsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deploymentName := "go-test-deployment"
	check, err := deploymentsClient.CheckExistence(testsuite.ctx, testsuite.resourceGroupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().False(check.Success)

	tmp, err := unmarshalTemplate(template)

	// create deployment
	testsuite.Require().NoError(err)
	pollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     to.Ptr(armresources.DeploymentModeIncremental),
				Template: tmp,
				Parameters: map[string]interface{}{
					"location": map[string]string{
						"value": "West US",
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	resp, err := testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *resp.Name)

	// get
	getResp, err := deploymentsClient.Get(testsuite.ctx, testsuite.resourceGroupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *getResp.Name)

	// list by resource group
	listPager := deploymentsClient.ListByResourceGroup(testsuite.resourceGroupName, nil)
	testsuite.Require().True(listPager.More())

	// what if
	whatPoller, err := deploymentsClient.BeginWhatIf(
		testsuite.ctx,
		testsuite.resourceGroupName,
		deploymentName,
		armresources.DeploymentWhatIf{
			Properties: &armresources.DeploymentWhatIfProperties{
				Mode:     to.Ptr(armresources.DeploymentModeIncremental),
				Template: tmp,
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	whatResp, err := testutil.PollForTest(testsuite.ctx, whatPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("InvalidTemplate", *whatResp.Error.Code)

	// validate
	vPoller, err := deploymentsClient.BeginValidate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     to.Ptr(armresources.DeploymentModeIncremental),
				Template: tmp,
				Parameters: map[string]interface{}{
					"location": map[string]string{
						"value": "West US",
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vResp, err := testutil.PollForTest(testsuite.ctx, vPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(to.Ptr(armresources.DeploymentModeIncremental), vResp.Properties.Mode)

	// export template
	exportTemplate, err := deploymentsClient.ExportTemplate(testsuite.ctx, testsuite.resourceGroupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(exportTemplate)

	// delete deployment
	delPoller, err := deploymentsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)
}

func (testsuite *DeploymentsClientTestSuite) TestDeploymentsAtScope() {
	// check deployment existence
	deploymentsClient, err := armresources.NewDeploymentsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deploymentName := "go-test-deployment-scope"
	scopeResource := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v", testsuite.subscriptionID, testsuite.resourceGroupName)
	check, err := deploymentsClient.CheckExistenceAtScope(testsuite.ctx, scopeResource, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().False(check.Success)

	temp, err := unmarshalTemplate(template)
	testsuite.Require().NoError(err)

	// create deployment at scope
	pollerResp, err := deploymentsClient.BeginCreateOrUpdateAtScope(
		testsuite.ctx,
		scopeResource,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     to.Ptr(armresources.DeploymentModeIncremental),
				Template: temp,
				Parameters: map[string]interface{}{
					"location": map[string]string{
						"value": "West US",
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	resp, err := testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *resp.Name)

	// get deployment at scope
	getResp, err := deploymentsClient.GetAtScope(testsuite.ctx, scopeResource, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *getResp.Name)

	// list deployment at scope
	listPager := deploymentsClient.ListAtScope(scopeResource, nil)
	testsuite.Require().True(listPager.More())

	vPoller, err := deploymentsClient.BeginValidateAtScope(
		testsuite.ctx,
		scopeResource,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:     to.Ptr(armresources.DeploymentModeIncremental),
				Template: temp,
				Parameters: map[string]interface{}{
					"location": map[string]string{
						"value": "West US",
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vResp, err := testutil.PollForTest(testsuite.ctx, vPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(armresources.DeploymentModeIncremental, *vResp.Properties.Mode)

	// export template
	exportTemplate, err := deploymentsClient.ExportTemplateAtScope(testsuite.ctx, scopeResource, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(exportTemplate)

	// delete deployment
	delPoller, err := deploymentsClient.BeginDeleteAtScope(testsuite.ctx, scopeResource, deploymentName, nil)
	testsuite.Require().NoError(err)
	_, err = delPoller.PollUntilDone(testsuite.ctx, 10*time.Second)
	testsuite.Require().NoError(err)
}

func (testsuite *DeploymentsClientTestSuite) TestDeploymentsAtManagementGroupScope() {
	// create management group
	managementGroupsClient, err := armmanagementgroups.NewClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	groupName := "00000000-0000-0000-0000-000000000000000"
	mgPoller, err := managementGroupsClient.BeginCreateOrUpdate(
		testsuite.ctx,
		groupName,
		armmanagementgroups.CreateManagementGroupRequest{
			Name: to.Ptr(groupName),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	mgResp, err := testutil.PollForTest(testsuite.ctx, mgPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(groupName, *mgResp.Name)

	// create deployment
	deploymentsClient, err := armresources.NewDeploymentsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deploymentName := "go-test-deploymentMG"
	pollerResp, err := deploymentsClient.BeginCreateOrUpdateAtManagementGroupScope(
		testsuite.ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Location: to.Ptr("West US"),
			Properties: &armresources.DeploymentProperties{
				Mode: to.Ptr(armresources.DeploymentModeIncremental),
				TemplateLink: &armresources.TemplateLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
				ParametersLink: &armresources.ParametersLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	resp, err := testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *resp.Name)

	// check
	check, err := deploymentsClient.CheckExistenceAtManagementGroupScope(testsuite.ctx, groupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().True(check.Success)

	// get deployment
	getResp, err := deploymentsClient.GetAtManagementGroupScope(testsuite.ctx, groupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *getResp.Name)

	// list deployment
	listPager := deploymentsClient.ListAtManagementGroupScope(groupName, nil)
	testsuite.Require().True(listPager.More())

	// validate deployment
	validatePoller, err := deploymentsClient.BeginValidateAtManagementGroupScope(
		testsuite.ctx,
		groupName,
		deploymentName,
		armresources.ScopedDeployment{
			Location: to.Ptr("West US"),
			Properties: &armresources.DeploymentProperties{
				Mode: to.Ptr(armresources.DeploymentModeIncremental),
				TemplateLink: &armresources.TemplateLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
				ParametersLink: &armresources.ParametersLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	validateResp, err := validatePoller.PollUntilDone(testsuite.ctx, 10*time.Second)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(validateResp.Properties)

	// export template deployment
	exportResp, err := deploymentsClient.ExportTemplateAtManagementGroupScope(testsuite.ctx, groupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(exportResp.Template)

	// delete template deployment
	delPoller, err := deploymentsClient.BeginDeleteAtManagementGroupScope(testsuite.ctx, groupName, deploymentName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)

	// delete management group
	delManagementGroup, err := managementGroupsClient.BeginDelete(testsuite.ctx, groupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delManagementGroup)
	testsuite.Require().NoError(err)
}

func (testsuite *DeploymentsClientTestSuite) TestDeploymentsAtTenantScope() {
	// check deployment existence
	deploymentsClient, err := armresources.NewDeploymentsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deploymentName := "go-test-deployment-at-tenant"
	resp, err := deploymentsClient.CheckExistenceAtTenantScope(testsuite.ctx, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().False(resp.Success)

	// list deployment
	listResp := deploymentsClient.ListAtTenantScope(nil)
	testsuite.Require().True(listResp.More())
}

func (testsuite *DeploymentsClientTestSuite) TestDeploymentsAtSubscriptionScope() {
	// check deployment existence
	deploymentsClient, err := armresources.NewDeploymentsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deploymentName := "go-test-at-subscription"
	check, err := deploymentsClient.CheckExistenceAtSubscriptionScope(testsuite.ctx, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().False(check.Success)

	// create deployment at subscription scope
	pollerResp, err := deploymentsClient.BeginCreateOrUpdateAtSubscriptionScope(
		testsuite.ctx,
		deploymentName,
		armresources.Deployment{
			Location: to.Ptr("West US"),
			Properties: &armresources.DeploymentProperties{
				Mode: to.Ptr(armresources.DeploymentModeIncremental),
				TemplateLink: &armresources.TemplateLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
				ParametersLink: &armresources.ParametersLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	resp, err := testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *resp.Name)

	// get deployment
	getResp, err := deploymentsClient.GetAtSubscriptionScope(testsuite.ctx, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(deploymentName, *getResp.Name)

	// list deployment
	listResp := deploymentsClient.ListAtSubscriptionScope(nil)
	testsuite.Require().True(listResp.More())

	// what if deployment
	whatIfPoller, err := deploymentsClient.BeginWhatIfAtSubscriptionScope(
		testsuite.ctx,
		deploymentName,
		armresources.DeploymentWhatIf{
			Location: to.Ptr("West US"),
			Properties: &armresources.DeploymentWhatIfProperties{
				Mode: to.Ptr(armresources.DeploymentModeIncremental),
				TemplateLink: &armresources.TemplateLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
				ParametersLink: &armresources.ParametersLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, whatIfPoller)
	testsuite.Require().NoError(err)

	// validate deployment
	validatePoller, err := deploymentsClient.BeginValidateAtSubscriptionScope(
		testsuite.ctx,
		deploymentName,
		armresources.Deployment{
			Location: to.Ptr("West US"),
			Properties: &armresources.DeploymentProperties{
				Mode: to.Ptr(armresources.DeploymentModeIncremental),
				TemplateLink: &armresources.TemplateLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
				ParametersLink: &armresources.ParametersLink{
					URI: to.Ptr("https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, validatePoller)
	testsuite.Require().NoError(err)

	// export template deployment
	exportResp, err := deploymentsClient.ExportTemplateAtSubscriptionScope(testsuite.ctx, deploymentName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(exportResp.Template)

	// delete deployment
	delPoller, err := deploymentsClient.BeginDeleteAtSubscriptionScope(testsuite.ctx, deploymentName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)
}

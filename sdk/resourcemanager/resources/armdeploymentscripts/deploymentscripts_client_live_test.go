//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdeploymentscripts_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeploymentscripts"
	"github.com/stretchr/testify/suite"
)

type DeploymentScriptsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *DeploymentScriptsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/resources/armdeploymentscripts/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *DeploymentScriptsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDeploymentScriptsClient(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	suite.Run(t, new(DeploymentScriptsClientTestSuite))
}

func (testsuite *DeploymentScriptsClientTestSuite) TestDeploymentScriptsCRUD() {
	// create identity
	userAssignedIdentitiesClient, err := armmsi.NewUserAssignedIdentitiesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	identityName := "go-test-identity"
	identityResp, err := userAssignedIdentitiesClient.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		identityName,
		armmsi.Identity{
			Location: to.Ptr(testsuite.location),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(identityName, *identityResp.Name)

	// create deployment script
	fmt.Println("Call operation: DeploymentScripts_Create")
	deploymentScriptsClient, err := armdeploymentscripts.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scriptName := "go-test-script"
	dsPollerResp, err := deploymentScriptsClient.BeginCreate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scriptName,
		&armdeploymentscripts.AzurePowerShellScript{
			Identity: &armdeploymentscripts.ManagedServiceIdentity{
				Type: to.Ptr(armdeploymentscripts.ManagedServiceIdentityTypeUserAssigned),
				UserAssignedIdentities: map[string]*armdeploymentscripts.UserAssignedIdentity{
					*identityResp.ID: {},
				},
			},
			Kind:     to.Ptr(armdeploymentscripts.ScriptTypeAzurePowerShell),
			Location: to.Ptr(testsuite.location),
			Properties: &armdeploymentscripts.AzurePowerShellScriptProperties{
				RetentionInterval:   to.Ptr("PT26H"),
				PrimaryScriptURI:    to.Ptr("https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"),
				Arguments:           to.Ptr("-name \"John Dole\""),
				Timeout:             to.Ptr("PT30M"),
				AzPowerShellVersion: to.Ptr("3.0"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	dsResp, err := testutil.PollForTest(testsuite.ctx, dsPollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(scriptName, *dsResp.GetDeploymentScript().Name)

	// get deployment scripts
	fmt.Println("Call operation: DeploymentScripts_Get")
	getResp, err := deploymentScriptsClient.Get(testsuite.ctx, testsuite.resourceGroupName, scriptName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(scriptName, *getResp.GetDeploymentScript().Name)

	// get log
	fmt.Println("Call operation: DeploymentScripts_GetLogs")
	getLogResp, err := deploymentScriptsClient.GetLogs(testsuite.ctx, testsuite.resourceGroupName, scriptName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(1, len(getLogResp.Value))

	// getLogsDefault
	fmt.Println("Call operation: DeploymentScripts_GetLogsDefault")
	getLogDefaultResp, err := deploymentScriptsClient.GetLogsDefault(testsuite.ctx, testsuite.resourceGroupName, scriptName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("default", *getLogDefaultResp.Name)

	// update deployment script
	fmt.Println("Call operation: DeploymentScripts_Update")
	updateResp, err := deploymentScriptsClient.Update(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scriptName,
		&armdeploymentscripts.ClientUpdateOptions{
			DeploymentScript: &armdeploymentscripts.DeploymentScriptUpdateParameter{
				Tags: map[string]*string{
					"test": to.Ptr("live"),
				},
			},
		},
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("live", *updateResp.GetDeploymentScript().Tags["test"])

	// list deployment script by subscription
	fmt.Println("Call operation: DeploymentScripts_ListBySubscription")
	listBySubscription := deploymentScriptsClient.NewListBySubscriptionPager(nil)
	testsuite.Require().True(listBySubscription.More())

	// list deployment script by resource group
	fmt.Println("Call operation: DeploymentScripts_ListByResourceGroup")
	listByResourceGroup := deploymentScriptsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	testsuite.Require().True(listByResourceGroup.More())

	// delete deployment script
	fmt.Println("Call operation: DeploymentScripts_Delete")
	_, err = deploymentScriptsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, scriptName, nil)
	testsuite.Require().NoError(err)
}

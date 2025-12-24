// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armoperationalinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v3"
	"github.com/stretchr/testify/suite"
)

type IntelligencePacksTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	workspaceName     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *IntelligencePacksTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.workspaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "oiintelliagencepack", 25, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *IntelligencePacksTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestIntelligencePacksTestSuite(t *testing.T) {
	suite.Run(t, new(IntelligencePacksTestSuite))
}

func (testsuite *IntelligencePacksTestSuite) Prepare() {
	var err error
	// From step Workspaces_Create
	fmt.Println("Call operation: Workspaces_CreateOrUpdate")
	workspacesClient, err := armoperationalinsights.NewWorkspacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	workspacesClientCreateOrUpdateResponsePoller, err := workspacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, armoperationalinsights.Workspace{
		Location: to.Ptr(testsuite.location),
		Properties: &armoperationalinsights.WorkspaceProperties{
			RetentionInDays: to.Ptr[int32](30),
			SKU: &armoperationalinsights.WorkspaceSKU{
				Name: to.Ptr(armoperationalinsights.WorkspaceSKUNameEnumPerGB2018),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, workspacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.OperationalInsights/workspaces/intelligencePacks/Disable
func (testsuite *IntelligencePacksTestSuite) TestIntelligencePack() {
	var err error
	// From step IntelligencePacks_List
	fmt.Println("Call operation: IntelligencePacks_List")
	intelligencePacksClient, err := armoperationalinsights.NewIntelligencePacksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = intelligencePacksClient.List(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, nil)
	testsuite.Require().NoError(err)

	// From step IntelligencePacks_Enable
	fmt.Println("Call operation: IntelligencePacks_Enable")
	_, err = intelligencePacksClient.Enable(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "ChangeTracking", nil)
	testsuite.Require().NoError(err)

	// From step IntelligencePacks_Disable
	fmt.Println("Call operation: IntelligencePacks_Disable")
	_, err = intelligencePacksClient.Disable(testsuite.ctx, testsuite.resourceGroupName, testsuite.workspaceName, "ChangeTracking", nil)
	testsuite.Require().NoError(err)
}

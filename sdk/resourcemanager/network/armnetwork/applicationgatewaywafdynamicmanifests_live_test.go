// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v8"
	"github.com/stretchr/testify/suite"
)

type ApplicationGatewayWafDynamicManifestsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApplicationGatewayWafDynamicManifestsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ApplicationGatewayWafDynamicManifestsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApplicationGatewayWafDynamicManifestsTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationGatewayWafDynamicManifestsTestSuite))
}

// Microsoft.Network/locations/{location}/applicationGatewayWafDynamicManifests
func (testsuite *ApplicationGatewayWafDynamicManifestsTestSuite) TestApplicationGatewayWafDynamicManifests() {
	var err error
	// From step ApplicationGatewayWafDynamicManifests_Get
	fmt.Println("Call operation: ApplicationGatewayWafDynamicManifests_Get")
	applicationGatewayWafDynamicManifestsClient, err := armnetwork.NewApplicationGatewayWafDynamicManifestsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	applicationGatewayWafDynamicManifestsClientNewGetPager := applicationGatewayWafDynamicManifestsClient.NewGetPager(testsuite.location, nil)
	for applicationGatewayWafDynamicManifestsClientNewGetPager.More() {
		_, err := applicationGatewayWafDynamicManifestsClientNewGetPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

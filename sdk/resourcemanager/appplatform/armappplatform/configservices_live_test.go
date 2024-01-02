//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armappplatform_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appplatform/armappplatform/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/stretchr/testify/suite"
)

type ConfigservicesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	serviceName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ConfigservicesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/appplatform/armappplatform/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicen", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ConfigservicesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestConfigservicesTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigservicesTestSuite))
}

func (testsuite *ConfigservicesTestSuite) Prepare() {
	var err error
	// From step Services_CreateOrUpdate
	fmt.Println("Call operation: Services_CreateOrUpdate")
	servicesClient, err := armappplatform.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	servicesClientCreateOrUpdateResponsePoller, err := servicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.ServiceResource{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armappplatform.ClusterResourceProperties{},
		SKU: &armappplatform.SKU{
			Name: to.Ptr("S0"),
			Tier: to.Ptr("Standard"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/configServers/default
func (testsuite *ConfigservicesTestSuite) TestConfigServers() {
	var err error
	// From step ConfigServers_UpdatePut
	fmt.Println("Call operation: ConfigServers_UpdatePut")
	configServersClient, err := armappplatform.NewConfigServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configServersClientUpdatePutResponsePoller, err := configServersClient.BeginUpdatePut(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.ConfigServerResource{
		Properties: &armappplatform.ConfigServerProperties{
			ConfigServer: &armappplatform.ConfigServerSettings{
				GitProperty: &armappplatform.ConfigServerGitProperty{
					Label: to.Ptr("main"),
					SearchPaths: []*string{
						to.Ptr("/")},
					URI: to.Ptr("https://github.com/Azure/azure-sdk-for-go.git"),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configServersClientUpdatePutResponsePoller)
	testsuite.Require().NoError(err)

	// From step ConfigServers_Get
	fmt.Println("Call operation: ConfigServers_Get")
	_, err = configServersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step ConfigServers_UpdatePatch
	fmt.Println("Call operation: ConfigServers_UpdatePatch")
	configServersClientUpdatePatchResponsePoller, err := configServersClient.BeginUpdatePatch(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.ConfigServerResource{
		Properties: &armappplatform.ConfigServerProperties{
			ConfigServer: &armappplatform.ConfigServerSettings{
				GitProperty: &armappplatform.ConfigServerGitProperty{
					Label: to.Ptr("main"),
					SearchPaths: []*string{
						to.Ptr("/")},
					URI: to.Ptr("https://github.com/Azure/azure-sdk-for-go.git"),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configServersClientUpdatePatchResponsePoller)
	testsuite.Require().NoError(err)

	// From step ConfigServers_Validate
	fmt.Println("Call operation: ConfigServers_Validate")
	configServersClientValidateResponsePoller, err := configServersClient.BeginValidate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.ConfigServerSettings{
		GitProperty: &armappplatform.ConfigServerGitProperty{
			Label: to.Ptr("main"),
			SearchPaths: []*string{
				to.Ptr("/")},
			URI: to.Ptr("https://github.com/Azure/azure-sdk-for-go.git"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configServersClientValidateResponsePoller)
	testsuite.Require().NoError(err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armchaos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/chaos/armchaos/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type OperationsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *OperationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *OperationsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(OperationsTestSuite))
}

// Microsoft.Chaos/operations
func (testsuite *OperationsTestSuite) TestOperation() {
	var err error
	// From step Operations_ListAll
	fmt.Println("Call operation: Operations_ListAll")
	operationsClient, err := armchaos.NewOperationStatusesClient(testsuite.subscriptionId, testsuite.cred, nil)
	testsuite.Require().NoError(err)
	_, err = operationsClient.Get(testsuite.ctx, testsuite.location, testsuite.subscriptionId, nil)
	testsuite.Require().NoError(err)
}

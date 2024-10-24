// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgezones_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgezones/armedgezones"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type EdgezonesTestSuite struct {
	suite.Suite

	ctx              context.Context
	cred             azcore.TokenCredential
	options          *arm.ClientOptions
	clientFactory    *armedgezones.ClientFactory
	armEndpoint      string
	extendedZoneName string
	location         string
	subscriptionId   string
}

func (testsuite *EdgezonesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	var err error
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.extendedZoneName = "losangeles"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.clientFactory, err = armedgezones.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
}

func (testsuite *EdgezonesTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestEdgezonesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgezonesTestSuite))
}

// Microsoft.EdgeZones/extendedZones/{extendedZoneName}
func (testsuite *EdgezonesTestSuite) TestExtendedZones() {
	var err error

	// From step ExtendedZones_Register
	fmt.Println("Call operation: ExtendedZones_Register")
	extendedZonesClient := testsuite.clientFactory.NewExtendedZonesClient()
	_, err = extendedZonesClient.Register(testsuite.ctx, testsuite.extendedZoneName, nil)
	testsuite.Require().NoError(err)

	// From step ExtendedZones_ListBySubscription
	fmt.Println("Call operation: ExtendedZones_ListBySubscription")
	extendedZonesClientNewListBySubscriptionPager := extendedZonesClient.NewListBySubscriptionPager(nil)
	for extendedZonesClientNewListBySubscriptionPager.More() {
		_, err := extendedZonesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ExtendedZones_Get
	fmt.Println("Call operation: ExtendedZones_Get")
	_, err = extendedZonesClient.Get(testsuite.ctx, testsuite.extendedZoneName, nil)
	testsuite.Require().NoError(err)

	// From step ExtendedZones_Unregister
	fmt.Println("Call operation: ExtendedZones_Unregister")
	_, err = extendedZonesClient.Unregister(testsuite.ctx, testsuite.extendedZoneName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EdgeZones/operations
func (testsuite *EdgezonesTestSuite) TestOperations() {
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient := testsuite.clientFactory.NewOperationsClient()
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		result, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		_ = result
		break
	}
}

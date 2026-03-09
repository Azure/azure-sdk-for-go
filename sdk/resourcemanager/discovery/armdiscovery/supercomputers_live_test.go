// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type SupercomputersTestSuite struct {
	suite.Suite
	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	location           string
	resourceGroupName  string
	subscriptionId     string
	supercomputerName  string
}

func (testsuite *SupercomputersTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.supercomputerName = "test-sc-go01"
}

func (testsuite *SupercomputersTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestSupercomputersTestSuite(t *testing.T) {
	suite.Run(t, new(SupercomputersTestSuite))
}

// Test listing supercomputers by subscription
func (testsuite *SupercomputersTestSuite) TestSupercomputersListBySubscription() {
	fmt.Println("Call operation: Supercomputers_ListBySubscription")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewSupercomputersClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

// Test creating a supercomputer
func (testsuite *SupercomputersTestSuite) TestSupercomputersCreateOrUpdate() {
	fmt.Println("Call operation: Supercomputers_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	miID := "/subscriptions/" + testsuite.subscriptionId + "/resourcegroups/" + testsuite.resourceGroupName + "/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity"
	subnetID := "/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/virtualNetworks/newapiv/subnets/default"

	supercomputersClient := clientFactory.NewSupercomputersClient()
	poller, err := supercomputersClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.supercomputerName,
		armdiscovery.Supercomputer{
			Location: to.Ptr(testsuite.location),
			Properties: &armdiscovery.SupercomputerProperties{
				SubnetID: to.Ptr(subnetID),
				Identities: &armdiscovery.SupercomputerIdentities{
					ClusterIdentity: &armdiscovery.Identity{
						ID: to.Ptr(miID),
					},
					KubeletIdentity: &armdiscovery.Identity{
						ID: to.Ptr(miID),
					},
					WorkloadIdentities: map[string]*armdiscovery.UserAssignedIdentity{
						miID: {},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	sc, err := poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(sc.ID)
	fmt.Println("Created supercomputer:", *sc.Name)
}

// Test listing supercomputers by resource group
func (testsuite *SupercomputersTestSuite) TestSupercomputersListByResourceGroup() {
	fmt.Println("Call operation: Supercomputers_ListByResourceGroup")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewSupercomputersClient().NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break // Just verify first page
	}
}

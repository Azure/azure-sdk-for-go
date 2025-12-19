// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type VirtualMachineImageTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *VirtualMachineImageTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *VirtualMachineImageTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestVirtualMachineImageTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualMachineImageTestSuite))
}

// Microsoft.Compute/locations/publishers
func (testsuite *VirtualMachineImageTestSuite) TestVirtualMachineImages() {
	offer := "office-365"
	publisherName := "MicrosoftWindowsDesktop"
	skus := "win11-22h2-avd-m365"
	version := "22621.1105.230110"
	var err error
	// From step VirtualMachineImages_ListPublishers
	fmt.Println("Call operation: VirtualMachineImages_ListPublishers")
	virtualMachineImagesClient, err := armcompute.NewVirtualMachineImagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = virtualMachineImagesClient.ListPublishers(testsuite.ctx, testsuite.location, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachineImages_ListOffers
	fmt.Println("Call operation: VirtualMachineImages_ListOffers")
	_, err = virtualMachineImagesClient.ListOffers(testsuite.ctx, testsuite.location, publisherName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachineImages_ListSkus
	fmt.Println("Call operation: VirtualMachineImages_ListSKUs")
	_, err = virtualMachineImagesClient.ListSKUs(testsuite.ctx, testsuite.location, publisherName, offer, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachineImages_List
	fmt.Println("Call operation: VirtualMachineImages_List")
	_, err = virtualMachineImagesClient.List(testsuite.ctx, testsuite.location, publisherName, offer, skus, &armcompute.VirtualMachineImagesClientListOptions{Expand: nil,
		Top:     nil,
		Orderby: nil,
	})
	testsuite.Require().NoError(err)

	// From step VirtualMachineImages_Get
	fmt.Println("Call operation: VirtualMachineImages_Get")
	_, err = virtualMachineImagesClient.Get(testsuite.ctx, testsuite.location, publisherName, offer, skus, version, nil)
	testsuite.Require().NoError(err)
}

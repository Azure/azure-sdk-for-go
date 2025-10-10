//go:build go1.18
// +build go1.18

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type VirtualMachineExtensionImageTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *VirtualMachineExtensionImageTestSuite) SetupSuite() {
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

func (testsuite *VirtualMachineExtensionImageTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVirtualMachineExtensionImageTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualMachineExtensionImageTestSuite))
}

// Microsoft.Compute/locations/{location}/publishers/{publisherName}/artifacttypes/vmextension/types
func (testsuite *VirtualMachineExtensionImageTestSuite) TestVirtualMachineExtensionImages() {
	publisherName := "Microsoft.Compute"
	typeParam := "CustomScriptExtension"
	version := "1.9"
	var err error
	// From step VirtualMachineExtensionImages_ListTypes
	fmt.Println("Call operation: VirtualMachineExtensionImages_ListTypes")
	virtualMachineExtensionImagesClient, err := armcompute.NewVirtualMachineExtensionImagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = virtualMachineExtensionImagesClient.ListTypes(testsuite.ctx, testsuite.location, publisherName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachineExtensionImages_ListVersions
	fmt.Println("Call operation: VirtualMachineExtensionImages_ListVersions")
	_, err = virtualMachineExtensionImagesClient.ListVersions(testsuite.ctx, testsuite.location, publisherName, typeParam, &armcompute.VirtualMachineExtensionImagesClientListVersionsOptions{Filter: nil,
		Top:     nil,
		Orderby: nil,
	})
	testsuite.Require().NoError(err)

	// From step VirtualMachineExtensionImages_Get
	fmt.Println("Call operation: VirtualMachineExtensionImages_Get")
	_, err = virtualMachineExtensionImagesClient.Get(testsuite.ctx, testsuite.location, publisherName, typeParam, version, nil)
	testsuite.Require().NoError(err)
}

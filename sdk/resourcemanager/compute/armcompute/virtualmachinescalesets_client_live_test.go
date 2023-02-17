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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/suite"
)

type VirtualMachineScaleSetsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *VirtualMachineScaleSetsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/compute/armcompute/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *VirtualMachineScaleSetsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVirtualMachineScaleSetsClient(t *testing.T) {
	suite.Run(t, new(VirtualMachineScaleSetsClientTestSuite))
}

func (testsuite *VirtualMachineScaleSetsClientTestSuite) TestVirtualMachineScaleSetsCRUD() {
	// create virtual network and subnet
	vnClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vnName := "go-test-network"
	subName := "go-test-subnet"
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		armnetwork.VirtualNetwork{
			Location: to.Ptr(testsuite.location),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.Ptr("10.1.0.0/16"),
					},
				},
				Subnets: []*armnetwork.Subnet{
					{
						Name: to.Ptr(subName),
						Properties: &armnetwork.SubnetPropertiesFormat{
							AddressPrefix: to.Ptr("10.1.0.0/24"),
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vnResp, err := testutil.PollForTest(testsuite.ctx, vnPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vnName, *vnResp.Name)

	// create virtual machine scale set
	fmt.Println("Call operation: VirtualMachineScaleSets_CreateOrUpdate")
	vmssClient, err := armcompute.NewVirtualMachineScaleSetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vmssName := "go-test-vmss"
	vmssPoller, err := vmssClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmssName,
		armcompute.VirtualMachineScaleSet{
			Location: to.Ptr(testsuite.location),
			SKU: &armcompute.SKU{
				//Name:     to.Ptr("Basic_A0"), //armcompute.VirtualMachineSizeTypesBasicA0
				Name:     to.Ptr("Standard_A0"), //armcompute.VirtualMachineSizeTypesBasicA0
				Capacity: to.Ptr[int64](1),
			},
			Properties: &armcompute.VirtualMachineScaleSetProperties{
				Overprovision: to.Ptr(false),
				UpgradePolicy: &armcompute.UpgradePolicy{
					Mode: to.Ptr(armcompute.UpgradeModeManual),
					AutomaticOSUpgradePolicy: &armcompute.AutomaticOSUpgradePolicy{
						EnableAutomaticOSUpgrade: to.Ptr(false),
						DisableAutomaticRollback: to.Ptr(false),
					},
				},
				VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
					OSProfile: &armcompute.VirtualMachineScaleSetOSProfile{
						ComputerNamePrefix: to.Ptr("vmss"),
						AdminUsername:      to.Ptr("sample-user"),
						AdminPassword:      to.Ptr("Password01!@#"),
					},
					StorageProfile: &armcompute.VirtualMachineScaleSetStorageProfile{
						ImageReference: &armcompute.ImageReference{
							Offer:     to.Ptr("WindowsServer"),
							Publisher: to.Ptr("MicrosoftWindowsServer"),
							SKU:       to.Ptr("2019-Datacenter"),
							Version:   to.Ptr("latest"),
						},
					},
					NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
						NetworkInterfaceConfigurations: []*armcompute.VirtualMachineScaleSetNetworkConfiguration{
							{
								Name: to.Ptr(vmssName),
								Properties: &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{
									Primary:            to.Ptr(true),
									EnableIPForwarding: to.Ptr(true),
									IPConfigurations: []*armcompute.VirtualMachineScaleSetIPConfiguration{
										{
											Name: to.Ptr(vmssName),
											Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
												Subnet: &armcompute.APIEntityReference{
													ID: to.Ptr(*vnResp.Properties.Subnets[0].ID),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vmssResp, err := testutil.PollForTest(testsuite.ctx, vmssPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vmssName, *vmssResp.Name)

	// update
	fmt.Println("Call operation: VirtualMachineScaleSets_Update")
	updatePollerResp, err := vmssClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmssName,
		armcompute.VirtualMachineScaleSetUpdate{
			Tags: map[string]*string{
				"test": to.Ptr("live"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	updateResp, err := testutil.PollForTest(testsuite.ctx, updatePollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("live", *updateResp.Tags["test"])

	// get
	fmt.Println("Call operation: VirtualMachineScaleSets_Get")
	getResp, err := vmssClient.Get(testsuite.ctx, testsuite.resourceGroupName, vmssName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vmssName, *getResp.Name)

	// list
	fmt.Println("Call operation: VirtualMachineScaleSets_List")
	listResp := vmssClient.NewListPager(testsuite.resourceGroupName, nil)
	testsuite.Require().True(listResp.More())

	// delete
	fmt.Println("Call operation: VirtualMachineScaleSets_Delete")
	delPoller, err := vmssClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vmssName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)
}

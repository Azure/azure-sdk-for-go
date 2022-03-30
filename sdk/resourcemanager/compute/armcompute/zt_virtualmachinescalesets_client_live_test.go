//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
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
	vnClient := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	vnName := "go-test-network"
	subName := "go-test-subnet"
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		armnetwork.VirtualNetwork{
			Location: to.StringPtr(testsuite.location),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
				Subnets: []*armnetwork.Subnet{
					{
						Name: to.StringPtr(subName),
						Properties: &armnetwork.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("10.1.0.0/24"),
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var vnResp armnetwork.VirtualNetworksClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vnPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if vnPoller.Poller.Done() {
				vnResp, err = vnPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		vnResp, err = vnPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(vnName, *vnResp.Name)

	// create virtual machine scale set
	vmssClient := armcompute.NewVirtualMachineScaleSetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	vmssName := "go-test-vmss"
	vmssPoller, err := vmssClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmssName,
		armcompute.VirtualMachineScaleSet{
			Location: to.StringPtr(testsuite.location),
			SKU: &armcompute.SKU{
				//Name:     to.StringPtr("Basic_A0"), //armcompute.VirtualMachineSizeTypesBasicA0
				Name:     to.StringPtr("Standard_A0"), //armcompute.VirtualMachineSizeTypesBasicA0
				Capacity: to.Int64Ptr(1),
			},
			Properties: &armcompute.VirtualMachineScaleSetProperties{
				Overprovision: to.BoolPtr(false),
				UpgradePolicy: &armcompute.UpgradePolicy{
					Mode: armcompute.UpgradeModeManual.ToPtr(),
					AutomaticOSUpgradePolicy: &armcompute.AutomaticOSUpgradePolicy{
						EnableAutomaticOSUpgrade: to.BoolPtr(false),
						DisableAutomaticRollback: to.BoolPtr(false),
					},
				},
				VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
					OSProfile: &armcompute.VirtualMachineScaleSetOSProfile{
						ComputerNamePrefix: to.StringPtr("vmss"),
						AdminUsername:      to.StringPtr("sample-user"),
						AdminPassword:      to.StringPtr("Password01!@#"),
					},
					StorageProfile: &armcompute.VirtualMachineScaleSetStorageProfile{
						ImageReference: &armcompute.ImageReference{
							Offer:     to.StringPtr("WindowsServer"),
							Publisher: to.StringPtr("MicrosoftWindowsServer"),
							SKU:       to.StringPtr("2019-Datacenter"),
							Version:   to.StringPtr("latest"),
						},
					},
					NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
						NetworkInterfaceConfigurations: []*armcompute.VirtualMachineScaleSetNetworkConfiguration{
							{
								Name: to.StringPtr(vmssName),
								Properties: &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{
									Primary:            to.BoolPtr(true),
									EnableIPForwarding: to.BoolPtr(true),
									IPConfigurations: []*armcompute.VirtualMachineScaleSetIPConfiguration{
										{
											Name: to.StringPtr(vmssName),
											Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
												Subnet: &armcompute.APIEntityReference{
													ID: to.StringPtr(*vnResp.Properties.Subnets[0].ID),
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
	var vmssResp armcompute.VirtualMachineScaleSetsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vmssPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if vmssPoller.Poller.Done() {
				vmssResp, err = vmssPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		vmssResp, err = vmssPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(vmssName, *vmssResp.Name)

	// update
	updatePollerResp, err := vmssClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmssName,
		armcompute.VirtualMachineScaleSetUpdate{
			Tags: map[string]*string{
				"test": to.StringPtr("live"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var updateResp armcompute.VirtualMachineScaleSetsClientUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = updatePollerResp.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if updatePollerResp.Poller.Done() {
				updateResp, err = updatePollerResp.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		updateResp, err = updatePollerResp.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal("live", *updateResp.Tags["test"])

	// get
	getResp, err := vmssClient.Get(testsuite.ctx, testsuite.resourceGroupName, vmssName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vmssName, *getResp.Name)

	// list
	listResp := vmssClient.List(testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(listResp.Err())

	// delete
	delPoller, err := vmssClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vmssName, nil)
	testsuite.Require().NoError(err)
	var delResp armcompute.VirtualMachineScaleSetsClientDeleteResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = delPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if delPoller.Poller.Done() {
				delResp, err = delPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		delResp, err = delPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}

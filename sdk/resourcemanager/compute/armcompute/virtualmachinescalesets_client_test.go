//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVirtualMachineScaleSetsClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "createVMSS", "westus2")
	rgName := *rg.Name
	defer clean()

	// create virtual network and subnet
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID, cred, opt)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	subName, err := createRandomName(t, "subnet")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("eastus"),
			},
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
	require.NoError(t, err)
	vnResp, err := vnPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vnResp.Name, vnName)

	// create virtual machine scale set
	vmssClient := armcompute.NewVirtualMachineScaleSetsClient(subscriptionID, cred, opt)
	vmssName, err := createRandomName(t, "vmss")
	require.NoError(t, err)
	vmssPoller, err := vmssClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vmssName,
		armcompute.VirtualMachineScaleSet{
			Resource: armcompute.Resource{
				Location: to.StringPtr("eastus"),
			},
			SKU: &armcompute.SKU{
				Name:     to.StringPtr("Basic_A0"), //armcompute.VirtualMachineSizeTypesBasicA0
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
	vmssResp, err := vmssPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vmssResp.Name, vmssName)
}

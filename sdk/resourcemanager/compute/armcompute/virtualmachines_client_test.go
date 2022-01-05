//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVirtualMachinesClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "createVM", "westus2")
	rgName := *rg.Name
	defer clean()

	//create virtual machine
	_, _ = createVirtualMachineTest(t, cred, opt, subscriptionID, rgName)
}

func TestVirtualMachinesClient_BeginDelete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "deleteVM", "westus")
	rgName := *rg.Name
	defer clean()

	// create virtual machine
	vmClient, vm := createVirtualMachineTest(t, cred, opt, subscriptionID, rgName)

	// delete virtual machine
	delPoller, err := vmClient.BeginDelete(context.Background(), rgName, *vm.Name, nil)
	require.NoError(t, err)
	delResp, err := delPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, delResp.RawResponse.StatusCode, 200)
}

func TestVirtualMachinesClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "getVM", "westus")
	rgName := *rg.Name
	defer clean()

	// create virtual machine
	vmClient, vm := createVirtualMachineTest(t, cred, opt, subscriptionID, rgName)

	// virtual machine get
	resp, err := vmClient.Get(context.Background(), rgName, *vm.Name, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.Name, *vm.Name)
}

func TestVirtualMachinesClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "listVM", "westus")
	rgName := *rg.Name
	defer clean()

	// create virtual machine
	vmClient, _ := createVirtualMachineTest(t, cred, opt, subscriptionID, rgName)

	// virtual machine list
	resp := vmClient.List(rgName, nil)
	require.Equal(t, resp.NextPage(context.Background()), true)
}

func TestVirtualMachinesClient_BeginUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "updateVM", "westus")
	rgName := *rg.Name
	defer clean()

	// create virtual machine
	vmClient, vm := createVirtualMachineTest(t, cred, opt, subscriptionID, rgName)

	// virtual machine update
	updatePoller, err := vmClient.BeginUpdate(
		context.Background(),
		rgName,
		*vm.Name,
		armcompute.VirtualMachineUpdate{
			UpdateResource: armcompute.UpdateResource{
				Tags: map[string]*string{
					"tag": to.StringPtr("value"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	updateResp, err := updatePoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *updateResp.Name, *vm.Name)
}

func createVirtualMachineTest(t *testing.T, cred azcore.TokenCredential, opt *arm.ClientOptions, subscriptionID, rgName string) (*armcompute.VirtualMachinesClient, *armcompute.VirtualMachine) {
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID, cred, opt)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus2"),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
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

	// create subnet
	subClient := armnetwork.NewSubnetsClient(subscriptionID, cred, opt)
	subName, err := createRandomName(t, "subnet")
	require.NoError(t, err)
	subPoller, err := subClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.1.10.0/24"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	subResp, err := subPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *subResp.Name, subName)

	// create public ip address
	ipClient := armnetwork.NewPublicIPAddressesClient(subscriptionID, cred, opt)
	ipName, err := createRandomName(t, "ip")
	require.NoError(t, err)
	ipPoller, err := ipClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		ipName,
		armnetwork.PublicIPAddress{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus2"),
			},
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(), // Static or Dynamic
			},
		},
		nil,
	)
	require.NoError(t, err)
	ipResp, err := ipPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *ipResp.Name, ipName)

	// create network security group
	nsgClient := armnetwork.NewNetworkSecurityGroupsClient(subscriptionID, cred, opt)
	nsgName, err := createRandomName(t, "nsg")
	require.NoError(t, err)
	nsgPoller, err := nsgClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		nsgName,
		armnetwork.NetworkSecurityGroup{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus2"),
			},
			Properties: &armnetwork.NetworkSecurityGroupPropertiesFormat{
				SecurityRules: []*armnetwork.SecurityRule{
					{
						Name: to.StringPtr("sample_inbound_22"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("*"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("22"),
							Protocol:                 armnetwork.SecurityRuleProtocolTCP.ToPtr(),
							Access:                   armnetwork.SecurityRuleAccessAllow.ToPtr(),
							Priority:                 to.Int32Ptr(100),
							Description:              to.StringPtr("sample network security group inbound port 22"),
							Direction:                armnetwork.SecurityRuleDirectionInbound.ToPtr(),
						},
					},
					// outbound
					{
						Name: to.StringPtr("sample_outbound_22"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("*"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("22"),
							Protocol:                 armnetwork.SecurityRuleProtocolTCP.ToPtr(),
							Access:                   armnetwork.SecurityRuleAccessAllow.ToPtr(),
							Priority:                 to.Int32Ptr(100),
							Description:              to.StringPtr("sample network security group outbound port 22"),
							Direction:                armnetwork.SecurityRuleDirectionOutbound.ToPtr(),
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	nsgResp, err := nsgPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *nsgResp.Name, nsgName)

	// create network interface
	nicClient := armnetwork.NewNetworkInterfacesClient(subscriptionID, cred, opt)
	nicName, err := createRandomName(t, "nic")
	require.NoError(t, err)
	nicPoller, err := nicClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		nicName,
		armnetwork.NetworkInterface{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus2"),
			},
			Properties: &armnetwork.NetworkInterfacePropertiesFormat{
				//NetworkSecurityGroup:
				IPConfigurations: []*armnetwork.NetworkInterfaceIPConfiguration{
					{
						Name: to.StringPtr("ipConfig"),
						Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
							Subnet: &armnetwork.Subnet{
								SubResource: armnetwork.SubResource{
									ID: to.StringPtr(*subResp.ID),
								},
							},
							PublicIPAddress: &armnetwork.PublicIPAddress{
								Resource: armnetwork.Resource{
									ID: to.StringPtr(*ipResp.ID),
								},
							},
						},
					},
				},
				NetworkSecurityGroup: &armnetwork.NetworkSecurityGroup{
					Resource: armnetwork.Resource{
						ID: to.StringPtr(*nsgResp.ID),
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	nicResp, err := nicPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *nicResp.Name, nicName)

	// create virtual machine
	vmClient := armcompute.NewVirtualMachinesClient(subscriptionID, cred, opt)
	vmName, err := createRandomName(t, "vm")
	require.NoError(t, err)
	diskName, err := createRandomName(t, "disk")
	require.NoError(t, err)
	vmPoller, err := vmClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vmName,
		armcompute.VirtualMachine{
			Resource: armcompute.Resource{
				Location: to.StringPtr("westus2"),
			},
			Identity: &armcompute.VirtualMachineIdentity{
				Type: armcompute.ResourceIdentityTypeNone.ToPtr(),
			},
			Properties: &armcompute.VirtualMachineProperties{
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Offer:     to.StringPtr("WindowsServer"),
						Publisher: to.StringPtr("MicrosoftWindowsServer"),
						SKU:       to.StringPtr("2019-Datacenter"),
						Version:   to.StringPtr("latest"),
					},
					OSDisk: &armcompute.OSDisk{
						Name:         to.StringPtr(diskName),
						CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
						Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
						ManagedDisk: &armcompute.ManagedDiskParameters{
							StorageAccountType: armcompute.StorageAccountTypesStandardLRS.ToPtr(), // OSDisk type Standard/Premium HDD/SSD
						},
					},
				},
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: armcompute.VirtualMachineSizeTypesStandardF2S.ToPtr(), // VM size include vCPUs,RAM,Data Disks,Temp storage.
				},
				OSProfile: &armcompute.OSProfile{
					ComputerName:  to.StringPtr("sample-compute"),
					AdminUsername: to.StringPtr("sample-user"),
					AdminPassword: to.StringPtr("Password01!@#"),
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							SubResource: armcompute.SubResource{
								ID: to.StringPtr(*nicResp.ID),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	vmResp, err := vmPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vmResp.Name, vmName)
	return vmClient, &vmResp.VirtualMachine
}

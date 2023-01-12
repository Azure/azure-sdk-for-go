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

type VirtualMachinesClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *VirtualMachinesClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/compute/armcompute/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *VirtualMachinesClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVirtualMachinesClient(t *testing.T) {
	suite.Run(t, new(VirtualMachinesClientTestSuite))
}

func (testsuite *VirtualMachinesClientTestSuite) TestVirtualMachineCRUD() {
	// create virtual network
	vnClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vnName := "go-test-network"
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
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vnResp, err := testutil.PollForTest(testsuite.ctx, vnPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*vnResp.Name, vnName)

	// create subnet
	subClient, err := armnetwork.NewSubnetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	subName := "go-test-subnet"
	subPoller, err := subClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.Ptr("10.1.10.0/24"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	subResp, err := testutil.PollForTest(testsuite.ctx, subPoller)
	testsuite.Require().NoError(err)
	subnetID := *subResp.ID
	testsuite.Require().Equal(*subResp.Name, subName)

	// create public ip address
	ipClient, err := armnetwork.NewPublicIPAddressesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	ipName := "go-test-ip"
	testsuite.Require().NoError(err)
	ipPoller, err := ipClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		ipName,
		armnetwork.PublicIPAddress{
			Location: to.Ptr(testsuite.location),
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic), // Static or Dynamic
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	ipResp, err := testutil.PollForTest(testsuite.ctx, ipPoller)
	testsuite.Require().NoError(err)
	publicIPAddressID := *ipResp.ID
	testsuite.Require().Equal(*ipResp.Name, ipName)

	// create network security group
	nsgClient, err := armnetwork.NewSecurityGroupsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	nsgName := "go-test-nsg"
	testsuite.Require().NoError(err)
	nsgPoller, err := nsgClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		nsgName,
		armnetwork.SecurityGroup{
			Location: to.Ptr(testsuite.location),
			Properties: &armnetwork.SecurityGroupPropertiesFormat{
				SecurityRules: []*armnetwork.SecurityRule{
					{
						Name: to.Ptr("sample_inbound_22"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
							SourcePortRange:          to.Ptr("*"),
							DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
							DestinationPortRange:     to.Ptr("22"),
							Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
							Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
							Priority:                 to.Ptr[int32](100),
							Description:              to.Ptr("sample network security group inbound port 22"),
							Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
						},
					},
					// outbound
					{
						Name: to.Ptr("sample_outbound_22"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
							SourcePortRange:          to.Ptr("*"),
							DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
							DestinationPortRange:     to.Ptr("22"),
							Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
							Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
							Priority:                 to.Ptr[int32](100),
							Description:              to.Ptr("sample network security group outbound port 22"),
							Direction:                to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	nsgResp, err := testutil.PollForTest(testsuite.ctx, nsgPoller)
	testsuite.Require().NoError(err)
	networkSecurityGroupID := *nsgResp.ID
	testsuite.Require().Equal(*nsgResp.Name, nsgName)

	// create network interface
	nicClient, err := armnetwork.NewInterfacesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	nicName := "go-test-nic"
	testsuite.Require().NoError(err)
	nicPoller, err := nicClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		nicName,
		armnetwork.Interface{
			Location: to.Ptr(testsuite.location),
			Properties: &armnetwork.InterfacePropertiesFormat{
				//NetworkSecurityGroup:
				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
					{
						Name: to.Ptr("ipConfig"),
						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
							Subnet: &armnetwork.Subnet{
								ID: to.Ptr(subnetID),
							},
							PublicIPAddress: &armnetwork.PublicIPAddress{
								ID: to.Ptr(publicIPAddressID),
							},
						},
					},
				},
				NetworkSecurityGroup: &armnetwork.SecurityGroup{
					ID: to.Ptr(networkSecurityGroupID),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	nicResp, err := testutil.PollForTest(testsuite.ctx, nicPoller)
	testsuite.Require().NoError(err)
	networkInterfaceID := *nicResp.ID
	testsuite.Require().Equal(*nicResp.Name, nicName)

	// create virtual machine
	fmt.Println("Call operation: VirtualMachines_CreateOrUpdate")
	vmClient, err := armcompute.NewVirtualMachinesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vmName := "go-test-vm"
	testsuite.Require().NoError(err)
	diskName := "go-test-disk"
	testsuite.Require().NoError(err)
	vmPoller, err := vmClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmName,
		armcompute.VirtualMachine{
			Location: to.Ptr(testsuite.location),
			Identity: &armcompute.VirtualMachineIdentity{
				Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
			},
			Properties: &armcompute.VirtualMachineProperties{
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Offer:     to.Ptr("WindowsServer"),
						Publisher: to.Ptr("MicrosoftWindowsServer"),
						SKU:       to.Ptr("2019-Datacenter"),
						Version:   to.Ptr("latest"),
					},
					OSDisk: &armcompute.OSDisk{
						Name:         to.Ptr(diskName),
						CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
						Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
						ManagedDisk: &armcompute.ManagedDiskParameters{
							StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
						},
					},
				},
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: to.Ptr(armcompute.VirtualMachineSizeTypesStandardF2S), // VM size include vCPUs,RAM,Data Disks,Temp storage.
				},
				OSProfile: &armcompute.OSProfile{
					ComputerName:  to.Ptr("sample-compute"),
					AdminUsername: to.Ptr("sample-user"),
					AdminPassword: to.Ptr("Password01!@#"),
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							ID: to.Ptr(networkInterfaceID),
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vmResp, err := testutil.PollForTest(testsuite.ctx, vmPoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*vmResp.Name, vmName)

	// virtual machine update
	fmt.Println("Call operation: VirtualMachines_Update")
	updatePoller, err := vmClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmName,
		armcompute.VirtualMachineUpdate{
			Tags: map[string]*string{
				"tag": to.Ptr("value"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	updateResp, err := testutil.PollForTest(testsuite.ctx, updatePoller)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*updateResp.Name, vmName)

	// virtual machine get
	fmt.Println("Call operation: VirtualMachines_Get")
	resp, err := vmClient.Get(testsuite.ctx, testsuite.resourceGroupName, vmName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*resp.Name, vmName)

	// virtual machine list
	fmt.Println("Call operation: VirtualMachines_List")
	vmList := vmClient.NewListPager(testsuite.resourceGroupName, nil)
	testsuite.Require().Equal(vmList.More(), true)

	// delete virtual machine
	fmt.Println("Call operation: VirtualMachines_Delete")
	delPoller, err := vmClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vmName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, delPoller)
	testsuite.Require().NoError(err)
}

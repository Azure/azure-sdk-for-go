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
	vnClient := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	vnName := "go-test-network"
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
	testsuite.Require().Equal(*vnResp.Name, vnName)

	// create subnet
	subClient := armnetwork.NewSubnetsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	subName := "go-test-subnet"
	subPoller, err := subClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.1.10.0/24"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var subResp armnetwork.SubnetsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = subPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if subPoller.Poller.Done() {
				subResp, err = subPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		subResp, err = subPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	subnetID := *subResp.ID
	testsuite.Require().Equal(*subResp.Name, subName)

	// create public ip address
	ipClient := armnetwork.NewPublicIPAddressesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	ipName := "go-test-ip"
	testsuite.Require().NoError(err)
	ipPoller, err := ipClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		ipName,
		armnetwork.PublicIPAddress{
			Location: to.StringPtr(testsuite.location),
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(), // Static or Dynamic
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var ipResp armnetwork.PublicIPAddressesClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = ipPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if ipPoller.Poller.Done() {
				ipResp, err = ipPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		ipResp, err = ipPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	publicIPAddressID := *ipResp.ID
	testsuite.Require().Equal(*ipResp.Name, ipName)

	// create network security group
	nsgClient := armnetwork.NewSecurityGroupsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	nsgName := "go-test-nsg"
	testsuite.Require().NoError(err)
	nsgPoller, err := nsgClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		nsgName,
		armnetwork.SecurityGroup{
			Location: to.StringPtr(testsuite.location),
			Properties: &armnetwork.SecurityGroupPropertiesFormat{
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
	testsuite.Require().NoError(err)
	var nsgResp armnetwork.SecurityGroupsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = nsgPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if nsgPoller.Poller.Done() {
				nsgResp, err = nsgPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		nsgResp, err = nsgPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	networkSecurityGroupID := *nsgResp.ID
	testsuite.Require().Equal(*nsgResp.Name, nsgName)

	// create network interface
	nicClient := armnetwork.NewInterfacesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	nicName := "go-test-nic"
	testsuite.Require().NoError(err)
	nicPoller, err := nicClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		nicName,
		armnetwork.Interface{
			Location: to.StringPtr(testsuite.location),
			Properties: &armnetwork.InterfacePropertiesFormat{
				//NetworkSecurityGroup:
				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
					{
						Name: to.StringPtr("ipConfig"),
						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
							Subnet: &armnetwork.Subnet{
								ID: to.StringPtr(subnetID),
							},
							PublicIPAddress: &armnetwork.PublicIPAddress{
								ID: to.StringPtr(publicIPAddressID),
							},
						},
					},
				},
				NetworkSecurityGroup: &armnetwork.SecurityGroup{
					ID: to.StringPtr(networkSecurityGroupID),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var nicResp armnetwork.InterfacesClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = nicPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if nicPoller.Poller.Done() {
				nicResp, err = nicPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		nicResp, err = nicPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	networkInterfaceID := *nicResp.ID
	testsuite.Require().Equal(*nicResp.Name, nicName)

	// create virtual machine
	vmClient := armcompute.NewVirtualMachinesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	vmName := "go-test-vm"
	testsuite.Require().NoError(err)
	diskName := "go-test-disk"
	testsuite.Require().NoError(err)
	vmPoller, err := vmClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmName,
		armcompute.VirtualMachine{
			Location: to.StringPtr(testsuite.location),
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
							ID: to.StringPtr(networkInterfaceID),
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var vmResp armcompute.VirtualMachinesClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vmPoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if vmPoller.Poller.Done() {
				vmResp, err = vmPoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		vmResp, err = vmPoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(*vmResp.Name, vmName)

	// virtual machine update
	updatePoller, err := vmClient.BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vmName,
		armcompute.VirtualMachineUpdate{
			Tags: map[string]*string{
				"tag": to.StringPtr("value"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	var updateResp armcompute.VirtualMachinesClientUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = updatePoller.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if updatePoller.Poller.Done() {
				updateResp, err = updatePoller.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		updateResp, err = updatePoller.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(*updateResp.Name, vmName)

	// virtual machine get
	resp, err := vmClient.Get(testsuite.ctx, testsuite.resourceGroupName, vmName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(*resp.Name, vmName)

	// virtual machine list
	vmList := vmClient.List(testsuite.resourceGroupName, nil)
	testsuite.Require().Equal(vmList.NextPage(testsuite.ctx), true)

	// delete virtual machine
	delPoller, err := vmClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, vmName, nil)
	testsuite.Require().NoError(err)
	var delResp armcompute.VirtualMachinesClientDeleteResponse
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
	testsuite.Require().Equal(delResp.RawResponse.StatusCode, 200)
}

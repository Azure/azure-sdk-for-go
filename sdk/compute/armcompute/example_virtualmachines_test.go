//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
)

func ExampleVirtualMachinesClient_BeginCreateOrUpdate() {
	// replace with your own value
	vmName := "<VM name>"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		vmName,
		armcompute.VirtualMachine{
			Resource: armcompute.Resource{
				Name:     to.StringPtr(vmName),
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armcompute.VirtualMachineProperties{
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: armcompute.VirtualMachineSizeTypesBasicA0.ToPtr(),
				},
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Publisher: to.StringPtr("<publisher>"),
						Offer:     to.StringPtr("<offer>"),
						SKU:       to.StringPtr("<sku>"),
						Version:   to.StringPtr("latest"),
					},
				},
				OSProfile: &armcompute.OSProfile{
					ComputerName:  to.StringPtr(vmName),
					AdminUsername: to.StringPtr("<username>"),
					AdminPassword: to.StringPtr("<password>"),
					LinuxConfiguration: &armcompute.LinuxConfiguration{
						SSH: &armcompute.SSHConfiguration{
							PublicKeys: []*armcompute.SSHPublicKey{
								{
									Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", "<username>")),
									KeyData: to.StringPtr("<SSH key data"),
								},
							},
						},
					},
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							SubResource: armcompute.SubResource{
								// call armnetwork.NetworkInterfacesClient.Get to retreive an existing NIC and see the ID
								ID: to.StringPtr("<NIC ID>"),
							},
							Properties: &armcompute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM ID: %v", *resp.VirtualMachine.ID)
}

func ExampleVirtualMachinesClient_BeginCreateOrUpdate_withDisk() {
	vmName := "<VM name>"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		vmName,
		armcompute.VirtualMachine{
			Resource: armcompute.Resource{
				Name:     to.StringPtr(vmName),
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armcompute.VirtualMachineProperties{
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: armcompute.VirtualMachineSizeTypesBasicA0.ToPtr(),
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							SubResource: armcompute.SubResource{
								ID: to.StringPtr("<NIC ID>"),
							},
							Properties: &armcompute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
				OSProfile: &armcompute.OSProfile{
					ComputerName:  to.StringPtr(vmName),
					AdminUsername: to.StringPtr("<username>"),
					AdminPassword: to.StringPtr("<password>"),
					LinuxConfiguration: &armcompute.LinuxConfiguration{
						DisablePasswordAuthentication: to.BoolPtr(false),
					},
				},
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Publisher: to.StringPtr("<publisher>"),
						Offer:     to.StringPtr("<offer>"),
						SKU:       to.StringPtr("<sku>"),
						Version:   to.StringPtr("latest"),
					},
					OSDisk: &armcompute.OSDisk{
						CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
						DiskSizeGB:   to.Int32Ptr(64),
					},
					DataDisks: []*armcompute.DataDisk{
						{
							CreateOption: armcompute.DiskCreateOptionTypesAttach.ToPtr(),
							Lun:          to.Int32Ptr(0),
							ManagedDisk: &armcompute.ManagedDiskParameters{
								SubResource: armcompute.SubResource{
									// call armcompute.DisksClient.Get to retreive an existing disk and see the ID
									ID: to.StringPtr("<disk ID>"),
								},
							},
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM ID: %v", *resp.VirtualMachine.ID)
}

func ExampleVirtualMachinesClient_BeginCreateOrUpdate_withLoadBalancer() {
	vmName := "<VM name>"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		vmName,
		armcompute.VirtualMachine{
			Resource: armcompute.Resource{
				Name:     to.StringPtr(vmName),
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armcompute.VirtualMachineProperties{
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: armcompute.VirtualMachineSizeTypesStandardA0.ToPtr(),
				},
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Publisher: to.StringPtr("<publisher>"),
						Offer:     to.StringPtr("<offer>"),
						SKU:       to.StringPtr("<sku>"),
						Version:   to.StringPtr("latest"),
					},
				},
				OSProfile: &armcompute.OSProfile{
					ComputerName:  to.StringPtr(vmName),
					AdminUsername: to.StringPtr("azureuser"),
					AdminPassword: to.StringPtr("password!1delete"),
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							SubResource: armcompute.SubResource{
								// get the NIC ID by calling armnetwork.NetworkInterfacesClient.Get and retreiving the ID from the desired NIC instance
								ID: to.StringPtr("<NIC ID>"),
							},
							Properties: &armcompute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
				AvailabilitySet: &armcompute.SubResource{
					// get the availability set ID by calling armcompute.AvailabilitySetsClient.Get and retreiving the ID from the desired availability set instance
					ID: to.StringPtr("<availability set ID>"),
				},
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM ID: %v", *resp.VirtualMachine.ID)
}

func ExampleVirtualMachinesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<VM name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM ID: %s", *resp.VirtualMachine.ID)
}

func ExampleVirtualMachinesClient_BeginDeallocate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginDeallocate(context.Background(), "<resource group name>", "<VM name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to deallocate the vm: %v", err)
	}
}

func ExampleVirtualMachinesClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginUpdate(
		context.Background(),
		"<resource group name>",
		"<VM name>",
		armcompute.VirtualMachineUpdate{
			UpdateResource: armcompute.UpdateResource{
				Tags: map[string]*string{
					"who rocks": to.StringPtr("golang"),
					"where":     to.StringPtr("on azure"),
				},
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to deallocate the vm: %v", err)
	}
	log.Printf("ID of the updated vm: %s", *resp.VirtualMachine.ID)
}

func ExampleVirtualMachinesClient_BeginStart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginStart(
		context.Background(),
		"<resource group name>",
		"<VM name>",
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to start the vm: %v", err)
	}
}

func ExampleVirtualMachinesClient_BeginRestart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginRestart(
		context.Background(),
		"<resource group name>",
		"<VM name>",
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to restart the vm: %v", err)
	}
}

func ExampleVirtualMachinesClient_BeginPowerOff() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginPowerOff(
		context.Background(),
		"<resource group name>",
		"<VM name>",
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to power off the vm: %v", err)
	}
}

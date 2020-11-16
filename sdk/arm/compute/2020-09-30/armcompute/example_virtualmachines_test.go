// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func ExampleVirtualMachinesOperations_BeginCreateOrUpdate() {
	// replace with your own value
	vmName := "<VM name>"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
							PublicKeys: &[]armcompute.SSHPublicKey{
								{
									Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", "<username>")),
									KeyData: to.StringPtr("<SSH key data"),
								},
							},
						},
					},
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: &[]armcompute.NetworkInterfaceReference{
						{
							SubResource: armcompute.SubResource{
								ID: to.StringPtr("<NIC ID>"), // call armnetwork.NetworkInterfacesClient.Get to retreive an existing NIC and see the ID
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

func ExampleVirtualMachinesOperations_BeginCreateOrUpdate_withDisk() {
	vmName := "<VM name>"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
					NetworkInterfaces: &[]armcompute.NetworkInterfaceReference{
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
						DiskSizeGb:   to.Int32Ptr(64),
					},
					DataDisks: &[]armcompute.DataDisk{
						{
							CreateOption: armcompute.DiskCreateOptionTypesAttach.ToPtr(),
							Lun:          to.Int32Ptr(0),
							ManagedDisk: &armcompute.ManagedDiskParameters{
								SubResource: armcompute.SubResource{
									ID: to.StringPtr("<disk ID>"), // call armcompute.DisksClient.Get to retreive an existing disk and see the ID
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

func ExampleVirtualMachinesOperations_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<VM name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM ID: %s", *resp.VirtualMachine.ID)
}

func ExampleVirtualMachinesOperations_BeginDeallocate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginDeallocate(context.Background(), "<resource group name>", "<VM name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to deallocate the vm: %v", err)
	}
}

func ExampleVirtualMachinesOperations_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginUpdate(
		context.Background(),
		"<resource group name>",
		"<VM name>",
		armcompute.VirtualMachineUpdate{
			UpdateResource: armcompute.UpdateResource{
				Tags: &map[string]string{
					"who rocks": "golang",
					"where":     "on azure",
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

func ExampleVirtualMachinesOperations_BeginStart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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

func ExampleVirtualMachinesOperations_BeginRestart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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

func ExampleVirtualMachinesOperations_BeginPowerOff() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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

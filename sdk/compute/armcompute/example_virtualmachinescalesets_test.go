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

func ExampleVirtualMachineScaleSetsClient_BeginCreateOrUpdate() {
	vmssName := "<VM scale set name>"
	username := "<username>"
	password := "<password>"
	sshKeyData := "<SSH key data>"
	offer := "<offer>"
	publisher := "<publisher>"
	sku := "<sku>"
	// Retreive the subnet ID of an existing subnet by calling armnetwork.SubnetsClient.Get to retreive the
	// the subnet instance that will be assigned to the virtual machine scale set.
	subnetID := "<subnet ID>"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<VM scale set name>",
		armcompute.VirtualMachineScaleSet{
			Resource: armcompute.Resource{
				Name:     to.StringPtr("<VM scale set name>"),
				Location: to.StringPtr("<Azure location>"),
			},
			SKU: &armcompute.SKU{
				Name:     to.StringPtr(string(armcompute.VirtualMachineSizeTypesBasicA0)),
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
						ComputerNamePrefix: to.StringPtr(vmssName),
						AdminUsername:      to.StringPtr(username),
						AdminPassword:      to.StringPtr(password),
						LinuxConfiguration: &armcompute.LinuxConfiguration{
							SSH: &armcompute.SSHConfiguration{
								PublicKeys: []*armcompute.SSHPublicKey{
									{
										Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)),
										KeyData: to.StringPtr(sshKeyData),
									},
								},
							},
						},
					},
					StorageProfile: &armcompute.VirtualMachineScaleSetStorageProfile{
						ImageReference: &armcompute.ImageReference{
							Offer:     to.StringPtr(offer),
							Publisher: to.StringPtr(publisher),
							SKU:       to.StringPtr(sku),
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
													ID: to.StringPtr(subnetID),
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
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM scale set ID: %v", *resp.VirtualMachineScaleSet.ID)
}

func ExampleVirtualMachineScaleSetsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<VM scale set name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM scale set ID: %s", *resp.VirtualMachineScaleSet.ID)
}

func ExampleVirtualMachineScaleSetsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginUpdate(
		context.Background(),
		"<resource group name>",
		"<VM scale set name>",
		armcompute.VirtualMachineScaleSetUpdate{
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
		log.Fatalf("failed to update VM scale set: %v", err)
	}
	log.Printf("VM scale set ID: %v", *resp.VirtualMachineScaleSet.ID)
}

func ExampleVirtualMachineScaleSetsClient_BeginDeallocate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginDeallocate(context.Background(), "<resource group name>", "<VM scale set name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to deallocate VM scale set: %v", err)
	}
}

func ExampleVirtualMachineScaleSetsClient_BeginStart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginStart(context.Background(), "<resource group name>", "<VM scale set name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to start VM scale set: %v", err)
	}
}

func ExampleVirtualMachineScaleSetsClient_BeginRestart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginRestart(context.Background(), "<resource group name>", "<VM scale set name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to restart VM scale set: %v", err)
	}
}

func ExampleVirtualMachineScaleSetsClient_BeginPowerOff() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginPowerOff(context.Background(), "<resource group name>", "<VM scale set name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to stop VM scale set: %v", err)
	}
}

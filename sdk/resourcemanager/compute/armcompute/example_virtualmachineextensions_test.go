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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

func ExampleVirtualMachineExtensionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineExtensionsClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<VM name>",
		"<VM extension name>",
		armcompute.VirtualMachineExtension{
			Resource: armcompute.Resource{
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armcompute.VirtualMachineExtensionProperties{
				AutoUpgradeMinorVersion: to.BoolPtr(true),
				ProtectedSettings: map[string]interface{}{
					"AADClientSecret": "<client secret>",
					"Passphrase":      "yourPassPhrase",
				},
				Publisher: to.StringPtr("Microsoft.Azure.Security"),
				Settings: map[string]interface{}{
					"AADClientID":               "<client ID>",
					"EncryptionOperation":       "EnableEncryption",
					"KeyEncryptionAlgorithm":    "RSA-OAEP",
					"KeyEncryptionKeyAlgorithm": "<key ID>",
					"KeyVaultURL":               fmt.Sprintf("https://%s.%s/", "<vault name>", "<keyvault DNS suffix>"),
					"SequenceVersion":           "<UUID string>",
					"VolumeType":                "ALL",
				},
				Type:               to.StringPtr("AzureDiskEncryptionForLinux"),
				TypeHandlerVersion: to.StringPtr("0.1"),
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
	log.Printf("VM extension ID: %v", *resp.VirtualMachineExtension.ID)
}

func ExampleVirtualMachineExtensionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineExtensionsClient("<subscription ID>", cred, nil)
	resp, err := client.Get(context.Background(), "<resource group name>", "<VM name>", "<VM extension name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("VM extension ID: %s", *resp.VirtualMachineExtension.ID)
}

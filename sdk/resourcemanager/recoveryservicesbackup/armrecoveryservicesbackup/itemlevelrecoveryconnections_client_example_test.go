//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armrecoveryservicesbackup_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3751704f5318f1175875c94b66af769db917f2d3/specification/recoveryservicesbackup/resource-manager/Microsoft.RecoveryServices/stable/2023-01-01/examples/AzureIaasVm/Provision_Ilr.json
func ExampleItemLevelRecoveryConnectionsClient_Provision() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicesbackup.NewItemLevelRecoveryConnectionsClient("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Provision(ctx, "PySDKBackupTestRsVault", "PythonSDKBackupTestRg", "Azure", "iaasvmcontainer;iaasvmcontainerv2;pysdktestrg;pysdktestv2vm1", "vm;iaasvmcontainerv2;pysdktestrg;pysdktestv2vm1", "1", armrecoveryservicesbackup.ILRRequestResource{
		Properties: &armrecoveryservicesbackup.IaasVMILRRegistrationRequest{
			ObjectType:                to.Ptr("IaasVMILRRegistrationRequest"),
			InitiatorName:             to.Ptr("Hello World"),
			RecoveryPointID:           to.Ptr("38823086363464"),
			RenewExistingRegistration: to.Ptr(true),
			VirtualMachineID:          to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/pysdktestrg/providers/Microsoft.Compute/virtualMachines/pysdktestv2vm1"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3751704f5318f1175875c94b66af769db917f2d3/specification/recoveryservicesbackup/resource-manager/Microsoft.RecoveryServices/stable/2023-01-01/examples/AzureIaasVm/Revoke_Ilr.json
func ExampleItemLevelRecoveryConnectionsClient_Revoke() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicesbackup.NewItemLevelRecoveryConnectionsClient("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Revoke(ctx, "PySDKBackupTestRsVault", "PythonSDKBackupTestRg", "Azure", "iaasvmcontainer;iaasvmcontainerv2;pysdktestrg;pysdktestv2vm1", "vm;iaasvmcontainerv2;pysdktestrg;pysdktestv2vm1", "1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

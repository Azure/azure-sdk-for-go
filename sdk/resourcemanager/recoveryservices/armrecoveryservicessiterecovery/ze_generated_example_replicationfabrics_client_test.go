//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservicessiterecovery_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_List.json
func ExampleReplicationFabricsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
			return
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_Get.json
func ExampleReplicationFabricsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"<fabric-name>",
		&armrecoveryservicessiterecovery.ReplicationFabricsClientGetOptions{Filter: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_Create.json
func ExampleReplicationFabricsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreate(ctx,
		"<fabric-name>",
		armrecoveryservicessiterecovery.FabricCreationInput{
			Properties: &armrecoveryservicessiterecovery.FabricCreationInputProperties{
				CustomDetails: &armrecoveryservicessiterecovery.FabricSpecificCreationInput{
					InstanceType: to.Ptr("<instance-type>"),
				},
			},
		},
		&armrecoveryservicessiterecovery.ReplicationFabricsClientBeginCreateOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_Purge.json
func ExampleReplicationFabricsClient_BeginPurge() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginPurge(ctx,
		"<fabric-name>",
		&armrecoveryservicessiterecovery.ReplicationFabricsClientBeginPurgeOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_CheckConsistency.json
func ExampleReplicationFabricsClient_BeginCheckConsistency() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCheckConsistency(ctx,
		"<fabric-name>",
		&armrecoveryservicessiterecovery.ReplicationFabricsClientBeginCheckConsistencyOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_ReassociateGateway.json
func ExampleReplicationFabricsClient_BeginReassociateGateway() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginReassociateGateway(ctx,
		"<fabric-name>",
		armrecoveryservicessiterecovery.FailoverProcessServerRequest{
			Properties: &armrecoveryservicessiterecovery.FailoverProcessServerRequestProperties{
				ContainerName:         to.Ptr("<container-name>"),
				SourceProcessServerID: to.Ptr("<source-process-server-id>"),
				TargetProcessServerID: to.Ptr("<target-process-server-id>"),
				UpdateType:            to.Ptr("<update-type>"),
				VMsToMigrate: []*string{
					to.Ptr("Vm1"),
					to.Ptr("Vm2")},
			},
		},
		&armrecoveryservicessiterecovery.ReplicationFabricsClientBeginReassociateGatewayOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_Delete.json
func ExampleReplicationFabricsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"<fabric-name>",
		&armrecoveryservicessiterecovery.ReplicationFabricsClientBeginDeleteOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/recoveryservicessiterecovery/resource-manager/Microsoft.RecoveryServices/stable/2022-02-01/examples/ReplicationFabrics_RenewCertificate.json
func ExampleReplicationFabricsClient_BeginRenewCertificate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicessiterecovery.NewReplicationFabricsClient("<resource-name>",
		"<resource-group-name>",
		"<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginRenewCertificate(ctx,
		"<fabric-name>",
		armrecoveryservicessiterecovery.RenewCertificateInput{
			Properties: &armrecoveryservicessiterecovery.RenewCertificateInputProperties{
				RenewCertificateType: to.Ptr("<renew-certificate-type>"),
			},
		},
		&armrecoveryservicessiterecovery.ReplicationFabricsClientBeginRenewCertificateOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

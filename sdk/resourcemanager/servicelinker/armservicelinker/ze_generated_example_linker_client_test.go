//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armservicelinker_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicelinker/armservicelinker"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/servicelinker/resource-manager/Microsoft.ServiceLinker/stable/2022-05-01/examples/LinkList.json
func ExampleLinkerClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armservicelinker.NewLinkerClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	pager := client.NewListPager("<resource-uri>",
		nil)
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/servicelinker/resource-manager/Microsoft.ServiceLinker/stable/2022-05-01/examples/Link.json
func ExampleLinkerClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armservicelinker.NewLinkerClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.Get(ctx,
		"<resource-uri>",
		"<linker-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/servicelinker/resource-manager/Microsoft.ServiceLinker/stable/2022-05-01/examples/PutLink.json
func ExampleLinkerClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armservicelinker.NewLinkerClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	poller, err := client.BeginCreateOrUpdate(ctx,
		"<resource-uri>",
		"<linker-name>",
		armservicelinker.LinkerResource{
			Properties: &armservicelinker.LinkerProperties{
				AuthInfo: &armservicelinker.SecretAuthInfo{
					AuthType: to.Ptr(armservicelinker.AuthTypeSecret),
					Name:     to.Ptr("<name>"),
					SecretInfo: &armservicelinker.ValueSecretInfo{
						SecretType: to.Ptr(armservicelinker.SecretTypeRawValue),
						Value:      to.Ptr("<value>"),
					},
				},
				TargetService: &armservicelinker.AzureResource{
					Type: to.Ptr(armservicelinker.TypeAzureResource),
					ID:   to.Ptr("<id>"),
				},
			},
		},
		&armservicelinker.LinkerClientBeginCreateOrUpdateOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/servicelinker/resource-manager/Microsoft.ServiceLinker/stable/2022-05-01/examples/DeleteLink.json
func ExampleLinkerClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armservicelinker.NewLinkerClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	poller, err := client.BeginDelete(ctx,
		"<resource-uri>",
		"<linker-name>",
		&armservicelinker.LinkerClientBeginDeleteOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
		return
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/servicelinker/resource-manager/Microsoft.ServiceLinker/stable/2022-05-01/examples/PatchLink.json
func ExampleLinkerClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armservicelinker.NewLinkerClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	poller, err := client.BeginUpdate(ctx,
		"<resource-uri>",
		"<linker-name>",
		armservicelinker.LinkerPatch{
			Properties: &armservicelinker.LinkerProperties{
				AuthInfo: &armservicelinker.ServicePrincipalSecretAuthInfo{
					AuthType:    to.Ptr(armservicelinker.AuthTypeServicePrincipalSecret),
					ClientID:    to.Ptr("<client-id>"),
					PrincipalID: to.Ptr("<principal-id>"),
					Secret:      to.Ptr("<secret>"),
				},
				TargetService: &armservicelinker.AzureResource{
					Type: to.Ptr(armservicelinker.TypeAzureResource),
					ID:   to.Ptr("<id>"),
				},
			},
		},
		&armservicelinker.LinkerClientBeginUpdateOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/servicelinker/resource-manager/Microsoft.ServiceLinker/stable/2022-05-01/examples/ValidateLinkSuccess.json
func ExampleLinkerClient_BeginValidate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armservicelinker.NewLinkerClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	poller, err := client.BeginValidate(ctx,
		"<resource-uri>",
		"<linker-name>",
		&armservicelinker.LinkerClientBeginValidateOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/servicelinker/resource-manager/Microsoft.ServiceLinker/stable/2022-05-01/examples/GetConfigurations.json
func ExampleLinkerClient_ListConfigurations() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armservicelinker.NewLinkerClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.ListConfigurations(ctx,
		"<resource-uri>",
		"<linker-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armapimanagement_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListGatewayCertificateAuthorities.json
func ExampleGatewayCertificateAuthorityClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGatewayCertificateAuthorityClient().NewListByServicePager("rg1", "apimService1", "gw1", &armapimanagement.GatewayCertificateAuthorityClientListByServiceOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.GatewayCertificateAuthorityCollection = armapimanagement.GatewayCertificateAuthorityCollection{
		// 	Value: []*armapimanagement.GatewayCertificateAuthorityContract{
		// 		{
		// 			Name: to.Ptr("cert1"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/gateways/certificateAuthorities"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/certificateAuthorities/cert1"),
		// 			Properties: &armapimanagement.GatewayCertificateAuthorityContractProperties{
		// 				IsTrusted: to.Ptr(false),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("cert2"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/gateways/certificateAuthorities"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/certificateAuthorities/cert2"),
		// 			Properties: &armapimanagement.GatewayCertificateAuthorityContractProperties{
		// 				IsTrusted: to.Ptr(true),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementHeadGatewayCertificateAuthority.json
func ExampleGatewayCertificateAuthorityClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewGatewayCertificateAuthorityClient().GetEntityTag(ctx, "rg1", "apimService1", "gw1", "cert1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetGatewayCertificateAuthority.json
func ExampleGatewayCertificateAuthorityClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGatewayCertificateAuthorityClient().Get(ctx, "rg1", "apimService1", "gw1", "cert1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GatewayCertificateAuthorityContract = armapimanagement.GatewayCertificateAuthorityContract{
	// 	Name: to.Ptr("cert1"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/gateways/certificateAuthorities"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/certificateAuthorities/cert1"),
	// 	Properties: &armapimanagement.GatewayCertificateAuthorityContractProperties{
	// 		IsTrusted: to.Ptr(true),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateGatewayCertificateAuthority.json
func ExampleGatewayCertificateAuthorityClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGatewayCertificateAuthorityClient().CreateOrUpdate(ctx, "rg1", "apimService1", "gw1", "cert1", armapimanagement.GatewayCertificateAuthorityContract{
		Properties: &armapimanagement.GatewayCertificateAuthorityContractProperties{
			IsTrusted: to.Ptr(false),
		},
	}, &armapimanagement.GatewayCertificateAuthorityClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GatewayCertificateAuthorityContract = armapimanagement.GatewayCertificateAuthorityContract{
	// 	Name: to.Ptr("cert1"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/gateways/certificateAuthorities"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/certificateAuthorities/cert1"),
	// 	Properties: &armapimanagement.GatewayCertificateAuthorityContractProperties{
	// 		IsTrusted: to.Ptr(false),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeleteGatewayCertificateAuthority.json
func ExampleGatewayCertificateAuthorityClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewGatewayCertificateAuthorityClient().Delete(ctx, "rg1", "apimService1", "gw1", "default", "*", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

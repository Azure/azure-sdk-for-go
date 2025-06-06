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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListGatewayHostnameConfigurations.json
func ExampleGatewayHostnameConfigurationClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGatewayHostnameConfigurationClient().NewListByServicePager("rg1", "apimService1", "gw1", &armapimanagement.GatewayHostnameConfigurationClientListByServiceOptions{Filter: nil,
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
		// page.GatewayHostnameConfigurationCollection = armapimanagement.GatewayHostnameConfigurationCollection{
		// 	Value: []*armapimanagement.GatewayHostnameConfigurationContract{
		// 		{
		// 			Name: to.Ptr("default"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/gateways/hostnameConfigurations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/hostnameConfigurations/default"),
		// 			Properties: &armapimanagement.GatewayHostnameConfigurationContractProperties{
		// 				CertificateID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/certificates/cert1"),
		// 				Hostname: to.Ptr("*"),
		// 				NegotiateClientCertificate: to.Ptr(false),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("default"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/gateways/hostnameConfigurations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/hostnameConfigurations/hostname1"),
		// 			Properties: &armapimanagement.GatewayHostnameConfigurationContractProperties{
		// 				CertificateID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/certificates/cert2"),
		// 				Hostname: to.Ptr("foo.bar.com"),
		// 				NegotiateClientCertificate: to.Ptr(true),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementHeadGatewayHostnameConfiguration.json
func ExampleGatewayHostnameConfigurationClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewGatewayHostnameConfigurationClient().GetEntityTag(ctx, "rg1", "apimService1", "gw1", "default", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetGatewayHostnameConfiguration.json
func ExampleGatewayHostnameConfigurationClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGatewayHostnameConfigurationClient().Get(ctx, "rg1", "apimService1", "gw1", "default", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GatewayHostnameConfigurationContract = armapimanagement.GatewayHostnameConfigurationContract{
	// 	Name: to.Ptr("default"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/gateways/hostnameConfigurations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/hostnameConfigurations/default"),
	// 	Properties: &armapimanagement.GatewayHostnameConfigurationContractProperties{
	// 		CertificateID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/certificates/cert1"),
	// 		Hostname: to.Ptr("*"),
	// 		NegotiateClientCertificate: to.Ptr(false),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateGatewayHostnameConfiguration.json
func ExampleGatewayHostnameConfigurationClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGatewayHostnameConfigurationClient().CreateOrUpdate(ctx, "rg1", "apimService1", "gw1", "default", armapimanagement.GatewayHostnameConfigurationContract{
		Properties: &armapimanagement.GatewayHostnameConfigurationContractProperties{
			CertificateID:              to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/certificates/cert1"),
			Hostname:                   to.Ptr("*"),
			HTTP2Enabled:               to.Ptr(true),
			NegotiateClientCertificate: to.Ptr(false),
			Tls10Enabled:               to.Ptr(false),
			Tls11Enabled:               to.Ptr(false),
		},
	}, &armapimanagement.GatewayHostnameConfigurationClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GatewayHostnameConfigurationContract = armapimanagement.GatewayHostnameConfigurationContract{
	// 	Name: to.Ptr("default"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/gateways/hostnameConfigurations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/gateways/gw1/hostnameConfigurations/default"),
	// 	Properties: &armapimanagement.GatewayHostnameConfigurationContractProperties{
	// 		CertificateID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/certificates/cert1"),
	// 		Hostname: to.Ptr("*"),
	// 		HTTP2Enabled: to.Ptr(true),
	// 		NegotiateClientCertificate: to.Ptr(false),
	// 		Tls10Enabled: to.Ptr(false),
	// 		Tls11Enabled: to.Ptr(false),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeleteGatewayHostnameConfiguration.json
func ExampleGatewayHostnameConfigurationClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewGatewayHostnameConfigurationClient().Delete(ctx, "rg1", "apimService1", "gw1", "default", "*", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// x-ms-original-file: specification/network/resource-manager/Microsoft.Network/stable/2021-05-01/examples/FirewallPolicySignatureOverridesPatch.json
func ExampleFirewallPolicyIdpsSignaturesOverridesClient_Patch() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armnetwork.NewFirewallPolicyIdpsSignaturesOverridesClient("<subscription-id>", cred, nil)
	res, err := client.Patch(ctx,
		"<resource-group-name>",
		"<firewall-policy-name>",
		armnetwork.SignaturesOverrides{
			Name: to.StringPtr("<name>"),
			Type: to.StringPtr("<type>"),
			ID:   to.StringPtr("<id>"),
			Properties: &armnetwork.SignaturesOverridesProperties{
				Signatures: map[string]*string{
					"2000105": to.StringPtr("Off"),
					"2000106": to.StringPtr("Deny"),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.FirewallPolicyIdpsSignaturesOverridesClientPatchResult)
}

// x-ms-original-file: specification/network/resource-manager/Microsoft.Network/stable/2021-05-01/examples/FirewallPolicySignatureOverridesPut.json
func ExampleFirewallPolicyIdpsSignaturesOverridesClient_Put() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armnetwork.NewFirewallPolicyIdpsSignaturesOverridesClient("<subscription-id>", cred, nil)
	res, err := client.Put(ctx,
		"<resource-group-name>",
		"<firewall-policy-name>",
		armnetwork.SignaturesOverrides{
			Name: to.StringPtr("<name>"),
			Type: to.StringPtr("<type>"),
			ID:   to.StringPtr("<id>"),
			Properties: &armnetwork.SignaturesOverridesProperties{
				Signatures: map[string]*string{
					"2000105": to.StringPtr("Off"),
					"2000106": to.StringPtr("Deny"),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.FirewallPolicyIdpsSignaturesOverridesClientPutResult)
}

// x-ms-original-file: specification/network/resource-manager/Microsoft.Network/stable/2021-05-01/examples/FirewallPolicySignatureOverridesGet.json
func ExampleFirewallPolicyIdpsSignaturesOverridesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armnetwork.NewFirewallPolicyIdpsSignaturesOverridesClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<firewall-policy-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.FirewallPolicyIdpsSignaturesOverridesClientGetResult)
}

// x-ms-original-file: specification/network/resource-manager/Microsoft.Network/stable/2021-05-01/examples/FirewallPolicySignatureOverridesList.json
func ExampleFirewallPolicyIdpsSignaturesOverridesClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armnetwork.NewFirewallPolicyIdpsSignaturesOverridesClient("<subscription-id>", cred, nil)
	res, err := client.List(ctx,
		"<resource-group-name>",
		"<firewall-policy-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.FirewallPolicyIdpsSignaturesOverridesClientListResult)
}

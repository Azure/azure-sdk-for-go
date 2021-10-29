//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func ExampleNetworkSecurityGroupsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<network security group name>",
		armnetwork.NetworkSecurityGroup{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr("<network security group name>"),
				Location: to.StringPtr("<Azure location>"),
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	log.Printf("network security group ID: %v", *resp.NetworkSecurityGroup.ID)
}

func ExampleNetworkSecurityGroupsClient_BeginCreateOrUpdate_withSSHandHTTPSrules() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<network security group name>",
		armnetwork.NetworkSecurityGroup{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr("<network security group name>"),
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armnetwork.NetworkSecurityGroupPropertiesFormat{
				SecurityRules: []*armnetwork.SecurityRule{
					{
						Name: to.StringPtr("allow_ssh"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							Protocol:                 armnetwork.SecurityRuleProtocolTCP.ToPtr(),
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("1-65535"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("22"),
							Access:                   armnetwork.SecurityRuleAccessAllow.ToPtr(),
							Direction:                armnetwork.SecurityRuleDirectionInbound.ToPtr(),
							Priority:                 to.Int32Ptr(100),
						},
					},
					{
						Name: to.StringPtr("allow_https"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							Protocol:                 armnetwork.SecurityRuleProtocolTCP.ToPtr(),
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("1-65535"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("443"),
							Access:                   armnetwork.SecurityRuleAccessAllow.ToPtr(),
							Direction:                armnetwork.SecurityRuleDirectionInbound.ToPtr(),
							Priority:                 to.Int32Ptr(200),
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
		log.Fatalf("failed to create resource: %v", err)
	}
	log.Printf("network security group ID: %v", *resp.NetworkSecurityGroup.ID)
}

func ExampleNetworkSecurityGroupsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient("<subscription ID>", cred, nil)
	resp, err := client.Get(context.Background(), "<resource group name>", "<network security group name>", nil)
	if err != nil {
		log.Fatalf("failed to get resource: %v", err)
	}
	log.Printf("network security group ID: %v", *resp.NetworkSecurityGroup.ID)
}

func ExampleNetworkSecurityGroupsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient("<subscription ID>", cred, nil)
	resp, err := client.BeginDelete(context.Background(), "<resource group name>", "<network security group name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
}

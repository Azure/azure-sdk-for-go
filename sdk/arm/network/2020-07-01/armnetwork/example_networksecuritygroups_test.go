// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func ExampleNetworkSecurityGroupsOperations_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("network security group ID: %v", *resp.NetworkSecurityGroup.ID)
}

func ExampleNetworkSecurityGroupsOperations_BeginCreateOrUpdate_withSSHandHTTPSrules() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
				SecurityRules: &[]armnetwork.SecurityRule{
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
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("network security group ID: %v", *resp.NetworkSecurityGroup.ID)
}

func ExampleNetworkSecurityGroupsOperations_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<network security group name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("network security group ID: %v", *resp.NetworkSecurityGroup.ID)
}

func ExampleNetworkSecurityGroupsOperations_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.BeginDelete(context.Background(), "<resource group name>", "<network security group name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
}

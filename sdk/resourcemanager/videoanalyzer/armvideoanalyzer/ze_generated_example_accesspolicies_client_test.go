//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armvideoanalyzer_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/videoanalyzer/armvideoanalyzer"
)

// x-ms-original-file: specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/access-policy-list.json
func ExampleAccessPoliciesClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armvideoanalyzer.NewAccessPoliciesClient("<subscription-id>", cred, nil)
	pager := client.List("<resource-group-name>",
		"<account-name>",
		&armvideoanalyzer.AccessPoliciesClientListOptions{Top: to.Int32Ptr(2)})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/access-policy-get.json
func ExampleAccessPoliciesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armvideoanalyzer.NewAccessPoliciesClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<access-policy-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.AccessPoliciesClientGetResult)
}

// x-ms-original-file: specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/access-policy-create.json
func ExampleAccessPoliciesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armvideoanalyzer.NewAccessPoliciesClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdate(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<access-policy-name>",
		armvideoanalyzer.AccessPolicyEntity{
			Properties: &armvideoanalyzer.AccessPolicyProperties{
				Authentication: &armvideoanalyzer.JwtAuthentication{
					Type: to.StringPtr("<type>"),
					Audiences: []*string{
						to.StringPtr("audience1")},
					Claims: []*armvideoanalyzer.TokenClaim{
						{
							Name:  to.StringPtr("<name>"),
							Value: to.StringPtr("<value>"),
						},
						{
							Name:  to.StringPtr("<name>"),
							Value: to.StringPtr("<value>"),
						}},
					Issuers: []*string{
						to.StringPtr("issuer1"),
						to.StringPtr("issuer2")},
					Keys: []armvideoanalyzer.TokenKeyClassification{
						&armvideoanalyzer.RsaTokenKey{
							Type: to.StringPtr("<type>"),
							Kid:  to.StringPtr("<kid>"),
							Alg:  armvideoanalyzer.AccessPolicyRsaAlgo("RS256").ToPtr(),
							E:    to.StringPtr("<e>"),
							N:    to.StringPtr("<n>"),
						},
						&armvideoanalyzer.EccTokenKey{
							Type: to.StringPtr("<type>"),
							Kid:  to.StringPtr("<kid>"),
							Alg:  armvideoanalyzer.AccessPolicyEccAlgo("ES256").ToPtr(),
							X:    to.StringPtr("<x>"),
							Y:    to.StringPtr("<y>"),
						}},
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.AccessPoliciesClientCreateOrUpdateResult)
}

// x-ms-original-file: specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/access-policy-delete.json
func ExampleAccessPoliciesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armvideoanalyzer.NewAccessPoliciesClient("<subscription-id>", cred, nil)
	_, err = client.Delete(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<access-policy-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/access-policy-patch.json
func ExampleAccessPoliciesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armvideoanalyzer.NewAccessPoliciesClient("<subscription-id>", cred, nil)
	res, err := client.Update(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<access-policy-name>",
		armvideoanalyzer.AccessPolicyEntity{
			Properties: &armvideoanalyzer.AccessPolicyProperties{
				Authentication: &armvideoanalyzer.JwtAuthentication{
					Type: to.StringPtr("<type>"),
					Audiences: []*string{
						to.StringPtr("audience1")},
					Claims: []*armvideoanalyzer.TokenClaim{
						{
							Name:  to.StringPtr("<name>"),
							Value: to.StringPtr("<value>"),
						},
						{
							Name:  to.StringPtr("<name>"),
							Value: to.StringPtr("<value>"),
						}},
					Issuers: []*string{
						to.StringPtr("issuer1"),
						to.StringPtr("issuer2")},
					Keys: []armvideoanalyzer.TokenKeyClassification{
						&armvideoanalyzer.RsaTokenKey{
							Type: to.StringPtr("<type>"),
							Kid:  to.StringPtr("<kid>"),
							Alg:  armvideoanalyzer.AccessPolicyRsaAlgo("RS256").ToPtr(),
							E:    to.StringPtr("<e>"),
							N:    to.StringPtr("<n>"),
						},
						&armvideoanalyzer.EccTokenKey{
							Type: to.StringPtr("<type>"),
							Kid:  to.StringPtr("<kid>"),
							Alg:  armvideoanalyzer.AccessPolicyEccAlgo("ES256").ToPtr(),
							X:    to.StringPtr("<x>"),
							Y:    to.StringPtr("<y>"),
						}},
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.AccessPoliciesClientUpdateResult)
}

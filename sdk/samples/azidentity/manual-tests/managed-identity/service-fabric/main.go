// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	// NOTE: the service fabric cluster used for testing uses a self-signed certificate,
	// this configuration is only used in a development environment. This is an insecure
	// configuration.
	cl := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: cl,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Calling GetToken()...")
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://vault.azure.net"}})
	if err != nil {
		panic(err)
	}
	fmt.Println("Success! Token received.")
}

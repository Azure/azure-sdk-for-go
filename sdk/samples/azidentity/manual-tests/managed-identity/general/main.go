// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var clientID = flag.String("clientID", "", "optional client ID of a user assigned identity. Mutually exclusive with resID.")
var printToken = flag.Bool("printToken", false, "print the acquired token")
var resID = flag.String("resID", "", "optional resource ID of a user assigned identity. Mutually exclusive with clientID.")
var scope = flag.String("scope", "https://management.core.windows.net//.default", "optional scope for access token")

func main() {
	flag.Parse()
	opts := &azidentity.ManagedIdentityCredentialOptions{}
	if *clientID != "" {
		if *resID != "" {
			panic(`"clientID" and "resID" are mutually exclusive`)
		}
		opts.ID = azidentity.ClientID(*clientID)
	} else if *resID != "" {
		opts.ID = azidentity.ResourceID(*resID)
	}
	cred, err := azidentity.NewManagedIdentityCredential(opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("Calling GetToken()...")
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{*scope}})
	if err != nil {
		panic(err)
	}
	fmt.Println("Success! Token received.")
	if *printToken {
		fmt.Println(tk.Token)
	}
}

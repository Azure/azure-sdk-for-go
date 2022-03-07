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

var c = flag.String("c", "", "optional client ID of a user assigned identity. Mutually exclusive with r.")
var p = flag.Bool("p", false, "print the acquired token")
var r = flag.String("r", "", "optional resource ID of a user assigned identity. Mutually exclusive with c.")
var s = flag.String("s", "https://management.core.windows.net//.default", "optional scope for access token")

func main() {
	flag.Parse()
	opts := &azidentity.ManagedIdentityCredentialOptions{}
	if *c != "" {
		if *r != "" {
			panic(`"c" and "r" are mutually exclusive`)
		}
		opts.ID = azidentity.ClientID(*c)
	} else if *r != "" {
		opts.ID = azidentity.ResourceID(*r)
	}
	cred, err := azidentity.NewManagedIdentityCredential(opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("Calling GetToken()...")
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{*s}})
	if err != nil {
		panic(err)
	}
	fmt.Println("Success! Token received.")
	if *p {
		fmt.Println(tk.Token)
	}
}

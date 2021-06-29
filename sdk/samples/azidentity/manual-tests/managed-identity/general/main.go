// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	cred, err := azidentity.NewManagedIdentityCredential("", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Calling GetToken()...")
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{"https://vault.azure.net"}})
	if err != nil {
		panic(err)
	}
	fmt.Println("Success! Token received.")
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentitysamples

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/jongio/azidext/go/azidext"
)

// Please note that the examples in this file use the Azure SDK for Go V1 code base with
// the azidentity package from the Azure SDK for Go V2. New applications should simply
// use only V2 modules, such as github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources

var (
	clientID     = os.Getenv("AZURE_CLIENT_ID")
	clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	tenantID     = os.Getenv("AZURE_TENANT_ID")
)

// ExampleNewDefaultAzureCredential demonstrates an adapter for azidentity.DefaultAzureCredential that enables using the credential with the Azure SDK for Go V1.
// NewDefaultAzureCredentialAdapter can replace auth.NewAuthorizerFromEnvironment(), because DefaultAzureCredential similarly reads environment variables for configuration.
// See azidentity documentation for more information about DefaultAzureCredential.
func ExampleNewDefaultAzureCredential() {
	// Passing nil configures the credential for Azure Public Cloud. To run in another cloud, such as Azure Government,
	// specify the Azure Resource Manager scope for that cloud in azidext.DefaultAzureCredentialOptions
	a, err := azidext.NewDefaultAzureCredentialAdapter(nil)
	if err != nil {
		panic(err)
	}

	client := subscription.NewSubscriptionsClient()
	client.Authorizer = a

	subscriptions, err := client.List(context.Background())
	if err != nil {
		panic(err)
	}
	subNames := []string{}
	for _, s := range subscriptions.Values() {
		subNames = append(subNames, *s.DisplayName)
	}
	// Output:
}

// ExampleNewClientSecretCredential demonstrates NewAzureIdentityCredentialAdapter, which adapts any azidentity credential type for use with the Azure SDK for Go V1.
func ExampleNewClientSecretCredential() {
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	// This scope is correct for Azure Resource Manager in Azure Public Cloud. Other clouds require a different scope.
	a := azidext.NewTokenCredentialAdapter(cred, []string{"https://management.azure.com//.default"})
	if err != nil {
		panic(err)
	}

	client := subscription.NewSubscriptionsClient()
	client.Authorizer = a

	subscriptions, err := client.List(context.Background())
	if err != nil {
		panic(err)
	}
	subNames := []string{}
	for _, s := range subscriptions.Values() {
		subNames = append(subNames, *s.DisplayName)
	}
	// Output:
}

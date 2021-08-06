// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentitysamples

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"

	"github.com/jongio/azidext/go/azidext"
)

// Please note that the examples in this file are using the Azure SDK for Go V1 code base, along
// with azidentity package from the Azure SDK for Go V2.
// The adapter in the azidext package provides a simple way to integrate azidentity credentials
// as authorizers for the V1 code base.

const (
	groupName = "samplegroup"
)

// Environment variables required for EnvironmentCredential to work and/or DefaultAzureCredential
var (
	clientID       = os.Getenv("AZURE_CLIENT_ID")
	clientSecret   = os.Getenv("AZURE_CLIENT_SECRET")
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	tenantID       = os.Getenv("AZURE_TENANT_ID")
)

var location = os.Getenv("AZURE_LOCATION")

// ExampleNewDefaultAzureCredential for using the DefaultAzureCredential through the NewDefaultAzureCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
// NewDefaultAzureCredentialAdapter should be used to replace auth.NewAuthorizerFromEnvironment(). DefaultAzureCredential, similarly to NewAuthorizerFromEnvironment, checks for
// environment variables that can construct ClientSecretCredentials, ClientCertificateCredentials, UsernamePasswordCredentials, ManagedIdentityCredentials and AzureCLICredentials.
func ExampleNewDefaultAzureCredential() {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	// call azidext.NewDefaultAzureCredentialAdapter in order to get an authorizer with a DefaultAzureCredential
	// leave azidext.DefaultAzureCredentialOptions as nil to get the default scope for management APIs.
	// The default scope is: https://management.azure.com//.default.
	// NOTE: Scopes define the set of resources and permissions that the credential will have assigned to it.
	// 		 To read more about scopes, see: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-permissions-and-consent
	a, err := azidext.NewDefaultAzureCredentialAdapter(nil)
	if err != nil {
		panic("failed to get credential")
	}
	groupsClient.Authorizer = a
	// use the groups client with the azidentity credential in the authorizer
	group, err := groupsClient.CreateOrUpdate(
		context.Background(),
		groupName,
		resources.Group{
			Location: to.StringPtr(location),
		})
	if err != nil {
		panic(err)
	}
	fmt.Println(*group.Name)
	// Output:
	// samplegroup
}

// ExampleNewClientSecretCredential for using the ClientSecretCredential with the NewAzureIdentityCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
// NewAzureIdentityCredentialAdapter can take any credential type defined in azidentity and convert it to an authorizer that is compatible with the Azure SDK for Go
// V1 implementation. For a list of the credentials that azidentity includes, please see: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity.
// NewClientSecretCredential can be used in place of auth.NewClientCredentialsConfig().
func ExampleNewClientSecretCredential() {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	// instantiate a new ClientSecretCredential as specified in the documentation
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	// call azidext.NewAzureIdentityCredentialAdapter with the azidentity credential and necessary scope
	// NOTE: Scopes define the set of resources and permissions that the credential will have assigned to it.
	// 		 To read more about scopes, see: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-permissions-and-consent
	a := azidext.NewAzureIdentityCredentialAdapter(
		cred,
		azcore.AuthenticationPolicyOptions{
			Options: azcore.TokenRequestOptions{
				Scopes: []string{"https://management.azure.com/.default"}}})
	if err != nil {
		panic("failed to get credential")
	}
	// assign the authorizer to your client's authorizer
	groupsClient.Authorizer = a
	// perform an operation with the complete client
	list, err := groupsClient.ListComplete(context.Background(), "", nil)
	if err != nil {
		panic(err)
	}
	for list.NotDone() {
		fmt.Println(*list.Value().Name)
		err = list.Next()
		if err != nil {
			panic(err.Error())
		}
	}
	// Output:
	// samplegroup
}

// ExampleNewEnvironmentCredential for using the EnvironmentCredential with the NewAzureIdentityCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
// NewAzureIdentityCredentialAdapter can take any credential type defined in azidentity and convert it to an authorizer that is compatible with the Azure SDK for Go
// V1 implementation. For a list of the credentials that azidentity includes, please see: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity.
// NewEnvironmentCredential can be used in place of auth.NewAuthorizerFromEnvironment(). An important distinction is that NewEnvironmentCredential does not include Managed
// Identity credential, for a credential that also checks the environment for Managed Identity credential use the NewDefaultAzureCredentialAdapter. Alternatively,
// create a custom credential chain with NewChainedTokenCredential and add all desired token credentials to try into the chain.
func ExampleNewEnvironmentCredential() {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		panic(err)
	}
	// call azidext.NewAzureIdentityCredentialAdapter with the azidentity credential and necessary scopes
	// NOTE: Scopes define the set of resources and/or permissions that the credential will have assigned to it.
	// 		 To read more about scopes, see: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-permissions-and-consent
	a := azidext.NewAzureIdentityCredentialAdapter(
		cred,
		azcore.AuthenticationPolicyOptions{
			Options: azcore.TokenRequestOptions{
				Scopes: []string{"https://management.azure.com/.default"}}})
	if err != nil {
		panic("failed to get credential")
	}
	// assign the authorizer to your client's authorizer
	groupsClient.Authorizer = a
	// perform an operation with the complete client
	_, err = groupsClient.Delete(context.Background(), groupName)
	if err == nil {
		fmt.Println("Delete in progress..")
	}
	// Output:
	// Delete in progress..
}

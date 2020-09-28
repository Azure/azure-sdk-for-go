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

var (
	location  = os.Getenv("AZURE_LOCATION")
	userAgent = "azidentitysample"
)

// ExampleGroupsClientWithDefaultAzureCredential for using the DefaultAzureCredential through the NewDefaultAzureCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
func ExampleGroupsClientWithDefaultAzureCredential() {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	// call azidext.NewDefaultAzureCredentialAdapter in order to get an authorizer with a DefaultAzureCredential
	// NOTE: Scopes define the set of resource and permissions that the credential will have assigned to it.
	// 		 To read more about scopes, see: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-permissions-and-consent
	a, err := azidext.NewDefaultAzureCredentialAdapter(
		&azidext.DefaultAzureCredentialOptions{
			AuthenticationPolicy: &azcore.AuthenticationPolicyOptions{
				Options: azcore.TokenRequestOptions{
					Scopes: []string{"https://management.azure.com/.default"}}}})
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

// ExampleGroupsClientWithClientSecretCredential for using the ClientSecretCredential with the NewAzureIdentityCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
func ExampleGroupsClientWithClientSecretCredential() {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	// instantiate a new ClientSecretCredential as specified in the documentation
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	// call azidext.NewAzureIdentityCredentialAdapter with the azidentity credential and necessary scope
	// NOTE: Scopes define the set of resource and permissions that the credential will have assigned to it.
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
		list.Next()
	}
	// Output:
	// samplegroup
}

// ExampleGroupsClientWithEnvironmentCredential for using the EnvironmentCredential with the NewAzureIdentityCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
func ExampleGroupsClientWithEnvironmentCredential() {
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

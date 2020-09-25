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

// Example for using the DefaultAzureCredential through the NewDefaultAzureCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
func getGroupsClientWithDefaultAzureCredential() resources.GroupsClient {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	a, err := azidext.NewDefaultAzureCredentialAdapter(nil)
	if err != nil {
		panic("failed to get credential")
	}
	groupsClient.Authorizer = a
	groupsClient.AddToUserAgent(userAgent)
	return groupsClient
}

// Example for using the ClientSecretCredential with the NewAzureIdentityCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
func getGroupsClientWithClientSecretCredential() resources.GroupsClient {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	a := azidext.NewAzureIdentityCredentialAdapter(cred, azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{"https://management.azure.com/.default"}}})
	if err != nil {
		panic("failed to get credential")
	}
	groupsClient.Authorizer = a
	groupsClient.AddToUserAgent(userAgent)
	return groupsClient
}

// Example for using the EnvironmentCredential with the NewAzureIdentityCredentialAdapter and assigning the credential to the
// SDK V1 authorizer.
func getGroupsClientWithEnvironmentCredential() resources.GroupsClient {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		panic(err)
	}
	a := azidext.NewAzureIdentityCredentialAdapter(cred, azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{"https://management.azure.com/.default"}}})
	if err != nil {
		panic("failed to get credential")
	}
	groupsClient.Authorizer = a
	groupsClient.AddToUserAgent(userAgent)
	return groupsClient
}

// ExampleSDKV1ResourcesCreateOrUpdate creates a new resource group using the DefaultAzureCredential through the
// getGroupsClientWithDefaultAzureCredential function.
func ExampleSDKV1ResourcesCreateOrUpdate() {
	// CreateGroup creates a new resource group named by env var
	groupsClient := getGroupsClientWithDefaultAzureCredential()
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

// ExampleSDKV1ResourcesListGroups gets an iterator that gets all resource groups in the subscription using the
// ClientSecretCredential through the getGroupsClientWithClientSecretCredential function.
func ExampleSDKV1ResourcesListGroups() {
	groupsClient := getGroupsClientWithClientSecretCredential()
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

// ExampleSDKV1ResourcesDeleteGroup removes the resource group using the EnvironmentCredential through the
// getGroupsClientWithEnvironmentCredential function.
func ExampleSDKV1ResourcesDeleteGroup() {
	groupsClient := getGroupsClientWithEnvironmentCredential()
	_, err := groupsClient.Delete(context.Background(), groupName)
	if err == nil {
		fmt.Println("Delete in progress..")
	}
	// Output:
	// Delete in progress..
}

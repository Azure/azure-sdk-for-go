package examples

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"fmt"
	"os"

	"bytes"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// SDK represents a working authenticated SDK structure.
type SDK struct {

	// Authentication holds misc. authentication information for interacting with the SDK.
	Authentication *Authentication
}

// Authentication represents the values needed to authenticate with the ARM API.
type Authentication struct {

	// ClientID is the `appId` from an Azure Service Principal.
	ClientID string

	// ClientSecret is the `password` from an Azure Service Principal.
	ClientSecret string

	// SubscriptionID is the unique Azure subscription ID to use.
	SubscriptionID string

	// TenantID is the `tenant` from an Azure Service Principal.
	TenantID string

	// HashMap is the data value some parts of the SDK will expect, and is stored here as convenience.
	HashMap map[string]string

	// AuthenticatedToken is the generated token that is needed to authenticate with each client.
	AuthenticatedToken *adal.ServicePrincipalToken
}

func NewSDK() (*SDK, error) {

	// -------------------------------------------------------------------
	// Validation that the environment is configured correctly.
	//
	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("empty $AZURE_CLIENT_ID")
	}
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("empty $AZURE_CLIENT_SECRET")
	}
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		return nil, fmt.Errorf("empty $AZURE_SUBSCRIPTION_ID")
	}
	tenantID := os.Getenv("AZURE_TENANT_ID")
	if tenantID == "" {
		return nil, fmt.Errorf("empty $AZURE_TENANT_ID")
	}

	// Here we build the SDK based on a user's configuration.
	sdk := &SDK{
		Authentication: &Authentication{
			ClientID:       clientID,
			ClientSecret:   clientSecret,
			SubscriptionID: subscriptionID,
			TenantID:       tenantID,
			HashMap: map[string]string{
				"AZURE_CLIENT_ID":       clientID,
				"AZURE_CLIENT_SECRET":   clientSecret,
				"AZURE_SUBSCRIPTION_ID": subscriptionID,
				"AZURE_TENANT_ID":       tenantID,
			},
		},
	}

	// Generate an authenticated token.
	authenticatedToken, err := helpers.NewServicePrincipalTokenFromCredentials(sdk.Authentication.HashMap, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}
	sdk.Authentication.AuthenticatedToken = authenticatedToken
	return sdk, nil
}

// E will check if an error is nil, and exit the program in the case of an error.
// Use this as convenience to simplify code in various examples.
func E(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}

// S will take an arbitrary string, and return a newly created pointer to this string.
// Use this as convenience to simplify code in various examples.
func S(i string) *string {
	return &i
}

// I32 will take an arbitrary int, and return a newly created pointer to an int32.
// Use this as convenience to simplify code in various examples.
func I32(i int) *int32 {
	i32 := int32(i)
	return &i32
}

// F64 will take an arbitrary int, and return a newly created pointer to a float64.
// Use this as convenience to simplify code in various examples.
func F64(i int) *float64 {
	f64 := float64(i)
	return &f64
}

// PrettyJSON will take an arbitrarry interface{} and attempt to print the data structure out to STDOUT.
// Use this as convenience to simplify code in various examples.
func PrettyJSON(data interface{}) error {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	fmt.Println(buffer.String())
	return nil
}

// EnsureResourceGroup will ensure that a resource group of an arbitrary name and location exists.
// Use this as convenience to simplify code in various examples.
func (sdk *SDK) EnsureResourceGroup(resourceGroupName, location string) error {
	resourceGroup := resources.NewGroupsClient(sdk.Authentication.SubscriptionID)
	resourceGroup.Authorizer = autorest.NewBearerAuthorizer(sdk.Authentication.AuthenticatedToken)

	// ---------------------------------
	// Resource Group Parameters
	parameters := resources.Group{
		Name:     &resourceGroupName,
		Location: &location,
	}

	newResourceGroup, err := resourceGroup.CreateOrUpdate(resourceGroupName, parameters)
	if err != nil {
		return err
	}
	fmt.Printf("Ensured resource group [%s] exists\n", *newResourceGroup.Name)
	return nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armnetwork_test

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var (
	// The name of the public IP addresses created in Azure
	ip = os.Getenv("AZURE_IP")
	// Azure location where the resource will be created
	location = os.Getenv("AZURE_LOCATION")
	// Azure resource group to retrieve and create resources
	resourceGroupName = os.Getenv("AZURE_RESOURCE_GROUP")
	// The subscription ID where the resource group exists
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	// The Virtual Network that will be used in the examples
	vnetName = os.Getenv("AZURE_VNET")
	// The subnet to be used in the samples
	subnetName = os.Getenv("AZURE_SUBNET")
)

// returns a credential that can be used to authenticate with Azure Active Directory
func getCredential() azcore.Credential {
	// NewEnvironmentCredential() will read various environment vars
	// to obtain a credential.  see the documentation for more info.
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		panic(err)
	}
	return cred
}

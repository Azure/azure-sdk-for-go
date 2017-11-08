package main

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
	"github.com/Azure/azure-sdk-for-go/arm/containerinstance"
	"github.com/Azure/azure-sdk-for-go/arm/examples"
	"github.com/Azure/go-autorest/autorest"
)

var (
	resourceGroupName     = "myresourcegroup"
	resourceGroupLocation = "eastus"
	containerGroupname    = "mycontainergroup"
)

func main() {

	// -----------------------------------------------------------------------------------------------------------------
	// Bootstrap the SDK and ensure a resource group is in place.-
	sdk, err := examples.NewSDK()
	examples.E(err)
	err = sdk.EnsureResourceGroup(resourceGroupName, resourceGroupLocation)
	examples.E(err)

	//------------------------------------------------------------------------------------------------------------------
	// Container Instance
	containerInstance := containerinstance.NewContainerGroupsClient(sdk.Authentication.SubscriptionID)
	containerInstance.Authorizer = autorest.NewBearerAuthorizer(sdk.Authentication.AuthenticatedToken)

	_, err = containerInstance.Delete(resourceGroupName, containerGroupname)
	examples.E(err)
	fmt.Printf("Deleted container group [%s]!\n", containerGroupname)

}

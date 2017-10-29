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

	// -------------------------------------------------------------------------------------------------------------
	// Container Instance Parameters
	parameters := containerinstance.ContainerGroup{
		Name:     &containerGroupname,
		Location: &resourceGroupLocation,
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{

			// Define the OsType for the container group, in this case Linux.
			OsType: containerinstance.Linux,

			// Containers is a set of Container{} structs each defined as below.
			Containers: &[]containerinstance.Container{
				{

					Name: examples.S("golang-aci-example"),

					// ---------------------------------------------------------------------------------------------
					//
					// Container Properties
					//
					// This is what ultimately will define a container for ACI.
					//
					//
					// FIELD                TYPE                             ENCODING
					// Image                *string                          `json:"image,omitempty"`
					// Command              *[]string                        `json:"command,omitempty"`
					// Ports                *[]ContainerPort                 `json:"ports,omitempty"`
					// EnvironmentVariables *[]EnvironmentVariable           `json:"environmentVariables,omitempty"`
					// InstanceView         *ContainerPropertiesInstanceView `json:"instanceView,omitempty"`
					// Resources            *ResourceRequirements            `json:"resources,omitempty"`
					// VolumeMounts         *[]VolumeMount                   `json:"volumeMounts,omitempty"`
					//
					ContainerProperties: &containerinstance.ContainerProperties{

						// Use to "golang:latest" docker container.
						Image: examples.S("golang:latest"),

						// Run the three shell commands listed below.
						Command: &[]string{"go version", "echo $GOPATH", "sleep 100"},

						// Open ports to the container.
						Ports: &[]containerinstance.ContainerPort{
							{
								// Expose port 8080 for the container
								Port: examples.I32(8080),
							},
						},

						// Define environmental variables for the container.
						EnvironmentVariables: &[]containerinstance.EnvironmentVariable{
							{

								// Set the environmental variable $MY_KEY="MY_VALUE"
								Name:  examples.S("MY_KEY"),
								Value: examples.S("MY_VALUE"),
							},
							{

								// Set the environmental variable $MY_OTHER_KEY="MY_OTHER_VALUE".
								Name:  examples.S("MY_OTHER_KEY"),
								Value: examples.S("MY_OTHER_VALUE"),
							},
						},

						// Set system resource constraints for the container.
						Resources: &containerinstance.ResourceRequirements{

							// Define resource limits  for the container.
							Limits: &containerinstance.ResourceLimits{

								// Set a limit of 1GB of memory for the container.
								MemoryInGB: examples.F64(1),

								// Set a request of 1 core for the CPU for the container.
								CPU: examples.F64(1),
							},

							// Define resource requests (nice-to-haves) for the container.
							Requests: &containerinstance.ResourceRequests{

								// Set a request of 1gb of memory for the container.
								MemoryInGB: examples.F64(1),

								// Set a request of 1 core for the CPU for the container.
								CPU: examples.F64(1),
							},
						},
					},
				},
			},
		},
	}
	newContainerGroup, err := containerInstance.CreateOrUpdate(resourceGroupName, containerGroupname, parameters)
	examples.E(err)
	fmt.Printf("Created container group [%s]!\n", containerGroupname)
	examples.PrettyJSON(newContainerGroup)
}

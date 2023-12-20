# Azure Automanage Module for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage)

The `armautomanage` module provides operations for working with Azure Automanage.

[Source code](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/automanage/armautomanage)

# Getting started

## Prerequisites

- an [Azure subscription](https://azure.microsoft.com/free/)
- Go 1.18 or above (You could download and install the latest version of Go from [here](https://go.dev/doc/install). It will replace the existing Go on your machine. If you want to install multiple Go versions on the same machine, you could refer this [doc](https://go.dev/doc/manage-install).)

## Install and import required packages

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Install the Azure Automanage and Azure Identity modules:

```sh
go get "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage"
go get "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
```

Import the Azure Automanage and Azure Identity modules:

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage"
)
```

## Authorization

When creating a client, you will need to provide a credential for authenticating with Azure Automanage.  The `azidentity` module provides facilities for various ways of authenticating with Azure including client/secret, certificate, managed identity, and more.

```go
cred, err := azidentity.NewDefaultAzureCredential(nil)
```

For more information on authentication, please see the documentation for `azidentity` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity).

## Client Factory

Azure Automanage module consists of one or more clients. We provide a client factory which could be used to create any client in this module.

```go
clientFactory, err := armautomanage.NewClientFactory(<subscription ID>, cred, nil)
```

You can use `ClientOptions` in package `github.com/Azure/azure-sdk-for-go/sdk/azcore/arm` to set endpoint to connect with public and sovereign clouds as well as Azure Stack. For more information, please see the documentation for `azcore` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore).

```go
options := arm.ClientOptions {
    ClientOptions: azcore.ClientOptions {
        Cloud: cloud.AzureChina,
    },
}
clientFactory, err := armautomanage.NewClientFactory(<subscription ID>, cred, &options)
```

## Clients

A client groups a set of related APIs, providing access to its functionality.  Create one or more clients to access the APIs you require using client factory.

```go
reportsClient := clientFactory.NewReportsClient()
configProfilesClient := clientFactory.NewReportsClient()
assignmentClient := clientFactory.NewReportsClient()
```

## Fakes

The fake package contains types used for constructing in-memory fake servers used in unit tests.
This allows writing tests to cover various success/error conditions without the need for connecting to a live service.

Please see https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/samples/fakes for details and examples on how to use fakes.

## Create or Update a Custom Automanage Configuration Profile

To update a profile, provide a value for all properties as if you were creating a configuration profile (ID, Name, Type, Location, Properties, Tags)

```go
configuration := map[string]interface{}{
    "Antimalware/Enable":                false,
    "AzureSecurityCenter/Enable":        true,
    "Backup/Enable":                     false,
    "BootDiagnostics/Enable":            true,
    "ChangeTrackingAndInventory/Enable": true,
    "GuestConfiguration/Enable":         true,
    "LogAnalytics/Enable":               true,
    "UpdateManagement/Enable":           true,
    "VMInsights/Enable":                 true,
}

properties := armautomanage.ConfigurationProfileProperties{
    Configuration: configuration,
}

location := "eastus"
environment := "dev"

// tags may be omitted 
tags := make(map[string]*string)
tags["environment"] = &environment

profile := armautomanage.ConfigurationProfile{
    Location:   &location,
    Properties: &properties,
    Tags:       tags,
}

newProfile, err := configProfilesClient.CreateOrUpdate(context.Background(), configurationProfileName, "resourceGroupName", profile, nil)
```


## Get an Automanage Configuration Profile

```go
profile, err := configProfilesClient.Get(context.Background(), "configurationProfileName", "resourceGroupName", nil)
data, err := json.MarshalIndent(profile, "", "   ")

fmt.Println(string(data))
```


## Delete an Automanage Configuration Profile

```go
_, err := configProfilesClient.Delete(context.Background(), "resourceGroupName", "configurationProfileName", nil)
```


## Get an Automanage Profile Assignment

```go
assignment, err := assignmentClient.Get(context.Background(), "resourceGroupName", "default", "vmName", nil)
data, err := json.MarshalIndent(assignment, "", "   ")
fmt.Println(string(data))
```


## Create an Assignment between a VM and an Automanage Best Practices Production Configuration Profile

```go
configProfileId := "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"

properties := armautomanage.ConfigurationProfileAssignmentProperties{
    ConfigurationProfile: &configProfileId,
}

assignment := armautomanage.ConfigurationProfileAssignment{
    Properties: &properties,
}

// assignment name must be 'default'
newAssignment, err = assignmentClient.CreateOrUpdate(context.Background(), "default", "resourceGroupName", "vmName", assignment, nil)
```


## Create an Assignment between a VM and a Custom Automanage Configuration Profile

```go
configProfileId := "/subscriptions/<subscription ID>/resourceGroups/resourceGroupName/providers/Microsoft.Automanage/configurationProfiles/configurationProfileName"

properties := armautomanage.ConfigurationProfileAssignmentProperties{
    ConfigurationProfile: &configProfileId,
}

assignment := armautomanage.ConfigurationProfileAssignment{
    Properties: &properties,
}

// assignment name must be 'default'
newAssignment, err = assignmentClient.CreateOrUpdate(context.Background(), "default", "resourceGroupName", "vmName", assignment, nil)
```


## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Automanage` label.

# Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any
additional questions or comments.
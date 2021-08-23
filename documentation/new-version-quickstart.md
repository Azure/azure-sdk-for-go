
Getting Started - New Azure Go SDK
=============================================================

We are excited to announce that a new set of management libraries are
now production-ready. Those packages share a number of new features
such as Azure Identity support, HTTP pipeline, error-handling.,etc, and
they also follow the new Azure SDK guidelines which create easy-to-use
APIs that are idiomatic, compatible, and dependable.

You can find the full list of those new libraries
[here](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk)

In this basic quickstart guide, we will walk you through how to
authenticate to Azure and start interacting with Azure resources. There are several possible approaches to
authentication. This document illustrates the most common scenario.

Migration from older versions of Azure management libraries for Go
------------------------------------------------------------------

If you are an existing user of the older version of Azure management library for Go (packages that are located under [`/services`](https://github.com/Azure/azure-sdk-for-go/tree/main/services)), and you are looking for a migration guide to upgrade to the latest version of the SDK, please refer to [this migration guide](https://aka.ms/azsdk/go/mgmt/migration) for detailed instructions.

Prerequisites
-------------

You will need the following values to authenticate to Azure

-   **Subscription ID**
-   **Client ID**
-   **Client Secret**
-   **Tenant ID**

These values can be obtained from the portal, here's the instructions:

### Get Subscription ID

1.  Login into your Azure account
2.  Select Subscriptions in the left sidebar
3.  Select whichever subscription is needed
4.  Click on Overview
5.  Copy the Subscription ID

### Get Client ID / Client Secret / Tenant ID

For information on how to get Client ID, Client Secret, and Tenant ID,
please refer to [this document](https://docs.microsoft.com/azure/active-directory/develop/howto-create-service-principal-portal)

### Setting Environment Variables

After you obtained the values, you need to set the following values as
your environment variables

-   `AZURE_CLIENT_ID`
-   `AZURE_CLIENT_SECRET`
-   `AZURE_TENANT_ID`
-   `AZURE_SUBSCRIPTION_ID`

To set the following environment variables on your development system:

Windows (Note: Administrator access is required)

1.  Open the Control Panel
2.  Click System Security, then System
3.  Click Advanced system settings on the left
4.  Inside the System Properties window, click the `Environment Variables…` button.
5.  Click on the property you would like to change, then click the `Edit…` button. If the property name is not listed, then click the `New…` button.

Linux-based OS :

    export AZURE_CLIENT_ID="__CLIENT_ID__"
    export AZURE_CLIENT_SECRET="__CLIENT_SECRET__"
    export AZURE_TENANT_ID="__TENANT_ID__"
    export AZURE_SUBSCRIPTION_ID="__SUBSCRIPTION_ID__"

Install the package
-------------------

This project uses Go modules for versioning and dependency management.

As an example, to install the Azure Compute module, you would run :

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute
```
We also recommend installing other packages for authentication and core functionalities :

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/armcore
go get github.com/Azure/azure-sdk-for-go/sdk/azcore
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
go get github.com/Azure/azure-sdk-for-go/sdk/to
```

Authentication
--------------

Once the environment is setup, all you need to do is to create an authenticated client. Before creating a client, you will first need to authenticate to Azure. In specific, you will need to provide a credential for authenticating with the Azure service.  The `azidentity` module provides facilities for various ways of authenticating with Azure including client/secret, certificate, managed identity, and more.

Our default option is to use **DefaultAzureCredential** which will make use of the environment variables we have set and take care of the authentication flow for us.

```go
cred, err := azidentity.NewDefaultAzureCredential(nil)
```

For more details on how authentication works in `azidentity`, please see the documentation for `azidentity` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity).


Connecting to Azure 
-------------------

Once you have a credential, create a connection to the desired ARM endpoint.  The `armcore` module provides facilities for connecting with ARM endpoints including public and sovereign clouds as well as Azure Stack.

```go
con := armcore.NewDefaultConnection(cred, nil)
```

For more information on ARM connections, please see the documentation for `armcore` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/armcore](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/armcore).

Creating a Resource Management Client
-------------------------------------

Once you have a connection to ARM, you will need to decide what service to use and create a client to connect to that service. In this section, we will use `Compute` as our target service. The Compute modules consist of one or more clients. A client groups a set of related APIs, providing access to its functionality within the specified subscription. You will need to create one or more clients to access the APIs you require using your `armcore.Connection`.

To show an example, we will create a client to manage Virtual Machines. The code to achieve this task would be:

```go
client := armcompute.NewVirtualMachinesClient(con, "<subscription ID>")
```
You can use the same pattern to connect with other Azure services that you are using. For example, in order to manage Virtual Network resources, you would install the Network package and create a `VirtualNetwork` Client:

```go
client := armnetwork.NewVirtualNetworksClient(acon, "<subscription ID>")
```

Interacting with Azure Resources
--------------------------------

Now that we are authenticated and have created our sub-resource clients, we can use our client to make API calls. For resource management scenarios, most of our cases are centered around creating / updating / reading / deleting Azure resources. Those scenarios correspond to what we call "operations" in Azure. Once you are sure of which operations you want to call, you can then implement the operation call using the management client we just created in previous section.

To write the concrete code for the API call, you might need to look up the information of request parameters, types, and response body for a certain opertaion. We recommend using the following site for SDK reference:

- [Official Go docs for new Azure Go SDK packages](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk) - This web-site contains the complete SDK references for each released package as well as embedded code snippets for some operation

To see the reference for a certain package, you can either click into each package on the web-site, or directly add the SDK path to the end of URL. For example, to see the reference for Azure Compute package, you can use [https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute). Certain development tool or IDE has features that allow you to directly look up API definitions as well.

Let's illustrate the SDK usage by a few quick examples. In the following sample. we are going to create a resource group using the SDK. To achieve this scenario, we can take the follow steps

- **Step 1** : Decide which client we want to use, in our case, we know that it's related to Resource Group so our choice is the [ResourceGroupsClient](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resources/armresources#ResourceGroupsClient)
- **Step 2** : Find out which operation is responsible for creating a resource group. By locating the client in previous step, we are able to see all the functions under `ResourceGroupsClient`, and we can see [the `CreateOrUpdate` function](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resources/armresources#ResourceGroupsClient.CreateOrUpdate) is what need. 
- **Step 3** : Using the information about this operation, we can then fill in the required parameters, and implement it using the Go SDK. If we need extra information on what those parameters mean, we can also use the [Azure service documentation](https://docs.microsoft.com/azure/?product=featured) on Microsoft Docs

Let's show our what final code looks like

Example: Creating a Resource Group
---------------------------------

***Import the packages***
```go
import (
    "contexts"
    "github.com/Azure/azure-sdk-for-go/sdk/armcore"
    "github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/to"
)
```

***Define some global variables***
```go
var (
    ctx                 = context.Background()
    subscriptionId      = os.Getenv("AZURE_SUBSCRIPTION_ID")
    location            = "westus2"
    resourceGroupName   = "resourceGroupName"
    interval            = 5 * time.Second
)
```

***Write a function to create a resource group***
```go
func createResourceGroup(ctx context.Context, connection *armcore.Connection) (armresources.ResourceGroupResponse, error) {
	rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)

	param := armresources.ResourceGroup{
		Location: to.StringPtr(location),
	}

	return rgClient.CreateOrUpdate(context.Background(), resourceGroupName, param, nil)
}
```

***Invoking the `createResourceGroup` function in main***
```go
func main() {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        log.Fatalf("authentication failure: %+v", err)
    }
    conn := armcore.NewDefaultConnection(cred, &armcore.ConnectionOptions{
        Logging: azcore.LogOptions{
            IncludeBody: true,
        },
    })
    
    resourceGroup, err := createResourceGroup(ctx, conn)
    if err != nil {
        log.Fatalf("cannot create resource group: %+v", err)
    }
    log.Printf("Resource Group %s created", *resourceGroup.ResourceGroup.ID)
}
```

Let's demonstrate management client's usage by showing additional samples

Example: Managing Resource Groups
---------------------------------

***Update a resource group***

```go
func updateResourceGroup(ctx context.Context, connection *armcore.Connection) (armresources.ResourceGroupResponse, error) {
    rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)
    
    update := armresources.ResourceGroupPatchable{
        Tags: map[string]*string{
            "new": to.StringPtr("tag"),
        },
    }
    return rgClient.Update(ctx, resourceGroupName, update, nil)
}
```

***List all resource groups***

```go
func listResourceGroups(ctx context.Context, connection *armcore.Connection) ([]*armresources.ResourceGroup, error) {
    rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)
    
    pager := rgClient.List(nil)
    
    var resourceGroups []*armresources.ResourceGroup
    for pager.NextPage(ctx) {
        resp := pager.PageResponse()
        if resp.ResourceGroupListResult != nil {
            resourceGroups = append(resourceGroups, resp.ResourceGroupListResult.Value...)
        }
    }
    return resourceGroups, pager.Err()
}
```

***Delete a resource group***

```go
func deleteResourceGroup(ctx context.Context, connection *armcore.Connection) error {
    rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)
    
    poller, err := rgClient.BeginDelete(ctx, resourceGroupName, nil)
    if err != nil {
        return err
    }
    if _, err := poller.PollUntilDone(ctx, interval); err != nil {
        return err
    }
    return nil
}
```

***Invoking the update, list and delete of resource group in the main function***
```go
func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("authentication failure: %+v", err)
	}
	conn := armcore.NewDefaultConnection(cred, &armcore.ConnectionOptions{
		Logging: azcore.LogOptions{
			IncludeBody: true,
		},
	})


    resourceGroup, err := createResourceGroup(ctx, conn)
    if err != nil {
        log.Fatalf("cannot create resource group: %+v", err)
    }
    log.Printf("Resource Group %s created", *resourceGroup.ResourceGroup.ID)

	updatedRG, err := updateResourceGroup(ctx, conn)
	if err != nil {
        log.Fatalf("cannot update resource group: %+v", err)
	}
	log.Printf("Resource Group %s updated", *updatedRG.ResourceGroup.ID)

	rgList, err := listResourceGroups(ctx, conn)
	if err != nil {
        log.Fatalf("cannot list resource group: %+v", err)
	}
	log.Printf("We totally have %d resource groups", len(rgList))

	if err := deleteResourceGroup(ctx, conn); err != nil {
        log.Fatalf("cannot delete resource group: %+v", err)
	}
	log.Printf("Resource Group deleted")
})
```

Example: Managing Virtual Machines
---------------------------------
In addition to resource groups, we will also use Virtual Machine as an example and show how to manage how to create a Virtual Machine which involves three Azure services (Resource Group, Network and Compute)

Due to the complexity of this scenario, please [click here](https://aka.ms/azsdk/go/mgmt/samples) for the complete sample.

Long Running Operations
-----------------------
In the samples above, you might notice that some operations have a ``Begin`` prefix (for example, ``BeginDelete``). This indicates the operation is a Long-Running Operation (LRO). For resource management libraries, this kind of operation is quite common since certain resource operations may take a while to finish. When you need to use those LROs, you will need to use a poller and keep polling for the result until it is done. To illustrate this pattern, here is an example

```go
poller, err := client.BeginCreate(context.Background(), "resource_identifier", "additonal_parameter")
if err != nil {
	// handle error...
}
resp, err = poller.PollUntilDone(context.Background(), 5 * time.Second)
if err != nil {
	// handle error...
}
fmt.Printf("LRO done")
// dealing with `resp`
```

Note that you will need to pass a polling interval to ```PollUntilDone``` and tell the poller how often it should try to get the status. This number is usually small but it's best to consult the [Azure service documentation](https://docs.microsoft.com/azure/?product=featured) on best practices and recommdend intervals for your specific use cases.

For more advanced usage of LRO and design guidelines of LRO, please visit [this documentation here](https://azure.github.io/azure-sdk/golang_introduction.html#methods-invoking-long-running-operations)

## Code Samples

More code samples for using the management library for Go SDK can be found in the following locations
- [Go SDK Code Samples](https://aka.ms/azsdk/go/mgmt/samples)
- Example files under each package. For example, examples for Network packages can be [found here](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/network/armnetwork/example_networkinterfaces_test.go)

Need help?
----------

-   File an issue via [Github
    Issues](https://github.com/Azure/azure-sdk-for-go/issues)
-   Check [previous
    questions](https://stackoverflow.com/questions/tagged/azure+go)
    or ask new ones on StackOverflow using azure and Go tags.

Contributing
------------

For details on contributing to this repository, see the contributing
guide.

This project welcomes contributions and suggestions. Most contributions
require you to agree to a Contributor License Agreement (CLA) declaring
that you have the right to, and actually do, grant us the rights to use
your contribution. For details, visit <https://cla.microsoft.com>.

When you submit a pull request, a CLA-bot will automatically determine
whether you need to provide a CLA and decorate the PR appropriately
(e.g., label, comment). Simply follow the instructions provided by the
bot. You will only need to do this once across all repositories using
our CLA.

This project has adopted the Microsoft Open Source Code of Conduct. For
more information see the Code of Conduct FAQ or contact
[opencode@microsoft.com](mailto:opencode@microsoft.com) with any questions or comments.

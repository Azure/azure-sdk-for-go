
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
4.  Inside the System Properties window, click the Environment
    Variables… button.
5.  Click on the property you would like to change, then click the Edit…
    button. If the property name is not listed, then click the New…
    button.

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

Now that we are authenticated and have created our sub-resource clients, we can use our client to make API calls. Let's demonstrate management client's usage by showing various concrete examples

Example: Managing Resource Groups
---------------------------------
In this example, we will show how to manage Resource Groups.

***Import the packages***
```go
// insert code
```

***Create a resource group***

```go
// insert code
```

***Update a resource group***

```go
// insert code
```


***List all resource groups***

```go
// insert code
```

***Delete a resource group***

```go
// insert code
```

Example: Managing Network Resources
---------------------------------
In this example, we will show to manage Network related resources.

***Import the packages***
```go
// insert code
```
***Creating a Virtual Network***
```go
// insert code
```

***Updating a Virtual Network***
```go
// insert code
```

***List all Virtual Networks***
```go
// insert code
```

Example: Managing Virtual Machines
---------------------------------
In this example, we will show how to manage Virtual Machines

***Import the packages***
```go
// insert code
```

***Creating a Virtual Machine***
```go
// insert code
```

***Updating a Virtual Machine***
```go
// insert code
```

***List all Virtual Machines***
```go
// insert code
```

***Delete a Virtual Machine***
```go
// insert code
```

## Code Samples

More code samples for using the management library for Go SDK can be found in the following locations
- [Go SDK Code Samples](https://github.com/Azure-Samples/azure-sdk-for-go-samples)

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
<opencode@microsoft.com> with any additional questions or comments.

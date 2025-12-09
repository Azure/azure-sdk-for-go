// Copyright 2017 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

/*
Package aztables can access an Azure Storage or CosmosDB account.

The aztables package is capable of:
  - Creating, deleting, and listing tables in an account
  - Creating, deleting, updating, and querying entities in a table account
  - Creating Shared Access Signatures for authentication

# Creating the Client

The Azure Data Tables library allows you to interact with two types of resources:
* the tables in your account
* the entities within those tables.
Interaction with these resources starts with an instance of a client.
To create a client object, you will need the account's table service endpoint
URL and a credential that allows you to access the account.

	cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
	handle(err)
	serviceClient, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net/", cred, nil)
	handle(err)

# Types of Credentials

The clients support different forms of authentication. The aztables library supports
any of the `azcore.TokenCredential` interfaces, authorization via a Connection String,
or authorization with a Shared Access Signature token.

# Using a Shared Key

To use an account shared key (aka account key or access key), provide the key as a string.
This can be found in your storage account in the Azure Portal under the "Access Keys" section.

Use the key as the credential parameter to authenticate the client:

	cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
	handle(err)
	serviceClient, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net/", cred, nil)
	handle(err)

Using a Connection String
Depending on your use case and authorization method, you may prefer to initialize a client instance with a connection string instead of providing the account URL and credential separately. To do this, pass the
connection string to the client's `from_connection_string` class method. The connection string can be found in your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or with the following Azure CLI command:

	connStr := "DefaultEndpointsProtocol=https;AccountName=<my_account_name>;AccountKey=<my_account_key>;EndpointSuffix=core.windows.net"
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)

Using a Shared Access Signature
To use a shared access signature (SAS) token, provide the token at the end of your service URL.
You can generate a SAS token from the Azure Portal under Shared Access Signature or use the
ServiceClient.GetAccountSASToken or Client.GetTableSASToken() functions.

	cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
	handle(err)
	service, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net", cred, nil)
	handle(err)

	resources := aztables.AccountSASResourceTypes{Service: true}
	permission := aztables.AccountSASPermissions{Read: true}
	start := time.Date(2021, time.August, 21, 1, 1, 0, 0, time.UTC)
	expiry := time.Date(2022, time.August, 21, 1, 1, 0, 0, time.UTC)
	sasUrl, err := service.GetAccountSASToken(resources, permission, start, expiry)
	handle(err)

	sasService, err := aztables.NewServiceClient(sasUrl, azcore.AnonymousCredential(), nil)
	handle(err)

# Key Concepts

Common uses of the Table service included:
* Storing TBs of structured data capable of serving web scale applications
* Storing datasets that do not require complex joins, foreign keys, or stored procedures and can be de-normalized for fast access
* Quickly querying data using a clustered index
* Accessing data using the OData protocol and LINQ filter expressions

The following components make up the Azure Data Tables Service:
* The account
* A table within the account, which contains a set of entities
* An entity within a table, as a dictionary

The Azure Data Tables client library for Go allows you to interact with each of these components
through the use of a dedicated client object.

Two different clients are provided to interact with the various components of the Table Service:
1. **`ServiceClient`** -
  - Get and set account setting
  - Query, create, and delete tables within the account.
  - Get a `Client` to access a specific table using the `NewClient` method.

2. **`Client`** -
  - Interacts with a specific table (which need not exist yet).
  - Create, delete, query, and upsert entities within the specified table.
  - Create or delete the specified table itself.

Entities are similar to rows. An entity has a PartitionKey, a RowKey, and a set of properties.
A property is a name value pair, similar to a column. Every entity in a table does not need to
have the same properties. Entities are returned as JSON, allowing developers to use JSON
marshalling and unmarshalling techniques. Additionally, you can use the aztables.EDMEntity to
ensure proper round-trip serialization of all properties.

	aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "pencils",
			RowKey: "id-003",
		},
		Properties: map[string]any{
			"Product": "Ticonderoga Pencils",
			"Price": 5.00,
			"Count": aztables.EDMInt64(12345678901234),
			"ProductGUID": aztables.EDMGUID("some-guid-value"),
			"DateReceived": aztables.EDMDateTime(time.Date{....}),
			"ProductCode": aztables.EDMBinary([]byte{"somebinaryvalue"})
		}
	}

# More Examples

The following sections provide several code snippets covering some of the most common Table tasks, including:

* Creating a table
* Creating entities
* Querying entities

# Creating a Table

Create a table in your account and get a `Client` to perform operations on the newly created table:

	cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
	handle(err)
	service, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net", cred, nil)
	handle(err)
	resp, err := service.CreateTable("myTable")

Creating Entities

	cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
	handle(err)
	service, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net", cred, nil)
	handle(err)

	myEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "001234",
			RowKey: "RedMarker",
		},
		Properties: map[string]any{
			"Stock": 15,
			"Price": 9.99,
			"Comments": "great product",
			"OnSale": true,
			"ReducedPrice": 7.99,
			"PurchaseDate": aztables.EDMDateTime(time.Date(2021, time.August, 21, 1, 1, 0, 0, time.UTC)),
			"BinaryRepresentation": aztables.EDMBinary([]byte{"Bytesliceinfo"})
		}
	}
	marshalled, err := json.Marshal(myEntity)
	handle(err)

	client, err := service.NewClient("myTable")
	handle(err)

	resp, err := client.AddEntity(context.Background(), marshalled, nil)
	handle(err)

Querying entities

	cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
	handle(err)
	client, err := aztables.NewClient("https://myAccountName.table.core.windows.net/myTableName", cred, nil)
	handle(err)

	filter := "PartitionKey eq 'markers' or RowKey eq 'id-001'"
	options := &ListEntitiesOptions{
		Filter: &filter,
		Select: to.Ptr("RowKey,Value,Product,Available"),
		Top: to.Ptr(int32(((15),
	}

	pager := client.List(options)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		fmt.Printf("Received: %v entities\n", len(resp.Entities))

		for _, entity := range resp.Entities {
			var myEntity aztables.EDMEntity
			err = json.Unmarshal(entity, &myEntity)
			handle(err)

			fmt.Printf("Received: %v, %v, %v, %v\n", myEntity.Properties["RowKey"], myEntity.Properties["Value"], myEntity.Properties["Product"], myEntity.Properties["Available"])
		}
	}

	if pager.Err() != nil {
		// handle error...
	}
*/
package aztables

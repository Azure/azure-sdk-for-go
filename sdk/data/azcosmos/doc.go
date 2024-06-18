// Copyright 2021 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

/*
Package azcosmos implements the client to interact with the Azure Cosmos DB SQL API.

The azcosmos package is capable of:
  - Creating, deleting, and reading databases in an account
  - Creating, deleting, updating, and reading containers in a database
  - Creating, deleting, replacing, upserting, and reading items in a container

# Creating the Client

Types of Credentials
The clients support different forms of authentication. The azcosmos library supports
authorization via Azure Active Directory or an account key.

Using Azure Active Directory
To create a client, you can use any of the TokenCredential implementations provided by `azidentity`.

	cred, err := azidentity.NewClientSecretCredential("tenantId", "clientId", "clientSecret")
	handle(err)
	client, err := azcosmos.NewClient("myAccountEndpointURL", cred, nil)
	handle(err)

Using account keys
To create a client, you will need the account's endpoint URL and a key credential.

	cred, err := azcosmos.NewKeyCredential("myAccountKey")
	handle(err)
	client, err := azcosmos.NewClientWithKey("myAccountEndpointURL", cred, nil)
	handle(err)

Using connection string
To create a client, you will need the account's connection string.

	client, err := azcosmos.NewClientFromConnectionString("myConnectionString", nil)
	handle(err)

# Key Concepts

The following are relevant concepts for the usage of the client:
  - A client is a connection to an Azure Cosmos DB account.
  - An account can have multiple databases, and the client allows you to create, read, and delete databases.
  - A database can have multiple containers, and the client allows you to create, read, update, and delete containers, and to modify throughput provision.
  - Information is stored as items inside containers and the client allows you to create, read, update, and delete items in containers.

# More Examples

The following sections provide several code snippets covering some of the most common Table tasks, including:
  - Creating a database
  - Creating a container
  - Creating, reading, and deleting items
  - Querying items
  - Using Transactional Batch

# Creating a database

Create a database and obtain a `DatabaseClient` to perform operations on your newly created database.

	cred, err := azcosmos.NewKeyCredential("myAccountKey")
	handle(err)
	client, err := azcosmos.NewClientWithKey("myAccountEndpointURL", cred, nil)
	handle(err)
	databaseProperties := azcosmos.DatabaseProperties{ID: "myDatabase"}
	response, err := client.CreateDatabase(context, databaseProperties, nil)
	handle(err)
	database, err := azcosmos.NewDatabase("myDatabase")
	handle(err)

# Creating a container

Create a container on an existing database and obtain a `ContainerClient` to perform operations on your newly created container.

	cred, err := azcosmos.NewKeyCredential("myAccountKey")
	handle(err)
	client, err := azcosmos.NewClientWithKey("myAccountEndpointURL", cred, nil)
	handle(err)
	database := azcosmos.NewDatabase("myDatabase")
	properties := azcosmos.ContainerProperties{
		ID: "myContainer",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/myPartitionKeyProperty"},
		},
	}

	throughput := azcosmos.NewManualThroughputProperties(400)
	response, err := database.CreateContainer(context, properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	handle(err)
	container, err := database.NewContainer("myContainer")
	handle(err)

Creating, reading, and deleting items

	item := map[string]string{
		"id":    "1",
		"myPartitionKeyProperty": "myPartitionKeyValue",
		"otherValue": 10
	}
	marshalled, err := json.Marshal(item)
	handle(err)

	pk := azcosmos.NewPartitionKeyString("myPartitionKeyValue")
	itemResponse, err := container.CreateItem(context, pk, marshalled, nil)
	handle(err)

	id := "1"
	itemResponse, err = container.ReadItem(context, pk, id, nil)
	handle(err)

	var itemResponseBody map[string]string
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	handle(err)

	itemResponseBody["value"] = "3"
	marshalledReplace, err := json.Marshal(itemResponseBody)
	handle(err)

	itemResponse, err = container.ReplaceItem(context, pk, id, marshalledReplace, nil)
	handle(err)

	itemResponse, err = container.DeleteItem(context, pk, id, nil)
	handle(err)

Querying items

	pk := azcosmos.NewPartitionKeyString("myPartitionKeyValue")
	queryPager := container.NewQueryItemsPager("select * from docs c", pk, nil)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context)
		if err != nil {
			handle(err)
		}

		for _, item := range queryResponse.Items {
			var itemResponseBody map[string]interface{}
			json.Unmarshal(item, &itemResponseBody)
		}
	}

Querying items with parametrized queries

	opt := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{"@value", "2"},
		},
	}
	pk := azcosmos.NewPartitionKeyString("myPartitionKeyValue")
	queryPager := container.NewQueryItemsPager("select * from docs c where c.value = @value", pk, &opt)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context)
		if err != nil {
			handle(err)
		}

		for _, item := range queryResponse.Items {
			var itemResponseBody map[string]interface{}
			json.Unmarshal(item, &itemResponseBody)
		}
	}

Using Transactional batch

	pk := azcosmos.NewPartitionKeyString("myPartitionKeyValue")
	batch := container.NewTransactionalBatch(pk)

	item := map[string]string{
		"id":    "1",
		"myPartitionKeyProperty": "myPartitionKeyValue",
		"otherValue": 10
	}
	marshalled, err := json.Marshal(item)
	handle(err)

	batch.CreateItem(marshalled, nil)
	batch.ReadItem("otherExistingId", nil)
	batch.DeleteItem("yetAnotherExistingId", nil)

	batchResponse, err  := container.ExecuteTransactionalBatch(context, batch, nil)
	handle(err)

	if batchResponse.Success {
		// Transaction succeeded
		// We can inspect the individual operation results
		for index, operation := range batchResponse.OperationResults {
			fmt.Printf("Operation %v completed with status code %v consumed %v RU", index, operation.StatusCode, operation.RequestCharge)
			if index == 1 {
				// Read operation would have body available
				var itemResponseBody map[string]string
				err = json.Unmarshal(operation.ResourceBody, &itemResponseBody)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		// Transaction failed, look for the offending operation
		for index, operation := range batchResponse.OperationResults {
			if operation.StatusCode != http.StatusFailedDependency {
				fmt.Printf("Transaction failed due to operation %v which failed with status code %v", index, operation.StatusCode)
			}
		}
	}
*/
package azcosmos

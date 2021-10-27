// Copyright 2021 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

/*
Package azcosmos implements the client to interact with the Azure Cosmos DB SQL API.

The azcosmos package is capable of:
	- Creating, deleting, and reading databases in an account
	- Creating, deleting, updating, and reading containers in a database
	- Creating, deleting, replacing, upserting, and reading items in a container

Creating the Client

To create a client, you will need the account's endpoint URL and a key credential.

	cred, err := azcosmos.NewKeyCredential("myAccountKey")
	handle(err)
	client, err := azcosmos.NewClientWithKey("myAccountEndpointURL", cred, nil)
	handle(err)


Key Concepts

The following are relevant concepts for the usage of the client:
	- A client is a connection to an Azure Cosmos DB account.
	- An account can have multiple databases, and the client allows you to create, read, and delete databases.
	- A database can have multiple containers, and the client allows you to create, read, update, and delete containers, and to modify throughput provision.
	- Information is stored as items inside containers and the client allows you to create, read, update, and delete items in containers.


More Examples

The following sections provide several code snippets covering some of the most common Table tasks, including:
	- Creating a database
	- Creating a container
	- Creating, reading, and deleting items


Creating a database

Create a database and obtain a `DatabaseClient` to perform operations on your newly created database.

	cred, err := azcosmos.NewKeyCredential("myAccountKey")
	handle(err)
	client, err := azcosmos.NewClientWithKey("myAccountEndpointURL", cred, nil)
	handle(err)
	database := azcosmos.DatabaseProperties{ID: "myDatabase"}
	response, err := client.CreateDatabase(context, database, nil)
	handle(err)
	database, err := azcosmos.NewDatabase("myDatabase")
	handle(err)


Creating a container

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
*/
package azcosmos

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package checkpointstore_test

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpointstore"
)

func Example() {
	cs := os.Getenv("CHECKPOINTSTORE_STORAGE_CONNECTION_STRING")
	containerName := os.Getenv("CHECKPOINTSTORE_STORAGE_CONTAINER_NAME")

	// Create the checkpoint store
	// NOTE: the container you pass in 'containerName' must already be created before the checkpoint
	// store starts.
	checkpointStore, err := checkpointstore.NewFromConnectionString(cs, containerName, nil)

	if err != nil {
		panic(err)
	}

	_, err = checkpointStore.ClaimOwnership(context.TODO(), nil, nil)

	if err != nil {
		panic(err)
	}

	// consumers := sync.Map{}
}

// func DistributedConsumer(chan<- ) {

// }

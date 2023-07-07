// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
)

func ExampleNewClientWithSharedKeyCredential() {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	sharedKey := os.Getenv("EVENTGRID_KEY")

	client, err := azeventgrid.NewClientWithSharedKeyCredential(endpoint, sharedKey, nil)

	if err != nil {
		panic(err)
	}

	_ = client // ignore
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func CreateDeleteFromServiceClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", accountName)

	service, err := aztables.NewServiceClient(serviceURL, cred, nil)
	check(err)

	// Create a table
	_, err = service.CreateTable(context.Background(), "fromServiceClient", nil)
	check(err)

	// Delete a table
	_, err = service.DeleteTable(context.Background(), "fromServiceClient", nil)
	check(err)
}

func CreateDeleteFromTableClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "fromTableClient")
	client, err := aztables.NewClient(serviceURL, cred, nil)
	check(err)

	// Create a table
	_, err = client.Create(context.Background(), nil)
	check(err)

	// Delete a table
	_, err = client.Delete(context.Background(), nil)
	check(err)
}

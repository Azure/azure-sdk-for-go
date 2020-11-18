// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/storage/2019-06-01/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleUsagesOperations_ListByLocation() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewUsagesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.ListByLocation(context.Background(), "<Azure location>", nil)
	if err != nil {
		log.Fatalf("failed to delete account: %v", err)
	}
	for _, u := range *resp.UsageListResult.Value {
		log.Printf("usage: %v, limit: %v, current value: %v", *u.Name.Value, *u.Limit, *u.CurrentValue)
	}
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpostgresqlflexibleservers_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v4"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/10925e3dec73699b950f256576cd6983947faaa3/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/LogFilesListByServer.json
func ExampleLogFilesClient_NewListByServerPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewLogFilesClient().NewListByServerPager("testrg", "postgresqltestsvc1", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.LogFileListResult = armpostgresqlflexibleservers.LogFileListResult{
		// 	Value: []*armpostgresqlflexibleservers.LogFile{
		// 		{
		// 			Name: to.Ptr("postgresql-slow-postgresqltestsvc1-2018022823.log"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/flexibleServers/logFiles"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/flexibleServers/postgresqltestsvc1/logFiles/postgresql-slow-postgresqltestsvc1-2018022823.log"),
		// 			Properties: &armpostgresqlflexibleservers.LogFileProperties{
		// 				Type: to.Ptr("slowlog"),
		// 				CreatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "1-01-01T00:00:00.000Z"); return t}()),
		// 				LastModifiedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-01T06:09:20.000Z"); return t}()),
		// 				SizeInKb: to.Ptr[int64](1),
		// 				URL: to.Ptr("https://wasd2prodwus1afse42.file.core.windows.net/833c99b2f36c47349e5554b903fe0440/serverlogs/postgresql-slow-postgresqltestsvc1-2018022823.log?sv=2015-04-05&sr=f&se=2018-03-01T07%3A12%3A13Z&sp=r"),
		// 			},
		// 	}},
		// }
	}
}

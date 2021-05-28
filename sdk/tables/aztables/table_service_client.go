// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	LegacyCosmosTableDomain = ".table.cosmosdb."
	CosmosTableDomain       = ".table.cosmos."
)

// A TableServiceClient represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type TableServiceClient struct {
	client  *tableClient
	service *serviceClient
	cred    SharedKeyCredential
}

// NewTableServiceClient creates a TableClient object using the specified URL and request policy pipeline.
func NewTableServiceClient(serviceURL string, cred azcore.Credential, options TableClientOptions) (*TableServiceClient, error) {
	conOptions := options.getConnectionOptions()
	if isCosmosEndpoint(serviceURL) {
		conOptions.PerCallPolicies = []azcore.Policy{CosmosPatchTransformPolicy{}}
	}
	con := newConnection(serviceURL, cred, conOptions)
	c, _ := cred.(*SharedKeyCredential)
	return &TableServiceClient{client: &tableClient{con}, service: &serviceClient{con}, cred: *c}, nil
}

// Gets a TableClient affinitzed to the specified table name and initialized with the same serviceURL and credentials as this TableServiceClient
func (t *TableServiceClient) NewTableClient(tableName string) *TableClient {
	return &TableClient{client: t.client, cred: t.cred, Name: tableName, service: t}
}

// Creates a table with the specified name
func (t *TableServiceClient) Create(ctx context.Context, name string) (TableResponseResponse, error) {
	resp, err := t.client.Create(ctx, TableProperties{&name}, new(TableCreateOptions), new(QueryOptions))
	if err == nil {
		tableResp := resp.(TableResponseResponse)
		return tableResp, nil
	}
	return TableResponseResponse{}, err
}

// Deletes a table by name
func (t *TableServiceClient) Delete(ctx context.Context, name string) (TableDeleteResponse, error) {
	return t.client.Delete(ctx, name, nil)
}

// Queries the tables using the specified QueryOptions
func (t *TableServiceClient) QueryTables(queryOptions QueryOptions) TableQueryResponsePager {
	return &tableQueryResponsePager{client: t.client, queryOptions: &queryOptions, tableQueryOptions: new(TableQueryOptions)}
}

func isCosmosEndpoint(url string) bool {
	isCosmosEmulator := strings.Index(url, "localhost") >= 0 && strings.Index(url, "8902") >= 0
	return isCosmosEmulator ||
		strings.Index(url, CosmosTableDomain) >= 0 ||
		strings.Index(url, LegacyCosmosTableDomain) >= 0
}

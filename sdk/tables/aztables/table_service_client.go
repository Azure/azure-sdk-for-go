// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
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
func NewTableServiceClient(serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableServiceClient, error) {
	var con *connection
	conOptions := options.getConnectionOptions()
	if isCosmosEndpoint(serviceURL) {
		p := azcore.NewPipeline(options.HTTPClient,
			azcore.NewTelemetryPolicy(conOptions.telemetryOptions()),
			CosmosPatchTransformPolicy{},
			azcore.NewRetryPolicy(&options.Retry),
			cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
			azcore.NewLogPolicy(&conOptions.Logging))
		con = newConnectionWithPipeline(serviceURL, p)

	} else {
		con = newConnection(serviceURL, cred, conOptions)
	}
	c, _ := cred.(*SharedKeyCredential)
	return &TableServiceClient{client: &tableClient{con}, service: &serviceClient{con}, cred: *c}, nil
}

// Gets a TableClient affinitzed to the specified table name and initialized with the same serviceURL and credentials as this TableServiceClient
func (t *TableServiceClient) GetTableClient(tableName string) *TableClient {
	return &TableClient{client: t.client, cred: t.cred, name: tableName, service: t}
}

// Creates a table with the specified name
func (t *TableServiceClient) Create(ctx context.Context, name string) (*TableResponseResponse, *runtime.ResponseError) {
	var r *TableResponseResponse = nil
	resp, err := t.client.Create(ctx, TableProperties{&name}, new(TableCreateOptions), new(QueryOptions))
	if err == nil {
		tableResp := resp.(TableResponseResponse)
		r = &tableResp
	}
	return r, convertErr(err)
}

// Deletes a table by name
func (t *TableServiceClient) Delete(ctx context.Context, name string) (*TableDeleteResponse, *runtime.ResponseError) {
	resp, err := t.client.Delete(ctx, name, nil)
	return &resp, convertErr(err)
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

func convertErr(err error) *runtime.ResponseError {
	var e *runtime.ResponseError
	if err == nil || !errors.As(err, &e) {
		return nil
	} else {
		return e
	}
}

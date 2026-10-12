// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestCreateDatabaseRejectsEmptyID(t *testing.T) {
	client := newTestClient(t)
	_, err := client.CreateDatabase(context.Background(), DatabaseProperties{}, nil)
	require.ErrorContains(t, err, "DatabaseProperties.ID")
}

func TestCreateContainerRejectsEmptyID(t *testing.T) {
	client := newTestClient(t)
	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	_, err = database.CreateContainer(context.Background(), ContainerProperties{
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}},
	}, nil)
	require.ErrorContains(t, err, "ContainerProperties.ID")
}

func TestCreateContainerRejectsMissingPartitionKeyPaths(t *testing.T) {
	client := newTestClient(t)
	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	_, err = database.CreateContainer(context.Background(), ContainerProperties{ID: "items"}, nil)
	require.ErrorContains(t, err, "PartitionKeyDefinition.Paths")
}

func TestManagementReadsAllowTokenCredentialClient(t *testing.T) {
	client, err := NewClient("https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	container, err := database.NewContainer("items")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = database.Read(ctx, nil)
	require.ErrorIs(t, err, context.Canceled)

	_, err = container.Read(ctx, nil)
	require.ErrorIs(t, err, context.Canceled)
}

func TestManagementMutationsRejectTokenCredentialClient(t *testing.T) {
	client, err := NewClient("https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	container, err := database.NewContainer("items")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	validContainer := ContainerProperties{ID: "items", PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}}}

	_, err = client.CreateDatabase(ctx, DatabaseProperties{ID: "db"}, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = database.Delete(ctx, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = database.CreateContainer(ctx, validContainer, nil)
	require.ErrorContains(t, err, "NewClientWithKey")

	_, err = container.Delete(ctx, nil)
	require.ErrorContains(t, err, "NewClientWithKey")
}

func TestDatabaseResponseParsesProperties(t *testing.T) {
	response, err := newDatabaseResponse(ItemResponse{}, []byte(`{"id":"db","_etag":"\"abc\"","_rid":"rid","_self":"self","_ts":1735689600}`))
	require.NoError(t, err)
	require.Equal(t, "db", response.DatabaseProperties.ID)
	require.Equal(t, "rid", response.DatabaseProperties.ResourceID)
	require.NotNil(t, response.DatabaseProperties.ETag)
	require.Equal(t, "\"abc\"", string(*response.DatabaseProperties.ETag))
	require.Equal(t, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), response.DatabaseProperties.LastModified)
}

func TestDatabaseResponseEmptyBodyIsZeroValue(t *testing.T) {
	response, err := newDatabaseResponse(ItemResponse{Response: Response{StatusCode: 204}}, nil)
	require.NoError(t, err)
	require.Equal(t, DatabaseProperties{}, response.DatabaseProperties)
	require.Equal(t, 204, response.StatusCode)
}

func TestContainerResponseParsesProperties(t *testing.T) {
	body := []byte(`{
		"id":"items",
		"partitionKey":{"paths":["/pk"],"kind":"Hash","version":2},
		"indexingPolicy":{"automatic":true,"indexingMode":"consistent"},
		"uniqueKeyPolicy":{"uniqueKeys":[{"paths":["/email"]}]},
		"_ts":1735689600
	}`)
	response, err := newContainerResponse(ItemResponse{}, body)
	require.NoError(t, err)
	require.Equal(t, "items", response.ContainerProperties.ID)
	require.Equal(t, []string{"/pk"}, response.ContainerProperties.PartitionKeyDefinition.Paths)
	require.Equal(t, 2, response.ContainerProperties.PartitionKeyDefinition.Version)
	require.Equal(t, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), response.ContainerProperties.LastModified)
	require.NotNil(t, response.ContainerProperties.IndexingPolicy)
	require.True(t, response.ContainerProperties.IndexingPolicy.Automatic)
	require.NotNil(t, response.ContainerProperties.UniqueKeyPolicy)
	require.Equal(t, []string{"/email"}, response.ContainerProperties.UniqueKeyPolicy.UniqueKeys[0].Paths)
}

func TestPartitionKeyDefinitionMarshalInfersKindFromPathCount(t *testing.T) {
	tests := []struct {
		name        string
		def         PartitionKeyDefinition
		wantKind    PartitionKeyKind
		wantVersion int
	}{
		{"explicit values are preserved", PartitionKeyDefinition{Kind: PartitionKeyKindMultiHash, Paths: []string{"/a"}, Version: 1}, PartitionKeyKindMultiHash, 1},
		{"single path infers hash and version two", PartitionKeyDefinition{Paths: []string{"/pk"}}, PartitionKeyKindHash, 2},
		{"multiple paths infer multi-hash and version two", PartitionKeyDefinition{Paths: []string{"/a", "/b"}}, PartitionKeyKindMultiHash, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.def)
			require.NoError(t, err)

			var decoded struct {
				Kind    PartitionKeyKind `json:"kind"`
				Version int              `json:"version"`
			}
			require.NoError(t, json.Unmarshal(data, &decoded))
			require.Equal(t, tt.wantKind, decoded.Kind)
			require.Equal(t, tt.wantVersion, decoded.Version)
		})
	}
}

func TestManagementResponseDecodeFailurePreservesMetadata(t *testing.T) {
	diagnostics := &Diagnostics{AttemptCount: 3}
	nativeResponse := ItemResponse{
		Response: Response{
			Diagnostics:   diagnostics,
			StatusCode:    201,
			SubStatus:     42,
			AttemptCount:  3,
			RequestCharge: 4.5,
			ActivityID:    "activity",
		},
		SessionToken: "session",
		ETag:         azcore.ETag(`"etag"`),
	}
	body := []byte(`{"id":`)

	response, err := newDatabaseResponse(nativeResponse, body)
	require.Equal(t, DatabaseResponse{}, response)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
	require.Same(t, diagnostics, cosmosErr.Diagnostics)
	require.Equal(t, 201, cosmosErr.StatusCode)
	require.Equal(t, 42, cosmosErr.SubStatus)
	require.Equal(t, uint32(3), cosmosErr.AttemptCount)
	require.Equal(t, 4.5, cosmosErr.RequestCharge)
	require.Equal(t, "activity", cosmosErr.ActivityID)
	require.Equal(t, SessionToken("session"), cosmosErr.SessionToken)
	require.Equal(t, azcore.ETag(`"etag"`), cosmosErr.ETag)
	require.Equal(t, body, cosmosErr.Body)
	require.True(t, cosmosErr.FromWire)
	var syntaxErr *json.SyntaxError
	require.True(t, errors.As(err, &syntaxErr))

	containerResponse, err := newContainerResponse(nativeResponse, body)
	require.Equal(t, ContainerResponse{}, containerResponse)
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
	require.Equal(t, body, cosmosErr.Body)
}

func TestForceContentResponseOnWriteOverridesFalse(t *testing.T) {
	disabled := false
	options := OperationOptions{EnableContentResponseOnWrite: &disabled}

	forceContentResponseOnWrite(&options)

	require.NotNil(t, options.EnableContentResponseOnWrite)
	require.True(t, *options.EnableContentResponseOnWrite)
}

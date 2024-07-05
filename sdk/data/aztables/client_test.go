// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

var services = []string{"storage", "cosmos"}

func TestServiceErrors(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name:   "Client.CreateTable",
				Status: tracing.SpanStatusError,
			}))
			defer delete()

			// Create a duplicate table to produce an error
			_, err := client.CreateTable(ctx, nil)
			require.Error(t, err)
			var httpErr *azcore.ResponseError
			require.ErrorAs(t, err, &httpErr)
			require.Equal(t, string(TableAlreadyExists), httpErr.ErrorCode)
			require.Contains(t, PossibleTableErrorCodeValues(), TableErrorCode(httpErr.ErrorCode))
		})
	}
}

func TestCreateTable(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, false, NewSpanValidator(t, SpanMatcher{
				Name: "Client.Delete",
			}))
			defer delete()

			_, err := client.CreateTable(ctx, nil)

			require.NoError(t, err)
		})
	}
}

type mdforAddGet struct {
	Metadata string `json:"odata.metadata"`
	Type     string `json:"odata.type"` // only for full metadata
}

func TestAddEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name: "Client.AddEntity",
			}))
			defer delete()

			simpleEntity := createSimpleEntity(1, "partition")

			marshalledEntity, err := json.Marshal(simpleEntity)
			require.NoError(t, err)
			resp, err := client.AddEntity(ctx, marshalledEntity, nil)
			require.NoError(t, err)
			require.NotEmpty(t, resp.Value)
			var md mdforAddGet
			require.NoError(t, json.Unmarshal(resp.Value, &md))
			require.NotEmpty(t, md.Metadata)
			require.Empty(t, md.Type)
		})
	}
}

func TestAddComplexEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			entity := createComplexEntity(1, "partition")

			marshalledEntity, err := json.Marshal(entity)
			require.NoError(t, err)

			resp, err := client.AddEntity(ctx, marshalledEntity, &AddEntityOptions{
				Format: to.Ptr(MetadataFormatFull),
			})
			require.NoError(t, err)
			require.NotEmpty(t, resp.Value)
			var md mdforAddGet
			require.NoError(t, json.Unmarshal(resp.Value, &md))
			require.NotEmpty(t, md.Metadata)
			if service == "storage" {
				// cosmos doesn't send full metadata
				require.NotEmpty(t, md.Type)
			}
		})
	}
}

func TestDeleteEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name: "Client.DeleteEntity",
			}))
			defer delete()

			simpleEntity := createSimpleEntity(1, "partition")

			marshalledEntity, err := json.Marshal(simpleEntity)
			require.NoError(t, err)
			_, err = client.AddEntity(ctx, marshalledEntity, nil)
			require.NoError(t, err)
			_, delErr := client.DeleteEntity(ctx, simpleEntity.PartitionKey, simpleEntity.RowKey, nil)
			require.Nil(t, delErr)
		})
	}
}

func TestDeleteEntityWithETag(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			simpleEntity := createSimpleEntity(1, "partition")
			simpleEntity2 := createSimpleEntity(2, "partition")

			marshalledEntity, err := json.Marshal(simpleEntity)
			require.NoError(t, err)
			resp, err := client.AddEntity(ctx, marshalledEntity, nil)
			require.NoError(t, err)
			oldETag := resp.ETag

			marshalledEntity, err = json.Marshal(simpleEntity2)
			require.NoError(t, err)
			resp, err = client.AddEntity(ctx, marshalledEntity, nil)
			require.NoError(t, err)
			newETag := resp.ETag

			_, err = client.DeleteEntity(ctx, simpleEntity2.PartitionKey, simpleEntity2.RowKey, &DeleteEntityOptions{IfMatch: &oldETag})
			require.Error(t, err)
			var httpErr *azcore.ResponseError
			require.ErrorAs(t, err, &httpErr)
			require.Contains(t, PossibleTableErrorCodeValues(), TableErrorCode(httpErr.ErrorCode))

			_, err = client.DeleteEntity(ctx, simpleEntity.PartitionKey, simpleEntity.RowKey, &DeleteEntityOptions{IfMatch: &oldETag})
			require.NoError(t, err)

			_, err = client.DeleteEntity(ctx, simpleEntity2.PartitionKey, simpleEntity2.RowKey, &DeleteEntityOptions{IfMatch: &newETag})
			require.NoError(t, err)
		})
	}
}

func TestMergeEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name: "Client.GetEntity",
			}))
			defer delete()

			entityToCreate := createSimpleEntity(1, "partition")
			marshalled, err := json.Marshal(entityToCreate)
			require.NoError(t, err)

			_, err = client.AddEntity(ctx, marshalled, nil)
			require.NoError(t, err)

			filter := "RowKey eq '1'"
			listOptions := &ListEntitiesOptions{Filter: &filter}

			preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey, nil)
			require.NoError(t, err)
			require.NotEmpty(t, preMerge.Value)
			var md mdforAddGet
			require.NoError(t, json.Unmarshal(preMerge.Value, &md))
			require.NotEmpty(t, md.Metadata)
			require.Empty(t, md.Type)

			var unMarshalledPreMerge map[string]any
			err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
			require.NoError(t, err)

			var mapEntity map[string]any
			err = json.Unmarshal(marshalled, &mapEntity)
			require.NoError(t, err)
			mapEntity["MergeProperty"] = "foo"

			reMarshalled, err := json.Marshal(mapEntity)
			require.NoError(t, err)

			_, updateErr := client.UpdateEntity(ctx, reMarshalled, &UpdateEntityOptions{UpdateMode: UpdateModeMerge})
			require.Nil(t, updateErr)

			var qResp ListEntitiesResponse
			pager := client.NewListEntitiesPager(listOptions)
			for pager.More() {
				qResp, err = pager.NextPage(ctx)
				require.NoError(t, err)
			}
			require.NotEmpty(t, qResp.Entities)
			postMerge := qResp.Entities[0]
			var unmarshalledPostMerge map[string]any
			err = json.Unmarshal(postMerge, &unmarshalledPostMerge)
			require.NoError(t, err)

			require.Equal(t, unmarshalledPostMerge["PartitionKey"], unMarshalledPreMerge["PartitionKey"])
			require.Equal(t, unmarshalledPostMerge["MergeProperty"], "foo")

			_, ok := unMarshalledPreMerge["MergeProperty"]
			require.False(t, ok)
		})
	}
}

func TestMergeEntityDoesNotExist(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name:   "Client.UpdateEntity",
				Status: tracing.SpanStatusError,
			}))
			defer delete()

			entityToCreate := createSimpleEntity(1, "partition")
			marshalled, err := json.Marshal(entityToCreate)
			require.NoError(t, err)

			_, updateErr := client.UpdateEntity(ctx, marshalled, &UpdateEntityOptions{UpdateMode: UpdateModeMerge})
			require.Error(t, updateErr)
			var httpErr *azcore.ResponseError
			require.ErrorAs(t, updateErr, &httpErr)
			require.Equal(t, string(ResourceNotFound), httpErr.ErrorCode)
			require.Contains(t, PossibleTableErrorCodeValues(), TableErrorCode(httpErr.ErrorCode))
		})
	}
}

func TestInsertEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name: "Client.UpsertEntity",
			}))
			defer delete()

			// 1. Create Basic Entity
			entityToCreate := createSimpleEntityWithRowKey(1, "parti'tion", "one'")
			marshalled, err := json.Marshal(entityToCreate)
			require.NoError(t, err)

			_, err = client.UpsertEntity(ctx, marshalled, &UpsertEntityOptions{UpdateMode: UpdateModeReplace})
			require.NoError(t, err)

			filter := "RowKey eq '1'"
			list := &ListEntitiesOptions{Filter: &filter}

			// 2. Query for basic Entity
			preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey, &GetEntityOptions{
				Format: to.Ptr(MetadataFormatFull),
			})
			require.NoError(t, err)
			require.NotEmpty(t, preMerge.Value)
			var md mdforAddGet
			require.NoError(t, json.Unmarshal(preMerge.Value, &md))
			require.NotEmpty(t, md.Metadata)
			if service == "storage" {
				// cosmos doesn't send full metadata
				require.NotEmpty(t, md.Type)
			}

			var unMarshalledPreMerge map[string]any
			err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
			require.NoError(t, err)

			// 3. Create same entity without Bool property, add "MergeProperty" prop
			mapEntity := createSimpleEntityNoBool(1, "partition")
			mapEntity["MergeProperty"] = "foo"

			reMarshalled, err := json.Marshal(mapEntity)
			require.NoError(t, err)

			// 4. Replace Entity with "bool"-less entity
			_, err = client.UpsertEntity(ctx, reMarshalled, &UpsertEntityOptions{UpdateMode: UpdateModeReplace})
			require.Nil(t, err)

			// 5. Query for new entity
			var qResp ListEntitiesResponse
			pager := client.NewListEntitiesPager(list)
			for pager.More() {
				qResp, err = pager.NextPage(ctx)
				require.NoError(t, err)
			}
			postMerge := qResp.Entities[0]
			var unmarshalledPostMerge map[string]any
			err = json.Unmarshal(postMerge, &unmarshalledPostMerge)
			require.NoError(t, err)

			// 6. Make assertions
			require.Less(t, len(unmarshalledPostMerge), len(unMarshalledPreMerge))
			require.Equal(t, unmarshalledPostMerge["MergeProperty"], "foo")

			_, ok := unmarshalledPostMerge["Bool"]
			require.Falsef(t, ok, "Bool property should not be available in the merged entity")
		})
	}
}
func TestInsertEntityTwice(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			// 1. Create Basic Entity
			entityToCreate := createSimpleEntity(1, "partition")
			marshalled, err := json.Marshal(entityToCreate)
			require.NoError(t, err)

			_, err = client.UpsertEntity(ctx, marshalled, &UpsertEntityOptions{UpdateMode: UpdateModeReplace})
			require.NoError(t, err)

			_, err = client.UpsertEntity(ctx, marshalled, &UpsertEntityOptions{UpdateMode: UpdateModeReplace})
			require.NoError(t, err)
		})
	}
}

type mdForListEntities struct {
	Timestamp time.Time `json:"Timestamp"`
	ID        string    `json:"odata.id"` // only for full metadata
}

func TestQuerySimpleEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name: "Pager[ListEntitiesResponse].NextPage",
			}))
			defer delete()

			// Add 5 entities
			entitiesToCreate := createSimpleEntities(5, "partition")
			for _, e := range *entitiesToCreate {
				marshalledEntity, err := json.Marshal(e)
				require.NoError(t, err)
				_, err = client.AddEntity(ctx, marshalledEntity, nil)
				require.NoError(t, err)
			}

			filter := "RowKey lt '5'"
			list := &ListEntitiesOptions{Filter: &filter}
			expectedCount := 4

			var resp ListEntitiesResponse
			pager := client.NewListEntitiesPager(list)
			for pager.More() {
				var err error
				resp, err = pager.NextPage(ctx)
				require.NoError(t, err)
				require.Equal(t, len(resp.Entities), expectedCount)
			}

			for i, e := range resp.Entities {
				var mapModel map[string]any
				err := json.Unmarshal(e, &mapModel)
				require.NoError(t, err)

				_, ok := mapModel[timestamp]
				require.True(t, ok)

				_, ok = mapModel[etagOData]
				require.True(t, ok)

				var b basicTestEntity
				err = json.Unmarshal(e, &b)
				require.NoError(t, err)

				require.Equal(t, b.PartitionKey, "partition")
				require.Equal(t, b.RowKey, fmt.Sprint(i+1))
				require.Equal(t, b.String, (*entitiesToCreate)[i].String)
				require.Equal(t, b.Integer, (*entitiesToCreate)[i].Integer)
				require.Equal(t, b.Bool, (*entitiesToCreate)[i].Bool)

				var md mdForListEntities
				require.NoError(t, json.Unmarshal(e, &md))
				require.False(t, md.Timestamp.IsZero())
				require.Empty(t, md.ID)
			}
		})
	}
}

func TestQueryComplexEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			// Add 5 entities
			entitiesToCreate := createComplexEntities(5, "partition")
			for _, e := range entitiesToCreate {
				marshalledEntity, err := json.Marshal(e)
				require.NoError(t, err)
				_, err = client.AddEntity(ctx, marshalledEntity, nil)
				require.NoError(t, err)
			}

			filter := "RowKey lt '5'"
			expectedCount := 4
			options := &ListEntitiesOptions{
				Filter: &filter,
				Format: to.Ptr(MetadataFormatFull),
			}

			pager := client.NewListEntitiesPager(options)
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				require.Equal(t, expectedCount, len(resp.Entities))

				for idx, entity := range resp.Entities {
					model := complexTestEntity{}
					err := json.Unmarshal(entity, &model)
					require.NoError(t, err)

					require.Equal(t, model.PartitionKey, "partition")
					require.Equal(t, model.RowKey, (entitiesToCreate)[idx].RowKey)
					require.Equal(t, model.Integer, (entitiesToCreate)[idx].Integer)
					require.Equal(t, model.String, (entitiesToCreate)[idx].String)
					require.Equal(t, model.Bool, (entitiesToCreate)[idx].Bool)
					require.Equal(t, model.Float, (entitiesToCreate)[idx].Float)
					require.Equal(t, model.DateTime, (entitiesToCreate)[idx].DateTime)
					require.Equal(t, model.Byte, (entitiesToCreate)[idx].Byte)

					var md mdForListEntities
					require.NoError(t, json.Unmarshal(entity, &md))
					require.False(t, md.Timestamp.IsZero())
					if service == "storage" {
						// cosmos doesn't send full metadata
						require.NotEmpty(t, md.ID)
					}
				}
			}
		})
	}
}

func TestInvalidEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			badEntity := map[string]any{
				"Value":  10,
				"String": "stringystring",
			}

			badEntityMarshalled, err := json.Marshal(badEntity)
			require.NoError(t, err)
			_, err = client.AddEntity(ctx, badEntityMarshalled, nil)

			require.NotNil(t, err)
			require.Contains(t, err.Error(), errPartitionKeyRowKeyError.Error())
		})
	}
}

func TestContinuationTokens(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			err := insertNEntities("contToken", 10, client)
			require.NoError(t, err)

			pager := client.NewListEntitiesPager(&ListEntitiesOptions{Top: to.Ptr(int32(1))})
			var pkContToken string
			var rkContToken string
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				require.Equal(t, 1, len(resp.Entities))
				require.NotNil(t, resp.NextPartitionKey)
				require.NotNil(t, resp.NextRowKey)
				pkContToken = *resp.NextPartitionKey
				rkContToken = *resp.NextRowKey
				break
			}

			require.NotNil(t, pkContToken)
			require.NotNil(t, rkContToken)

			newPager := client.NewListEntitiesPager(&ListEntitiesOptions{
				NextPartitionKey: &pkContToken,
				NextRowKey:       &rkContToken,
			})
			count := 0
			for newPager.More() {
				resp, err := newPager.NextPage(ctx)
				require.NoError(t, err)
				count += len(resp.Entities)
			}
			require.Equal(t, 9, count)
		})
	}
}

func TestContinuationTokensFilters(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			err := insertNEntities("contToken", 10, client)
			require.NoError(t, err)

			pager := client.NewListEntitiesPager(&ListEntitiesOptions{
				Top:    to.Ptr(int32(1)),
				Filter: to.Ptr("Value le 5"),
			})
			var pkContToken string
			var rkContToken string
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				require.Equal(t, 1, len(resp.Entities))
				require.NotNil(t, resp.NextPartitionKey)
				require.NotNil(t, resp.NextRowKey)
				pkContToken = *resp.NextPartitionKey
				rkContToken = *resp.NextRowKey
				break
			}

			require.NotNil(t, pkContToken)
			require.NotNil(t, rkContToken)

			newPager := client.NewListEntitiesPager(&ListEntitiesOptions{
				NextPartitionKey: &pkContToken,
				NextRowKey:       &rkContToken,
				Filter:           to.Ptr("Value le 5"),
			})
			count := 0
			for newPager.More() {
				resp, err := newPager.NextPage(ctx)
				require.NoError(t, err)
				count += len(resp.Entities)
			}
			require.Equal(t, 4, count)
		})
	}
}

func TestClientConstructor(t *testing.T) {
	// Test NewClient, which is not used by recording infra
	client, err := NewClient("https://fakeaccount.table.core.windows.net/", credential.Fake{}, nil)
	require.NoError(t, err)
	require.NotNil(t, client.client)

	// Test NewClientWithNoCredential, which is also not used
	client2, err := NewClientWithNoCredential("https://fakeaccount.table.core.windows.net/", nil)
	require.NoError(t, err)
	require.NotNil(t, client2.client)
}

func TestAzurite(t *testing.T) {
	// quick and dirty make sure azurite is running
	req, err := http.NewRequest("POST", "http://localhost:10002", nil)
	require.NoError(t, err)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Skip("Skipping Azurite test, azurite is not running")
	}

	connStr := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;"
	svc, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)

	name, err := createRandomName(t, "Table")
	require.NoError(t, err)
	_, err = svc.CreateTable(ctx, name, nil)
	defer func() {
		_, err = svc.DeleteTable(ctx, name, nil)
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	client := svc.NewClient(name)
	entity := EDMEntity{
		Entity: Entity{
			PartitionKey: "pencils",
			RowKey:       "id-003",
		},
		Properties: map[string]any{
			"Product":      "Ticonderoga Pencils",
			"Price":        5.00,
			"Count":        EDMInt64(12345678901234),
			"ProductGUID":  EDMGUID("some-guid-value"),
			"DateReceived": EDMDateTime(time.Now()),
			"ProductCode":  EDMBinary([]byte("somebinaryvalue")),
		},
	}

	data, err := json.Marshal(entity)
	require.NoError(t, err)

	_, err = client.AddEntity(ctx, data, nil)
	require.NoError(t, err)

	count := 0
	pager := client.NewListEntitiesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		require.NoError(t, err)
		count += len(resp.Entities)
	}
	require.Equal(t, 1, count)
}

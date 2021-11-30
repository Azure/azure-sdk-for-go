// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

var services = []string{"storage", "cosmos"}

func TestServiceErrors(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			// Create a duplicate table to produce an error
			_, err := client.Create(ctx, nil)
			require.Error(t, err)
		})
	}
}

func TestCreateTable(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, false)
			defer delete()

			resp, err := client.Create(ctx, nil)

			require.NoError(t, err)
			require.NotNil(t, resp.RawResponse)
		})
	}
}

func TestAddEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			simpleEntity := createSimpleEntity(1, "partition")

			marshalledEntity, err := json.Marshal(simpleEntity)
			require.NoError(t, err)
			_, err = client.AddEntity(ctx, marshalledEntity, nil)
			require.NoError(t, err)
		})
	}
}

func TestAddComplexEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			entity := createComplexEntity(1, "partition")

			marshalledEntity, err := json.Marshal(entity)
			require.NoError(t, err)

			_, err = client.AddEntity(ctx, marshalledEntity, nil)
			require.NoError(t, err)
		})
	}
}

func TestDeleteEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
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
			client, delete := initClientTest(t, service, true)
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
			client, delete := initClientTest(t, service, true)
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

			var unMarshalledPreMerge map[string]interface{}
			err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
			require.NoError(t, err)

			var mapEntity map[string]interface{}
			err = json.Unmarshal(marshalled, &mapEntity)
			require.NoError(t, err)
			mapEntity["MergeProperty"] = "foo"

			reMarshalled, err := json.Marshal(mapEntity)
			require.NoError(t, err)

			_, updateErr := client.UpdateEntity(ctx, reMarshalled, &UpdateEntityOptions{UpdateMode: MergeEntity})
			require.Nil(t, updateErr)

			var qResp ListEntitiesPage
			pager := client.List(listOptions)
			for pager.NextPage(ctx) {
				qResp = pager.PageResponse()
			}
			postMerge := qResp.Entities[0]
			var unmarshalledPostMerge map[string]interface{}
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
			client, delete := initClientTest(t, service, true)
			defer delete()

			entityToCreate := createSimpleEntity(1, "partition")
			marshalled, err := json.Marshal(entityToCreate)
			require.NoError(t, err)

			_, updateErr := client.UpdateEntity(ctx, marshalled, &UpdateEntityOptions{UpdateMode: MergeEntity})
			require.Error(t, updateErr)
		})
	}
}

func TestInsertEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			// 1. Create Basic Entity
			entityToCreate := createSimpleEntity(1, "partition")
			marshalled, err := json.Marshal(entityToCreate)
			require.NoError(t, err)

			_, err = client.InsertEntity(ctx, marshalled, &InsertEntityOptions{UpdateMode: ReplaceEntity})
			require.NoError(t, err)

			filter := "RowKey eq '1'"
			list := &ListEntitiesOptions{Filter: &filter}

			// 2. Query for basic Entity
			preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey, nil)
			require.NoError(t, err)

			var unMarshalledPreMerge map[string]interface{}
			err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
			require.NoError(t, err)

			// 3. Create same entity without Bool property, add "MergeProperty" prop
			mapEntity := createSimpleEntityNoBool(1, "partition")
			mapEntity["MergeProperty"] = "foo"

			reMarshalled, err := json.Marshal(mapEntity)
			require.NoError(t, err)

			// 4. Replace Entity with "bool"-less entity
			_, err = client.InsertEntity(ctx, reMarshalled, &InsertEntityOptions{UpdateMode: ReplaceEntity})
			require.Nil(t, err)

			// 5. Query for new entity
			var qResp ListEntitiesPage
			pager := client.List(list)
			for pager.NextPage(ctx) {
				qResp = pager.PageResponse()
			}
			postMerge := qResp.Entities[0]
			var unmarshalledPostMerge map[string]interface{}
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
			client, delete := initClientTest(t, service, true)
			defer delete()

			// 1. Create Basic Entity
			entityToCreate := createSimpleEntity(1, "partition")
			marshalled, err := json.Marshal(entityToCreate)
			require.NoError(t, err)

			_, err = client.InsertEntity(ctx, marshalled, &InsertEntityOptions{UpdateMode: ReplaceEntity})
			require.NoError(t, err)

			_, err = client.InsertEntity(ctx, marshalled, &InsertEntityOptions{UpdateMode: ReplaceEntity})
			require.NoError(t, err)
		})
	}
}

func TestQuerySimpleEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
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

			var resp ListEntitiesPage
			pager := client.List(list)
			for pager.NextPage(ctx) {
				resp = pager.PageResponse()
				require.Equal(t, len(resp.Entities), expectedCount)
			}

			for i, e := range resp.Entities {
				var mapModel map[string]interface{}
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
			}
		})
	}
}

func TestQueryComplexEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			// Add 5 entities
			entitiesToCreate := createComplexEntities(5, "partition")
			for _, e := range *entitiesToCreate {
				marshalledEntity, err := json.Marshal(e)
				require.NoError(t, err)
				_, err = client.AddEntity(ctx, marshalledEntity, nil)
				require.NoError(t, err)
			}

			filter := "RowKey lt '5'"
			expectedCount := 4
			options := &ListEntitiesOptions{Filter: &filter}

			var resp ListEntitiesPage
			pager := client.List(options)
			for pager.NextPage(ctx) {
				resp = pager.PageResponse()
				require.Equal(t, expectedCount, len(resp.Entities))

				for idx, entity := range resp.Entities {
					model := complexTestEntity{}
					err := json.Unmarshal(entity, &model)
					require.NoError(t, err)

					require.Equal(t, model.PartitionKey, "partition")
					require.Equal(t, model.RowKey, (*entitiesToCreate)[idx].RowKey)
					require.Equal(t, model.Integer, (*entitiesToCreate)[idx].Integer)
					require.Equal(t, model.String, (*entitiesToCreate)[idx].String)
					require.Equal(t, model.Bool, (*entitiesToCreate)[idx].Bool)
					require.Equal(t, model.Float, (*entitiesToCreate)[idx].Float)
					require.Equal(t, model.DateTime, (*entitiesToCreate)[idx].DateTime)
					require.Equal(t, model.Byte, (*entitiesToCreate)[idx].Byte)
				}
			}
		})
	}
}

func TestInvalidEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			badEntity := map[string]interface{}{
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
			client, delete := initClientTest(t, service, true)
			defer delete()

			err := insertNEntities("contToken", 10, client)
			require.NoError(t, err)

			pager := client.List(&ListEntitiesOptions{Top: to.Int32Ptr(1)})
			var pkContToken *string
			var rkContToken *string
			for pager.NextPage(ctx) {
				require.Equal(t, 1, len(pager.PageResponse().Entities))
				pkContToken = pager.NextPagePartitionKey()
				rkContToken = pager.NextPageRowKey()
				break
			}

			require.NoError(t, pager.Err())
			require.NotNil(t, pkContToken)
			require.NotNil(t, rkContToken)

			newPager := client.List(&ListEntitiesOptions{
				PartitionKey: pkContToken,
				RowKey:       rkContToken,
			})
			count := 0
			for newPager.NextPage(ctx) {
				count += len(newPager.PageResponse().Entities)
			}

			require.NoError(t, pager.Err())
			require.Equal(t, 9, count)
		})
	}
}

func TestContinuationTokensFilters(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			err := insertNEntities("contToken", 10, client)
			require.NoError(t, err)

			pager := client.List(&ListEntitiesOptions{
				Top:    to.Int32Ptr(1),
				Filter: to.StringPtr("Value le 5"),
			})
			var pkContToken *string
			var rkContToken *string
			for pager.NextPage(ctx) {
				require.Equal(t, 1, len(pager.PageResponse().Entities))
				pkContToken = pager.NextPagePartitionKey()
				rkContToken = pager.NextPageRowKey()
				break
			}

			require.NoError(t, pager.Err())
			require.NotNil(t, pkContToken)
			require.NotNil(t, rkContToken)

			newPager := client.List(&ListEntitiesOptions{
				PartitionKey: pkContToken,
				RowKey:       rkContToken,
				Filter:       to.StringPtr("Value le 5"),
			})
			count := 0
			for newPager.NextPage(ctx) {
				count += len(newPager.PageResponse().Entities)
			}

			require.NoError(t, pager.Err())
			require.Equal(t, 4, count)
		})
	}
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
	client, err := svc.CreateTable(ctx, name, nil)
	defer func() {
		_, err = client.Delete(ctx, nil)
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	entity := EDMEntity{
		Entity: Entity{
			PartitionKey: "pencils",
			RowKey:       "id-003",
		},
		Properties: map[string]interface{}{
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
	pager := client.List(nil)
	for pager.NextPage(ctx) {
		count += len(pager.PageResponse().Entities)
	}
	require.NoError(t, pager.Err())
	require.Equal(t, 1, count)
}

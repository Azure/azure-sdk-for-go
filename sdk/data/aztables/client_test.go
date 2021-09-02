// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

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

			resp, err := client.Create(context.Background(), nil)

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

func TestMergeEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			require := require.New(t) // Remove this later

			entityToCreate := createSimpleEntity(1, "partition")
			marshalled := marshalBasicEntity(entityToCreate, require)

			_, err := client.AddEntity(ctx, *marshalled, nil)
			require.NoError(err)

			filter := "RowKey eq '1'"
			listOptions := &ListEntitiesOptions{Filter: &filter}

			preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey, nil)
			require.NoError(err)

			var unMarshalledPreMerge map[string]interface{}
			err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
			require.NoError(err)

			var mapEntity map[string]interface{}
			err = json.Unmarshal(*marshalled, &mapEntity)
			require.NoError(err)
			mapEntity["MergeProperty"] = "foo"

			reMarshalled, err := json.Marshal(mapEntity)
			require.NoError(err)

			_, updateErr := client.UpdateEntity(ctx, reMarshalled, &UpdateEntityOptions{UpdateMode: MergeEntity})
			require.Nil(updateErr)

			var qResp ListEntitiesPage
			pager := client.List(listOptions)
			for pager.NextPage(ctx) {
				qResp = pager.PageResponse()
			}
			postMerge := qResp.Entities[0]
			var unmarshalledPostMerge map[string]interface{}
			err = json.Unmarshal(postMerge, &unmarshalledPostMerge)
			require.NoError(err)

			require.Equal(unmarshalledPostMerge["PartitionKey"], unMarshalledPreMerge["PartitionKey"])
			require.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

			_, ok := unMarshalledPreMerge["MergeProperty"]
			require.False(ok)
		})
	}
}

func TestInsertEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			require := require.New(t)

			// 1. Create Basic Entity
			entityToCreate := createSimpleEntity(1, "partition")
			marshalled := marshalBasicEntity(entityToCreate, require)

			_, err := client.InsertEntity(ctx, *marshalled, &InsertEntityOptions{UpdateMode: ReplaceEntity})
			require.NoError(err)

			filter := "RowKey eq '1'"
			list := &ListEntitiesOptions{Filter: &filter}

			// 2. Query for basic Entity
			preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey, nil)
			require.NoError(err)

			var unMarshalledPreMerge map[string]interface{}
			err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
			require.NoError(err)

			// 3. Create same entity without Bool property, add "MergeProperty" prop
			mapEntity := createSimpleEntityNoBool(1, "partition")
			mapEntity["MergeProperty"] = "foo"

			reMarshalled, err := json.Marshal(mapEntity)
			require.NoError(err)

			// 4. Replace Entity with "bool"-less entity
			_, err = client.InsertEntity(ctx, reMarshalled, &InsertEntityOptions{UpdateMode: ReplaceEntity})
			require.Nil(err)

			// 5. Query for new entity
			var qResp ListEntitiesPage
			pager := client.List(list)
			for pager.NextPage(ctx) {
				qResp = pager.PageResponse()
			}
			postMerge := qResp.Entities[0]
			var unmarshalledPostMerge map[string]interface{}
			err = json.Unmarshal(postMerge, &unmarshalledPostMerge)
			require.NoError(err)

			// 6. Make assertions
			require.Less(len(unmarshalledPostMerge), len(unMarshalledPreMerge))
			require.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

			_, ok := unmarshalledPostMerge["Bool"]
			require.Falsef(ok, "Bool property should not be available in the merged entity")
		})
	}
}

func TestQuerySimpleEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()
			require := require.New(t)

			// Add 5 entities
			entitiesToCreate := createSimpleEntities(5, "partition")
			for _, e := range *entitiesToCreate {
				marshalledEntity, err := json.Marshal(e)
				require.NoError(err)
				_, err = client.AddEntity(ctx, marshalledEntity, nil)
				require.NoError(err)
			}

			filter := "RowKey lt '5'"
			list := &ListEntitiesOptions{Filter: &filter}
			expectedCount := 4

			var resp ListEntitiesPage
			pager := client.List(list)
			for pager.NextPage(ctx) {
				resp = pager.PageResponse()
				require.Equal(len(resp.Entities), expectedCount)
			}

			for i, e := range resp.Entities {
				var mapModel map[string]interface{}
				err := json.Unmarshal(e, &mapModel)
				require.NoError(err)

				_, ok := mapModel[timestamp]
				require.True(ok)

				_, ok = mapModel[etagOData]
				require.True(ok)

				var b basicTestEntity
				err = json.Unmarshal(e, &b)
				require.NoError(err)

				require.Equal(b.PartitionKey, "partition")
				require.Equal(b.RowKey, fmt.Sprint(i+1))
				require.Equal(b.String, (*entitiesToCreate)[i].String)
				require.Equal(b.Integer, (*entitiesToCreate)[i].Integer)
				require.Equal(b.Bool, (*entitiesToCreate)[i].Bool)
			}
		})
	}
}

func TestQueryComplexEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			require := require.New(t)

			// Add 5 entities
			entitiesToCreate := createComplexEntities(5, "partition")
			for _, e := range *entitiesToCreate {
				marshalledEntity, err := json.Marshal(e)
				require.NoError(err)
				_, err = client.AddEntity(ctx, marshalledEntity, nil)
				require.NoError(err)
			}

			filter := "RowKey lt '5'"
			expectedCount := 4
			options := &ListEntitiesOptions{Filter: &filter}

			var resp ListEntitiesPage
			pager := client.List(options)
			for pager.NextPage(ctx) {
				resp = pager.PageResponse()
				require.Equal(expectedCount, len(resp.Entities))

				for idx, entity := range resp.Entities {
					model := complexTestEntity{}
					err := json.Unmarshal(entity, &model)
					require.NoError(err)

					require.Equal(model.PartitionKey, "partition")
					require.Equal(model.RowKey, (*entitiesToCreate)[idx].RowKey)
					require.Equal(model.Integer, (*entitiesToCreate)[idx].Integer)
					require.Equal(model.String, (*entitiesToCreate)[idx].String)
					require.Equal(model.Bool, (*entitiesToCreate)[idx].Bool)
					require.Equal(model.Float, (*entitiesToCreate)[idx].Float)
					require.Equal(model.DateTime, (*entitiesToCreate)[idx].DateTime)
					require.Equal(model.Byte, (*entitiesToCreate)[idx].Byte)
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

			require := require.New(t)

			badEntity := map[string]interface{}{
				"Value":  10,
				"String": "stringystring",
			}

			badEntityMarshalled, err := json.Marshal(badEntity)
			require.NoError(err)
			_, err = client.AddEntity(ctx, badEntityMarshalled, nil)

			require.NotNil(err)
			require.Contains(err.Error(), errPartitionKeyRowKeyError.Error())
		})
	}
}

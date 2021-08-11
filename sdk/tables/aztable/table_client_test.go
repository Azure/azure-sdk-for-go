// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type tableClientLiveTests struct {
	suite.Suite
	endpointType EndpointType
	mode         recording.RecordMode
}

// Hookup to the testing framework
func TestTableClient_Storage(t *testing.T) {
	storage := tableClientLiveTests{endpointType: StorageEndpoint, mode: recording.Playback /* change to Record to re-record tests */}
	suite.Run(t, &storage)
}

// Hookup to the testing framework
func TestTableClient_Cosmos(t *testing.T) {
	cosmos := tableClientLiveTests{endpointType: CosmosEndpoint, mode: recording.Playback /* change to Record to re-record tests */}
	suite.Run(t, &cosmos)
}

func (s *tableClientLiveTests) TestServiceErrors() {
	client, delete := s.init(true)
	defer delete()

	// Create a duplicate table to produce an error
	_, err := client.Create(ctx)
	var svcErr *runtime.ResponseError
	errors.As(err, &svcErr)
	require.Equal(s.T(), svcErr.RawResponse().StatusCode, http.StatusConflict)
}

func (s *tableClientLiveTests) TestCreateTable() {
	require := require.New(s.T())
	client, delete := s.init(false)
	defer delete()

	resp, err := client.Create(ctx)

	require.NoError(err)
	require.Equal(*resp.TableResponse.TableName, client.Name)
}

func (s *tableClientLiveTests) TestAddEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	marshalledEntity, err := json.Marshal((*entitiesToCreate)[0])
	require.NoError(err)
	resp, err := client.AddEntity(ctx, marshalledEntity)
	require.NoError(err)
	require.NotNil(resp)
}

func (s *tableClientLiveTests) TestAddComplexEntity() {
	require := require.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(1, "partition")

	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		require.NoError(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		var svcErr *runtime.ResponseError
		errors.As(err, &svcErr)
		require.Nilf(err, getStringFromBody(svcErr))
	}
}

func (s *tableClientLiveTests) TestDeleteEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	simpleEntity := createSimpleEntity(1, "partition")

	marshalledEntity, err := json.Marshal(simpleEntity)
	require.NoError(err)
	_, err = client.AddEntity(ctx, marshalledEntity)
	require.NoError(err)
	_, delErr := client.DeleteEntity(ctx, simpleEntity.PartitionKey, simpleEntity.RowKey, nil)
	require.Nil(delErr)
}

func (s *tableClientLiveTests) TestMergeEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entityToCreate := createSimpleEntity(1, "partition")
	marshalled := marshalBasicEntity(entityToCreate, require)

	_, err := client.AddEntity(ctx, *marshalled)
	require.NoError(err)

	filter := "RowKey eq '1'"
	listOptions := &ListOptions{Filter: &filter}

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

	_, updateErr := client.UpdateEntity(ctx, reMarshalled, nil, MergeEntity)
	require.Nil(updateErr)

	var qResp TableEntityQueryByteResponseResponse
	pager := client.List(listOptions)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unmarshalledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unmarshalledPostMerge)
	require.NoError(err)

	require.Equal(unmarshalledPostMerge["PartitionKey"], unMarshalledPreMerge["PartitionKey"])
	require.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

	_, ok := unMarshalledPreMerge["MergeProperty"]
	require.False(ok)
}

func (s *tableClientLiveTests) TestInsertEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// 1. Create Basic Entity
	entityToCreate := createSimpleEntity(1, "partition")
	marshalled := marshalBasicEntity(entityToCreate, require)

	_, err := client.InsertEntity(ctx, *marshalled, ReplaceEntity)
	require.NoError(err)

	filter := "RowKey eq '1'"
	list := &ListOptions{Filter: &filter}

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

	// 4. Replace Entity with "bool"-less entity
	_, err = client.InsertEntity(ctx, reMarshalled, ReplaceEntity)
	require.Nil(err)

	// 5. Query for new entity
	var qResp TableEntityQueryByteResponseResponse
	pager := client.List(list)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unmarshalledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unmarshalledPostMerge)

	// 6. Make assertions
	require.Less(len(unmarshalledPostMerge), len(unMarshalledPreMerge))
	require.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

	_, ok := unmarshalledPostMerge["Bool"]
	require.Falsef(ok, "Bool property should not be available in the merged entity")
}

func (s *tableClientLiveTests) TestQuerySimpleEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createSimpleEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		require.NoError(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		require.NoError(err)
	}

	filter := "RowKey lt '5'"
	list := &ListOptions{Filter: &filter}
	expectedCount := 4

	var resp TableEntityQueryByteResponseResponse
	pager := client.List(list)
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		require.Equal(len(resp.TableEntityQueryResponse.Value), expectedCount)
	}

	for i, e := range resp.TableEntityQueryResponse.Value {
		var mapModel map[string]interface{}
		err := json.Unmarshal(e, &mapModel)
		require.NoError(err)

		_, ok := mapModel[timestamp]
		require.True(ok)

		_, ok = mapModel[etagOdata]
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
}

func (s *tableClientLiveTests) TestQueryComplexEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createComplexEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		require.NoError(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		require.NoError(err)
	}

	filter := "RowKey lt '5'"
	expectedCount := 4
	options := &ListOptions{Filter: &filter}

	var resp TableEntityQueryByteResponseResponse
	pager := client.List(options)
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		require.Equal(expectedCount, len(resp.TableEntityQueryResponse.Value))

		for idx, entity := range resp.TableEntityQueryResponse.Value {
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
}

func (s *tableClientLiveTests) TestInvalidEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	badEntity := map[string]interface{}{
		"Value":  10,
		"String": "stringystring",
	}

	badEntityMarshalled, err := json.Marshal(badEntity)
	_, err = client.AddEntity(ctx, badEntityMarshalled)

	require.NotNil(err)
	require.Contains(err.Error(), partitionKeyRowKeyError.Error())
}

// setup the test environment
func (s *tableClientLiveTests) BeforeTest(suite string, test string) {
	recordedTestSetup(s.T(), s.T().Name(), s.endpointType, s.mode)
}

// teardown the test context
func (s *tableClientLiveTests) AfterTest(suite string, test string) {
	recordedTestTeardown(s.T().Name())
}

func (s *tableClientLiveTests) init(createTable bool) (*TableClient, func()) {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, _ := getTableName(context)
	client := context.client.NewTableClient(tableName)
	if createTable {
		_, err := client.Create(ctx)
		// fmt.Println("CREATE ERROR: ", err.Error())
		if err != nil {
			var svcErr *runtime.ResponseError
			errors.As(err, &svcErr)
			require.FailNow(getStringFromBody(svcErr))
		}
	}
	return client, func() {
		_, err := client.Delete(ctx, nil)
		if err != nil {
			fmt.Printf("Error deleting table. %v\n", err.Error())
		}
	}
}

func getStringFromBody(e *runtime.ResponseError) string {
	if e == nil {
		return "Error is nil"
	}
	r := e.RawResponse()
	body := bytes.Buffer{}
	b := r.Body
	b.Close()
	if b != nil {
		_, err := body.ReadFrom(b)
		if err != nil {
			return "<emtpy body>"
		}
		_ = ioutil.NopCloser(&body)
	}
	return body.String()
}

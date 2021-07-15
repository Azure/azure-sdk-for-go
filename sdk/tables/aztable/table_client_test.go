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
	"github.com/stretchr/testify/assert"
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
	assert.Equal(s.T(), svcErr.RawResponse().StatusCode, http.StatusConflict)
}

func (s *tableClientLiveTests) TestCreateTable() {
	assert := assert.New(s.T())
	client, delete := s.init(false)
	defer delete()

	resp, err := client.Create(ctx)

	assert.Nil(err)
	assert.Equal(*resp.TableResponse.TableName, client.Name)
}

func (s *tableClientLiveTests) TestAddEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	marshalledEntity, err := json.Marshal((*entitiesToCreate)[0])
	assert.Nil(err)
	resp, err := client.AddEntity(ctx, marshalledEntity)
	assert.Nil(err)
	assert.NotNil(resp)
}

func (s *tableClientLiveTests) TestAddComplexEntity() {
	assert := assert.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(1, "partition")

	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		assert.Nil(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		var svcErr *runtime.ResponseError
		errors.As(err, &svcErr)
		assert.Nilf(err, getStringFromBody(svcErr))
	}
}

func (s *tableClientLiveTests) TestDeleteEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	simpleEntity := createSimpleEntity(1, "partition")

	marshalledEntity, err := json.Marshal(simpleEntity)
	assert.Nil(err)
	_, err = client.AddEntity(ctx, marshalledEntity)
	assert.Nil(err)
	_, delErr := client.DeleteEntity(ctx, simpleEntity.PartitionKey, simpleEntity.RowKey, nil)
	assert.Nil(delErr)
}

func (s *tableClientLiveTests) TestMergeEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entityToCreate := createSimpleEntity(1, "partition")
	marshalled := marshalBasicEntity(entityToCreate, assert)

	_, err := client.AddEntity(ctx, *marshalled)
	assert.Nil(err)

	filter := "RowKey eq '1'"
	queryOptions := &QueryOptions{Filter: &filter}

	preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey, nil)
	assert.Nil(err)

	var unMarshalledPreMerge map[string]interface{}
	err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
	assert.Nil(err)

	var mapEntity map[string]interface{}
	err = json.Unmarshal(*marshalled, &mapEntity)
	assert.Nil(err)
	mapEntity["MergeProperty"] = "foo"

	reMarshalled, err := json.Marshal(mapEntity)
	assert.Nil(err)

	_, updateErr := client.UpdateEntity(ctx, reMarshalled, nil, Merge)
	assert.Nil(updateErr)

	var qResp TableEntityQueryByteResponseResponse
	pager := client.Query(queryOptions)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unmarshalledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unmarshalledPostMerge)
	assert.Nil(err)

	assert.Equal(unmarshalledPostMerge["PartitionKey"], unMarshalledPreMerge["PartitionKey"])
	assert.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

	_, ok := unMarshalledPreMerge["MergeProperty"]
	assert.False(ok)
}

func (s *tableClientLiveTests) TestUpsertEntity() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// 1. Create Basic Entity
	entityToCreate := createSimpleEntity(1, "partition")
	marshalled := marshalBasicEntity(entityToCreate, assert)

	_, err := client.UpsertEntity(ctx, *marshalled, Replace)
	assert.Nil(err)

	filter := "RowKey eq '1'"
	query := &QueryOptions{Filter: &filter}

	// 2. Query for basic Entity
	preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey, nil)
	assert.Nil(err)

	var unMarshalledPreMerge map[string]interface{}
	err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
	assert.Nil(err)

	// 3. Create same entity without Bool property, add "MergeProperty" prop
	mapEntity := createSimpleEntityNoBool(1, "partition")
	mapEntity["MergeProperty"] = "foo"

	reMarshalled, err := json.Marshal(mapEntity)

	// 4. Replace Entity with "bool"-less entity
	_, err = client.UpsertEntity(ctx, reMarshalled, Replace)
	require.Nil(err)

	// 5. Query for new entity
	var qResp TableEntityQueryByteResponseResponse
	pager := client.Query(query)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unmarshalledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unmarshalledPostMerge)

	// 6. Make assertions
	assert.Less(len(unmarshalledPostMerge), len(unMarshalledPreMerge))
	assert.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

	_, ok := unmarshalledPostMerge["Bool"]
	assert.Falsef(ok, "Bool property should not be available in the merged entity")
}

func (s *tableClientLiveTests) TestQuerySimpleEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createSimpleEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		assert.Nil(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		assert.Nil(err)
	}

	filter := "RowKey lt '5'"
	query := &QueryOptions{Filter: &filter}
	expectedCount := 4

	var resp TableEntityQueryByteResponseResponse
	pager := client.Query(query)
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		assert.Equal(len(resp.TableEntityQueryResponse.Value), expectedCount)
	}

	for i, e := range resp.TableEntityQueryResponse.Value {
		var mapModel map[string]interface{}
		err := json.Unmarshal(e, &mapModel)
		assert.Nil(err)

		_, ok := mapModel[timestamp]
		assert.True(ok)

		_, ok = mapModel[etagOdata]
		assert.True(ok)

		var b basicTestEntity
		err = json.Unmarshal(e, &b)
		assert.Nil(err)

		assert.Equal(b.PartitionKey, "partition")
		assert.Equal(b.RowKey, fmt.Sprint(i+1))
		assert.Equal(b.String, (*entitiesToCreate)[i].String)
		assert.Equal(b.Integer, (*entitiesToCreate)[i].Integer)
		assert.Equal(b.Bool, (*entitiesToCreate)[i].Bool)
	}
}

func (s *tableClientLiveTests) TestQueryComplexEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createComplexEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		assert.Nil(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		assert.Nil(err)
	}

	filter := "RowKey lt '5'"
	expectedCount := 4
	query := &QueryOptions{Filter: &filter}

	var resp TableEntityQueryByteResponseResponse
	pager := client.Query(query)
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		assert.Equal(expectedCount, len(resp.TableEntityQueryResponse.Value))

		for idx, entity := range resp.TableEntityQueryResponse.Value {
			model := complexTestEntity{}
			err := json.Unmarshal(entity, &model)
			assert.Nil(err)

			assert.Equal(model.PartitionKey, "partition")
			assert.Equal(model.RowKey, (*entitiesToCreate)[idx].RowKey)
			assert.Equal(model.Integer, (*entitiesToCreate)[idx].Integer)
			assert.Equal(model.String, (*entitiesToCreate)[idx].String)
			assert.Equal(model.Bool, (*entitiesToCreate)[idx].Bool)
			assert.Equal(model.Float, (*entitiesToCreate)[idx].Float)
			assert.Equal(model.DateTime, (*entitiesToCreate)[idx].DateTime)
			assert.Equal(model.Byte, (*entitiesToCreate)[idx].Byte)
		}

	}
}

func (s *tableClientLiveTests) TestInvalidEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	badEntity := map[string]interface{}{
		"Value":  10,
		"String": "stringystring",
	}

	badEntityMarshalled, err := json.Marshal(badEntity)
	_, err = client.AddEntity(ctx, badEntityMarshalled)

	assert.NotNil(err)
	assert.Contains(err.Error(), partitionKeyRowKeyError.Error())
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
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, _ := getTableName(context)
	client := context.client.NewTableClient(tableName)
	if createTable {
		_, err := client.Create(ctx)
		if err != nil {
			var svcErr *runtime.ResponseError
			errors.As(err, &svcErr)
			assert.FailNow(getStringFromBody(svcErr))
		}
	}
	return client, func() {
		client.Delete(ctx)
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
		b = ioutil.NopCloser(&body)
	}
	return body.String()
}

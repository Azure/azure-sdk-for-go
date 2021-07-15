package aztable

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *tableClientLiveTests) TestAddBasicEntity() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	basicEntity := basicTestEntity{
		Entity: Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Integer: 10,
		String:  "abcdef",
		Bool:    true,
	}

	marshalled, err := json.Marshal(basicEntity)
	require.Nil(err)
	_, err = client.AddEntity(ctx, marshalled)
	require.Nil(err)

	resp, err := client.GetEntity(ctx, "pk001", "rk001")
	require.Nil(err)

	receivedEntity := basicTestEntity{}
	err = json.Unmarshal(resp.Value, &receivedEntity)
	require.Nil(err)
	assert.Equal(receivedEntity.PartitionKey, "pk001")
	assert.Equal(receivedEntity.RowKey, "rk001")

	queryString := "PartitionKey eq 'pk001'"
	queryOptions := QueryOptions{Filter: &queryString}
	pager := client.Query(&queryOptions)
	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		for _, e := range resp.TableEntityQueryResponse.Value {
			err = json.Unmarshal(e, &receivedEntity)
			assert.Nil(err)
			assert.Equal(receivedEntity.PartitionKey, "pk001")
			assert.Equal(receivedEntity.RowKey, "rk001")
			count += 1
		}
	}

	assert.Equal(count, 1)
}

type AnnotatedEntity struct {
	Entity
	Large               string    `json:"large"`
	LargeAnnotation     string    `json:"large@odata.type"`
	FloatType           float64   `json:"FloatType"`
	FloatTypeAnnotation string    `json:"FloatType@odata.type"`
	DateType            time.Time `json:"DateType"`
	DateTypeAnnotation  string    `json:"DateType@odata.type"`
	Stringy             string    `json:"Stringy"`
	StringyAnnotation   string    `json:"Stringy@odata.type"`
	Bool                bool      `json:"Bool"`
	BoolAnnotation      string    `json:"Bool@odata.type"`
	Small               int32     `json:"Small"`
	SmallAnnotation     string    `json:"Small@odata.type"`
	Binary              []byte    `json:"Binary"`
	BinaryAnnotation    string    `json:"Binary@odata.type"`
}

func createAnnotatedEntity(count int, pk string) AnnotatedEntity {
	return AnnotatedEntity{
		Entity: Entity{
			PartitionKey: pk,
			RowKey:       fmt.Sprint(count),
		},
		Large:               "1125899906842624",
		LargeAnnotation:     "Edm.Int64",
		FloatType:           math.Pow(2, 33),
		FloatTypeAnnotation: "Edm.Double",
		DateType:            time.Date(2021, time.April, 1, 1, 1, 1, 1, time.UTC),
		DateTypeAnnotation:  "Edm.DateTime",
		Stringy:             "somestring",
		StringyAnnotation:   "Edm.String",
		Bool:                true,
		BoolAnnotation:      "Edm.Boolean",
		Small:               10,
		SmallAnnotation:     "Edm.Int32",
		Binary:              []byte("binary"),
		BinaryAnnotation:    "Edm.Binary",
	}
}

func (s *tableClientLiveTests) TestAddAnnotatedEntity() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	annotatedEntity := createAnnotatedEntity(1, "partition")

	marshalled, err := json.Marshal(annotatedEntity)
	require.Nil(err)
	_, err = client.AddEntity(ctx, marshalled)
	require.Nil(err)

	resp, err := client.GetEntity(ctx, "partition", fmt.Sprint(1))
	require.Nil(err)
	receivedEntity := AnnotatedEntity{}
	err = json.Unmarshal(resp.Value, &receivedEntity)
	require.Nil(err)
	assert.Equal(receivedEntity.PartitionKey, annotatedEntity.PartitionKey)
	assert.Equal(receivedEntity.RowKey, annotatedEntity.RowKey)
	// assert.Equal(receivedEntity.Large, annotatedEntity.Large)
	// assert.Equal(receivedEntity.BigIntAnnotation, annotatedEntity.BigIntAnnotation)

	queryString := "PartitionKey eq 'partition'"
	queryOptions := QueryOptions{Filter: &queryString}
	pager := client.Query(&queryOptions)
	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		for _, e := range resp.TableEntityQueryResponse.Value {
			err = json.Unmarshal(e, &receivedEntity)
			require.Nil(err)
			assert.Equal(receivedEntity.PartitionKey, "partition")
			assert.Equal(receivedEntity.RowKey, "1")
			count += 1
		}
	}

	assert.Equal(count, 1)
}

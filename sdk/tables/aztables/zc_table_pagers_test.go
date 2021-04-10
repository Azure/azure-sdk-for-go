// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/stretchr/testify/assert"
)

type pagerTests struct{}

func TestCastAndRemoveAnnotations(t *testing.T) {
	assert := assert.New(t)

	r := &http.Response{Body: closerFromString(complexPayload)}
	resp := azcore.Response{Response: r}

	var val map[string]interface{}
	err := resp.UnmarshalAsJSON(&val)
	assert.Nil(err)
	entity, err := castAndRemoveAnnotations(&val)
	assert.Nil(err)
	// assert all odata annotations are removed.
	for k, _ := range *entity {
		assert.NotContains(k, OdataType)
	}

	assert.IsType(time.Now(), (*entity)["SomeDateProperty"])
	assert.IsType([]byte{}, (*entity)["SomeBinaryProperty"])
	assert.IsType(float64(0), (*entity)["SomeDoubleProperty0"])
	// TODO: fix this
	// assert.IsType(int(0), (*entity)["SomeIntProperty"])
}

func TestToMap(t *testing.T) {
	assert := assert.New(t)

	ent := createComplexEntity()

	entMap, err := toMap(ent)
	assert.Nil(err)

	// Validate that we have all the @odata.type properties for types []byte, int64, float64, time.Time, and uuid
	for k, v := range odataHintProps {
		vv, ok := (*entMap)[odataType(k)]
		assert.Truef(ok, "Should have found map key of name '%s'", odataType(k))
		assert.Equal(v, vv)
	}

	// validate all the types were properly casted / converted
	assert.Equal(ent.PartitionKey, (*entMap)["PartitionKey"])
	assert.Equal(ent.RowKey, (*entMap)["RowKey"])
	ts, _ := time.Parse(ISO8601, (*entMap)["Timestamp"].(string))
	assert.Equal(ent.Timestamp.UTC().String(), ts.String())
	assert.Equal(base64.StdEncoding.EncodeToString(ent.SomeBinaryProperty), string((*entMap)["SomeBinaryProperty"].(string)))
	ts, _ = time.Parse(ISO8601, (*entMap)["SomeDateProperty"].(string))
	assert.Equal(ent.SomeDateProperty.UTC().String(), ts.String())
	assert.Equal(ent.SomeDoubleProperty0, (*entMap)["SomeDoubleProperty0"])
	assert.Equal(ent.SomeDoubleProperty1, (*entMap)["SomeDoubleProperty1"])
	var u uuid.UUID = ent.SomeGuidProperty
	assert.Equal(u.String(), (*entMap)["SomeGuidProperty"].(string))
	assert.Equal(strconv.FormatInt(ent.SomeInt64Property, 10), (*entMap)["SomeInt64Property"].(string))
	assert.Equal(ent.SomeIntProperty, (*entMap)["SomeIntProperty"])
	assert.Equal(ent.SomeStringProperty, (*entMap)["SomeStringProperty"])
	assert.Equal(*ent.SomePtrStringProperty, (*entMap)["SomePtrStringProperty"])
}

func TestEntitySerialization(t *testing.T) {
	assert := assert.New(t)

	ent := createComplexEntity()

	b, err := json.Marshal(ent)
	assert.Nil(err)
	assert.NotEmpty(b)
	s := string(b)
	//assert.FailNow(s)
	assert.NotEmpty(s)
}

func createComplexEntity() complexEntity {
	sp := "some pointer to string"
	var e = complexEntity{
		PartitionKey:          "partition",
		ETag:                  "*",
		RowKey:                "row",
		Timestamp:             time.Now(),
		SomeBinaryProperty:    []byte("some bytes"),
		SomeDateProperty:      time.Now(),
		SomeDoubleProperty0:   float64(1),
		SomeDoubleProperty1:   float64(1.2345),
		SomeGuidProperty:      uuid.New(),
		SomeInt64Property:     math.MaxInt64,
		SomeIntProperty:       42,
		SomeStringProperty:    "some string",
		SomePtrStringProperty: &sp}
	return e
}

func closerFromString(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}

var odataHintProps = map[string]string{
	"Timestamp":           EdmDateTime,
	"SomeBinaryProperty":  EdmBinary,
	"SomeDateProperty":    EdmDateTime,
	"SomeDoubleProperty0": EdmDouble,
	"SomeDoubleProperty1": EdmDouble,
	"SomeGuidProperty":    EdmGuid,
	"SomeInt64Property":   EdmInt64}

const complexPayload = "{\"odata.metadata\": \"https://jverazsdkprim.table.core.windows.net/$metadata#testtableifprd13i\"," +
	"\"value\": [" +
	"{\"odata.etag\": \"W/\\u0022datetime\\u00272021-03-23T18%3A29%3A15.9686039Z\\u0027\\u0022\"," +
	"\"PartitionKey\": \"somPartition\"," +
	"\"RowKey\": \"01\"," +
	"\"Timestamp\": \"2021-03-23T18:29:15.9686039Z\"," +
	"\"SomeBinaryProperty@odata.type\": \"Edm.Binary\"," +
	"\"SomeBinaryProperty\": \"AQIDBAU=\"," +
	"\"SomeDateProperty@odata.type\": \"Edm.DateTime\"," +
	"\"SomeDateProperty\": \"2020-01-01T01:02:00Z\"," +
	"\"SomeDoubleProperty0\": 1.0," +
	"\"SomeDoubleProperty1\": 1.5," +
	"\"SomeGuidProperty@odata.type\": \"Edm.Guid\"," +
	"\"SomeGuidProperty\": \"0d391d16-97f1-4b9a-be68-4cc871f90001\"," +
	"\"SomeInt64Property@odata.type\": \"Edm.Int64\"," +
	"\"SomeInt64Property\": \"1\"," +
	"\"SomeIntProperty\": 1," +
	"\"SomeStringProperty\": \"This is table entity number 01\"  }]  }"

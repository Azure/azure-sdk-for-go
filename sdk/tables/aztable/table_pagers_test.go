// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func TestCastAndRemoveAnnotations(t *testing.T) {
	assert := assert.New(t)

	r := &http.Response{Body: closerFromString(complexPayload)}
	resp := azcore.Response{Response: r}

	var val map[string]interface{}
	err := resp.UnmarshalAsJSON(&val)
	assert.Nil(err)
	err = castAndRemoveAnnotations(&val)
	assert.Nil(err)
	// assert all odata annotations are removed.
	for k := range val {
		assert.NotContains(k, OdataType)
	}

	assert.IsType(time.Now(), val["SomeDateProperty"])
	assert.IsType([]byte{}, val["SomeBinaryProperty"])
	assert.IsType(float64(0), val["SomeDoubleProperty0"])
	// TODO: fix this
	// assert.IsType(int(0), (*entity)["SomeIntProperty"])
}

func TestToOdataAnnotatedDictionary(t *testing.T) {
	assert := assert.New(t)

	var val = createComplexEntityMap()
	err := toOdataAnnotatedDictionary(&val)
	assert.Nil(err)
	// assert all odata annotations are removed.
	for k := range odataHintProps {
		_, ok := val[k]
		assert.Truef(ok, fmt.Sprintf("map does not contain %s", k))
		iSuffix := strings.Index(k, OdataType)
		if iSuffix > 0 {
			// Get the name of the property that this odataType key describes.
			valueKey := k[0:iSuffix]
			if !strings.Contains(valueKey, "SomeDoubleProperty") {
				assert.IsTypef("", val[valueKey], fmt.Sprintf("should be type string %s", valueKey))
			}
		}
		_, ok = val[odataType(k)]
		assert.Truef(ok, fmt.Sprintf("map does not contain %s", odataType(k)))
	}
}

func BenchmarkUnMarshal_AsJson_CastAndRemove_Map(b *testing.B) {
	assert := assert.New(b)
	b.ReportAllocs()
	bt := []byte(complexPayload)
	for i := 0; i < b.N; i++ {
		var val = make(map[string]interface{})
		err := json.Unmarshal(bt, &val)
		if err != nil {
			panic(err)
		}
		castAndRemoveAnnotations(&val)
		assert.Equal("somePartition", val["PartitionKey"])
	}
}

func BenchmarkUnMarshal_FromMap_Entity(b *testing.B) {
	assert := assert.New(b)

	bt := []byte(complexPayload)
	for i := 0; i < b.N; i++ {
		var val = make(map[string]interface{})
		err := json.Unmarshal(bt, &val)
		if err != nil {
			panic(err)
		}
		result := complexEntity{}
		err = EntityMapAsModel(val, &result)
		assert.Nil(err)
		assert.Equal("somePartition", result.PartitionKey)
	}
}

func BenchmarkMarshal_Entity_ToMap_ToOdataDict_Map(b *testing.B) {
	ent := createComplexEntity()
	for i := 0; i < b.N; i++ {
		m, _ := toMap(ent)
		toOdataAnnotatedDictionary(m)
		_, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMarshal_Map_ToOdataDict_Map(b *testing.B) {
	ent := createComplexEntityMap()
	for i := 0; i < b.N; i++ {
		toOdataAnnotatedDictionary(&ent)
		_, err := json.Marshal(ent)
		if err != nil {
			panic(err)
		}
	}
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
	assert.Equal(base64.StdEncoding.EncodeToString(ent.SomeBinaryProperty), string((*entMap)["SomeBinaryProperty"].(string)))
	ts, _ := time.Parse(ISO8601, (*entMap)["SomeDateProperty"].(string))
	assert.Equal(ent.SomeDateProperty.UTC().Format(ISO8601), ts.Format(ISO8601))
	assert.Equal(ent.SomeDoubleProperty0, (*entMap)["SomeDoubleProperty0"])
	assert.Equal(ent.SomeDoubleProperty1, (*entMap)["SomeDoubleProperty1"])
	var u uuid.UUID = ent.SomeGuidProperty
	assert.Equal(u.String(), (*entMap)["SomeGuidProperty"].(string))
	assert.Equal(strconv.FormatInt(ent.SomeInt64Property, 10), (*entMap)["SomeInt64Property"].(string))
	assert.Equal(ent.SomeIntProperty, (*entMap)["SomeIntProperty"])
	assert.Equal(ent.SomeStringProperty, (*entMap)["SomeStringProperty"])
	assert.Equal(*ent.SomePtrStringProperty, (*entMap)["SomePtrStringProperty"])
}

func TestToMapWithMap(t *testing.T) {
	assert := assert.New(t)

	ent := createComplexEntityMap()

	entMap, err := toMap(ent)
	assert.Nil(err)

	// Validate that we have all the @odata.type properties for types []byte, int64, float64, time.Time, and uuid
	for k, v := range odataHintProps {
		vv, ok := (*entMap)[odataType(k)]
		assert.Truef(ok, "Should have found map key of name '%s'", odataType(k))
		assert.Equal(v, vv)
	}

	assert.Equal(&ent, entMap)
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

func TestDeserializeFromMap(t *testing.T) {
	assert := assert.New(t)

	expected := createComplexEntity()
	bt := []byte(complexPayload)
	var val = make(map[string]interface{})
	err := json.Unmarshal(bt, &val)
	assert.Nil(err)
	result := complexEntity{}
	// tt := reflect.TypeOf(complexEntity{})
	// err := fromMap(tt, getTypeValueMap(tt), &val, reflect.ValueOf(&result).Elem())
	err = EntityMapAsModel(val, &result)
	assert.Nil(err)
	assert.EqualValues(expected, result)
}

func createComplexEntity() complexEntity {
	sp := "some pointer to string"
	t, _ := time.Parse(ISO8601, "2021-03-23T18:29:15.9686039Z")
	t2, _ := time.Parse(ISO8601, "2020-01-01T01:02:00Z")
	b, _ := base64.StdEncoding.DecodeString("AQIDBAU=")
	var e = complexEntity{
		PartitionKey:          "somePartition",
		ETag:                  "W/\"datetime'2021-04-05T05%3A02%3A40.7371784Z'\"",
		RowKey:                "01",
		Timestamp:             t,
		SomeBinaryProperty:    b,
		SomeDateProperty:      t2,
		SomeDoubleProperty0:   float64(1.0),
		SomeDoubleProperty1:   float64(1.5),
		SomeGuidProperty:      uuid.Parse("0d391d16-97f1-4b9a-be68-4cc871f90001"),
		SomeInt64Property:     int64(math.MaxInt64),
		SomeIntProperty:       42,
		SomeStringProperty:    "This is table entity number 01",
		SomePtrStringProperty: &sp}
	return e
}

func createComplexEntityMap() map[string]interface{} {
	sp := "some pointer to string"
	t, _ := time.Parse(ISO8601, "2021-03-23T18:29:15.9686039Z")
	t2, _ := time.Parse(ISO8601, "2020-01-01T01:02:00Z")
	b, _ := base64.StdEncoding.DecodeString("AQIDBAU=")
	var e = map[string]interface{}{
		"PartitionKey":          "somePartition",
		"ETag":                  "W/\"datetime'2021-04-05T05%3A02%3A40.7371784Z'\"",
		"RowKey":                "01",
		"Timestamp":             t,
		"SomeBinaryProperty":    b,
		"SomeDateProperty":      t2,
		"SomeDoubleProperty0":   float64(1.0),
		"SomeDoubleProperty1":   float64(1.5),
		"SomeGuidProperty":      uuid.Parse("0d391d16-97f1-4b9a-be68-4cc871f90001"),
		"SomeInt64Property":     int64(math.MaxInt64),
		"SomeIntProperty":       42,
		"SomeStringProperty":    "This is table entity number 01",
		"SomePtrStringProperty": &sp}
	return e
}

func closerFromString(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}

var odataHintProps = map[string]string{
	"SomeBinaryProperty":  edmBinary,
	"SomeDateProperty":    edmDateTime,
	"SomeDoubleProperty0": edmDouble,
	"SomeDoubleProperty1": edmDouble,
	"SomeGuidProperty":    edmGuid,
	"SomeInt64Property":   edmInt64}

var complexPayload = "{\"odata.etag\": \"W/\\\"datetime'2021-04-05T05%3A02%3A40.7371784Z'\\\"\"," +
	"\"PartitionKey\": \"somePartition\"," +
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
	"\"SomeInt64Property\": \"" + strconv.FormatInt(math.MaxInt64, 10) + "\"," +
	"\"SomeIntProperty\": 42," +
	"\"SomeStringProperty\": \"This is table entity number 01\"," +
	"\"SomePtrStringProperty\": \"some pointer to string\"  }"

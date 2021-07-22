// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkUnMarshal_AsJson_CastAndRemove_Map(b *testing.B) {
	assert := assert.New(b)
	b.ReportAllocs()
	bt := []byte(complexPayload)
	for i := 0; i < b.N; i++ {
		var val = make(map[string]interface{})
		err := json.Unmarshal(bt, &val)
		assert.Nil(err)
		err = castAndRemoveAnnotations(&val)
		assert.Nil(err)
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

func check(e error) {
	if e != nil {
		panic(e)
	}
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

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

	assert.IsType((*entity)["SomeDateProperty"], time.Now())
}

func closerFromString(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}

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

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

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

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azlogs_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"
	"github.com/stretchr/testify/require"
)

type ComputerInfo struct {
	Time              time.Time
	Computer          string
	AdditionalContext AdditionalContext
}

type AdditionalContext struct {
	TestContextKey int
	CounterName    string
}

func generateLogs() []byte {
	var data []ComputerInfo
	for i := 0; i < 10; i++ {
		data = append(data, ComputerInfo{
			Time:              time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC),
			Computer:          "Computer" + strconv.Itoa(i),
			AdditionalContext: AdditionalContext{TestContextKey: i, CounterName: "AppMetric2"},
		})
	}
	data2, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return data2
}

func TestUpload(t *testing.T) {
	client := startTest(t)

	logs := generateLogs()

	res, err := client.Upload(context.Background(), ruleID, streamName, logs, nil)
	require.NoError(t, err)
	require.Empty(t, res)
}

func TestUploadWithGzip(t *testing.T) {
	client := startTest(t)

	logs := generateLogs()

	// gzip data
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(logs)
	require.NoError(t, err)
	err = zw.Close()
	require.NoError(t, err)

	res, err := client.Upload(context.Background(), ruleID, streamName, buf.Bytes(), &azlogs.UploadOptions{ContentEncoding: to.Ptr("gzip")})
	require.NoError(t, err)
	require.Empty(t, res)
}

func TestUploadWithError(t *testing.T) {
	client := startTest(t)

	logs := generateLogs()

	res, err := client.Upload(context.Background(), ruleID, "incorrect stream", logs, nil)
	require.Error(t, err)
	require.Empty(t, res)

	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "InvalidStream")
	require.Equal(t, httpErr.StatusCode, 400)
}

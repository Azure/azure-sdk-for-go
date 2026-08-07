// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testcommon

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	RecordingDirectory          = "sdk/storage/azdatalake/testdata"
	FileSystemPrefix            = "gofs"
	FilePrefix                  = "gotestfile"
	DirPrefix                   = "gotestdir"
	SubDirPrefix                = "gotestsubdir"
	DefaultData                 = "Godatalakedata"
	InvalidHeaderErrorSubstring = "invalid header field" // error thrown by the http client
)

func GenerateFileSystemName(testName string) string {
	return FileSystemPrefix + GenerateEntityName(testName)
}

func GenerateFileName(testName string) string {
	return FilePrefix + GenerateEntityName(testName)
}

func GenerateDirName(testName string) string {
	return DirPrefix + GenerateEntityName(testName)
}

func GenerateSubDirName(testName string) string {
	return SubDirPrefix + GenerateEntityName(testName)
}

func GenerateEntityName(testName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(testName), "/", ""), "test", "")
}

func SetupSuite(suite *suite.Suite) *recording.TestProxyInstance {
	proxy, err := recording.StartTestProxy(RecordingDirectory, nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	return proxy
}

func TearDownSuite(suite *suite.Suite, proxy *recording.TestProxyInstance) {
	err := recording.StopTestProxy(proxy)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func BeforeTest(t *testing.T, suite string, test string) {
	const blobURLRegex = `https://\S+\.blob\.core\.windows\.net`
	const dfsURLRegex = `https://\S+\.dfs\.core\.windows\.net`
	const tokenRegex = `(?:Bearer\s).*`

	require.NoError(t, recording.AddURISanitizer(FakeBlobStorageURL, blobURLRegex, nil))
	require.NoError(t, recording.AddURISanitizer(FakeDFSStorageURL, dfsURLRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source", FakeBlobStorageURL, blobURLRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source", FakeDFSStorageURL, dfsURLRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source-authorization", FakeToken, tokenRegex, nil))
	// we freeze request IDs and timestamps to avoid creating noisy diffs
	// NOTE: we can't freeze time stamps as that breaks some tests that use if-modified-since etc (maybe it can be fixed?)
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-request-id", "00000000-0000-0000-0000-000000000000", "", nil))
	// TODO: more freezing
	// testframework.AddHeaderRegexSanitizer("X-Ms-Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-request-id", "00000000-0000-0000-0000-000000000000", "", nil))
	// testframework.AddHeaderRegexSanitizer("Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	// TODO: more freezing
	// testframework.AddBodyRegexSanitizer("RequestId:00000000-0000-0000-0000-000000000000", `RequestId:\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`, nil)
	// testframework.AddBodyRegexSanitizer("Time:2022-08-11T00:21:56.4562741Z", `Time:\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d*)?Z`, nil)
	require.NoError(t, recording.Start(t, RecordingDirectory, nil))
}

func AfterTest(t *testing.T, suite string, test string) {
	require.NoError(t, recording.Stop(t, nil))
}

// GetRequiredEnv gets an environment variable by name and returns an error if it is not found
func GetRequiredEnv(name string) (string, error) {
	env, ok := os.LookupEnv(name)
	if ok {
		return env, nil
	} else {
		return "", errors.New("Required environment variable not set: " + name)
	}
}

func ValidateErrorCode(_require *require.Assertions, err error, code datalakeerror.StorageErrorCode) {
	_require.Error(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(string(code), responseErr.ErrorCode)
	} else {
		_require.Contains(err.Error(), code)
	}
}

func GetRelativeTimeFromAnchor(anchorTime *time.Time, amount time.Duration) time.Time {
	return anchorTime.Add(amount * time.Second)
}

const random64BString string = "2SDgZj6RkKYzJpu04sweQek4uWHO8ndPnYlZ0tnFS61hjnFZ5IkvIGGY44eKABov"

func GenerateData(sizeInBytes int) (io.ReadSeekCloser, []byte) {
	data := make([]byte, sizeInBytes)
	_len := len(random64BString)
	if sizeInBytes > _len {
		count := sizeInBytes / _len
		if sizeInBytes%_len != 0 {
			count++
		}
		copy(data, strings.Repeat(random64BString, count))
	} else {
		copy(data, random64BString)
	}
	return streaming.NopCloser(bytes.NewReader(data)), data
}

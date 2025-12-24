// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Contains common helpers for TESTS ONLY
package testcommon

import (
	"bytes"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

const (
	RecordingDirectory = "sdk/storage/azfile/testdata"
	SharePrefix        = "gos"
	DirectoryPrefix    = "godir"
	FilePrefix         = "gotestfile"
	FileDefaultData    = "GoFileDefaultData"
)

func GenerateShareName(testName string) string {
	return SharePrefix + GenerateEntityName(testName)
}

func GenerateEntityName(testName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(testName), "/", ""), "test", "")
}

func GenerateDirectoryName(testName string) string {
	return DirectoryPrefix + GenerateEntityName(testName)
}

func GenerateFileName(testName string) string {
	return FilePrefix + GenerateEntityName(testName)
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

func ValidateHTTPErrorCode(_require *require.Assertions, err error, code int) {
	_require.Error(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(responseErr.StatusCode, code)
	} else {
		_require.Equal(strings.Contains(err.Error(), strconv.Itoa(code)), true)
	}
}

func ValidateFileErrorCode(_require *require.Assertions, err error, code fileerror.Code) {
	_require.Error(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(string(code), responseErr.ErrorCode)
	} else {
		_require.Contains(err.Error(), code)
	}
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
	const urlRegex = `https://\S+\.file\.core\.windows\.net`
	const tokenRegex = `(?:Bearer\s).*`

	require.NoError(t, recording.AddURISanitizer(FakeStorageURL, urlRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source", FakeStorageURL, urlRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source-authorization", FakeToken, tokenRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-file-rename-source", FakeStorageURL, urlRegex, nil))
	// we freeze request IDs and timestamps to avoid creating noisy diffs
	// NOTE: we can't freeze time stamps as that breaks some tests that use if-modified-since etc (maybe it can be fixed?)
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-request-id", "00000000-0000-0000-0000-000000000000", "", nil))
	// TODO: more freezing
	require.NoError(t, recording.Start(t, RecordingDirectory, nil))
}

func AfterTest(t *testing.T, suite string, test string) {
	require.NoError(t, recording.Stop(t, nil))
}
